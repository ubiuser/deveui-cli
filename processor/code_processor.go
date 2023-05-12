package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/device"
)

type CodeProcessor struct {
	CodeRegistrationLimit int
	MaxConcurrentJobs     int
	LoraWAN               client.LoraWAN
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
			registeredDevice, err := registerDevice(cp.LoraWAN, ctx)
			if err == nil {
				cp.Device <- *registeredDevice
			} else {
				return err
			}
		}
	}
}

func registerDevice(loraWAN client.LoraWAN, ctx context.Context) (*device.Device, error) {
	device := device.NewDevice()
	identifier := device.Get()

	b := new(bytes.Buffer)
	reqBody := map[string]string{"Deveui": identifier}

	err := json.NewEncoder(b).Encode(&reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := loraWAN.DoPost(b, ctx)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return device, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}
