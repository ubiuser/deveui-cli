package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/NickGowdy/deveui-cli/device"
)

// Client used to communicate to external services
type Client interface {
	Do(*http.Request) (resp *http.Response, err error)
}

// LoraWAN used to communicate to LoRaWAN external system
type LoraWAN struct {
	baseURL string
	client  Client
}

func NewLoraWAN(baseURL string, client Client) *LoraWAN {
	return &LoraWAN{
		baseURL: baseURL,
		client:  client,
	}
}

const endpoint = "/sensor-onboarding-sample" // endpoint for saving DevEUI via LoRaWAN

// RegisterDevice registers new device using LoraWAN external service
func (l *LoraWAN) RegisterDevice(ctx context.Context) (*device.Device, error) {
	device := device.NewDevice()
	identifier := device.GetIdentifier()
	b := new(bytes.Buffer)
	reqBody := map[string]string{"Deveui": identifier}

	err := json.NewEncoder(b).Encode(&reqBody)
	if err != nil {
		return nil, err
	}

	fullUrl := l.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, "POST", fullUrl, b)
	if err != nil {
		return nil, err
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return device, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}
