package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bartekn/go-bip39"
	"github.com/btcsuite/btcd/chaincfg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// btc
var Net = &chaincfg.MainNetParams

// usdt
var PropertyId = 31

// eth
var Transfer = "transfer(address,uint256)"

func GetMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(data string, passphrase string) (string, error) {
	block, err := aes.NewCipher([]byte(createHash(passphrase)))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := gcm.Seal(nonce, nonce, []byte(data), nil)
	return hex.EncodeToString(cipherText), nil
}

func Decrypt(encryptedData string, passphrase string) (string, error) {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	data, err := hex.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func EncryptFile(filename string, data string, passphrase string) (bool, error) {
	f, err := os.Create(filename)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err.Error())
		}
	}()
	encryptedData, err := Encrypt(data, passphrase)
	if err != nil {
		return false, err
	}
	_, err = f.Write([]byte(encryptedData))
	if err != nil {
		return false, err
	}
	return true, nil
}

func DecryptFile(filename string, passphrase string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	decryptedData, err := Decrypt(string(data), passphrase)
	if err != nil {
		return "", err
	}
	return decryptedData, nil
}

func GetData(url string, resu interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	return creatResult(res, resu)
}

func PostData(url string, postData interface{}, resu interface{}) error {
	var data []byte
	var err error
	if postData != nil {
		data, err = json.Marshal(postData)
		if err != nil {
			return err
		}
	}
	bodyReader := bytes.NewBuffer(data)
	res, err := http.Post(url, "application/json", bodyReader)
	if err != nil {
		return err
	}
	return creatResult(res, resu)
}

func creatResult(resp *http.Response, resu interface{}) error {
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("statusCode :%v ; body :%v", resp.StatusCode, string(body))
	}
	if resu == nil {
		return nil
	}
	return json.Unmarshal(body, resu)
}
