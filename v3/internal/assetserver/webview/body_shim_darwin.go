//go:build darwin && !ios

package webview

import (
	"encoding/base64"
	"io"
	"net/http"
)

// WKWebView hands a wails:// scheme handler no body at all (neither HTTPBody
// nor HTTPBodyStream) when the page sends a binary body: ArrayBuffer, typed
// arrays, Blob, FormData. String bodies arrive. The document-start shim
// (pkg/application/darwin_body_fetch_shim.js) therefore sends binary bodies
// as base64 text, marked with these headers; both macOS backends undo it here
// so the asset server and the app see the original request.
const (
	bodyEncodingHeader = "X-Wails-Body-Encoding"
	bodyTypeHeader     = "X-Wails-Body-Type"
)

// unshimHeader restores the original headers of a shimmed request in place and
// reports whether its body is base64 encoded.
func unshimHeader(h http.Header) bool {
	if h.Get(bodyEncodingHeader) != "base64" {
		return false
	}
	if contentType := h.Get(bodyTypeHeader); contentType != "" {
		h.Set("Content-Type", contentType)
	} else {
		h.Del("Content-Type")
	}
	h.Del(bodyEncodingHeader)
	h.Del(bodyTypeHeader)
	h.Del("Content-Length") // of the encoded text
	return true
}

type decodedBody struct {
	io.Reader
	io.Closer
}

// unshimBody decodes a shimmed body.
func unshimBody(body io.ReadCloser) io.ReadCloser {
	if body == nil || body == http.NoBody {
		return body
	}
	return decodedBody{Reader: base64.NewDecoder(base64.StdEncoding, body), Closer: body}
}
