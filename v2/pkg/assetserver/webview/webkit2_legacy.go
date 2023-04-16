//go:build linux && !(webkit2_36 || webkit2_40)

package webview

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include "gtk/gtk.h"
#include "webkit2/webkit2.h"
*/
import "C"

import (
	"fmt"
	"io"
	"net/http"
	"unsafe"
)

const Webkit2MinMinorVersion = 0

func webkit_uri_scheme_request_get_http_method(_ *C.WebKitURISchemeRequest) string {
	return http.MethodGet
}

func webkit_uri_scheme_request_get_http_headers(_ *C.WebKitURISchemeRequest) http.Header {
	return http.Header{}
}

func webkit_uri_scheme_request_get_http_body(_ *C.WebKitURISchemeRequest) io.ReadCloser {
	return http.NoBody
}

func webkit_uri_scheme_request_finish(req *C.WebKitURISchemeRequest, code int, header http.Header, stream *C.GInputStream, streamLength int64) error {
	if code != http.StatusOK {
		return fmt.Errorf("StatusCodes not supported: %d - %s", code, http.StatusText(code))
	}

	cMimeType := C.CString(header.Get(HeaderContentType))
	C.webkit_uri_scheme_request_finish(req, stream, C.gint64(streamLength), cMimeType)
	C.free(unsafe.Pointer(cMimeType))
	return nil
}
