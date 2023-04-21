package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/codegenerator"
)

type CodeProcessor struct {
	CodeRegistrationLimit int
	MaxConcurrentJobs     int
	BaseUrl               string
	Client                client.Client
	RegisteredDevices     chan RegisterDevice
}

type RegisterDevice struct {
	Identifier string
	Code       string
}

// Worker attempts to register a valid DevEUI via external LoRaWAN API.
// If successfull, a RegisterDevice struct with it's Identifier and Code will be sent to the work channel.
//
// # Example
//
//	Identifier: 1CEB0080F074F750, Code: 4F750
//
// When an unexpected error occurs, return ctx.Err instead.
func (cp *CodeProcessor) Worker(ctx context.Context, work chan struct{}) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-work:
			saved, registeredDevice := registerDevice(cp.Client, cp.BaseUrl)
			if saved {
				cp.RegisteredDevices <- registeredDevice
			}
		}
	}
}

func registerDevice(client client.Client, url string) (bool, RegisterDevice) {
	hex, err := codegenerator.GenerateHexString()
	if err != nil {
		log.Print(err)
		return false, RegisterDevice{}
	}

	code, err := codegenerator.GenerateCode(hex)
	if err != nil {
		log.Print(err)
		return false, RegisterDevice{}
	}

	b := new(bytes.Buffer)
	reqBody := map[string]string{"Deveui": code}

	err = json.NewEncoder(b).Encode(&reqBody)
	if err != nil {
		log.Print(err)
	}

	resp, err := client.Post(url, "application/json", b)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, RegisterDevice{Code: code, Identifier: hex}
	} else {
		return false, RegisterDevice{}
	}
}
