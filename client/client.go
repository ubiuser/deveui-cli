package client

import (
	"io"
	"net/http"
)

// Generic client used to communicate to external services
type Client interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}
