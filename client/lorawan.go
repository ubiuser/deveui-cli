package client

import (
	"io"
	"net/http"
	"path"
)

// Generic client used to communicate to external services
type Client interface {
	Post(baseURL string, contentType string, body io.Reader) (resp *http.Response, err error)
}

// Client used to communicate to LoRaWAN external system
type LoraWAN struct {
	client Client
}

func NewLoraWAN(client Client) *LoraWAN {
	return &LoraWAN{
		client: client,
	}
}

const endpoint = "/sensor-onboarding-sample" // endpoint for saving DevEUI via LoRaWAN

// Send data via POST (HTTP) request
func (l *LoraWAN) Post(baseURL string, contentType string, body io.Reader) (resp *http.Response, err error) {
	fullUrl := path.Join(baseURL, endpoint)
	resp, err = l.client.Post(fullUrl, contentType, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
