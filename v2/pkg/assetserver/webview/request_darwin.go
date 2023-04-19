//go:build darwin

package webview

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework WebKit

#import <Foundation/Foundation.h>
#import <WebKit/WebKit.h>

static void URLSchemeTaskRetain(void *wkUrlSchemeTask) {
	id<WKURLSchemeTask> urlSchemeTask = (id<WKURLSchemeTask>) wkUrlSchemeTask;
	[urlSchemeTask retain];
}

static void URLSchemeTaskRelease(void *wkUrlSchemeTask) {
	id<WKURLSchemeTask> urlSchemeTask = (id<WKURLSchemeTask>) wkUrlSchemeTask;
	[urlSchemeTask release];
}

static const char * URLSchemeTaskRequestURL(void *wkUrlSchemeTask) {
	id<WKURLSchemeTask> urlSchemeTask = (id<WKURLSchemeTask>) wkUrlSchemeTask;
	@autoreleasepool {
		return [urlSchemeTask.request.URL.absoluteString UTF8String];
	}
}

static const char * URLSchemeTaskRequestMethod(void *wkUrlSchemeTask) {
	id<WKURLSchemeTask> urlSchemeTask = (id<WKURLSchemeTask>) wkUrlSchemeTask;
	@autoreleasepool {
		return [urlSchemeTask.request.HTTPMethod UTF8String];
	}
}

static const char * URLSchemeTaskRequestHeadersJSON(void *wkUrlSchemeTask) {
	id<WKURLSchemeTask> urlSchemeTask = (id<WKURLSchemeTask>) wkUrlSchemeTask;
	@autoreleasepool {
		NSData *headerData = [NSJSONSerialization dataWithJSONObject: urlSchemeTask.request.allHTTPHeaderFields options:0 error: nil];
		if (!headerData) {
			return nil;
		}

		NSString* headerString = [[[NSString alloc] initWithData:headerData encoding:NSUTF8StringEncoding] autorelease];
		const char * headerJSON = [headerString UTF8String];

		char * headersOut = malloc(strlen(headerJSON));
		strcpy(headersOut, headerJSON);
		return headersOut;
	}
}

static bool URLSchemeTaskRequestBodyBytes(void *wkUrlSchemeTask, const void **body, int *bodyLen) {
	id<WKURLSchemeTask> urlSchemeTask = (id<WKURLSchemeTask>) wkUrlSchemeTask;
	@autoreleasepool {
		if (!urlSchemeTask.request.HTTPBody) {
			return false;
		}

		*body = urlSchemeTask.request.HTTPBody.bytes;
		*bodyLen = urlSchemeTask.request.HTTPBody.length;
		return true;
	}
}

static bool URLSchemeTaskRequestBodyStreamOpen(void *wkUrlSchemeTask) {
	id<WKURLSchemeTask> urlSchemeTask = (id<WKURLSchemeTask>) wkUrlSchemeTask;
	@autoreleasepool {
		if (!urlSchemeTask.request.HTTPBodyStream) {
			return false;
		}

		[urlSchemeTask.request.HTTPBodyStream open];
		return true;
	}
}

static void URLSchemeTaskRequestBodyStreamClose(void *wkUrlSchemeTask) {
	id<WKURLSchemeTask> urlSchemeTask = (id<WKURLSchemeTask>) wkUrlSchemeTask;
	@autoreleasepool {
		if (!urlSchemeTask.request.HTTPBodyStream) {
			return;
		}

		[urlSchemeTask.request.HTTPBodyStream close];
	}
}

static int URLSchemeTaskRequestBodyStreamRead(void *wkUrlSchemeTask, void *buf, int bufLen) {
	id<WKURLSchemeTask> urlSchemeTask = (id<WKURLSchemeTask>) wkUrlSchemeTask;

	@autoreleasepool {
		NSInputStream *stream = urlSchemeTask.request.HTTPBodyStream;
		if (!stream) {
			return -2;
		}

		NSStreamStatus status = stream.streamStatus;
		if (status == NSStreamStatusAtEnd || !stream.hasBytesAvailable) {
			return 0;
		} else if (status != NSStreamStatusOpen) {
			return -3;
		}

		return [stream read:buf maxLength:bufLen];
	}
}
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

