package webview

import (
	"io"
	"net/http"
)

type Request interface {
	URL() (string, error)
	Method() (string, error)
	Header() (http.Header, error)
	Body() (io.ReadCloser, error)

	Response() ResponseWriter

	Close() error
}
