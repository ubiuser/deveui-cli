package processor

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
			return &http.Response{}, nil
		},
	}

	CodeProcessor := &CodeProcessor{
		Client:         client,
		CodeChannel:    codeChannel,
		SignalChannel:  signalChannel,
		RegisterNumber: 10,
	}

	reader := bufio.NewReader(os.Stdin)
	CodeProcessor.Start()

	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}

func (m *MockClient) Post(endpoint string, b *bytes.Buffer) (*http.Response, error) {
	return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(nil)),
			Status:     "200 OK"},
		nil
}
