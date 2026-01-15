package gosumsub

import (
	"io"
	"net/http"
	"net/url"
)

type request struct {
	Method   string
	Endpoint string
	Params   any
	Query    url.Values
	Header   http.Header
	Body     io.Reader
	FullURL  string
}
