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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"unsafe"
)

//export processURLRequest
func processURLRequest(ctx unsafe.Pointer, requestId C.ulonglong, url *C.char, method *C.char, headers *C.char, body unsafe.Pointer, bodyLen C.int, hasBodyStream C.int) {
	var bodyReader io.Reader
	if body != nil && bodyLen != 0 {
		bodyReader = bytes.NewReader(C.GoBytes(body, bodyLen))
	} else if hasBodyStream != 0 {
		bodyReader = &bodyStreamReader{id: requestId, ctx: ctx}
	}

	requestBuffer <- &wkWebViewRequest{
		id:      requestId,
		url:     C.GoString(url),
		method:  C.GoString(method),
		headers: C.GoString(headers),
		body:    bodyReader,
		ctx:     ctx,
	}
}

type wkWebViewRequest struct {
	id      C.ulonglong
	url     string
	method  string
	headers string
	body    io.Reader

	ctx unsafe.Pointer
}

func (r *wkWebViewRequest) GetHttpRequest() (*http.Request, error) {
	req, err := http.NewRequest(r.method, r.url, r.body)
	if err != nil {
		return nil, err
	}

	if r.headers != "" {
		var h map[string]string
		if err := json.Unmarshal([]byte(r.headers), &h); err != nil {
			return nil, fmt.Errorf("unable to unmarshal request headers: %s", err)
		}

		for k, v := range h {
			req.Header.Add(k, v)
		}
	}

	return req, nil
}

var _ io.Reader = &bodyStreamReader{}

type bodyStreamReader struct {
	id  C.ulonglong
	ctx unsafe.Pointer
}

// Read implements io.Reader
func (r *bodyStreamReader) Read(p []byte) (n int, err error) {
	var content unsafe.Pointer
	var contentLen int
	if p != nil {
		content = unsafe.Pointer(&p[0])
		contentLen = len(p)
	}

	res := C.ProcessURLRequestReadBodyStream(r.ctx, r.id, content, C.int(contentLen))
	if res > 0 {
		return int(res), nil
	}

	switch res {
	case 0:
		return 0, io.EOF
	case -1:
		return 0, fmt.Errorf("body: stream error")
	case -2:
		return 0, errRequestStopped
	case -3:
		return 0, fmt.Errorf("body: no stream defined")
	case -4:
		return 0, io.ErrClosedPipe
	default:
		return 0, fmt.Errorf("body: unknown error %d", res)
	}
}
