//go:build linux && purego && !gtk3 && !android

package webview

// CGO-free (purego) port of request_linux.go.

import (
	"io"
	"net/http"
	"unsafe"

	"github.com/ebitengine/purego"
)

// unrefRequestOnMain runs on the GTK main thread (scheduled via
// g_main_context_invoke) and drops the reference taken in NewRequest. Created
// once as a package variable: purego callback slots are a process-wide,
// never-freed resource, so a per-call NewCallback would leak them.
var unrefRequestOnMain = purego.NewCallback(func(data uintptr) uintptr {
	if data != 0 {
		g_object_unref(data)
	}
	return 0 // G_SOURCE_REMOVE
})

// releaseRequestOnMainThread schedules the WebKitURISchemeRequest unref on the
// GTK main context. Close() runs on the assetserver goroutine, and dropping
// what may be the last reference finalizes a WebKit GObject — only safe on the
// UI thread (see #5557).
func releaseRequestOnMainThread(request uintptr) {
	if request == 0 {
		return
	}
	g_main_context_invoke(0, unrefRequestOnMain, request)
}

func NewRequest(webKitURISchemeRequest unsafe.Pointer) Request {
	ensureWebviewLibs()

	webkitReq := uintptr(webKitURISchemeRequest)
	g_object_ref(webkitReq)

	req := &request{req: webkitReq}
	return newRequestFinalizer(req)
}

var _ Request = &request{}

type request struct {
	req uintptr

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
	releaseRequestOnMainThread(r.req)
	return err
}
