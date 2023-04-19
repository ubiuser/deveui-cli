package client

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type HttpClient struct {
	timeout     time.Duration
	baseUrl     string
	client      http.Client
	contentType string
}

func NewHttpClient(duration time.Duration, baseUrl string) *HttpClient {
	return &HttpClient{
		timeout: duration,
		baseUrl: baseUrl,
	}
}

func (h *HttpClient) Post(endpoint string, b *bytes.Buffer) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", h.baseUrl, endpoint)
	resp, err := h.client.Post(url, h.contentType, b)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return resp, nil
}
