package processor

import (
	"context"

	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/device"
)

type CodeProcessor struct {
	CodeRegistrationLimit int
	MaxConcurrentJobs     int
	LoraWAN               client.LoraWAN
	DeviceCh              chan device.Device
	DoneCh                chan struct{}
}

// Worker attempts to register a valid DevEUI via external LoRaWAN API.
// If successful, a RegisterDevice struct with its Identifier and Code will be sent to the work channel.
//
// # Example
//
//	Identifier: 1CEB0080F074F750, Code: 4F750
//
// When an unexpected error occurs, return ctx.Err instead.
func (cp *CodeProcessor) Worker(ctx context.Context, workCh chan device.Device) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-workCh:
			registeredDevice, err := cp.LoraWAN.Send(ctx)
			if err == nil {
				cp.DeviceCh <- *registeredDevice

			} else {
				return err
			}
		}
	}
}
