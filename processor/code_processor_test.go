package processor

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/device"
)

type MockClient struct {
	DoFunc func(*http.Request) (resp *http.Response, err error)
}

const (
	MAX_CONCURRENT_JOBS     = 2
	CODE_REGISTRATION_LIMIT = 10
)

func TestCanProcessCodes(t *testing.T) {
	mockClient := &MockClient{
		DoFunc: func(*http.Request) (resp *http.Response, err error) {
			return &http.Response{}, nil
		},
	}

	loraWAN := client.NewLoraWAN("www.example.com", mockClient)

	codeProcessor := &CodeProcessor{
		CodeRegistrationLimit: CODE_REGISTRATION_LIMIT,
		MaxConcurrentJobs:     MAX_CONCURRENT_JOBS,
		Device:                make(chan device.Device),
		LoraWAN:               *loraWAN,
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

		identifier := d.GetIdentifier()
		code := d.GetCode()

		if code == "" {
			t.Error("code should not be nil")
		}

		if identifier == "" {
			t.Error("identifier should not be nil")
		}

		if identifier[len(identifier)-5:] != code {
			t.Errorf("code should be last 5 characters of identifier, but is %s", code)
		}

		if len(identifier) != 16 {
			t.Errorf("identifier should be exactly 16 characters, but is %d", len(identifier))
		}

		if len(code) != 5 {
			t.Errorf("code should be exactly 5 characters, but is %d", len(code))
		}

		n += 1
		if n == CODE_REGISTRATION_LIMIT {
			break
		}
	}
}

func (m *MockClient) Do(*http.Request) (resp *http.Response, err error) {
	return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(nil)),
			Status:     "200 OK"},
		nil
}
