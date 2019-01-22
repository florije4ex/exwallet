package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"exwallet/utils"
	"github.com/bartekn/go-bip39"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
)

type Address struct {
	PrivateKey string `json:"privateKey"`
	Address    string `json:"address"`
}

type Tx struct {
	Inputs  []TXInput  `json:"inputs"`
	Outputs []TXOutput `json:"outputs"`
}

type TXInput struct {
	PrivateKey   string `json:"privateKey"`
	ScriptPubKey string `json:"scriptPubKey"`
	TxHash       string `json:"txHash"`
	VoutIndex    int32  `json:"voutIndex"`
}

type TXOutput struct {
	Address string `json:"address"`
	Amount  int64  `json:"amount"`
}

func GetBtcAddress() (string, error) {
	prvKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", err
	}
	wif, err := btcutil.NewWIF(prvKey, utils.Net, true)
	if err != nil {
		return "", err
	}
	pubKeySerial := prvKey.PubKey().SerializeCompressed()
	addressPubKey, err := btcutil.NewAddressPubKey(pubKeySerial, utils.Net)
	if err != nil {
		return "", err
	}
	address := Address{PrivateKey: wif.String(), Address: addressPubKey.EncodeAddress()}
	jsonRes, err := json.Marshal(address)
	if err != nil {
		return "", err
	}
	return string(jsonRes), nil
}

func GetAddrByPrvKey(prvKey string) (string, error) {
	wif, err := btcutil.DecodeWIF(prvKey)
	if err != nil {
		return "", err
	}
	pubKeySerial := wif.PrivKey.PubKey().SerializeCompressed()
	addressPubKey, err := btcutil.NewAddressPubKey(pubKeySerial, utils.Net)
	if err != nil {
		return "", err
	}
	address := Address{PrivateKey: wif.String(), Address: addressPubKey.EncodeAddress()}
	jsonRes, err := json.Marshal(address)
	if err != nil {
		return "", err
	}
	return string(jsonRes), nil
}

func GetHdBtcAddress(mnemonic string, index int32) (string, error) {
	seed := bip39.NewSeed(mnemonic, "")
	root, err := hdkeychain.NewMaster(seed, utils.Net)
	// m/44'/0'/0'/0/0
	purpose, err := root.Child(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		return "", err
	}
	coinType, err := purpose.Child(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return "", err
	}
	account, err := coinType.Child(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return "", err
	}
	change, err := account.Child(0)
	if err != nil {
		return "", err
	}
	prvExKey, err := change.Child(uint32(index))
	if err != nil {
		return "", err
	}
	prvKey, err := prvExKey.ECPrivKey()
	if err != nil {
		return "", err
	}
	wif, err := btcutil.NewWIF(prvKey, utils.Net, true)
	if err != nil {
		return "", err
	}
	pubKeySerial := prvKey.PubKey().SerializeCompressed()
	addressPubKey, err := btcutil.NewAddressPubKey(pubKeySerial, utils.Net)
	if err != nil {
		return "", err
	}
	address := Address{PrivateKey: wif.String(), Address: addressPubKey.EncodeAddress()}
	jsonRes, err := json.Marshal(address)
	if err != nil {
		return "", err
	}
	return string(jsonRes), nil
}

func SignBtcTx(txData string) (string, error) {
	var btcTx Tx
	err := json.Unmarshal([]byte(txData), &btcTx)
	if err != nil {
		return "", err
	}

	var inputs []*wire.TxIn
	for _, input := range btcTx.Inputs {
		hash, err := chainhash.NewHashFromStr(input.TxHash)
		if err != nil {
			return "", err
		}
		outPoint := wire.NewOutPoint(hash, uint32(input.VoutIndex))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		inputs = append(inputs, txIn)
	}

	var outputs []*wire.TxOut
	for _, output := range btcTx.Outputs {
		address, err := btcutil.DecodeAddress(output.Address, utils.Net)
		if err != nil {
			return "", err
		}
		pkScript, err := txscript.PayToAddrScript(address)
		if err != nil {
			return "", err
		}
		outputs = append(outputs, wire.NewTxOut(output.Amount, pkScript))
	}

	tx := &wire.MsgTx{
		Version:  wire.TxVersion,
		TxIn:     inputs,
		TxOut:    outputs,
		LockTime: 0,
	}

	for i, input := range btcTx.Inputs {
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
