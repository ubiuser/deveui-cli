package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/device"
)

type CodeProcessor struct {
	CodeRegistrationLimit int
	MaxConcurrentJobs     int
	BaseUrl               string
	Client                client.Client
	Device                chan device.Device
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
				cp.Device <- *registeredDevice
			}
		}
	}
}

func registerDevice(client client.Client, url string) (bool, *device.Device) {
	device := device.NewDevice()

	b := new(bytes.Buffer)
	reqBody := map[string]string{"Deveui": device.Identifier}

	err := json.NewEncoder(b).Encode(&reqBody)
	if err != nil {
		log.Print(err)
	}

	resp, err := client.Post(b)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, device
	} else {
		return false, nil
	}
}
