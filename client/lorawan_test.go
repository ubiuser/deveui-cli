package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/NickGowdy/deveui-cli/device"
)

func TestNewLoraWAN(t *testing.T) {
	t.Parallel()
	type args struct {
		baseURL string
		timeout time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    *LoraWAN
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "invalid-base-url",
			args: args{
				baseURL: ":invalid",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to join url")
			},
		},
		{
			name: "ok",
			args: args{
				baseURL: "base",
				timeout: 1 * time.Second,
			},
			want: &LoraWAN{
				fullURL: func() *url.URL {
					u, err := url.Parse(fmt.Sprintf("%s%s", "base", endpoint))
					require.NoError(t, err)

					return u
				}(),
				client: &http.Client{
					Timeout: 1 * time.Second,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewLoraWAN(tt.args.baseURL, tt.args.timeout)
			if !tt.wantErr(t, err, fmt.Sprintf("NewLoraWAN(%v, %v)", tt.args.baseURL, tt.args.timeout)) {
				return
			}
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestLoraWAN_RegisterDevice_Request(t *testing.T) {
	t.Parallel()

	const timeout = 10 * time.Second

	t.Run("server-side-checks", func(t *testing.T) {
		t.Parallel()

		newDevice, err := device.NewDevice()
		require.NoError(t, err)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, endpoint, r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			data, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			require.NoError(t, err)
			var body struct {
				Deveui string `json:"deveui"`
			}
			err = json.Unmarshal(data, &body)
			require.NoError(t, err)
			assert.Equal(t, newDevice.GetIdentifier(), body.Deveui)

			w.WriteHeader(http.StatusCreated)
		}))
		defer ts.Close()

		client, err := NewLoraWAN(ts.URL, timeout)
		require.NoError(t, err)

		err = client.RegisterDevice(context.Background(), newDevice)
		assert.NoError(t, err)
	})
}

func TestLoraWAN_RegisterDevice(t *testing.T) {
	t.Parallel()
	type fields struct {
		handler func(w http.ResponseWriter, r *http.Request)
	}
	tests := []struct {
		name    string
		fields  fields
		want    *device.Device
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error-status-code",
			fields: fields{
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnprocessableEntity)
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "request failed")
			},
		},
		{
			name: "ok-status-code",
			fields: fields{
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
			},
			want: &device.Device{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewRouter()
			r.Post(endpoint, tt.fields.handler)
			ts := httptest.NewServer(r)
			defer ts.Close()

			tsPath, err := url.JoinPath(ts.URL, endpoint)
			require.NoError(t, err)

			tsURL, err := url.Parse(tsPath)
			require.NoError(t, err)

			l := &LoraWAN{
				fullURL: tsURL,
				client:  ts.Client(),
			}

			tt.wantErr(t, l.RegisterDevice(context.Background(), &device.Device{}), "RegisterDevice()")
		})
	}
}
