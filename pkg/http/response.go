package http

import "net/http"

type Response interface {
	StatusCode() int
	Body() ([]byte, error)
	Header() http.Header
	ContentType() string
}
