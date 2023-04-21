package client

import (
	"fmt"
	"io"
	"net/http"
)

// Client used to communicate to LoRaWAN external system
type LoraWanClient struct {
	Client Client
}

const endpoint = "/sensor-onboarding-sample" // endpoint for saving DevEUI via LoRaWAN

// Send data via POST (HTTP) request
func (h *LoraWanClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	fullUrl := fmt.Sprintf("%s/%s", url, endpoint)
	resp, err = h.Client.Post(fullUrl, contentType, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
