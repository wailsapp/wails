//go:build linux && purego
// +build linux,purego

package webview

import (
	"io"
	"net/http"

	"github.com/ebitengine/purego"
)

// NewRequest creates as new WebViewRequest based on a pointer to an `WebKitURISchemeRequest`
//
// Please make sure to call Release() when finished using the request.
func NewRequest(webKitURISchemeRequest uintptr) Request {
	webkitReq := webKitURISchemeRequest
	req := &request{req: webkitReq}
	req.AddRef()
	return req
}

var _ Request = &request{}

type request struct {
	req uintptr

	header http.Header
	body   io.ReadCloser
	rw     *responseWriter
}

func (r *request) AddRef() error {
	var objectRef func(uintptr)
	purego.RegisterLibFunc(&objectRef, gtk, "g_object_ref")
	objectRef(r.req)
	return nil
}

func (r *request) Release() error {
	var objectUnref func(uintptr)
	purego.RegisterLibFunc(&objectUnref, gtk, "g_object_unref")
	objectUnref(r.req)
	return nil
}

func (r *request) URL() (string, error) {
	var getUri func(uintptr) string
	purego.RegisterLibFunc(&getUri, webkit, "webkit_uri_scheme_request_get_uri")
	return getUri(r.req), nil
}

func (r *request) Method() (string, error) {
	return webkit_uri_scheme_request_get_http_method(r.req), nil
}

func (r *request) Header() (http.Header, error) {
	if r.header != nil {
		return r.header, nil
	}

	r.header = webkit_uri_scheme_request_get_http_headers(r.req)
	return r.header, nil
}

func (r *request) Body() (io.ReadCloser, error) {
	if r.body != nil {
		return r.body, nil
	}

	// WebKit2GTK has currently no support for request bodies.
	r.body = http.NoBody

	return r.body, nil
}

func (r *request) Response() ResponseWriter {
	if r.rw != nil {
		return r.rw
	}

	r.rw = &responseWriter{req: r.req}
	return r.rw
}

func (r *request) Close() error {
	var err error
	if r.body != nil {
		err = r.body.Close()
	}
	r.Response().Finish()
	r.Release()
	return err
}
