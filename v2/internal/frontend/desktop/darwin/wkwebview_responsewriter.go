//go:build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit

#import "Application.h"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"errors"
	"net/http"
	"unsafe"
)

var (
	errRequestStopped = errors.New("request has been stopped")
)

type wkWebViewResponseWriter struct {
	r *wkWebViewRequest

	header      http.Header
	wroteHeader bool
}

func (rw *wkWebViewResponseWriter) Header() http.Header {
	if rw.header == nil {
		rw.header = http.Header{}
	}
	return rw.header
}

func (rw *wkWebViewResponseWriter) Write(buf []byte) (int, error) {
	rw.WriteHeader(http.StatusOK)

	var content unsafe.Pointer
	var contentLen int
	if buf != nil {
		content = unsafe.Pointer(&buf[0])
		contentLen = len(buf)
	}

	if !C.ProcessURLDidReceiveData(rw.r.ctx, rw.r.id, content, C.int(contentLen)) {
		return 0, errRequestStopped
	}
	return contentLen, nil
}

func (rw *wkWebViewResponseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.wroteHeader = true

	header := map[string]string{}
	for k := range rw.Header() {
		header[k] = rw.Header().Get(k)
	}
	headerData, _ := json.Marshal(header)

	var headers unsafe.Pointer
	var headersLen int
	if len(headerData) != 0 {
		headers = unsafe.Pointer(&headerData[0])
		headersLen = len(headerData)
	}

	C.ProcessURLDidReceiveResponse(rw.r.ctx, rw.r.id, C.int(code), headers, C.int(headersLen))
}

func (rw *wkWebViewResponseWriter) Close() {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusNotImplemented)
	}
	C.ProcessURLDidFinish(rw.r.ctx, rw.r.id)
}
