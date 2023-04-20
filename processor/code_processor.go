package processor

import (
	"crypto/rand"
	"log"
	"math/big"
	"sync"

	"github.com/NickGowdy/deveui-cli/channel"
	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/register"
)

const allowedChars = "ABCDEF0123456789"

type CodeProcessor struct {
	Client         client.Client
	CodeChannel    *channel.CodeChannel
	SignalChannel  *channel.SignalChannel
	RegisterNumber int
}

func (cp *CodeProcessor) Start() {
	go cp.CodeChannel.StartAndListen()
	go cp.SignalChannel.StartAndListen()

	process(cp)
}

func process(cp *CodeProcessor) {
	var i int
	var lock sync.Mutex
	for {

		hexStr, err := generateHexString(16)
		if err != nil {
			log.Print(err)
		}
		code := hexStr[len(hexStr)-5:]
		codeRegister := &register.CodeRegister{
			HttpClient: cp.Client,
			Code:       code,
		}

		go func(code string) {

			resp, err := codeRegister.RegisterCode()
			if err != nil {
				log.Print(err)
			}

			if resp.StatusCode == 200 {
				msg := channel.Message{
					Code:   code,
					Status: resp.Status,
				}

				cp.CodeChannel.Msgch <- msg
				lock.Lock()
				defer lock.Unlock()
				i++
			}

			if i == cp.RegisterNumber {
				close(cp.CodeChannel.Quitch)
			}
		}(code)
	}
}

func generateHexString(length int) (string, error) {
	max := big.NewInt(int64(len(allowedChars)))
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = allowedChars[n.Int64()]
	}
	return string(b), nil
}
