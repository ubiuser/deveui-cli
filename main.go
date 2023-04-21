package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/processor"
	"github.com/joho/godotenv"
)

const (
	MAX_CONCURRENT_JOBS     = /* Buffer limit for channel */ 10
	CODE_REGISTRATION_LIMIT = /* Maximum number of devices that will be registered */ 100
	TIMEOUT                 = /* Seconds */ 5
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	godotenv.Load(".env")
	baseurl := os.Getenv("BASE_URL")

	// setup client for requests
	httpClient := &http.Client{
		Timeout: time.Second * TIMEOUT,
	}
	loraWanClient := &client.LoraWanClient{
		Client: httpClient,
	}

	// setup processor to do work
	codeProcessor := &processor.CodeProcessor{
		MaxConcurrentJobs:     MAX_CONCURRENT_JOBS,
		BaseUrl:               baseurl,
		CodeRegistrationLimit: CODE_REGISTRATION_LIMIT,
		Client:                loraWanClient,
		RegisteredDevices:     make(chan processor.RegisterDevice),
	}

	work := make(chan struct{}, MAX_CONCURRENT_JOBS)
	go func() {
		for {
			work <- struct{}{}
		}
	}()

	// Spawn workers
	for j := 0; j < MAX_CONCURRENT_JOBS; j++ {
		go codeProcessor.Worker(ctx, work)
	}

	n := 0
	for d := range codeProcessor.RegisteredDevices {
		fmt.Printf("device: %d has identifier: %s and code: %s\n", n+1, d.Identifier, d.Code)
		n += 1
		if n == CODE_REGISTRATION_LIMIT {
			break
		}
	}

	// goroutine to listen for syscall.SIGTERM, syscall.SIGINT
	go func() {
		cancelChan := make(chan os.Signal, 1)
		signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
		go func() {
			for {
				time.Sleep(1000)
			}
		}()
		sig := <-cancelChan
		log.Printf("Caught signal %v", sig)
	}()
}
