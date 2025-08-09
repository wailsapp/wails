package assetserver

import (
	"bytes"
	"net/http"
)

type bodyRecorder struct {
	http.ResponseWriter
	doRecord func(code int, header http.Header) bool

	body        *bytes.Buffer
	code        int
	wroteHeader bool
}

func (rw *bodyRecorder) Write(buf []byte) (int, error) {
	rw.writeHeader(buf, http.StatusOK)
	if rw.body != nil {
		return rw.body.Write(buf)
	}
	return rw.ResponseWriter.Write(buf)
}

func (rw *bodyRecorder) WriteHeader(code int) {
	rw.writeHeader(nil, code)
}

func (rw *bodyRecorder) Code() int {
	return rw.code
}

func (rw *bodyRecorder) Body() *bytes.Buffer {
	return rw.body
}

func (rw *bodyRecorder) writeHeader(buf []byte, code int) {
	if rw.wroteHeader {
		return
	}

	if rw.doRecord != nil {
		header := rw.Header()
		if len(buf) != 0 {
			if _, hasType := header[HeaderContentType]; !hasType {
				header.Set(HeaderContentType, http.DetectContentType(buf))
			}
		}

		if rw.doRecord(code, header) {
			rw.body = bytes.NewBuffer(nil)
		}
	}

	if rw.body == nil {
		rw.ResponseWriter.WriteHeader(code)
	}

	rw.code = code
	rw.wroteHeader = true
}
