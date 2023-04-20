package client

import (
	"bytes"
	"net/http"
)

type Client interface {
	Post(endpoint string, b *bytes.Buffer) (*http.Response, error)
}
