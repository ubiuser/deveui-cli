package client

import (
	"io"
	"net/http"
	"path"
	"time"
)

// Client used to communicate to external services
type Client interface {
	Post(body io.Reader) (resp *http.Response, err error)
}

// LoraWAN used to communicate to LoRaWAN external system
type LoraWAN struct {
	client  http.Client
	baseURL string
}

func NewLoraWAN(baseURL string, timeout time.Duration) *LoraWAN {
	return &LoraWAN{
		client: http.Client{
			Timeout: timeout * time.Second,
		},
	}
}

const endpoint = "/sensor-onboarding-sample" // endpoint for saving DevEUI via LoRaWAN

// Send data via POST (HTTP) request
func (l *LoraWAN) Post(body io.Reader) (resp *http.Response, err error) {
	fullUrl := path.Join(l.baseURL, endpoint)
	resp, err = l.client.Post(fullUrl, "application/json", body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
