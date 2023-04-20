package processor

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/NickGowdy/deveui-cli/channel"
)

type MockClient struct {
	DoPost func(endpoint string, b *bytes.Buffer) (*http.Response, error)
}

func TestCanProcessCodes(t *testing.T) {
	codeChannel := &channel.CodeChannel{
		Msgch:  make(chan channel.Message, 10),
		Quitch: make(chan struct{}),
	}

	signalChannel := &channel.SignalChannel{}

	client := &MockClient{
		DoPost: func(endpoint string, b *bytes.Buffer) (*http.Response, error) {
			// do whatever you want
			return &http.Response{
				StatusCode: http.StatusOK,
			}, nil
		},
	}

	CodeProcessor := &CodeProcessor{
		Client:         client,
		CodeChannel:    codeChannel,
		SignalChannel:  signalChannel,
		RegisterNumber: 10,
	}

	CodeProcessor.Start()
}

func (m *MockClient) Post(endpoint string, b *bytes.Buffer) (*http.Response, error) {
	return &http.Response{}, nil
}
