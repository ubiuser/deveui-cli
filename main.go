package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/NickGowdy/deveui-cli/channel"
	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/processor"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	baseurl := os.Getenv("BASE_URL")

	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}

	limit, err := strconv.Atoi(os.Getenv("CODE_REGISTRATION_LIMIT"))
	if err != nil {
		log.Fatal(err)
	}

	codeChannel := &channel.CodeChannel{
		Msgch:  make(chan channel.Message, 10),
		Quitch: make(chan struct{}),
	}

	signalChannel := &channel.SignalChannel{}

	client := &client.HttpClient{
		BaseUrl: baseurl,
		Client:  http.Client{Timeout: time.Duration(time.Second * time.Duration(timeout))},
	}

	CodeProcessor := &processor.CodeProcessor{
		Client:         client,
		CodeChannel:    codeChannel,
		SignalChannel:  signalChannel,
		RegisterNumber: limit,
	}

	CodeProcessor.Start()
}
