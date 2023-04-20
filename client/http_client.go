package client

import (
	"bytes"
	"fmt"
	"net/http"
)

type HttpClient struct {
	BaseUrl string
	Client  http.Client
}

func (h *HttpClient) Post(endpoint string, b *bytes.Buffer) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", h.BaseUrl, endpoint)
	resp, err := h.Client.Post(url, "application/json", b)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
