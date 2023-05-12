package client

import (
	"context"
	"io"
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
	req, err := http.NewRequestWithContext(ctx, "POST", fullUrl, body)
	if err != nil {
		return nil, err
	}

	resp, err = l.client.Do(req)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}
