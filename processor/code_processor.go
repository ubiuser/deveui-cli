package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/NickGowdy/deveui-cli/codegenerator"
)

type CodeProcessor struct {
	CodeRegistrationLimit int32
	MaxConcurrentJobs     int
	BaseUrl               string
}

func (cp *CodeProcessor) Start() {

	waitChan := make(chan struct{}, cp.MaxConcurrentJobs)
	var count int32

	for count < cp.CodeRegistrationLimit {
		waitChan <- struct{}{}
		go func(ops int32) {
			saved := job()
			if saved {
				atomic.AddInt32(&count, 1)
			}

			<-waitChan

		}(count)
	}

	close(waitChan)
}

func job() bool {
	client := http.Client{Timeout: time.Second * 30}
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

	resp, err := client.Post("http://europe-west1-machinemax-dev-d524.cloudfunctions.net/sensor-onboarding-sample", "application/json", b)

	if err != nil {
		log.Print(err)
	}

	defer resp.Body.Close()

	fmt.Printf("%s\n", resp.Status)

	return resp.StatusCode == http.StatusOK
}
