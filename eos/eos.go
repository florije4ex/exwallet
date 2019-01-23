package eos

import (
	"encoding/hex"
	"encoding/json"
	eosGo "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"time"
)

type Account struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

type Tx struct {
	Action      string `json:"action"`
	ActionData  string `json:"actionData"`
	Actor       string `json:"actor"`
	ChainId     string `json:"chainId"`
	Contract    string `json:"contract"`
	Expiration  int    `json:"expiration"`
	HeadBlockId string `json:"headBlockId"`
	Permission  string `json:"permission"`
	PrivateKey  string `json:"privateKey"`
}

func GetEosAddress() (string, error) {
	key, err := ecc.NewRandomPrivateKey()
	if err != nil {
		return "", err
	}
	jsonRes, err := json.Marshal(Account{PrivateKey: key.String(), PublicKey: key.PublicKey().String()})
	if err != nil {
		return "", err
	}
	return string(jsonRes), nil
}

func GetPubKeyByPrvKey(prvKey string) (string, error) {
	key, err := ecc.NewPrivateKey(prvKey)
	if err != nil {
		return "", err
	}
	jsonRes, err := json.Marshal(Account{PrivateKey: key.String(), PublicKey: key.PublicKey().String()})
	if err != nil {
		return "", err
	}
	return string(jsonRes), nil
}

func SignTx(eosTx Tx) (string, error) {
	actionData, err := hex.DecodeString(eosTx.ActionData)
	if err != nil {
		return "", err
	}
	actions := []*eosGo.Action{
		{
			Account: eosGo.AN(eosTx.Contract),
			Name:    eosGo.ActN(eosTx.Action),
			Authorization: []eosGo.PermissionLevel{
				{
					Actor: eosGo.AN(eosTx.Actor), Permission: eosGo.PN(eosTx.Permission),
				},
			},
			ActionData: eosGo.NewActionDataFromHexData(actionData),
		},
	}
	chainId, err := hex.DecodeString(eosTx.ChainId)
	if err != nil {
		return "", err
	}
	headBlockId, err := hex.DecodeString(eosTx.HeadBlockId)
	if err != nil {
		return "", err
	}

	opts := eosGo.TxOptions{ChainID: eosGo.Checksum256(chainId), HeadBlockID: eosGo.Checksum256(headBlockId)}
	tx := eosGo.NewTransaction(actions, &opts)
	tx.SetExpiration(time.Duration(time.Duration(eosTx.Expiration) * time.Second))

	signedTx := eosGo.NewSignedTransaction(tx)
	rawTx, cfd, err := signedTx.PackedTransactionAndCFD()
	if err != nil {
		return "", err
	}
	sigDigest := eosGo.SigDigest(chainId, rawTx, cfd)

	prvKey, err := ecc.NewPrivateKey(eosTx.PrivateKey)
	if err != nil {
		return "", err
	}
	sig, err := prvKey.Sign(sigDigest)
	if err != nil {
		return "", err
	}
	signedTx.Signatures = append(signedTx.Signatures, sig)
	packed, err := signedTx.Pack(eosGo.CompressionNone)
	if err != nil {
		return "", err
	}
	packedTx, err := json.Marshal(packed)
	if err != nil {
		return "", err
	}
	return string(packedTx), nil
}
