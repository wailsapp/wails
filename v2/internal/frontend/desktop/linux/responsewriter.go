//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0 gio-unix-2.0

#include "gtk/gtk.h"
#include "webkit2/webkit2.h"
#include "gio/gunixinputstream.h"

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

	"github.com/wailsapp/wails/v2/internal/frontend/assetserver"
)

type webKitResponseWriter struct {
	req *C.WebKitURISchemeRequest

	header      http.Header
	wroteHeader bool

	w    io.WriteCloser
	wErr error
}

func (rw *webKitResponseWriter) Header() http.Header {
	if rw.header == nil {
		rw.header = http.Header{}
	}
	return rw.header
}

func (rw *webKitResponseWriter) Write(buf []byte) (int, error) {
	rw.WriteHeader(http.StatusOK)
	if rw.wErr != nil {
		return 0, rw.wErr
	}
	return rw.w.Write(buf)
}

func (rw *webKitResponseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.wroteHeader = true

	if code != http.StatusOK {
		// WebKitGTK stable < 2.36 API does not support response headers and response statuscodes
		rw.w = &nopCloser{io.Discard}
		rw.finishWithError(http.StatusText(code), code)
		return
	}

	// We can't use os.Pipe here, because that returns files with a finalizer for closing the FD. But the control over the
	// read FD is given to the InputStream and will be closed there.
	// Furthermore we especially don't want to have the FD_CLOEXEC
	rFD, w, err := pipe()
	if err != nil {
		rw.wErr = fmt.Errorf("Unable opening pipe: %s", err)
		rw.finishWithError(rw.wErr.Error(), http.StatusInternalServerError)
		return
	}
	rw.w = w

	cMimeType := C.CString(rw.Header().Get(assetserver.HeaderContentType))
	defer C.free(unsafe.Pointer(cMimeType))

	contentLength := int64(-1)
	if sLen := rw.Header().Get(assetserver.HeaderContentLength); sLen != "" {
		if pLen, _ := strconv.ParseInt(sLen, 10, 64); pLen > 0 {
			contentLength = pLen
		}
	}

	stream := C.g_unix_input_stream_new(C.int(rFD), gtkBool(true))
	C.webkit_uri_scheme_request_finish(rw.req, stream, C.gint64(contentLength), cMimeType)
	C.g_object_unref(C.gpointer(stream))
}

func (rw *webKitResponseWriter) Close() {
	if rw.w != nil {
		rw.w.Close()
	}
}

func (rw *webKitResponseWriter) finishWithError(message string, code int) {
	msg := C.CString(http.StatusText(code))
	gerr := C.g_error_new_literal(C.g_quark_from_string(msg), C.int(code), msg)
	C.webkit_uri_scheme_request_finish_error(rw.req, gerr)
	C.g_error_free(gerr)
	C.free(unsafe.Pointer(msg))
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
