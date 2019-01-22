package btc

import (
	"exwallet/utils"
	"github.com/shopspring/decimal"
	"testing"
)

func TestGetHdAddress(t *testing.T) {
	mnemonic, err := utils.GetMnemonic()
	if err != nil {
		t.Error(err)
	}
	t.Log(mnemonic)
	address, err := GetHdAddress(mnemonic, 0)
	if err != nil {
		t.Error(err)
	}
	t.Log(address)
}

func TestGetAddress(t *testing.T) {
	address, err := GetAddress()
	if err != nil {
		t.Error(err)
	}
	t.Log(address)
}

func TestGetAddrByPrvKey(t *testing.T) {
	// {"privateKey":"cSdLmMAHejoeEJmLxDojjaasftu5P8qvp3Gi6R82zvx9ZGgDGhMQ","address":"mpS3FS68vr3QhMC76Xr8t8rubMfRQemWiw"}
	address, err := GetAddrByPrvKey("cSdLmMAHejoeEJmLxDojjaasftu5P8qvp3Gi6R82zvx9ZGgDGhMQ")
	if err != nil {
		t.Error(err)
	}
	t.Log(address)
}

func TestSignTx(t *testing.T) {
	type jsonRequest struct {
		// {"jsonrpc": "2.0", "id": 1, "method": "generate", "params": [100] }
		Id      uint64        `json:"id"`
		Method  string        `json:"method"`
		JsonRpc string        `json:"jsonrpc"`
		Params  []interface{} `json:"params"`
	}
	type jsonResponse struct {
		// {"result":["73f2fa7c3043e547d9364b7e4211d593607cef11e42f1ecad9150afe167e8e16"],"error":null,"id":1}
		Result []interface{} `json:"result"`
		Error  string        `json:"error"`
		Id     uint64        `json:"id"`
	}
	type unspentItem struct {
		Txid          string          `json:"txid"`
		Vout          uint64          `json:"vout"`
		Address       string          `json:"address"`
		Account       string          `json:"account"`
		ScriptPubKey  string          `json:"scriptPubKey"`
		Amount        decimal.Decimal `json:"amount"`
		Confirmations uint64          `json:"confirmations"`
		Spendable     bool            `json:"spendable"`
		Solvable      bool            `json:"solvable"`
		Safe          bool            `json:"safe"`
	}
	type unspentResponse struct {
		Result []unspentItem `json:"result"`
		Error  string        `json:"error"`
		Id     uint64        `json:"id"`
	}
	type txResponse struct {
		Result string `json:"result"`
		Error  string `json:"error"`
		Id     uint64 `json:"id"`
	}
	rpcUrl := "http://localhost:8339"
	// get address and private key
	fromAddress, err := GetAddress()
	if err != nil {
		t.Error(err)
	}
	t.Log(fromAddress)
	toAddress, err := GetAddress()
	if err != nil {
		t.Error(err)
	}
	t.Log(toAddress)
	importkeyData := jsonRequest{Id: 1, Method: "importprivkey", JsonRpc: "2.0", Params: []interface{}{fromAddress.PrivateKey}}
	importkeyResp := jsonResponse{}
	err = utils.PostData(rpcUrl, importkeyData, &importkeyResp)
	if err != nil {
		t.Error(err.Error())
	}
	generateData := jsonRequest{Id: 1, Method: "generatetoaddress", JsonRpc: "2.0", Params: []interface{}{101, fromAddress.Address}}
	generateResp := jsonResponse{}
	err = utils.PostData(rpcUrl, generateData, &generateResp)
	if err != nil {
		t.Error(err.Error())
	}
	unspentData := jsonRequest{Id: 1, Method: "listunspent", JsonRpc: "2.0", Params: []interface{}{1, 9999999, []string{fromAddress.Address}}}
	unspentResp := unspentResponse{}
	err = utils.PostData(rpcUrl, unspentData, &unspentResp)
	if err != nil {
		t.Error(err.Error())
	}
	if len(unspentResp.Result) < 1 {
		t.Error("unspent is empty!")
	}
	unspent := unspentResp.Result[0]
	sendAmount := int64(100000000)
	feeAmount := int64(10000)
	changeAmount := unspent.Amount.Mul(decimal.New(1, 8)).Sub(decimal.New(sendAmount, 0)).Sub(decimal.New(feeAmount, 0)).IntPart()
	tx := Tx{
		Inputs:  []TxInput{{TxHash: unspent.Txid, VoutIndex: unspent.Vout, ScriptPubKey: unspent.ScriptPubKey, PrivateKey: fromAddress.PrivateKey}},
		Outputs: []TxOutput{{Address: fromAddress.Address, Amount: changeAmount}, {Address: toAddress.Address, Amount: sendAmount}},
	}
	t.Log(tx)
	signedTx, err := SignTx(tx)
	t.Log(signedTx)
	sendRawTxData := jsonRequest{Id: 1, Method: "sendrawtransaction", JsonRpc: "2.0", Params: []interface{}{signedTx}}
	sendRawTxResp := txResponse{}
	err = utils.PostData(rpcUrl, sendRawTxData, &sendRawTxResp)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(sendRawTxResp.Id, sendRawTxResp.Error, sendRawTxResp.Result)
}
