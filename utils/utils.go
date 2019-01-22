package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"github.com/bartekn/go-bip39"
	"github.com/btcsuite/btcd/chaincfg"
	"io"
	"io/ioutil"
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
