package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
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
	msgch  chan Message
	quitch chan struct{}
}

func (s *Server) StartAndListen() {
	for {
		select {
		// block here until someone is sending a message to the channel
		case msg := <-s.msgch:
			fmt.Printf("code: %s with status: %s\n", msg.Code, msg.Status)
		case <-s.quitch:
		default:

		}
	}
}

func main() {

	s := &Server{
		msgch: make(chan Message, 10),
	}
	client := client.NewHttpClient(30, "http://europe-west1-machinemax-dev-d524.cloudfunctions.net")

	go s.StartAndListen()
	// var i int
	for {
		time.Sleep(2000 * time.Millisecond)
		hexStr, err := generateHexString(16)
		if err != nil {
			log.Print(err)
		}

		code := hexStr[len(hexStr)-5:]
		go registerCode(code, client, s.msgch)
	}
}

func registerCode(code string, client *client.HttpClient, msgch chan Message) {

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

	msg := Message{
		Code:   code,
		Status: resp.Status,
	}

	msgch <- msg
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
