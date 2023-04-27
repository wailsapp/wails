//go:build linux && (webkit2_36 || webkit2_40) && purego

package webview

import (
	"net/http"
	"strings"

	"github.com/ebitengine/purego"
)

func webkit_uri_scheme_request_get_http_method(req uintptr) string {
	var getMethod func(uintptr) string
	purego.RegisterLibFunc(&getMethod, gtk, "webkit_uri_scheme_request_get_http_method")
	return strings.ToUpper(getMethod(req))
}

func webkit_uri_scheme_request_get_http_headers(req uintptr) http.Header {
	var getHeaders func(uintptr) uintptr
	purego.RegisterLibFunc(&getUri, webkit, "webkit_uri_scheme_request_get_http_headers")

	hdrs := getHeaders(req)

	var headersIterInit func(uintptr, uintptr) uintptr
	purego.RegisterLibFunc(&headersIterInit, gtk, "soup_message_headers_iter_init")

	// TODO: How do we get a struct?
	/*
	   typedef struct {
	   	SoupMessageHeaders *hdrs;
	           int index_common;
	   	int index_uncommon;
	   } SoupMessageHeadersIterReal;
	*/
	iter := make([]byte, 12)
	headersIterInit(&iter, hdrs)

	var iterNext func(uintptr, *string, *string) int
	purego.RegisterLibFunc(&iterNext, gtk, "soup_message_headers_iter_next")

	var name string
	var value string
	h := http.Header{}

	for iterNext(&iter, &name, &value) != 0 {
		h.Add(name, value)
	}

	return h
}

func webkit_uri_scheme_request_finish(req uintptr, code int, header http.Header, stream uintptr, streamLength int64) error {

	var newResponse func(uintptr, int64) string
	purego.RegisterLibFunc(&newResponse, webkit, "webkit_uri_scheme_response_new")
	var unRef func(uintptr)
	purego.RegisterLibFunc(&unRef, gtk, "g_object_unref")

	resp := newResponse(stream, streamLength)
	defer unRef(resp)

	var setStatus func(uintptr, int, string)
	purego.RegisterLibFunc(&unRef, webkit, "webkit_uri_scheme_response_set_status")

	setStatus(resp, code, cReason)

	var setContentType func(uintptr, string)
	purego.RegisterLibFunc(&unRef, webkit, "webkit_uri_scheme_response_set_content_type")

	setContentType(resp, header.Get(HeaderContentType))

	soup := gtk
	var soupHeadersNew func(int) uintptr
	purego.RegisterLibFunc(&unRef, soup, "soup_message_headers_new")
	var soupHeadersAppend func(uintptr, string, string)
	purego.RegisterLibFunc(&unRef, soup, "soup_message_headers_append")

	hdrs := soupHeadersNew(SOUP_MESSAGE_HEADERS_RESPONSE)
	for name, values := range header {
		for _, value := range values {
			soupHeadersAppend(hdrs, name, value)
		}
	}

	var setHttpHeaders func(uintptr, uintptr)
	purego.RegisterLibFunc(&unRef, webkit, "webkit_uri_scheme_response_set_http_headers")

	setHttpHeaders(resp, hdrs)
	var finishWithResponse func(uintptr, uintptr)
	purego.RegisterLibFunc(&unRef, webkit, "webkit_uri_scheme_request_finish_with_response")
	finishWithResponse(req, resp)

	return nil
}
