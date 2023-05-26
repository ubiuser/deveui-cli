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

	// T is a short alias for a blank struct, which has zero size
	T struct{}
)

func New(codeRegistrationLimit int, maxConcurrentJobs int, client Client) *Processor {
	return &Processor{
		codeRegistrationLimit: codeRegistrationLimit,
		maxConcurrentJobs:     maxConcurrentJobs,
		client:                client,
	}
}

func (p *Processor) printDevice(done chan<- T, ch <-chan *device.Device) {
	// count is not under race condition here
	count := 1
	for dev := range ch {
		log.Printf("registered device %d/%d: %s\n", count, p.codeRegistrationLimit, dev.String())
		count++
	}
	done <- T{}
}

func (p *Processor) Start(ctx context.Context, cancel context.CancelFunc) {
	done := make(chan T)
	devicePrinter := make(chan *device.Device)
	go p.printDevice(done, devicePrinter)

	workers := make(chan T, p.maxConcurrentJobs)
	for i := 0; i < p.maxConcurrentJobs; i++ {
		workers <- T{}
	}

	// NOTE: use a buffered channel to track successfully completed jobs
	// First, fill up the channel to max capacity, then during the loop take
	// one item until the channel is closed. The channel is closed when it becomes
	// empty.
	leftToDo := make(chan T, p.codeRegistrationLimit)
	for i := 0; i < p.codeRegistrationLimit; i++ {
		leftToDo <- T{}
	}

	// This will block when there are only in-flight jobs running and we wait for their outcome. In other words,
	// we don't start any new jobs until we have in-flight jobs and they have a chance to satisfy the requested
	// number of registrations.
	// When the channel is closed, we know that all jobs have been completed.
	for range leftToDo {
		select {
		case <-ctx.Done():
			// graceful shutdown
			log.Printf("processor shutting down")

			return
		default:
			if _, ok := <-workers; ok {
				go func() {
					defer func() {
						workers <- T{}

						if len(workers) == cap(workers) && len(leftToDo) == 0 {
							close(leftToDo)
							close(workers)
							close(devicePrinter)
						}
					}()

					newDevice, err := device.NewDevice()
					if err != nil {
						log.Printf("failed to create new device: %v\n", err)
						leftToDo <- T{}

						return
					}

					if err = p.client.RegisterDevice(ctx, newDevice); err != nil {
						log.Printf("failed to register device: %v\n", err)
						leftToDo <- T{}

						return
					}

					devicePrinter <- newDevice
				}()
			}
		}
	}

	<-done
	log.Println("all jobs have been completed")
	cancel()
}
