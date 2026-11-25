//go:build linux && purego && !android

package webview

// CGO-free (purego) port of responsewriter_linux.go.

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"syscall"
)

var (
	webviewAssetErrorQuarkOnce sync.Once
	webviewAssetErrorQuarkID   uint32
)

// webviewAssetErrorQuark returns a stable GError domain for asset-server
// failures. The domain string is interned exactly once (sync.Once), the purego
// equivalent of the cgo g_quark_from_static_string over a static string
// literal, so it never leaks. Interning the per-request error message instead
// (the previous behaviour) grew the global quark table unboundedly on
// long-running apps, since GQuarks are never freed.
func webviewAssetErrorQuark() uint32 {
	webviewAssetErrorQuarkOnce.Do(func() {
		webviewAssetErrorQuarkID = g_quark_from_string("wails-webview-assetserver")
	})
	return webviewAssetErrorQuarkID
}

type responseWriter struct {
	req uintptr

	header      http.Header
	wroteHeader bool
	finished    bool
	code        int

	w    io.WriteCloser
	wErr error
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

	rFD, w, err := pipe()
	if err != nil {
		rw.finishWithError(http.StatusInternalServerError, fmt.Errorf("unable to open pipe: %s", err))
		return
	}
	rw.w = w

	// webkitURISchemeRequestFinish wraps the read end of the pipe in a
	// GUnixInputStream and completes the request on the GTK main thread; the
	// stream is created and released there too. See webkit_linux_purego.go and
	// #5631.
	if err := webkitURISchemeRequestFinish(rw.req, code, rw.Header(), rFD, contentLength); err != nil {
		rw.finishWithError(http.StatusInternalServerError, fmt.Errorf("unable to finish request: %s", err))
		return
	}
}

func (rw *responseWriter) finishWithError(code int, err error) {
	if rw.w != nil {
		rw.w.Close()
		rw.w = &nopCloser{io.Discard}
	}
	rw.wErr = err

	msg := err.Error()

	// webkit_uri_scheme_request_finish_error touches the WebKit-owned request
	// on the GTK main loop; this runs on an asset-server worker goroutine, so
	// it must hop to the main thread. See mainthread_linux_purego.go and issue
	// #5631.
	invokeOnMainSync(func() {
		gerr := g_error_new_literal(webviewAssetErrorQuark(), int32(code), msg)
		webkit_uri_scheme_request_finish_error(rw.req, gerr)
		g_error_free(gerr)
	})
}

type nopCloser struct {
	io.Writer
}

func pipe() (r int, w *os.File, err error) {
	var p [2]int
	// O_CLOEXEC, as in the cgo backends: a child process started while a
	// response streams must not inherit the write end, or the reader never
	// sees EOF and the fetch stalls.
	e := syscall.Pipe2(p[0:], syscall.O_CLOEXEC)
	if e != nil {
		return 0, nil, fmt.Errorf("pipe2: %s", e)
	}

	return p[0], os.NewFile(uintptr(p[1]), "|1"), nil
}
