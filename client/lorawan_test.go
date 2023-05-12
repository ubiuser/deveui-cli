package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"reflect"
	"time"

	"net/http"
	"strings"
	"testing"
)

type MockClient struct {
	DoFunc func(*http.Request) (resp *http.Response, err error)
}

func TestLorawanClientHappyPath(t *testing.T) {
	mockClient := &MockClient{
		DoFunc: func(*http.Request) (resp *http.Response, err error) {
			return &http.Response{}, nil
		},
	}

	loraWAN := NewLoraWAN("www.example.com", mockClient)

	b := new(bytes.Buffer)
	reqBody := map[string]string{"Deveui": "Abcde"}

	_ = json.NewEncoder(b).Encode(&reqBody)

	ctx, cancel := context.WithCancel(context.Background())

	if cancel == nil {
		t.Errorf("cancel should not be nil but is: %v", cancel)
	}

	resp, err := loraWAN.DoPost(b, ctx)

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

func TestNewLoraWanClient(t *testing.T) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	t.Parallel()
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want *LoraWAN
	}{
		{
			name: "create-new-lorawan-client",
			args: args{
				timeout: 30,
			},
			want: &LoraWAN{
				baseURL: "https://www.example.com",
				client:  client,
			},
		},
	}
	for _, tt := range tests {
		tt := tt // it is important to capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // this makes sure that all cases from the table here are executed in parallel
			if got := NewLoraWAN("https://www.example.com", client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLoraWanClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (m *MockClient) Do(*http.Request) (resp *http.Response, err error) {
	b := new(bytes.Buffer)
	reqBody := true

	_ = json.NewEncoder(b).Encode(&reqBody)
	return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(b),
			Status:     "200 OK"},
		nil
}
