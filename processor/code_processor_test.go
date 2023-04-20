package processor

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

type MockClient struct {
	DoPost func(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

func TestCanProcessCodes(t *testing.T) {
	client := &MockClient{
		DoPost: func(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
			return &http.Response{}, nil
		},
	}

	CodeProcessor := &CodeProcessor{
		CodeRegistrationLimit: 0,
		MaxConcurrentJobs:     10,
		BaseUrl:               "http://www.mock-url.com",
		Client:                client,
	}

	reader := bufio.NewReader(os.Stdin)
	CodeProcessor.Start()

	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}

func (m *MockClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(nil)),
			Status:     "200 OK"},
		nil
}
