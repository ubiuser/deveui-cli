package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/NickGowdy/deveui-cli/channel"
	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/processor"
	"github.com/joho/godotenv"
)

const (
	MAX_CONCURRENT_JOBS     = /* Buffer limit for channel */ 10
	CODE_REGISTRATION_LIMIT = /* Maximum number of devices that will be registered */ 100
	TIMEOUT                 = /* Seconds */ 30000
)

func main() {
	godotenv.Load(".env")
	baseurl := os.Getenv("BASE_URL")

	// setup channel for listening to SIG cmds
	signalChannel := &channel.SignalChannel{}

	// setup client for requests
	httpClient := &http.Client{
		Timeout: time.Second * TIMEOUT,
	}
	loraWanClient := &client.LoraWanClient{
		Client: httpClient,
	}

	// setup processor to do work
	CodeProcessor := &processor.CodeProcessor{
		MaxConcurrentJobs:     MAX_CONCURRENT_JOBS,
		BaseUrl:               baseurl,
		CodeRegistrationLimit: CODE_REGISTRATION_LIMIT,
		Client:                loraWanClient,
	}

	go signalChannel.StartAndListen()
	registeredDevices := CodeProcessor.Process()

	for _, d := range *registeredDevices {
		for k, v := range d {
			fmt.Println(k, "value is", v)
		}
	}
}
