package channel

import (
	"fmt"
	"time"
)

type Message struct {
	Code   string
	Status string
}

type CodeChannel struct {
	Msgch  chan Message
	Quitch chan struct{}
}

func (cc *CodeChannel) StartAndListen() {
listening:
	for {
		select {
		case msg := <-cc.Msgch:
			fmt.Printf("code: %s with status: %s\n", msg.Code, msg.Status)
		case <-cc.Quitch:
			fmt.Print("shutting down...\n")
			time.Sleep(1000)
			// os.Exit(4)
			break listening
		}
	}
}
