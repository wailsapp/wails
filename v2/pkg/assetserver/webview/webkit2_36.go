//go:build linux && webkit2_36

package webview

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0 libsoup-2.4

#include "gtk/gtk.h"
#include "webkit2/webkit2.h"
#include "libsoup/soup.h"
*/
import "C"

import (
	"net/http"
	"strings"
	"unsafe"
)

const Webkit2MinMinorVersion = 36

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
