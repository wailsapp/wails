package assetserver

import (
	"net/http"
)

type contentTypeSniffer struct {
	rw http.ResponseWriter

	wroteHeader bool
}

func (rw *contentTypeSniffer) Header() http.Header {
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
