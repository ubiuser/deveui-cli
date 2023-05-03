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
	BaseUrl               string // seems like something the client should know about, not this processor
	Client                client.Client
	Device                chan device.Device
}

// Worker attempts to register a valid DevEUI via external LoRaWAN API.
// If successful, a RegisterDevice struct with its Identifier and Code will be sent to the work channel.
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

// I think you should store the devices in a map or something to be able to ensure their uniqueness.
// The reason is that even if the 16 char code is unique, the last 5 chars used for the lookup can be the same.
// This wasn't an explicit requirement, but from the description it would just make sense.

// if this was a method of CodeProcessor then you wouldn't need the input parameters
// The bool return parameter is redundant, you can just return the device pointer and check if it is nil.
// Even better, return the pointer AND an error.
func registerDevice(client client.Client, url string) (bool, *device.Device) {
	device := device.NewDevice()

	b := new(bytes.Buffer)
	reqBody := map[string]string{"Deveui": device.Identifier}

	err := json.NewEncoder(b).Encode(&reqBody)
	if err != nil {
		// if you just log the error here, the POST request will still hit the endpoint with invalid data
		log.Print(err)
	}

	resp, err := client.Post(url, "application/json", b)

	if err != nil {
		// this will kill the whole app on the first POST error
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, device
	} else {
		return false, nil
	}
}
