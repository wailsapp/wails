package assetserver

import (
	"net/http"
)

func newContentTypeSniffer(rw http.ResponseWriter) *contentTypeSniffer {
	return &contentTypeSniffer{
		rw:           rw,
		closeChannel: make(chan bool, 1),
	}
}

type contentTypeSniffer struct {
	rw           http.ResponseWriter
	status       int
	wroteHeader  bool
	closeChannel chan bool
}

func (rw contentTypeSniffer) Header() http.Header {
	return rw.rw.Header()
}

func (rw *contentTypeSniffer) Write(buf []byte) (int, error) {
	rw.writeHeader(buf)
	return rw.rw.Write(buf)
}

func (rw *contentTypeSniffer) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.rw.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *contentTypeSniffer) writeHeader(b []byte) {
	if rw.wroteHeader {
		return
	}

	m := rw.rw.Header()
	if _, hasType := m[HeaderContentType]; !hasType {
		m.Set(HeaderContentType, http.DetectContentType(b))
	}

	rw.WriteHeader(http.StatusOK)
}

// CloseNotify implements the http.CloseNotifier interface.
func (rw *contentTypeSniffer) CloseNotify() <-chan bool {
	return rw.closeChannel
}

func (rw *contentTypeSniffer) closeClient() {
	rw.closeChannel <- true
}

// Flush implements the http.Flusher interface.
func (rw *contentTypeSniffer) Flush() {
	if f, ok := rw.rw.(http.Flusher); ok {
		f.Flush()
	}
}
