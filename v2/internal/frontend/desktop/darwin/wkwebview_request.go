//go:build darwin

package darwin

/*
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
func processURLRequest(ctx unsafe.Pointer, requestId C.ulonglong, url *C.char, method *C.char, headers *C.char, body unsafe.Pointer, bodyLen C.int) {
	var goBody []byte
	if body != nil && bodyLen != 0 {
		goBody = C.GoBytes(body, bodyLen)
	}

	requestBuffer <- &wkWebViewRequest{
		id:      requestId,
		url:     C.GoString(url),
		method:  C.GoString(method),
		headers: C.GoString(headers),
		body:    goBody,
		ctx:     ctx,
	}
}

type wkWebViewRequest struct {
	id      C.ulonglong
	url     string
	method  string
	headers string
	body    []byte

	ctx unsafe.Pointer
}

func (r *wkWebViewRequest) GetHttpRequest() (*http.Request, error) {
	var body io.Reader
	if len(r.body) != 0 {
		body = bytes.NewReader(r.body)
	}

	req, err := http.NewRequest(r.method, r.url, body)
	if err != nil {
		return nil, err
	}

	if r.headers != "" {
		var h map[string]string
		if err := json.Unmarshal([]byte(r.headers), &h); err != nil {
			return nil, fmt.Errorf("Unable to unmarshal request headers: %s", err)
		}

		for k, v := range h {
			req.Header.Add(k, v)
		}
	}

	return req, nil
}
