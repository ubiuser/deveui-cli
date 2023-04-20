package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/codegenerator"
)

type CodeProcessor struct {
	CodeRegistrationLimit int
	MaxConcurrentJobs     int
	BaseUrl               string
	Client                client.Client
	RegisteredDevices     []RegisterDevice
}

type RegisterDevice struct {
	Identifier string
	Code       string
}

func (cp *CodeProcessor) Process() *[]RegisterDevice {
	waitChan := make(chan struct{}, cp.MaxConcurrentJobs)
	var count int
	m := sync.Mutex{}
	wg := sync.WaitGroup{}

	for count < cp.CodeRegistrationLimit {
		waitChan <- struct{}{}
		wg.Add(1)
		go func(innerCount int) {
			saved, registeredDevice := registerDevice(cp.Client, cp.BaseUrl)
			if saved {
				m.Lock()
				count++
				cp.RegisteredDevices = append(cp.RegisteredDevices, registeredDevice)
				defer m.Unlock()
			}

			<-waitChan

		}(count)
		wg.Done()
	}

	close(waitChan)
	return &cp.RegisteredDevices
}

func registerDevice(client client.Client, url string) (bool, RegisterDevice) {
	hex, err := codegenerator.GenerateHexString(16)
	if err != nil {
		log.Print(err)
		return false, RegisterDevice{}
	}

	code, err := codegenerator.Generate(hex)
	if err != nil {
		log.Print(err)
		return false, RegisterDevice{}
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
		return true, RegisterDevice{Code: code, Identifier: hex}
	} else {
		return false, RegisterDevice{}
	}
}
