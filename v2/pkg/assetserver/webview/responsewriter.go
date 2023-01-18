package webview

import (
	"net/http"
)

// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response for the WebView.
type ResponseWriter interface {
	http.ResponseWriter

	// Finish the response and flush all data.
	Finish() error
}
