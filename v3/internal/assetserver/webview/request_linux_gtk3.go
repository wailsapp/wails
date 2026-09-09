//go:build linux && cgo && gtk3 && !android

package webview

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.1 gio-unix-2.0

#include "gtk/gtk.h"
#include "webkit2/webkit2.h"

static gboolean unref_request_on_main(gpointer data) {
	if (data != NULL) {
		g_object_unref(data);
	}
	return G_SOURCE_REMOVE;
}

// releaseRequestOnMainThread schedules the WebKitURISchemeRequest unref on the
// GTK main context. Close() runs on the assetserver goroutine, and dropping
// what may be the last reference finalizes a WebKit GObject — only safe on the
// UI thread (see #5557).
static void releaseRequestOnMainThread(WebKitURISchemeRequest *request) {
	if (request == NULL) {
		return;
	}
	g_main_context_invoke(NULL, unref_request_on_main, request);
}
*/
import "C"

import (
	"io"
	"net/http"
	"unsafe"
)

// NewRequest creates as new WebViewRequest based on a pointer to an `WebKitURISchemeRequest`
func NewRequest(webKitURISchemeRequest unsafe.Pointer) Request {
	webkitReq := (*C.WebKitURISchemeRequest)(webKitURISchemeRequest)
	C.g_object_ref(C.gpointer(webkitReq))

	req := &request{req: webkitReq}
	return newRequestFinalizer(req)
}

var _ Request = &request{}

type request struct {
	req *C.WebKitURISchemeRequest

	header http.Header
	body   io.ReadCloser
	rw     *responseWriter
}

func (r *request) URL() (string, error) {
	// Reading the URI touches the WebKit-owned request on the GTK main loop;
	// this runs on a worker goroutine, so it must hop to the main thread.
	// See mainthread_linux.go and issue #5631.
	var uri string
	invokeOnMainSync(func() {
		uri = C.GoString(C.webkit_uri_scheme_request_get_uri(r.req))
	})
	return uri, nil
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

	r.body = webkit_uri_scheme_request_get_http_body(r.req)

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
	C.releaseRequestOnMainThread(r.req)
	return err
}
