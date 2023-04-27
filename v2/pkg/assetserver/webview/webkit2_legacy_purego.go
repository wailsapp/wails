//go:build linux && !(webkit2_36 || webkit2_40) && purego

package webview

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ebitengine/purego"
)

const Webkit2MinMinorVersion = 0

func webkit_uri_scheme_request_get_http_method(_ uintptr) string {
	return http.MethodGet
}

func webkit_uri_scheme_request_get_http_headers(_ uintptr) http.Header {
	return http.Header{}
}

func webkit_uri_scheme_request_get_http_body(_ uintptr) io.ReadCloser {
	return http.NoBody
}

func webkit_uri_scheme_request_finish(req uintptr, code int, header http.Header, stream uintptr, streamLength int64) error {
	if code != http.StatusOK {
		return fmt.Errorf("StatusCodes not supported: %d - %s", code, http.StatusText(code))
	}

	var requestFinish func(uintptr, uintptr, int64, string)
	purego.RegisterLibFunc(&requestFinish, webkit, "webkit_uri_scheme_request_finish")
	requestFinish(req, stream, streamLength, header.Get(HeaderContentType))
	return nil
}
