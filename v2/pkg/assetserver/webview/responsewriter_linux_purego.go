//go:build linux && purego
// +build linux,purego

package webview

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"syscall"

	"github.com/ebitengine/purego"
)

const (
	gtk3 = "libgtk-3.so"
	gtk4 = "libgtk-4.so"
)

var (
	gtk     uintptr
	webkit  uintptr
	version int
)

func init() {
	var err error
	// gtk, err = purego.Dlopen(gtk4, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	// if err == nil {
	// 	version = 4
	// 	return
	// }
	//	log.Println("Failed to open GTK4: Falling back to GTK3")
	gtk, err = purego.Dlopen(gtk3, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}
	version = 3

	var webkit4 string = "libwebkit2gtk-4.1.so"
	webkit, err = purego.Dlopen(webkit4, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}
}

type responseWriter struct {
	req uintptr

	header      http.Header
	wroteHeader bool
	finished    bool

	w    io.WriteCloser
	wErr error
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
	// TODO? Is this ever called? I don't think so!
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

	var newStream func(int, bool) uintptr
	purego.RegisterLibFunc(&newStream, gtk, "g_unix_input_stream_new")
	var unRef func(uintptr)
	purego.RegisterLibFunc(&unRef, gtk, "g_object_unref")
	stream := newStream(rFD, true)

	/*	var reqFinish func(uintptr, uintptr, uintptr, uintptr, int64) int
		purego.RegisterLibFunc(&reqFinish, webkit, "webkit_uri_scheme_request_finish")

		header := rw.Header()
			defer unRef(stream)
		if err := reqFinish(rw.req, code, header, stream, contentLength); err != nil {
			rw.finishWithError(http.StatusInternalServerError, fmt.Errorf("unable to finish request: %s", err))
		}
	*/
	if err := webkit_uri_scheme_request_finish(rw.req, code, rw.Header(), stream, contentLength); err != nil {
		rw.finishWithError(http.StatusInternalServerError, fmt.Errorf("unable to finish request: %s", err))
		return
	}
}

func (rw *responseWriter) Finish() {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusNotImplemented)
	}

	if rw.finished {
		return
	}
	rw.finished = true
	if rw.w != nil {
		rw.w.Close()
	}
}

func (rw *responseWriter) finishWithError(code int, err error) {
	if rw.w != nil {
		rw.w.Close()
		rw.w = &nopCloser{io.Discard}
	}
	rw.wErr = err

	var newLiteral func(uint32, string, int, string) uintptr // is this correct?
	purego.RegisterLibFunc(&newLiteral, gtk, "g_error_new_literal")
	var newQuark func(string) uintptr
	purego.RegisterLibFunc(&newQuark, gtk, "g_quark_from_string")
	var freeError func(uintptr)
	purego.RegisterLibFunc(&freeError, gtk, "g_error_free")
	var finishError func(uintptr, uintptr)
	purego.RegisterLibFunc(&finishError, webkit, "webkit_uri_scheme_request_finish_error")

	msg := string(err.Error())
	//gquark := newQuark(msg)
	gerr := newLiteral(1, msg, code, msg)
	finishError(rw.req, gerr)
	freeError(gerr)
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
