package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	DoPost func(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

// TestLorawanClientHappyPath doesn't seem to test anything really.
// You should use httptest.NewServer to test your Post method.
func TestLorawanClientHappyPath(t *testing.T) {
	//mockClient := &MockClient{
	//	DoPost: func(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	//		return &http.Response{}, nil
	//	},
	//}

	loraWanClient := LoraWanClient{
		//client: mockClient,
	}

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

// unused parameters should be removed or replaced with underscores
func (m *MockClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(nil)),
			Status:     "200 OK"},
		nil
}

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	m.Called(url, contentType, body)

	return
}

// This is a table test, my Goland IDE can generate the skeleton for methods easily.
func TestNewLoraWanClient(t *testing.T) {
	t.Parallel() // this makes sure that TestNewLoraWanClient is executed in parallel with other tests
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want *LoraWanClient
	}{
		{
			name: "create-new-lorawan-client",
			args: args{
				timeout: 30,
			},
			want: &LoraWanClient{
				client: &http.Client{
					Timeout: 30 * time.Second,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt // it is important to capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // this makes sure that all cases from the table here are executed in parallel
			if got := NewLoraWanClient(tt.args.timeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLoraWanClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
