//go:build darwin && !ios

package webview

import (
	"net/http"
	"unsafe"
)

func (rw *responseWriter) Code() int {
	return rw.code
}

func (rw *responseWriter) Header() http.Header {
	if rw.header == nil {
		rw.header = http.Header{}
	}
	return rw.header
}

func (rw *responseWriter) Finish() error {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusNotImplemented)
	}

	if rw.finished {
		return nil
	}
	rw.finished = true

	urlSchemeTaskDidFinish(rw.r.task)

	return nil
}

func (rw *responseWriter) Write(buf []byte) (int, error) {
	if rw.finished {
		return 0, errResponseFinished
	}

	rw.WriteHeader(http.StatusOK)

	var content unsafe.Pointer
	var contentLen int
	if len(buf) > 0 {
		content = unsafe.Pointer(&buf[0])
		contentLen = len(buf)
	}

	if !urlSchemeTaskDidReceiveData(rw.r.task, content, contentLen) {
		return 0, errRequestStopped
	}
	return contentLen, nil
}
