//go:build linux

package webview

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.1 libsoup-3.0

#include "gtk/gtk.h"
#include "webkit2/webkit2.h"
#include "libsoup/soup.h"
*/
import "C"

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unsafe"
)

const Webkit2MinMinorVersion = 40

func webkit_uri_scheme_request_get_http_method(req *C.WebKitURISchemeRequest) string {
	method := C.GoString(C.webkit_uri_scheme_request_get_http_method(req))
	return strings.ToUpper(method)
}

func webkit_uri_scheme_request_get_http_headers(req *C.WebKitURISchemeRequest) http.Header {
	hdrs := C.webkit_uri_scheme_request_get_http_headers(req)

	var iter C.SoupMessageHeadersIter
	C.soup_message_headers_iter_init(&iter, hdrs)

	var name *C.char
	var value *C.char

	h := http.Header{}
	for C.soup_message_headers_iter_next(&iter, &name, &value) != 0 {
		h.Add(C.GoString(name), C.GoString(value))
	}

	return h
}

func webkit_uri_scheme_request_finish(req *C.WebKitURISchemeRequest, code int, header http.Header, stream *C.GInputStream, streamLength int64) error {
	resp := C.webkit_uri_scheme_response_new(stream, C.gint64(streamLength))
	defer C.g_object_unref(C.gpointer(resp))

	cReason := C.CString(http.StatusText(code))
	C.webkit_uri_scheme_response_set_status(resp, C.guint(code), cReason)
	C.free(unsafe.Pointer(cReason))

	cMimeType := C.CString(header.Get(HeaderContentType))
	C.webkit_uri_scheme_response_set_content_type(resp, cMimeType)
	C.free(unsafe.Pointer(cMimeType))

	hdrs := C.soup_message_headers_new(C.SOUP_MESSAGE_HEADERS_RESPONSE)
	for name, values := range header {
		cName := C.CString(name)
		for _, value := range values {
			cValue := C.CString(value)
			C.soup_message_headers_append(hdrs, cName, cValue)
			C.free(unsafe.Pointer(cValue))
		}
		C.free(unsafe.Pointer(cName))
	}

	C.webkit_uri_scheme_response_set_http_headers(resp, hdrs)

	C.webkit_uri_scheme_request_finish_with_response(req, resp)
	return nil
}

func webkit_uri_scheme_request_get_http_body(req *C.WebKitURISchemeRequest) io.ReadCloser {
	stream := C.webkit_uri_scheme_request_get_http_body(req)
	if stream == nil {
		return http.NoBody
	}
	return &webkitRequestBody{stream: stream}
}

type webkitRequestBody struct {
	stream *C.GInputStream
	closed bool
}

// Read implements io.Reader
func (r *webkitRequestBody) Read(p []byte) (int, error) {
	if r.closed {
		return 0, io.ErrClosedPipe
	}

	var content unsafe.Pointer
	var contentLen int
	if p != nil {
		content = unsafe.Pointer(&p[0])
		contentLen = len(p)
	}

	var n C.gsize
	var gErr *C.GError
	res := C.g_input_stream_read_all(r.stream, content, C.gsize(contentLen), &n, nil, &gErr)
	if res == 0 {
		return 0, formatGError("stream read failed", gErr)
	} else if n == 0 {
		return 0, io.EOF
	}
	return int(n), nil
}

func (r *webkitRequestBody) Close() error {
	if r.closed {
		return nil
	}
	r.closed = true

	// https://docs.gtk.org/gio/method.InputStream.close.html
	// Streams will be automatically closed when the last reference is dropped, but you might want to call this function
	// to make sure resources are released as early as possible.
	var err error
	var gErr *C.GError
	if C.g_input_stream_close(r.stream, nil, &gErr) == 0 {
		err = formatGError("stream close failed", gErr)
	}
	C.g_object_unref(C.gpointer(r.stream))
	r.stream = nil
	return err
}

func formatGError(msg string, gErr *C.GError, args ...any) error {
	if gErr != nil && gErr.message != nil {
		msg += ": " + C.GoString(gErr.message)
		C.g_error_free(gErr)
	}
	return fmt.Errorf(msg, args...)
}
