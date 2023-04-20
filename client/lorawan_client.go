package client

import (
	"fmt"
	"io"
	"net/http"
)

type LoraWanClient struct {
	Client Client
}

const endpoint = "/sensor-onboarding-sample"

func (h *LoraWanClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	fullUrl := fmt.Sprintf("%s/%s", url, endpoint)
	resp, err = h.Client.Post(fullUrl, contentType, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
