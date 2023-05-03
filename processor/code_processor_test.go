package processor

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/NickGowdy/deveui-cli/device"
)

// for mocking in tests I highly recommend using the github.com/stretchr/testify/mock package

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

const (
	MAX_CONCURRENT_JOBS     = 2
	CODE_REGISTRATION_LIMIT = 10
)

func TestCanProcessCodes(t *testing.T) {
	client := &MockClient{
		DoPost: func(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
			return &http.Response{}, nil
		},
	}

	codeProcessor := &CodeProcessor{
		CodeRegistrationLimit: CODE_REGISTRATION_LIMIT,
		MaxConcurrentJobs:     MAX_CONCURRENT_JOBS,
		BaseUrl:               "http://www.mock-url.com",
		Client:                client,
		Device:                make(chan device.Device),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	work := make(chan struct{}, MAX_CONCURRENT_JOBS)

	go func() {
		for {
			work <- struct{}{}
		}
	}()

	// Spawn workers
	for j := 0; j < MAX_CONCURRENT_JOBS; j++ {
		go codeProcessor.Worker(ctx, work)
	}

	n := 0
	for d := range codeProcessor.Device {

		if d.Code == "" {
			t.Error("code should not be nil")
		}

		if d.Identifier == "" {
			t.Error("identifier should not be nil")
		}

		if d.Identifier[len(d.Identifier)-5:] != d.Code {
			t.Errorf("code should be last 5 characters of identifier, but is %s", d.Code)
		}

		if len(d.Identifier) != 16 {
			t.Errorf("identifier should be exactly 16 characters, but is %d", len(d.Identifier))
		}

		if len(d.Code) != 5 {
			t.Errorf("code should be exactly 5 characters, but is %d", len(d.Code))
		}

		n += 1
		if n == CODE_REGISTRATION_LIMIT {
			break
		}
	}
}
