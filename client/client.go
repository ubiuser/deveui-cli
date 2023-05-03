package client

import (
	"io"
	"net/http"
)

// I don't see any reason why this is in a separate file, it would be better in the other file (lorawan).

// Client generic client used to communicate to external services
type Client interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}
