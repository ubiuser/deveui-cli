package client

import (
	"io"
	"net/http"
)

// Client used to communicate to external services
type Client interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
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

// Send data via POST (HTTP) request
func (l *LoraWAN) DoPost(body io.Reader) (resp *http.Response, err error) {
	fullUrl := l.baseURL + endpoint
	resp, err = l.client.Post(fullUrl, "application/json", body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
