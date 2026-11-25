//go:build linux && purego && !android

package webview

// CGO-free (purego) port of request_linux.go.

import (
	"io"
	"net/http"
	"unsafe"
)

func NewRequest(webKitURISchemeRequest unsafe.Pointer) Request {
	ensureWebviewLibs()

	webkitReq := uintptr(webKitURISchemeRequest)
	lifetime := retainRequest(webKitURISchemeRequest)

	req := &request{req: webkitReq, lifetime: lifetime}
	return newRequestFinalizer(req)
}

var _ Request = &request{}

type request struct {
	lifetime *requestLifetime
	req      uintptr

	header http.Header
	body   io.ReadCloser
	rw     *responseWriter
}

func (r *request) URL() (string, error) {
	// Reading the URI touches the WebKit-owned request on the GTK main loop;
	// this runs on a worker goroutine, so it must hop to the main thread.
	// See mainthread_linux_purego.go and issue #5631.
	var uri string
	invokeOnMainSync(func() {
		uri = goString(webkit_uri_scheme_request_get_uri(r.req))
	})
	return uri, nil
}

func (r *request) Method() (string, error) {
	return webkitURISchemeRequestGetHTTPMethod(r.req), nil
}

func (r *request) Header() (http.Header, error) {
	if r.header != nil {
		return r.header, nil
	}

	r.header = webkitURISchemeRequestGetHTTPHeaders(r.req)
	return r.header, nil
}

func (r *request) Body() (io.ReadCloser, error) {
	if r.body != nil {
		return r.body, nil
	}

	r.body = webkitURISchemeRequestGetHTTPBody(r.req)

	return r.body, nil
}
