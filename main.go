package main

import (
	"net/http"
	"os"
	"time"

	"github.com/NickGowdy/deveui-cli/channel"
	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/processor"
	"github.com/joho/godotenv"
)

const (
	MAX_CONCURRENT_JOBS     = 10
	CODE_REGISTRATION_LIMIT = 10
	TIMEOUT                 = /* Seconds */ 30000
)

func main() {
	godotenv.Load(".env")
	baseurl := os.Getenv("BASE_URL")

	signalChannel := &channel.SignalChannel{}
	httpClient := &http.Client{
		Timeout: time.Second * TIMEOUT,
	}

	CodeProcessor := &processor.CodeProcessor{
		MaxConcurrentJobs:     MAX_CONCURRENT_JOBS,
		BaseUrl:               baseurl,
		CodeRegistrationLimit: CODE_REGISTRATION_LIMIT,
		Client: &client.LoraWanClient{
			Client: httpClient,
		},
	}

	go signalChannel.StartAndListen()

	CodeProcessor.Start()
}
