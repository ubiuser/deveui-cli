package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net/http"
	"time"
)

const allowedChars = "ABCDEF0123456789"

type Request struct {
	Deveui string `json:"deveui"`
}

type Response struct {
	Description string `json:"description"`
}

func main() {
	hexStr, err := generateHexString(16)
	code := hexStr[len(hexStr)-5:]

	if err != nil {
		log.Print(err)
	}

	client := http.Client{Timeout: time.Duration(time.Second * time.Duration(30))}
	b := new(bytes.Buffer)

	reqBody := Request{Deveui: code}

	err = json.NewEncoder(b).Encode(&reqBody)

	if err != nil {
		log.Print(err)
	}

	resp, err := client.Post("http://europe-west1-machinemax-dev-d524.cloudfunctions.net/sensor-onboarding-sample", "application/json", b)
	if err != nil {
		log.Print(err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}

	if resp.StatusCode == http.StatusUnprocessableEntity || resp.StatusCode == http.StatusOK {
		bodyString := string(bodyBytes)
		log.Print(bodyString)
	}

	if err != nil {
		log.Print(err)
	}
}

func generateHexString(length int) (string, error) {
	max := big.NewInt(int64(len(allowedChars)))
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = allowedChars[n.Int64()]
	}
	return string(b), nil
}
