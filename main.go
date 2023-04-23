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
	"github.com/NickGowdy/deveui-cli/device"
	"github.com/NickGowdy/deveui-cli/processor"
	"github.com/joho/godotenv"
)

const (
	MAX_CONCURRENT_JOBS     = /* Buffer limit for channel */ 10
	CODE_REGISTRATION_LIMIT = /* Maximum number of devices that will be registered */ 100
	TIMEOUT                 = /* Seconds */ 30
)

/*
deveui-cli is a Go CLI program.
It is used for concurrently generating unique 16-character (hex) identifier called a DevEUI.
These are generated by the program and registered via an external (LoRaWAN) API.

Usage:

	go run main.go (locally)

	go run deveui-cli (docker)

Once this program starts, it will listen to syscall.SIGTERM, syscall.SIGINT via a channel.
This is to handle any unexpected terminations of the program and to resume processing DevEUIs.
*/
func main() {
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
		Device:                make(chan device.Device),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	work := make(chan struct{}, MAX_CONCURRENT_JOBS)
	listener := make(chan os.Signal, 1)

	// goroutine to listen for syscall.SIGTERM, syscall.SIGINT
	go func() {
		signal.Notify(listener, syscall.SIGTERM, syscall.SIGINT)
		go func() {
			for {
				time.Sleep(1000)
			}
		}()
		sig := <-listener
		log.Printf("Caught signal %v", sig)
	}()

	go func() {
		for {
			work <- struct{}{}
		}
	}()

	// Spawn workers
	for job := 0; job < MAX_CONCURRENT_JOBS; job++ {
		go codeProcessor.Worker(ctx, work)
	}

	// stdout any registered devices and increment until CODE_REGISTRATION_LIMIT is reached.
	count := 0
	for d := range codeProcessor.Device {
		fmt.Printf("device: %d has identifier: %s and code: %s\n", count+1, d.Identifier, d.Code)
		count += 1
		if count == CODE_REGISTRATION_LIMIT {
			break
		}
	}

	close(work)
	close(listener)
}
