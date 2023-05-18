package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/NickGowdy/deveui-cli/device"
)

// LoraWAN used to communicate to LoRaWAN external system
type LoraWAN struct {
	fullURL *url.URL
	client  *http.Client
}

const endpoint = "/sensor-onboarding-sample" // endpoint for saving DevEUI via LoRaWAN

func NewLoraWAN(baseURL string, timeout time.Duration) (*LoraWAN, error) {
	path, err := url.JoinPath(baseURL, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to join url: %w", err)
	}

	fullURL, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	return &LoraWAN{
		fullURL: fullURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

// RegisterDevice registers a new device using LoraWAN external service
func (l *LoraWAN) RegisterDevice(ctx context.Context) (*device.Device, error) {
	newDevice := device.NewDevice()
	reqBody := map[string]string{"deveui": newDevice.GetIdentifier()}

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(&reqBody); err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", l.fullURL.String(), b)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}

	return newDevice, nil
}
