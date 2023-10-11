package server

import "net/http"

type Consumer struct{}

func (c Consumer) Header() http.Header {
	return http.Header{}
}

func (c Consumer) Write(data []byte) (int, error) {
	return len(data), nil
}

func (c Consumer) WriteHeader(statusCode int) {
}
