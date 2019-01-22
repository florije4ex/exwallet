package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

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

type Data struct {
	Name string `json:"name"`
}

func parseHandler(w http.ResponseWriter, r *http.Request) {
	location := Data{}
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading the body", err)
	}
	err = json.Unmarshal(jsn, &location)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}
	log.Printf("Received: %v\n", location)
	locationJson, err := json.Marshal(location)
	if err != nil {
		_, err := fmt.Fprintf(w, "Error: %s", err)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(locationJson)
	if err != nil {
		log.Fatal(err.Error())
	}

}

func server() {
	http.HandleFunc("/", parseHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}

//func client() {
//	locJson, err := json.Marshal(Data{Name: "florije4ex"})
//	req, err := http.NewRequest("POST", "http://localhost:8080", bytes.NewBuffer(locJson))
//	req.Header.Set("Content-Type", "application/json")
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Fatal(err)
//	}
//	body, err := ioutil.ReadAll(resp.Body)
//
//	fmt.Println("Response: ", string(body))
//	defer func() {
//		err := resp.Body.Close()
//		if err != nil {
//			log.Fatal(err.Error())
//		}
//	}()
//}

func TestPostData(t *testing.T) {
	go server()
	type Data struct {
		Name string `json:"name"`
	}
	data := Data{}
	err := PostData("http://localhost:8080", Data{Name: "florije4ex"}, &data)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(data.Name)
}
