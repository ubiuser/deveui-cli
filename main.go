package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"
)

const allowedChars = "ABCDEF0123456789"

type Request struct {
	Deveui string `json:"deveui"`
}

func main() {
	hexStr, err := generateHexString(16)

	if err != nil {
		log.Print(err)
	}

	fmt.Print(hexStr)

	client := http.Client{Timeout: time.Duration(time.Second * time.Duration(30))}
	b := new(bytes.Buffer)

	reqBody := Request{Deveui: hexStr}

	err = json.NewEncoder(b).Encode(&reqBody)

	if err != nil {
		log.Print(err)
	}

	resp, err := client.Post("http://europe-west1-machinemax-dev-d524.cloudfunctions.net/sensor-onboarding-sample", "application/json", b)

	if err != nil {
		log.Print(err)
	}

	fmt.Print(resp)
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
