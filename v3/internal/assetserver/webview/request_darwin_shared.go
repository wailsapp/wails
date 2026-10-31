//go:build darwin && !ios

package webview

import (
	"errors"
	"fmt"
	"io"
	"unsafe"
)

func (r *request) Response() ResponseWriter {
	if r.rw != nil {
		return r.rw
	}

	r.rw = &responseWriter{r: r}
	return r.rw
}

// Read implements io.Reader.
func (r *requestBodyStreamReader) Read(p []byte) (n int, err error) {
	var content unsafe.Pointer
	var contentLen int
	if len(p) > 0 {
		content = unsafe.Pointer(&p[0])
		contentLen = len(p)
	}

	res := r.readStream(content, contentLen)
	if res > 0 {
		return res, nil
	}

	switch res {
	case 0:
		return 0, io.EOF
	case -1:
		return 0, errors.New("body: stream error")
	case -2:
		return 0, errors.New("body: no stream defined")
	case -3:
		return 0, io.ErrClosedPipe
	default:
		return 0, fmt.Errorf("body: unknown error %d", res)
	}
}
