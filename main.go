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
	"time"

	"github.com/NickGowdy/deveui-cli/client"
)

const allowedChars = "ABCDEF0123456789"

type Request struct {
	Deveui string `json:"deveui"`
}

type Message struct {
	Code   string
	Status string
}

type Server struct {
	msgch   chan Message
	quitch  chan struct{}
	request int
}

func (s *Server) StartAndListen() {
listening:
	for {
		select {
		case msg := <-s.msgch:
			fmt.Printf("code: %s with status: %s\n", msg.Code, msg.Status)
		case <-s.quitch:
			fmt.Print("shutting down...")
			break listening
		}
	}
}

func main() {

	s := &Server{
		msgch:   make(chan Message, 10),
		quitch:  make(chan struct{}),
		request: 100,
	}
	client := &client.HttpClient{
		BaseUrl: "http://europe-west1-machinemax-dev-d524.cloudfunctions.net",
		Client:  http.Client{Timeout: time.Duration(time.Second * time.Duration(30000))},
	}

	go s.StartAndListen()

	var i int
	var lock sync.Mutex

	var wg sync.WaitGroup
	for {
		time.Sleep(500 * time.Millisecond)
		hexStr, err := generateHexString(16)
		if err != nil {
			log.Print(err)
		}

		code := hexStr[len(hexStr)-5:]
		wg.Add(1)
		go func(code string) {
			resp, err := registerCode(code, client)
			if err != nil {
				log.Print(err)
			}

			if resp.StatusCode == 200 {
				msg := Message{
					Code:   code,
					Status: resp.Status,
				}

				s.msgch <- msg
				lock.Lock()
				defer lock.Unlock()
				i++
			}

			if i == 2 {
				close(s.quitch)
			}

			wg.Done()
		}(code)
	}

}

func registerCode(code string, client *client.HttpClient) (*http.Response, error) {

	b := new(bytes.Buffer)
	reqBody := Request{Deveui: code}

	err := json.NewEncoder(b).Encode(&reqBody)
	if err != nil {
		log.Print(err)
	}

	resp, err := client.Post("sensor-onboarding-sample", b)
	if err != nil {
		log.Print(err)
	}

	defer resp.Body.Close()
	return resp, nil
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
