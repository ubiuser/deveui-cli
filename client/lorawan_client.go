package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// I'd be inclined to rename this to just lorawan.go, because it's already in the client package.
// comments in Go start with the name of the described entity

// LoraWanClient used to communicate to LoRaWAN external system
type LoraWanClient struct {
	client *http.Client // the internal field can be the direct type, the struct itself will implement the interface
}

const endpoint = "/sensor-onboarding-sample" // endpoint for saving DevEUI via LoRaWAN

// Unless you want to share the same http client with other code, it's better to hide it as an internal field.
// The users of the client package are not interested in it anyway, they only care about the interface methods.

// NewLoraWanClient creates a new LoraWanClient that implements the Client interface
func NewLoraWanClient(timeout time.Duration) *LoraWanClient {
	return &LoraWanClient{
		client: &http.Client{
			Timeout: timeout * time.Second,
		},
	}
}

// Post sends data via POST (HTTP) request
func (h *LoraWanClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	// Probably use url.JoinPath instead of fmt.Sprintf
	// Also, I think your url won't change throughout the app, so could be on the struct as an internal field, because
	// at the moment it's calculated with every call of the Post method.
	// contentType also doesn't change, could just be "application/json" as a const
	fullUrl := fmt.Sprintf("%s/%s", url, endpoint)
	resp, err = h.client.Post(fullUrl, contentType, body)
	if err != nil {
		// some people find it annoying, but I prefer wrapping errors, which can greatly help during debugging
		return nil, err
	}

	return resp, nil

	// Another way to further improve this is to use the http.Client.Do method that requires a http.Request.
	// The request can be created with the http.NewRequestWithContext method, so you can pass down the context
	// all the way. When the parent context is cancelled, the request will be cancelled too.
}
