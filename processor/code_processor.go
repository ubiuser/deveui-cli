package processor

import (
	"context"
	"log"

	"github.com/NickGowdy/deveui-cli/device"
)

type (
	// Client is the client interface the processor expects
	Client interface {
		RegisterDevice(ctx context.Context, newDevice *device.Device) error
	}

	Processor struct {
		codeRegistrationLimit int
		maxConcurrentJobs     int
		client                Client
	}
)

func New(codeRegistrationLimit int, maxConcurrentJobs int, client Client) *Processor {
	return &Processor{
		codeRegistrationLimit: codeRegistrationLimit,
		maxConcurrentJobs:     maxConcurrentJobs,
		client:                client,
	}
}

func (p *Processor) Start(ctx context.Context, cancel context.CancelFunc) {
	// workCh := make(chan struct{})
	count := 0

	for count < p.codeRegistrationLimit {
		newDevice, err := device.NewDevice()
		if err != nil {
			log.Printf("failed to create new device: %v\n", err)
			continue
		}

		if err = p.client.RegisterDevice(ctx, newDevice); err != nil {
			log.Printf("failed to register device: %v\n", err)
			continue
		}

		log.Println(newDevice)
		count++
	}

	// go func(ctx context.Context) {
	// 	for {
	// 		p.doWork(ctx, workCh)
	// 	}
	// }(ctx)

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	case <-workCh:
	// 		count++
	// 		if count == p.CodeRegistrationLimit {
	// 			cancel()
	// 			fmt.Printf("work complete \n")
	// 		}
	// 	}
	// }
}

// func (cp *Processor) doWork(ctx context.Context, workCh chan<- struct{}) {
// 	device, err := cp.LoraWAN.RegisterDevice(ctx)
// 	if err != nil {
// 		return
// 	} else {
// 		device.Print()
// 		workCh <- struct{}{}
// 	}
// }
