package main

import (
	"net/http"
	"time"

	"github.com/NickGowdy/deveui-cli/channel"
	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/processor"
)

func main() {

	codeChannel := &channel.CodeChannel{
		Msgch:  make(chan channel.Message, 10),
		Quitch: make(chan struct{}),
	}

	client := &client.HttpClient{
		BaseUrl: "http://europe-west1-machinemax-dev-d524.cloudfunctions.net",
		Client:  http.Client{Timeout: time.Duration(time.Second * time.Duration(30000))},
	}

	CodeProcessor := &processor.CodeProcessor{
		Client:         client,
		CodeChannel:    codeChannel,
		RegisterNumber: 100,
	}

	go codeChannel.StartAndListen()

	CodeProcessor.Process()

}
