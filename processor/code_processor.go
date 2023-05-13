package processor

import (
	"context"
	"log"

	"github.com/NickGowdy/deveui-cli/client"
)

type Processor struct {
	CodeRegistrationLimit int
	MaxConcurrentJobs     int
	LoraWAN               client.LoraWAN
}

func (p *Processor) Start(ctx context.Context, cancel context.CancelFunc) {
	// workCh := make(chan struct{})
	count := 0

	for count < p.CodeRegistrationLimit {
		device, err := p.LoraWAN.RegisterDevice(ctx)
		if err != nil {
			log.Print(err)
		} else {
			device.Print()
			count++
		}
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
