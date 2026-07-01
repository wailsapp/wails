//go:build linux && cgo && gtk3 && !android

package webview

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.1

#include "gtk/gtk.h"
#include "webkit2/webkit2.h"

// webview_asset_error_quark returns a stable GError domain for asset-server
// failures. The string literal is static storage, so g_quark_from_static_string
// interns it once and never leaks. Interning the per-request error message
// instead (the previous behaviour) grew the global quark table unboundedly on
// long-running apps, since GQuarks are never freed.
static GQuark webview_asset_error_quark(void) {
	return g_quark_from_static_string("wails-webview-assetserver");
}

*/
import "C"
import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

type responseWriter struct {
	req *C.WebKitURISchemeRequest

	header      http.Header
	wroteHeader bool
	finished    bool
	code        int

	w    io.WriteCloser
	wErr error
}

func (rw *responseWriter) Code() int {
	return rw.code
}

func (rw *responseWriter) Header() http.Header {
	if rw.header == nil {
		rw.header = http.Header{}
	}
	return rw.header
}

func (rw *responseWriter) Write(buf []byte) (int, error) {
	if rw.finished {
		return 0, errResponseFinished
	}

	rw.WriteHeader(http.StatusOK)
	if rw.wErr != nil {
		return 0, rw.wErr
	}
	return rw.w.Write(buf)
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.code = code
	if rw.wroteHeader || rw.finished {
		return
	}
	rw.wroteHeader = true

	contentLength := int64(-1)
	if sLen := rw.Header().Get(HeaderContentLength); sLen != "" {
		if pLen, _ := strconv.ParseInt(sLen, 10, 64); pLen > 0 {
			contentLength = pLen
		}
	}

	// We can't use os.Pipe here, because that returns files with a finalizer for closing the FD. But the control over the
	// read FD is given to the InputStream and will be closed there.
	// Furthermore we especially don't want to have the FD_CLOEXEC
	rFD, w, err := pipe()
	if err != nil {
		rw.finishWithError(http.StatusInternalServerError, fmt.Errorf("unable to open pipe: %s", err))
		return
	}
	rw.w = w

	// webkit_uri_scheme_request_finish wraps the read end of the pipe in a
	// GUnixInputStream and completes the request on the GTK main thread; the
	// stream is created and released there too. See webkit_linux_gtk3.go and #5631.
	if err := webkit_uri_scheme_request_finish(rw.req, code, rw.Header(), rFD, contentLength); err != nil {
		rw.finishWithError(http.StatusInternalServerError, fmt.Errorf("unable to finish request: %s", err))
		return
	}
}

func (rw *responseWriter) Finish() error {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusNotImplemented)
	}

	if rw.finished {
		return nil
	}
	rw.finished = true
	if rw.w != nil {
		rw.w.Close()
	}
	return nil
}

func (rw *responseWriter) finishWithError(code int, err error) {
	if rw.w != nil {
		rw.w.Close()
		rw.w = &nopCloser{io.Discard}
	}
	rw.wErr = err

	msg := C.CString(err.Error())
	defer C.free(unsafe.Pointer(msg))

	// webkit_uri_scheme_request_finish_error touches the WebKit-owned request
	// on the GTK main loop; this runs on an asset-server worker goroutine, so
	// it must hop to the main thread. See mainthread_linux.go and issue #5631.
	invokeOnMainSync(func() {
		gerr := C.g_error_new_literal(C.webview_asset_error_quark(), C.int(code), msg)
		C.webkit_uri_scheme_request_finish_error(rw.req, gerr)
		C.g_error_free(gerr)
	})
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func pipe() (r int, w *os.File, err error) {
	var p [2]int
	e := syscall.Pipe2(p[0:], 0)
	if e != nil {
		return 0, nil, fmt.Errorf("pipe2: %s", e)
	}

	return p[0], os.NewFile(uintptr(p[1]), "|1"), nil
}
