package usdt

import "testing"

func TestSignUsdtTx(t *testing.T) {
	signedTx, err := SignUsdtTx("{\"inputs\": [{\"txHash\": \"b282c0320aff1a3c8616987659ed656e4a57a71feadb467191d52fbe84330cc9\", \"voutIndex\": 1, \"privateKey\": \"cSCsRY5xjq1vHEkdFN2QLzSDAJEhigPQitD9CuGSv64JNsuC3EFc\", \"scriptPubKey\": \"76a91499dcf88a3c82edcb39b25f5bcf518cc967a9ab1388ac\"}], \"amount\": 100, \"outputs\": [{\"amount\": 500000000, \"address\": \"mmVmF9ayNyLFwZ6imW6CvywgX8wzGoHGhU\"}, {\"amount\": 546, \"address\": \"muYWPjD54vjYdCz6amkXha92EJWVQ7aSC1\"}], \"address\": \"muYWPjD54vjYdCz6amkXha92EJWVQ7aSC1\"}")
	if err != nil {
		t.Error(err)
	}
	t.Log(signedTx)
}
