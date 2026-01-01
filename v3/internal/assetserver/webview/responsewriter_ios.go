//go:build ios

package webview

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Foundation -framework WebKit

#import <Foundation/Foundation.h>
#import <WebKit/WebKit.h>

typedef void (^schemeTaskCaller)(id<WKURLSchemeTask>);

static bool urlSchemeTaskCall(void *wkUrlSchemeTask, schemeTaskCaller fn) {
    id<WKURLSchemeTask> urlSchemeTask = (__bridge id<WKURLSchemeTask>) wkUrlSchemeTask;
    if (urlSchemeTask == nil) {
        return false;
    }

	@autoreleasepool {
		@try {
			fn(urlSchemeTask);
		} @catch (NSException *exception) {
			// This is very bad to detect a stopped schemeTask this should be implemented in a better way
			// But it seems to be very tricky to not deadlock when keeping a lock curing executing fn()
			// It seems like those call switch the thread back to the main thread and then deadlocks when they reentrant want
			// to get the lock again to start another request or stop it.
			if ([exception.reason isEqualToString: @"This task has already been stopped"]) {
				return false;
			}

			@throw exception;
		}

		return true;
	}
}

static bool URLSchemeTaskDidReceiveData(void *wkUrlSchemeTask, void* data, int datalength) {
	return urlSchemeTaskCall(
		wkUrlSchemeTask,
		^(id<WKURLSchemeTask> urlSchemeTask) {
			NSData *nsdata = [NSData dataWithBytes:data length:datalength];
			[urlSchemeTask didReceiveData:nsdata];
	});
}

static bool URLSchemeTaskDidFinish(void *wkUrlSchemeTask) {
	return urlSchemeTaskCall(
		wkUrlSchemeTask,
		^(id<WKURLSchemeTask> urlSchemeTask) {
			[urlSchemeTask didFinish];
	});
}

static bool URLSchemeTaskDidReceiveResponse(void *wkUrlSchemeTask, int statusCode, void *headersString, int headersStringLength) {
	return urlSchemeTaskCall(
		wkUrlSchemeTask,
		^(id<WKURLSchemeTask> urlSchemeTask) {
			NSData *nsHeadersJSON = [NSData dataWithBytes:headersString length:headersStringLength];
            NSDictionary *headerFields = [NSJSONSerialization JSONObjectWithData:nsHeadersJSON options:NSJSONReadingMutableContainers error:nil];
            NSHTTPURLResponse *response = [[NSHTTPURLResponse alloc] initWithURL:urlSchemeTask.request.URL statusCode:statusCode HTTPVersion:nil headerFields:headerFields];

			[urlSchemeTask didReceiveResponse:response];
	});
}
*/
import "C"

import (
	"fmt"
	"net/http"
	"strings"
	"unsafe"

	json "github.com/goccy/go-json"
)

var _ ResponseWriter = &responseWriter{}

type responseWriter struct {
	r *request

	header      http.Header
	wroteHeader bool
	code        int

	finished bool
}

func (rw *responseWriter) Header() http.Header {
	if rw.header == nil {
		rw.header = http.Header{}
	}
	return rw.header
}

func (rw *responseWriter) Write(buf []byte) (int, error) {
	if rw.finished {
		return 0, errResponseFinished
	}

	rw.WriteHeader(http.StatusOK)

	// Debug logging for CSS files
	if url, err := rw.r.URL(); err == nil && (strings.Contains(url, ".css") || strings.Contains(url, "style")) {
		preview := string(buf)
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		fmt.Printf("🎨 CSS Write: URL=%s Size=%d Preview=%s\n", url, len(buf), preview)
	}

	var content unsafe.Pointer
	var contentLen int
	if buf != nil {
		content = unsafe.Pointer(&buf[0])
		contentLen = len(buf)
	}

	if !C.URLSchemeTaskDidReceiveData(rw.r.task, content, C.int(contentLen)) {
		return 0, errRequestStopped
	}
	return contentLen, nil
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.code = code
	if rw.wroteHeader || rw.finished {
		return
	}
	rw.wroteHeader = true

	header := map[string]string{}
	for k := range rw.Header() {
		header[k] = rw.Header().Get(k)
	}
	headerData, _ := json.Marshal(header)

	// Debug logging for CSS files
	if url, err := rw.r.URL(); err == nil && (strings.Contains(url, ".css") || strings.Contains(url, "style")) {
		fmt.Printf("🎨 CSS Response: URL=%s Code=%d Headers=%s\n", url, code, string(headerData))
	}

	var headers unsafe.Pointer
	var headersLen int
	if len(headerData) != 0 {
		headers = unsafe.Pointer(&headerData[0])
		headersLen = len(headerData)
	}

	C.URLSchemeTaskDidReceiveResponse(rw.r.task, C.int(code), headers, C.int(headersLen))
}

func (rw *responseWriter) Finish() error {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusNotImplemented)
	}

	if rw.finished {
		return nil
	}
	rw.finished = true

	C.URLSchemeTaskDidFinish(rw.r.task)

	return nil
}

func (rw *responseWriter) Code() int {
	return rw.code
}
