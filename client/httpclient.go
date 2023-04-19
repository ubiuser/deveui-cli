package client

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type HttpClient struct {
	baseUrl string
	client  http.Client
}

func NewHttpClient(duration time.Duration, baseUrl string) *HttpClient {
	return &HttpClient{
		client:  http.Client{Timeout: time.Duration(time.Second * time.Duration(duration))},
		baseUrl: baseUrl,
	}
}

func (h *HttpClient) Post(endpoint string, b *bytes.Buffer) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", h.baseUrl, endpoint)
	resp, err := h.client.Post(url, "application/json", b)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return resp, nil
}
