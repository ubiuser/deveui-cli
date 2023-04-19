package channel

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type SignalChannel struct{}

func (sc *SignalChannel) StartAndListen() {
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			time.Sleep(1000)
		}
	}()
	sig := <-cancelChan
	log.Printf("Caught signal %v", sig)
}
