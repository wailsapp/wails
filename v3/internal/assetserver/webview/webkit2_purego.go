//go:build linux && purego && !android

package webview

import (
	"net/http"
	"strings"
	"unsafe"

	"github.com/ebitengine/purego"
)

// SoupMessageHeadersType value for response headers
// (libsoup soup-message-headers.h: REQUEST = 0, RESPONSE = 1, MULTIPART = 2).
const SOUP_MESSAGE_HEADERS_RESPONSE = 1

// goString converts a NUL-terminated C string to a Go string. Needed for
// C out-parameters, where purego's automatic string conversion does not apply.
func goString(c uintptr) string {
	if c == 0 {
		return ""
	}
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&c))
	n := 0
	for *(*byte)(unsafe.Add(ptr, n)) != 0 {
		n++
	}
	return string(unsafe.Slice((*byte)(ptr), n))
}

func webkit_uri_scheme_request_get_http_method(req uintptr) string {
	var getMethod func(uintptr) string
	purego.RegisterLibFunc(&getMethod, webkit, "webkit_uri_scheme_request_get_http_method")
	return strings.ToUpper(getMethod(req))
}

func webkit_uri_scheme_request_get_http_headers(req uintptr) http.Header {
	var getHeaders func(uintptr) uintptr
	purego.RegisterLibFunc(&getHeaders, webkit, "webkit_uri_scheme_request_get_http_headers")

	h := http.Header{}

	hdrs := getHeaders(req)
	if hdrs == 0 {
		return h
	}

	// SoupMessageHeadersIter is declared as `gpointer dummy[3]` in libsoup;
	// allocate matching storage for the out-parameter.
	var iter [3]uintptr

	// The soup_* symbols live in libsoup, which is resolved through the
	// webkit handle (libwebkit2gtk links libsoup).
	var headersIterInit func(*[3]uintptr, uintptr)
	purego.RegisterLibFunc(&headersIterInit, webkit, "soup_message_headers_iter_init")
	headersIterInit(&iter, hdrs)

	var iterNext func(*[3]uintptr, *uintptr, *uintptr) int32
	purego.RegisterLibFunc(&iterNext, webkit, "soup_message_headers_iter_next")

	var name, value uintptr
	for iterNext(&iter, &name, &value) != 0 {
		h.Add(goString(name), goString(value))
	}

	return h
}

func webkit_uri_scheme_request_finish(req uintptr, code int, header http.Header, stream uintptr, streamLength int64) error {
	var newResponse func(uintptr, int64) uintptr
	purego.RegisterLibFunc(&newResponse, webkit, "webkit_uri_scheme_response_new")
	var unRef func(uintptr)
	purego.RegisterLibFunc(&unRef, gtk, "g_object_unref")

	resp := newResponse(stream, streamLength)
	defer unRef(resp)

	var setStatus func(uintptr, uint32, string)
	purego.RegisterLibFunc(&setStatus, webkit, "webkit_uri_scheme_response_set_status")
	setStatus(resp, uint32(code), http.StatusText(code))

	var setContentType func(uintptr, string)
	purego.RegisterLibFunc(&setContentType, webkit, "webkit_uri_scheme_response_set_content_type")
	setContentType(resp, header.Get(HeaderContentType))

	var soupHeadersNew func(int32) uintptr
	purego.RegisterLibFunc(&soupHeadersNew, webkit, "soup_message_headers_new")
	var soupHeadersAppend func(uintptr, string, string)
	purego.RegisterLibFunc(&soupHeadersAppend, webkit, "soup_message_headers_append")

	// Ownership of hdrs is transferred to the response by
	// webkit_uri_scheme_response_set_http_headers (transfer full), so we must
	// not unref it here — doing so frees the headers while WebKit/libsoup still
	// reference them, crashing in soup_message_headers_iter_next on render.
	hdrs := soupHeadersNew(SOUP_MESSAGE_HEADERS_RESPONSE)
	for name, values := range header {
		for _, value := range values {
			soupHeadersAppend(hdrs, name, value)
		}
	}

	var setHttpHeaders func(uintptr, uintptr)
	purego.RegisterLibFunc(&setHttpHeaders, webkit, "webkit_uri_scheme_response_set_http_headers")
	setHttpHeaders(resp, hdrs)

	var finishWithResponse func(uintptr, uintptr)
	purego.RegisterLibFunc(&finishWithResponse, webkit, "webkit_uri_scheme_request_finish_with_response")
	finishWithResponse(req, resp)

	return nil
}
