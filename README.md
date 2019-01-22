# Usage

This is a project for crypto currency development, the main language is golang.

## before
```text
1. set GOPATH
2. set GoLand vgo Proxy: https://goproxy.io or just export GOPROXY=https://goproxy.io
then sync packages will install dependence in $GOPATH/pkg
```

## bitcoin
```text
cd bitcoin
docker-compose up -d
test function TestSignTx

links:
https://bitcoincore.org/en/doc/0.16.3/
http://chainquery.com/bitcoin-api
```

## litecoin
```text
https://download.litecoin.org/litecoin-0.16.0/linux/litecoin-0.16.0-x86_64-linux-gnu.tar.gz

curl -X POST http://localhost:8545 -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
```