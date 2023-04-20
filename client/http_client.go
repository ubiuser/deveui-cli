package client

import (
	"fmt"
	"io"
	"net/http"
)

type HttpClient struct {
	BaseUrl string
	Client  Client
}

func (h *HttpClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	fullUrl := fmt.Sprintf("%s/%s", h.BaseUrl, url)
	resp, err = h.Client.Post(fullUrl, contentType, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
