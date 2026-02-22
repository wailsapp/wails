//go:build android

package webview

import (
	"bytes"
	"io"
	"net/http"
)

// Request interface for Android asset requests
// On Android, requests are handled via JNI from Java's WebViewAssetLoader

// androidRequest implements the Request interface for Android
type androidRequest struct {
	url     string
	method  string
	headers http.Header
	body    io.ReadCloser
	rw      *androidResponseWriter
}

// NewRequestFromJNI creates a new request from JNI parameters
func NewRequestFromJNI(url string, method string, headersJSON string) Request {
	return &androidRequest{
		url:     url,
		method:  method,
		headers: http.Header{},
		body:    http.NoBody,
	}
}

func (r *androidRequest) URL() (string, error) {
	return r.url, nil
}

func (r *androidRequest) Method() (string, error) {
	return r.method, nil
}

func (r *androidRequest) Header() (http.Header, error) {
	return r.headers, nil
}

func (r *androidRequest) Body() (io.ReadCloser, error) {
	return r.body, nil
}

func (r *androidRequest) Response() ResponseWriter {
	if r.rw == nil {
		r.rw = &androidResponseWriter{}
	}
	return r.rw
}

func (r *androidRequest) Close() error {
	if r.body != nil {
		return r.body.Close()
	}
	return nil
}

// androidResponseWriter implements ResponseWriter for Android
type androidResponseWriter struct {
	statusCode int
	headers    http.Header
	body       bytes.Buffer
	finished   bool
}

func (w *androidResponseWriter) Header() http.Header {
	if w.headers == nil {
		w.headers = http.Header{}
	}
	return w.headers
}

func (w *androidResponseWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

func (w *androidResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *androidResponseWriter) Finish() error {
	w.finished = true
	return nil
}

// Code returns the HTTP status code of the response
func (w *androidResponseWriter) Code() int {
	if w.statusCode == 0 {
		return 200
	}
	return w.statusCode
}

// GetResponseData returns the response data for JNI
func (w *androidResponseWriter) GetResponseData() []byte {
	return w.body.Bytes()
}
