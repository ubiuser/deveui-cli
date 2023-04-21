package processor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type MockClient struct {
	DoPost func(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

func (m *MockClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(nil)),
			Status:     "200 OK"},
		nil
}

func TestCanProcessCodes(t *testing.T) {
	client := &MockClient{
		DoPost: func(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
			return &http.Response{}, nil
		},
	}

	codeProcessor := &CodeProcessor{
		CodeRegistrationLimit: 10,
		MaxConcurrentJobs:     10,
		BaseUrl:               "http://www.mock-url.com",
		Client:                client,
		RegisteredDevices:     make(chan RegisterDevice),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	work := make(chan struct{}, 10)

	go func() {
		for {
			work <- struct{}{}
		}
	}()

	// Spawn workers
	for j := 0; j < 10; j++ {
		go codeProcessor.Worker(ctx, work)
	}

	n := 0
	for d := range codeProcessor.RegisteredDevices {
		fmt.Printf("device: %d has identifier: %s and code: %s\n", n+1, d.Identifier, d.Code)
		n += 1
		if n == 10 {
			break
		}
	}
}
