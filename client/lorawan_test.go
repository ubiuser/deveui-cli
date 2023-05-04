package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type MockClient struct {
	DoPost func(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

func TestLorawanHappyPath(t *testing.T) {
	mockClient := &MockClient{
		DoPost: func(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
			return &http.Response{}, nil
		},
	}

	loraWanClient := NewLoraWAN(mockClient)

	b := new(bytes.Buffer)
	reqBody := map[string]string{"Deveui": "Abcde"}

	_ = json.NewEncoder(b).Encode(&reqBody)

	resp, err := loraWanClient.Post("mock-url", "application/json", b)

	if err != nil {
		t.Errorf("err should be nil but is: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		t.Errorf("resp should be nil but is: %d", resp.StatusCode)
	}
}

func (m *MockClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(nil)),
			Status:     "200 OK"},
		nil
}
