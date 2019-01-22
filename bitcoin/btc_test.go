package btc

import (
	"exwallet/utils"
	"testing"
)

func TestGetHdBtcAddress(t *testing.T) {
	mnemonic, err := utils.GetMnemonic()
	if err != nil {
		t.Error(err)
	}
	t.Log(mnemonic)
	address, err := GetHdBtcAddress(mnemonic, 0)
	if err != nil {
		t.Error(err)
	}
	t.Log(address)
}

func TestGetBtcAddress(t *testing.T) {
	address, err := GetBtcAddress()
	if err != nil {
		t.Error(err)
	}
	t.Log(address)
}

func TestSignBtcTx(t *testing.T) {
	//cT2sMcN7iitjtK1fZ6EvfbE8omyQ5ZUXEewrDv5AjAgq9kd7Ra7N mw1ASgotu94p8NDMCmfyy8wctvFUcvABy1
	//signedTx, err := SignBtcTx("{\"inputs\":[{\"privateKey\":\"cSCsRY5xjq1vHEkdFN2QLzSDAJEhigPQitD9CuGSv64JNsuC3EFc\",\"scriptPubKey\":\"76a91499dcf88a3c82edcb39b25f5bcf518cc967a9ab1388ac\",\"txHash\":\"b282c0320aff1a3c8616987659ed656e4a57a71feadb467191d52fbe84330cc9\",\"voutIndex\":1}],\"outputs\":[{\"address\":\"mmVmF9ayNyLFwZ6imW6CvywgX8wzGoHGhU\",\"amount\":499990000},{\"address\":\"muYWPjD54vjYdCz6amkXha92EJWVQ7aSC1\",\"amount\":500000000}]}")
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(signedTx)

	// get key
	address, err := GetBtcAddress()
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
