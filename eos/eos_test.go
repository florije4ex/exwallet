package eos

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	eosGo "github.com/eoscanada/eos-go"
	"io"
	"net/http"
	"testing"
)

func TestGetEosAddress(t *testing.T) {
	address, err := GetEosAddress()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(address)
}

func TestGetPubKeyByPrvKey(t *testing.T) {
	address, err := GetPubKeyByPrvKey("5JbSYC6jwyEpz81BiPbBvi9hcoDpg5gsvVNxsvEuMDmyJPMk4FC")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(address == "{\"privateKey\":\"5JbSYC6jwyEpz81BiPbBvi9hcoDpg5gsvVNxsvEuMDmyJPMk4FC\",\"publicKey\":\"EOS8gyyQyiwhpu6Q9PiQSn6LWiA9XdVv9PZW4qiqJ55u9oomfhqrp\"}")
}

func TestSignEosTx(t *testing.T) {
	contentType := "application/json"
	// get abi encoded transfer data
	abiUrl := "http://127.0.0.1:8888/v1/chain/abi_json_to_bin"
	type transfer struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Quantity string `json:"quantity"`
		Memo     string `json:"memo"`
	}
	type transferAction struct {
		Action string   `json:"action"`
		Code   string   `json:"code"`
		Args   transfer `json:"args"`
	}
	transferData := transferAction{
		Code:   "eosio.token",
		Action: "transfer",
		Args: transfer{
			From:     "eosio.token",
			To:       "inita",
			Quantity: "1.0000 SYS",
			Memo:     "transfer 1 EOS from eosio.token to inita",
		},
	}
	transferJson, err := json.Marshal(transferData)
	if err != nil {
		t.Error(err.Error())
	}
	abiResp, err := http.Post(abiUrl, contentType, bytes.NewBuffer(transferJson))
	if err != nil {
		t.Error(err.Error())
	}
	var abiCnt bytes.Buffer
	_, err = io.Copy(&abiCnt, abiResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	//t.Log(abiCnt.String())
	var abi eosGo.ABIJSONToBinResp
	err = json.Unmarshal(abiCnt.Bytes(), &abi)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(abi.Binargs)
	// get info(head_block_id)
	infoUrl := "http://127.0.0.1:8888/v1/chain/get_info"
	infoResp, err := http.Post(infoUrl, contentType, nil)
	if err != nil {
		t.Error(err.Error())
	}
	var infoCnt bytes.Buffer
	_, err = io.Copy(&infoCnt, infoResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var info eosGo.InfoResp
	err = json.Unmarshal(infoCnt.Bytes(), &info)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(info.HeadBlockID)
	t.Log(info.ChainID)
	// construct tx struct
	eosTx := Tx{
		Action:      "transfer",
		ActionData:  abi.Binargs,
		Actor:       "eosio.token",
		ChainId:     hex.EncodeToString(info.ChainID),
		Contract:    "eosio.token",
		Expiration:  120,
		HeadBlockId: hex.EncodeToString(info.HeadBlockID),
		Permission:  "active",
		PrivateKey:  "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3",
	}
	signedTx, err := SignTx(eosTx)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(signedTx)
	pushUrl := "http://127.0.0.1:8888/v1/chain/push_transaction"
	pushResp, err := http.Post(pushUrl, contentType, bytes.NewBuffer([]byte(signedTx)))
	if err != nil {
		t.Error(err.Error())
	}
	var pushCnt bytes.Buffer
	_, err = io.Copy(&pushCnt, pushResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	if pushResp.StatusCode == http.StatusOK || pushResp.StatusCode == http.StatusAccepted {
		var push eosGo.PushTransactionFullResp
		err = json.Unmarshal(pushCnt.Bytes(), &push)
		if err != nil {
			t.Error(err.Error())
		}
		t.Log(push.StatusCode)
		t.Log(push.Processed.Status)
		t.Log(push.TransactionID)
	} else {
		t.Error(pushResp.Status)
		t.Error(pushCnt.String())
	}
}
