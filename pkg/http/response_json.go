package http

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	statusCode int
	data       interface{}
	header     http.Header
	next       bool
}

func NewJsonResponse(statusCode int, data interface{}) (response Response) {
	response = &jsonResponse{
		statusCode: statusCode,
		data:       data,
		header:     http.Header{},
	}
	return
}

func (jsonResponse *jsonResponse) StatusCode() int {
	return jsonResponse.statusCode
}

func (jsonResponse *jsonResponse) Body() ([]byte, error) {
	b, err := json.Marshal(jsonResponse.data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (jsonResponse *jsonResponse) Header() http.Header {
	return jsonResponse.header
}

func (jsonResponse *jsonResponse) ContentType() string {
	return "application/json; charset=utf-8"
}
