package client

import (
	"bytes"
	"encoding/json"
	"io"

	"net/http"
	"strings"
	"testing"
	"time"
)

type MockClient struct {
	DoPost func(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

func TestLorawanClientHappyPath(t *testing.T) {
	mockClient := &MockClient{
		DoPost: func(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
			return &http.Response{}, nil
		},
	}

	loraWanClient := NewLoraWAN("www.example.com", time.Microsecond*30000)

	b := new(bytes.Buffer)
	reqBody := map[string]string{"Deveui": "Abcde"}

	_ = json.NewEncoder(b).Encode(&reqBody)

	resp, err := loraWanClient.DoPost(mockClient, b)

	if err != nil {
		t.Errorf("err should be nil but is: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("resp should be nil but is: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	val := string(body)

	if strings.TrimSpace(val) != "true" {
		t.Errorf("body should equal true but is: %d", body)
	}
}

func (m *MockClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	b := new(bytes.Buffer)
	reqBody := true

	_ = json.NewEncoder(b).Encode(&reqBody)
	return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(b),
			Status:     "200 OK"},
		nil
}
