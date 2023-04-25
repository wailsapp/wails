//go:build linux && webkit2_36

package webview

/*
#cgo linux pkg-config: webkit2gtk-4.0

#include "webkit2/webkit2.h"
*/
import "C"

import (
	"io"
	"net/http"
)

const Webkit2MinMinorVersion = 36

func webkit_uri_scheme_request_get_http_body(_ *C.WebKitURISchemeRequest) io.ReadCloser {
	return http.NoBody
}
