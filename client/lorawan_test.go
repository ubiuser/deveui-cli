package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type MockClient struct {
	DoPost func(body io.Reader) (resp *http.Response, err error)
}

func TestLorawanHappyPath(t *testing.T) {
	mockClient := &MockClient{
		DoPost: func(body io.Reader) (resp *http.Response, err error) {
			return &http.Response{}, nil
		},
	}

	// loraWanClient := NewLoraWAN("mock-url", 30000)

	b := new(bytes.Buffer)
	reqBody := map[string]string{"Deveui": "Abcde"}

	_ = json.NewEncoder(b).Encode(&reqBody)

	resp, err := mockClient.Post(b)

	if err != nil {
		t.Errorf("err should be nil but is: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		t.Errorf("resp should be nil but is: %d", resp.StatusCode)
	}
}

func (m *MockClient) Post(body io.Reader) (resp *http.Response, err error) {
	return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(nil)),
			Status:     "200 OK"},
		nil
}
