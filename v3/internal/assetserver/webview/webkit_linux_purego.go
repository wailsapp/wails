//go:build linux && purego && !gtk3 && !android

package webview

// CGO-free (purego) port of webkit_linux.go. Every WebKit/GObject/libsoup
// function is dlopen(3)ed at runtime and bound with purego.RegisterFunc.
//
// Conventions (shared with the pkg/application purego backend):
//   - All C pointers are uintptr.
//   - gboolean is int32 on the C side; compare returns with != 0.
//   - C string ARGS are declared as Go string (purego marshals them).
//   - Const char* RETURNS must not be freed — declare uintptr and copy with
//     goString.

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"unsafe"
)

const Webkit2MinMinorVersion = 0

// SOUP_MESSAGE_HEADERS_RESPONSE from libsoup's SoupMessageHeadersType enum.
const soupMessageHeadersResponse = 1

// soupMessageHeadersIter mirrors SoupMessageHeadersIter: an opaque struct of
// three pointers, allocated by the caller and initialised by
// soup_message_headers_iter_init.
type soupMessageHeadersIter struct {
	_ [3]uintptr
}

// gError mirrors GError: guint32 domain, gint32 code, char *message.
type gError struct {
	domain  uint32
	code    int32
	message uintptr
}

// ----------------------------------------------------------------------------
// Library loading
// ----------------------------------------------------------------------------

var webviewLibsOnce sync.Once

var (
	// glib
	g_free              func(uintptr)
	g_error_free        func(gErr uintptr)
	g_error_new_literal func(domain uint32, code int32, message string) uintptr
	g_quark_from_string func(str string) uint32

	// gobject
	g_object_ref   func(obj uintptr) uintptr
	g_object_unref func(obj uintptr)

	// gio (GUnixInputStream lives in libgio's UNIX API)
	g_unix_input_stream_new func(fd int32, closeFD int32) uintptr
	g_input_stream_read_all func(stream, buffer, count, bytesRead, cancellable, gErr uintptr) int32
	g_input_stream_close    func(stream, cancellable, gErr uintptr) int32

	// webkit
	webkit_uri_scheme_request_get_uri              func(req uintptr) uintptr // const return, do NOT free
	webkit_uri_scheme_request_get_http_method      func(req uintptr) uintptr // const return, do NOT free
	webkit_uri_scheme_request_get_http_headers     func(req uintptr) uintptr
	webkit_uri_scheme_request_get_http_body        func(req uintptr) uintptr
	webkit_uri_scheme_response_new                 func(stream uintptr, streamLength int64) uintptr
	webkit_uri_scheme_response_set_status          func(resp uintptr, statusCode uint32, reasonPhrase string)
	webkit_uri_scheme_response_set_content_type    func(resp uintptr, contentType string)
	webkit_uri_scheme_response_set_http_headers    func(resp uintptr, headers uintptr)
	webkit_uri_scheme_request_finish_with_response func(req uintptr, resp uintptr)
	webkit_uri_scheme_request_finish_error         func(req uintptr, gErr uintptr)

	// soup
	soup_message_headers_new       func(headersType int32) uintptr
	soup_message_headers_append    func(hdrs uintptr, name string, value string)
	soup_message_headers_iter_init func(iter uintptr, hdrs uintptr)
	soup_message_headers_iter_next func(iter uintptr, name uintptr, value uintptr) int32
)

