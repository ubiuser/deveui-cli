package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/codegenerator"
)

type CodeProcessor struct {
	CodeRegistrationLimit int32
	MaxConcurrentJobs     int
	BaseUrl               string
	Client                client.Client
}

func (cp *CodeProcessor) Start() {
	waitChan := make(chan struct{}, cp.MaxConcurrentJobs)
	// client := http.Client{Timeout: time.Second * 30}
	var count int32

	for count < cp.CodeRegistrationLimit {
		waitChan <- struct{}{}
		go func(ops int32) {
			saved := job(cp.Client, cp.BaseUrl)
			if saved {
				atomic.AddInt32(&count, 1)
			}

			<-waitChan

		}(count)
	}

	close(waitChan)
}

func job(client client.Client, url string) bool {

	code, err := codegenerator.Generate()

	if err != nil {
		log.Print(err)
		return false
	}

	b := new(bytes.Buffer)
	reqBody := map[string]string{"Deveui": code}

	err = json.NewEncoder(b).Encode(&reqBody)
	if err != nil {
		log.Print(err)
	}

	resp, err := client.Post(url, "application/json", b)

	if err != nil {
		log.Print(err)
	}

	defer resp.Body.Close()

	fmt.Printf("%s\n", resp.Status)

	return resp.StatusCode == http.StatusOK
}
