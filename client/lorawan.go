package client

import (
	"context"
	"io"
	"log"
	"net/http"
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

// DoPost sends data via POST (HTTP) request
func (l *LoraWAN) DoPost(body io.Reader, ctx context.Context) (resp *http.Response, err error) {
	fullUrl := l.baseURL + endpoint
	req, err := http.NewRequest("POST", fullUrl, body)
	if err != nil {
		log.Fatalf("%v", err)
	}

	req = req.WithContext(ctx)

	resp, err = l.client.Do(req)
	if err != nil {
		log.Printf("%v", err)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}
