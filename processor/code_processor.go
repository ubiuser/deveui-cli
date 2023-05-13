package processor

import (
	"context"
	"fmt"

	"github.com/NickGowdy/deveui-cli/client"
)

type CodeProcessor struct {
	CodeRegistrationLimit int
	MaxConcurrentJobs     int
	LoraWAN               client.LoraWAN
}

func (cp *CodeProcessor) Start(ctx context.Context, cancel context.CancelFunc) {
	workCh := make(chan struct{})
	count := 0
	go func(ctx context.Context) {
		for {
			cp.doWork(ctx, workCh)
		}
	}(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-workCh:
			count++
			if count == cp.CodeRegistrationLimit {
				cancel()
				fmt.Printf("work complete \n")
			}
		}
	}
}

func (cp *CodeProcessor) doWork(ctx context.Context, workCh chan<- struct{}) {
	device, err := cp.LoraWAN.RegisterDevice(ctx)
	if err != nil {
		return
	} else {
		device.Print()
		workCh <- struct{}{}
	}
}
