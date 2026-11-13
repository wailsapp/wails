package assetserver

import (
	"net/http"
	"path/filepath"
)

type contentTypeSniffer struct {
	rw      http.ResponseWriter
	reqPath string

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
		ct := http.DetectContentType(b)
		if rw.reqPath != "" {
			if ext := filepath.Ext(rw.reqPath); ext != "" {
				if custom := getMimeTypeByExt(ext); custom != "" {
					ct = custom
				}
			}
		}
		m.Set(HeaderContentType, ct)
	}

	rw.WriteHeader(http.StatusOK)
}

func getMimeTypeByExt(ext string) string {
	if ct, ok := mimeTypesByExt[ext]; ok {
		return ct
	}
	return ""
}
