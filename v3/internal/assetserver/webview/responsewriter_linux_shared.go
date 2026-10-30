//go:build linux && !android

package webview

import "net/http"

func (nopCloser) Close() error { return nil }

func (rw *responseWriter) Code() int {
	return rw.code
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

// Flush implements the http.Flusher interface.
//
// Write goes into the pipe whose read end WebKitGTK consumes as a
// GUnixInputStream, so bytes are already on their way out and there is
// nothing held back here.
func (rw *responseWriter) Flush() {}

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
