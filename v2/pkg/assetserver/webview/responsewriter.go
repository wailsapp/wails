package webview

import (
	"errors"
	"net/http"
)

const (
	HeaderContentLength = "Content-Length"
	HeaderContentType   = "Content-Type"
)

var (
	errRequestStopped   = errors.New("request has been stopped")
	errResponseFinished = errors.New("response has been finished")
)

// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response for the WebView.
type ResponseWriter interface {
	http.ResponseWriter

	// Finish the response and flush all data.
	Finish() error
}
