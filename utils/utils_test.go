package utils

import "testing"

func TestGenMnemonic(t *testing.T) {
	mnemonic, err := GetMnemonic()
	if err != nil {
		t.Error(err)
	}
	t.Log(mnemonic)
}

func TestEncrypt(t *testing.T) {
	encryptedData, err := Encrypt("Hello, Wallet!", "helloworld")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(encryptedData)
}

func TestDecrypt(t *testing.T) {
	decryptedData, err := Decrypt("bf1a86482c9b7e646eb87dfdff8a53b857aa5817fa7b67528bd5bff2815f0030e528d105951d358b3001", "helloworld")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(decryptedData)
}

func TestEncryptFile(t *testing.T) {
	result, err := EncryptFile("encrypted.dat", "Hello, Wallet!", "helloworld")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(result)
}

func TestDecryptFile(t *testing.T) {
	result, err := DecryptFile("encrypted.dat", "helloworlds")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(result)
}
