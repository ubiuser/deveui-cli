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
	RegisteredDevices     []map[int]string
}

func (cp *CodeProcessor) Process() *[]map[int]string {
	waitChan := make(chan struct{}, cp.MaxConcurrentJobs)
	var count int32

	for count < cp.CodeRegistrationLimit {
		waitChan <- struct{}{}
		go func(innerCount int32) {
			saved, code := registerDevice(cp.Client, cp.BaseUrl)
			if saved {
				atomic.AddInt32(&count, 1)
				cp.RegisteredDevices = append(cp.RegisteredDevices, map[int]string{int(count): code})
			}

			<-waitChan

		}(count)
	}

	close(waitChan)
	return &cp.RegisteredDevices
}

func registerDevice(client client.Client, url string) (bool, string) {

	code, err := codegenerator.Generate()

	if err != nil {
		log.Print(err)
		return false, ""
	}

	b := new(bytes.Buffer)
	reqBody := map[string]string{"Deveui": code}

	err = json.NewEncoder(b).Encode(&reqBody)
	if err != nil {
		log.Print(err)
	}

	resp, err := client.Post(url, "application/json", b)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	fmt.Printf("%s\n", resp.Status)

	if resp.StatusCode == http.StatusOK {
		return true, code
	} else {
		return false, ""
	}
}