// ensureWebviewLibs lazily dlopens the GObject/Gio/WebKitGTK/libsoup libraries
// and binds every function this package needs. pkg/application has already
// loaded and validated these same libraries before any request can arrive, so
// a failure here (panic, see mainthread_linux_purego.go helpers) is a safety
// net, not primary UX.
func ensureWebviewLibs() {
	webviewLibsOnce.Do(func() {
		ensureGLib()

		libGObject := dlopenWebviewLib("libgobject-2.0.so.0", "libgobject-2.0.so")
		libGio := dlopenWebviewLib("libgio-2.0.so.0", "libgio-2.0.so")
		libWebKit := dlopenWebviewLib("libwebkitgtk-6.0.so.4", "libwebkitgtk-6.0.so")
		libSoup := dlopenWebviewLib("libsoup-3.0.so.0", "libsoup-3.0.so")

		mustRegisterWebviewFunc(&g_free, webviewLibGLib, "g_free")
		mustRegisterWebviewFunc(&g_error_free, webviewLibGLib, "g_error_free")
		mustRegisterWebviewFunc(&g_error_new_literal, webviewLibGLib, "g_error_new_literal")
		mustRegisterWebviewFunc(&g_quark_from_string, webviewLibGLib, "g_quark_from_string")

		mustRegisterWebviewFunc(&g_object_ref, libGObject, "g_object_ref")
		mustRegisterWebviewFunc(&g_object_unref, libGObject, "g_object_unref")

		mustRegisterWebviewFunc(&g_unix_input_stream_new, libGio, "g_unix_input_stream_new")
		mustRegisterWebviewFunc(&g_input_stream_read_all, libGio, "g_input_stream_read_all")
		mustRegisterWebviewFunc(&g_input_stream_close, libGio, "g_input_stream_close")

		mustRegisterWebviewFunc(&webkit_uri_scheme_request_get_uri, libWebKit, "webkit_uri_scheme_request_get_uri")
		mustRegisterWebviewFunc(&webkit_uri_scheme_request_get_http_method, libWebKit, "webkit_uri_scheme_request_get_http_method")
		mustRegisterWebviewFunc(&webkit_uri_scheme_request_get_http_headers, libWebKit, "webkit_uri_scheme_request_get_http_headers")
		mustRegisterWebviewFunc(&webkit_uri_scheme_request_get_http_body, libWebKit, "webkit_uri_scheme_request_get_http_body")
		mustRegisterWebviewFunc(&webkit_uri_scheme_response_new, libWebKit, "webkit_uri_scheme_response_new")
		mustRegisterWebviewFunc(&webkit_uri_scheme_response_set_status, libWebKit, "webkit_uri_scheme_response_set_status")
		mustRegisterWebviewFunc(&webkit_uri_scheme_response_set_content_type, libWebKit, "webkit_uri_scheme_response_set_content_type")
		mustRegisterWebviewFunc(&webkit_uri_scheme_response_set_http_headers, libWebKit, "webkit_uri_scheme_response_set_http_headers")
		mustRegisterWebviewFunc(&webkit_uri_scheme_request_finish_with_response, libWebKit, "webkit_uri_scheme_request_finish_with_response")
		mustRegisterWebviewFunc(&webkit_uri_scheme_request_finish_error, libWebKit, "webkit_uri_scheme_request_finish_error")

		mustRegisterWebviewFunc(&soup_message_headers_new, libSoup, "soup_message_headers_new")
		mustRegisterWebviewFunc(&soup_message_headers_append, libSoup, "soup_message_headers_append")
		mustRegisterWebviewFunc(&soup_message_headers_iter_init, libSoup, "soup_message_headers_iter_init")
		mustRegisterWebviewFunc(&soup_message_headers_iter_next, libSoup, "soup_message_headers_iter_next")
	})
}

// goString copies a NUL-terminated C string. The pointer is not freed — use
// this for const char* returns, which we must not free.
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

// ----------------------------------------------------------------------------
// WebKit URI scheme request helpers
// ----------------------------------------------------------------------------

func webkitURISchemeRequestGetHTTPMethod(req uintptr) string {
	// Reading request metadata touches the WebKit-owned request object, which
	// belongs to the GTK main loop; this runs on a worker goroutine, so it must
	// hop to the main thread. See mainthread_linux_purego.go and issue #5631.
	var method string
	invokeOnMainSync(func() {
		method = goString(webkit_uri_scheme_request_get_http_method(req))
	})
	return strings.ToUpper(method)
}

func webkitURISchemeRequestGetHTTPHeaders(req uintptr) http.Header {
	h := http.Header{}
	// Reading and iterating the request's libsoup headers touches WebKit-owned
	// state on the GTK main loop; this runs on a worker goroutine, so it must hop
	// to the main thread. See mainthread_linux_purego.go and issue #5631.
	invokeOnMainSync(func() {
		hdrs := webkit_uri_scheme_request_get_http_headers(req)

		var iter soupMessageHeadersIter
		soup_message_headers_iter_init(uintptr(unsafe.Pointer(&iter)), hdrs)

		var name uintptr
		var value uintptr

		for soup_message_headers_iter_next(uintptr(unsafe.Pointer(&iter)), uintptr(unsafe.Pointer(&name)), uintptr(unsafe.Pointer(&value))) != 0 {
			h.Add(goString(name), goString(value))
		}
	})
	return h
}

