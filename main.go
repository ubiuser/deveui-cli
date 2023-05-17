package main

import (
	"context"
	"net/http"
	"time"

	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/processor"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	BaseURL               string        `envconfig:"BASE_URL" default:"http://europe-west1-machinemax-dev-d524.cloudfunctions.net"`
	MaxConcurrent         int           `envconfig:"MAX_CONCURRENT" default:"10"`
	CodeRegistrationLimit int           `envconfig:"CODE_REGISTRATION_LIMIT" default:"100"`
	Timeout               time.Duration `envconfig:"TIMEOUT" default:"30s"`
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		panic("failed to process env vars: " + err.Error())
	}

	httpClient := &http.Client{
		Timeout: cfg.Timeout,
	}
	loraWAN := client.NewLoraWAN(cfg.BaseURL, httpClient)

	proc := &processor.Processor{
		CodeRegistrationLimit: cfg.CodeRegistrationLimit,
		MaxConcurrentJobs:     cfg.MaxConcurrent,
		LoraWAN:               *loraWAN,
	}

	ctx, cancel := context.WithCancel(context.Background())
	proc.Start(ctx, cancel)
}