// NewRequest creates as new WebViewRequest based on a pointer to an `id<WKURLSchemeTask>`
func NewRequest(wkURLSchemeTask unsafe.Pointer) Request {
	C.URLSchemeTaskRetain(wkURLSchemeTask)
	return newRequestFinalizer(&request{task: wkURLSchemeTask})
}

var _ Request = &request{}

type request struct {
	task unsafe.Pointer

	header http.Header
	body   io.ReadCloser
	rw     *responseWriter
}

func (r *request) URL() (string, error) {
	return C.GoString(C.URLSchemeTaskRequestURL(r.task)), nil
}

func (r *request) Method() (string, error) {
	return C.GoString(C.URLSchemeTaskRequestMethod(r.task)), nil
}

func (r *request) Header() (http.Header, error) {
	if r.header != nil {
		return r.header, nil
	}

	header := http.Header{}
	if cHeaders := C.URLSchemeTaskRequestHeadersJSON(r.task); cHeaders != nil {
		if headers := C.GoString(cHeaders); headers != "" {
			var h map[string]string
			if err := json.Unmarshal([]byte(headers), &h); err != nil {
				return nil, fmt.Errorf("unable to unmarshal request headers: %s", err)
			}

			for k, v := range h {
				header.Add(k, v)
			}
		}
		C.free(unsafe.Pointer(cHeaders))
	}
	r.header = header
	return header, nil
}

func (r *request) Body() (io.ReadCloser, error) {
	if r.body != nil {
		return r.body, nil
	}

	var body unsafe.Pointer
	var bodyLen C.int
	if C.URLSchemeTaskRequestBodyBytes(r.task, &body, &bodyLen) {
		if body != nil && bodyLen > 0 {
			r.body = io.NopCloser(bytes.NewReader(C.GoBytes(body, bodyLen)))
		} else {
			r.body = http.NoBody
		}
	} else if C.URLSchemeTaskRequestBodyStreamOpen(r.task) {
		r.body = &requestBodyStreamReader{task: r.task}
	}

	return r.body, nil
}

func (r *request) Response() ResponseWriter {
	if r.rw != nil {
		return r.rw
	}

	r.rw = &responseWriter{r: r}
	return r.rw
}

func (r *request) Close() error {
	var err error
	if r.body != nil {
		err = r.body.Close()
	}
	r.Response().Finish()
	C.URLSchemeTaskRelease(r.task)
	return err
}

var _ io.ReadCloser = &requestBodyStreamReader{}

type requestBodyStreamReader struct {
	task   unsafe.Pointer
	closed bool
}

// Read implements io.Reader
func (r *requestBodyStreamReader) Read(p []byte) (n int, err error) {
	var content unsafe.Pointer
	var contentLen int
	if p != nil {
		content = unsafe.Pointer(&p[0])
		contentLen = len(p)
	}

	res := C.URLSchemeTaskRequestBodyStreamRead(r.task, content, C.int(contentLen))
	if res > 0 {
		return int(res), nil
	}

	switch res {
	case 0:
		return 0, io.EOF
	case -1:
		return 0, fmt.Errorf("body: stream error")
	case -2:
		return 0, fmt.Errorf("body: no stream defined")
	case -3:
		return 0, io.ErrClosedPipe
	default:
		return 0, fmt.Errorf("body: unknown error %d", res)
	}
}

func (r *requestBodyStreamReader) Close() error {
	if r.closed {
		return nil
	}
	r.closed = true

	C.URLSchemeTaskRequestBodyStreamClose(r.task)
	return nil
}
