package channel

import "fmt"

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
			fmt.Print("shutting down...")
			break listening
		}
	}
}
