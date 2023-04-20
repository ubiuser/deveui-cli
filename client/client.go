package client

import (
	"io"
	"net/http"
)

type Client interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}