func webkitURISchemeRequestFinish(req uintptr, code int, header http.Header, rFD int, streamLength int64) error {
	// Completing the request touches WebKit/libsoup objects owned by the GTK
	// main loop, but this runs on an asset-server worker goroutine. WebKit2GTK
	// is not thread-safe, so the whole sequence must hop to the main thread.
	//
	// The response input stream is created and unref'd inside the same hop: it is
	// ref-taken by webkit_uri_scheme_response_new on the main thread, so creating
	// and releasing our reference here too keeps every refcount operation on a
	// single thread. Previously the stream was built and unref'd on the worker
	// while WebKit took its ref on the main thread, splitting the stream's
	// refcount across threads. See mainthread_linux_purego.go and issue #5631.
	invokeOnMainSync(func() {
		stream := g_unix_input_stream_new(int32(rFD), 1)
		defer g_object_unref(stream)

		resp := webkit_uri_scheme_response_new(stream, streamLength)
		defer g_object_unref(resp)

		webkit_uri_scheme_response_set_status(resp, uint32(code), http.StatusText(code))

		webkit_uri_scheme_response_set_content_type(resp, header.Get(HeaderContentType))

		// Ownership of hdrs is transferred to the response by
		// webkit_uri_scheme_response_set_http_headers (transfer full), so we must
		// not unref it here — doing so frees the headers while WebKit/libsoup still
		// reference them, crashing in soup_message_headers_iter_next on render.
		hdrs := soup_message_headers_new(soupMessageHeadersResponse)
		for name, values := range header {
			for _, value := range values {
				soup_message_headers_append(hdrs, name, value)
			}
		}

		webkit_uri_scheme_response_set_http_headers(resp, hdrs)

		webkit_uri_scheme_request_finish_with_response(req, resp)
	})
	return nil
}

func webkitURISchemeRequestGetHTTPBody(req uintptr) io.ReadCloser {
	// Fetching the request body stream touches the WebKit-owned request on the
	// GTK main loop; this runs on a worker goroutine, so it must hop to the main
	// thread. See mainthread_linux_purego.go and issue #5631.
	var stream uintptr
	invokeOnMainSync(func() {
		stream = webkit_uri_scheme_request_get_http_body(req)
	})
	if stream == 0 {
		return http.NoBody
	}
	return &webkitRequestBody{stream: stream}
}

type webkitRequestBody struct {
	stream uintptr
	closed bool
}

func (r *webkitRequestBody) Read(p []byte) (int, error) {
	if r.closed {
		return 0, io.ErrClosedPipe
	}

	// io.Reader allows a zero-length read; taking &p[0] on an empty slice would
	// panic, so return early before touching the backing array.
	if len(p) == 0 {
		return 0, nil
	}

	content := unsafe.Pointer(&p[0])
	contentLen := len(p)

	var n uintptr
	var gErr uintptr
	var res int32
	// Reading the WebKit-owned request body stream must happen on the GTK main
	// loop thread; this runs on a worker goroutine. See issue #5631.
	invokeOnMainSync(func() {
		res = g_input_stream_read_all(r.stream, uintptr(content), uintptr(contentLen), uintptr(unsafe.Pointer(&n)), 0, uintptr(unsafe.Pointer(&gErr)))
	})
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

	var err error
	var gErr uintptr
	// Closing and unref-ing the WebKit-owned request body stream finalizes a
	// GObject tied to the GTK main loop; this runs on a worker goroutine, so it
	// must hop to the main thread. See issue #5631.
	invokeOnMainSync(func() {
		if g_input_stream_close(r.stream, 0, uintptr(unsafe.Pointer(&gErr))) == 0 {
			err = formatGError("stream close failed", gErr)
		}
		g_object_unref(r.stream)
	})
	r.stream = 0
	return err
}

func formatGError(msg string, gErr uintptr, args ...any) error {
	if gErr != 0 {
		// GError layout: guint32 domain, gint32 code, char *message.
		e := (*gError)(unsafe.Pointer(gErr))
		if e.message != 0 {
			msg += ": " + goString(e.message)
			g_error_free(gErr)
		}
	}
	return fmt.Errorf(msg, args...)
}
