package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const MAX_CONCURRENT_JOBS = 10
const allowedChars = "ABCDEF0123456789"

func main() {
	waitChan := make(chan struct{}, MAX_CONCURRENT_JOBS)
	var count int32
	var wg sync.WaitGroup

	for count < 100 {
		wg.Add(1)
		waitChan <- struct{}{}

		go func(ops int32) {
			saved := job()
			if saved {
				atomic.AddInt32(&count, 1)
			}

			<-waitChan
			wg.Done()
		}(count)
	}

	close(waitChan)
}

func job() bool {
	hexStr, err := generateHexString(16)
	if err != nil {
		log.Print(err)
	}
	code := hexStr[len(hexStr)-5:]
	client := http.Client{Timeout: time.Second * 30}

	b := new(bytes.Buffer)
	reqBody := map[string]string{"Deveui": code}

	err = json.NewEncoder(b).Encode(&reqBody)
	if err != nil {
		log.Print(err)
	}

	resp, err := client.Post("http://europe-west1-machinemax-dev-d524.cloudfunctions.net/sensor-onboarding-sample", "application/json", b)

	if err != nil {
		log.Print(err)
	}

	defer resp.Body.Close()

	fmt.Printf("%s\n", resp.Status)

	return resp.StatusCode == http.StatusOK
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

// godotenv.Load(".env")

// baseurl := os.Getenv("BASE_URL")

// timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
// if err != nil {
// 	log.Fatal(err)
// }

// limit, err := strconv.Atoi(os.Getenv("CODE_REGISTRATION_LIMIT"))
// if err != nil {
// 	log.Fatal(err)
// }

// codeChannel := &channel.CodeChannel{
// 	Msgch:  make(chan channel.Message, 10),
// 	Quitch: make(chan struct{}),
// }

// signalChannel := &channel.SignalChannel{}

// client := &client.HttpClient{
// 	BaseUrl: baseurl,
// 	Client: &http.Client{
// 		Timeout: time.Duration(time.Second * time.Duration(timeout)),
// 		Transport: &http.Transport{
// 			MaxIdleConns:        10,
// 			MaxIdleConnsPerHost: 10,
// 		},
// 	},
// }

// CodeProcessor := &processor.CodeProcessor{
// 	Client:         client,
// 	CodeChannel:    codeChannel,
// 	SignalChannel:  signalChannel,
// 	RegisterNumber: limit,
// }

// CodeProcessor.Start()
