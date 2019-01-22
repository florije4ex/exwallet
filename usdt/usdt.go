package usdt

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

// https://github.com/OmniLayer/omnicore/blob/master/src/omnicore/doc/rpc-api.md#omni_sendissuancefixed
// https://steemit.com/usdt/@chaimyu/omni-usdt-raw-transaction
// https://github.com/OmniLayer/omnicore/wiki/Use-the-raw-transaction-API-to-create-a-Simple-Send-transaction
// https://omniexplorer.info/asset/31
// https://github.com/OmniLayer/spec#field-number-of-coins

var PropertyId = 31
var Network = &chaincfg.TestNet3Params

type TxInput struct {
	PrivateKey   string `json:"privateKey"`
	ScriptPubKey string `json:"scriptPubKey"`
	TxHash       string `json:"txHash"`
	VoutIndex    int32  `json:"voutIndex"`
}

type TxOutput struct {
	Address string `json:"address"`
	Amount  int64  `json:"amount"`
}

type Tx struct {
	Address string     `json:"address"`
	Amount  int64      `json:"amount"`
	Inputs  []TxInput  `json:"inputs"`
	Outputs []TxOutput `json:"outputs"`
}

func getPayload(value int64) string {
	prefix := "omni"
	txVersion := 0
	txType := 0
	payload := fmt.Sprintf("%x%04x%04x%08x%016x", prefix, txVersion, txType, PropertyId, value)
	return payload
}

func SignUsdtTx(txData string) (string, error) {
	var usdtTx Tx
	err := json.Unmarshal([]byte(txData), &usdtTx)
	if err != nil {
		return "", err
	}

	var inputs []*wire.TxIn
	for _, input := range usdtTx.Inputs {
		hash, err := chainhash.NewHashFromStr(input.TxHash)
		if err != nil {
			return "", err
		}
		outPoint := wire.NewOutPoint(hash, uint32(input.VoutIndex))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		inputs = append(inputs, txIn)
	}

	var outputs []*wire.TxOut
	for _, output := range usdtTx.Outputs {
		address, err := btcutil.DecodeAddress(output.Address, Network)
		if err != nil {
			return "", err
		}
		pkScript, err := txscript.PayToAddrScript(address)
		if err != nil {
			return "", err
		}
		outputs = append(outputs, wire.NewTxOut(output.Amount, pkScript))
	}

	// add omni payload
	payload := getPayload(usdtTx.Amount)
	pkScript, err := txscript.NullDataScript([]byte(payload))
	if err != nil {
		return "", err
	}
	outputs = append(outputs, wire.NewTxOut(int64(0), pkScript))

	tx := &wire.MsgTx{
		Version:  wire.TxVersion,
		TxIn:     inputs,
		TxOut:    outputs,
		LockTime: 0,
	}

	for i, input := range usdtTx.Inputs {
		pkScript, _ := hex.DecodeString(input.ScriptPubKey)
		wif, err := btcutil.DecodeWIF(input.PrivateKey)
		if err != nil {
			return "", err
		}
		script, err := txscript.SignatureScript(tx, i, pkScript, txscript.SigHashAll, wif.PrivKey, true)
		if err != nil {
			return "", err
		}
		inputs[i].SignatureScript = script
	}

	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
	}
	txHex := hex.EncodeToString(buf.Bytes())

	return txHex, nil
}
