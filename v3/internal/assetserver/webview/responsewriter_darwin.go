//go:build darwin && !ios && !purego

package webview

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework WebKit

#import <Foundation/Foundation.h>
#import <WebKit/WebKit.h>

typedef void (^schemeTaskCaller)(id<WKURLSchemeTask>);

static bool urlSchemeTaskCall(void *wkUrlSchemeTask, schemeTaskCaller fn) {
	id<WKURLSchemeTask> urlSchemeTask = (id<WKURLSchemeTask>) wkUrlSchemeTask;
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
			NSDictionary *headerFields = [NSJSONSerialization JSONObjectWithData:nsHeadersJSON options: NSJSONReadingMutableContainers error: nil];
        	NSHTTPURLResponse *response = [[[NSHTTPURLResponse alloc] initWithURL:urlSchemeTask.request.URL statusCode:statusCode HTTPVersion:@"HTTP/1.1" headerFields:headerFields] autorelease];

			[urlSchemeTask didReceiveResponse:response];
    	});
}
*/
import "C"

import (
	"net/http"
	"unsafe"

	"encoding/json"
)

var _ ResponseWriter = &responseWriter{}

type responseWriter struct {
	r *request

	header      http.Header
	wroteHeader bool
	code        int

	finished bool
}

// Flush implements the http.Flusher interface.
//
// Write hands each chunk straight to the URL scheme task via
// didReceiveData, so there is nothing buffered on this side and a flush has
// no work to do. Declaring it anyway is what lets http.ResponseController
// report success truthfully to a streaming handler.
func (rw *responseWriter) Flush() {}

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

	var headers unsafe.Pointer
	var headersLen int
	if len(headerData) != 0 {
		headers = unsafe.Pointer(&headerData[0])
		headersLen = len(headerData)
	}

	C.URLSchemeTaskDidReceiveResponse(rw.r.task, C.int(code), headers, C.int(headersLen))
}

// urlSchemeTaskDidReceiveData hands data to the WKURLSchemeTask. It reports
// false when the task has already been stopped.
func urlSchemeTaskDidReceiveData(task unsafe.Pointer, data unsafe.Pointer, length int) bool {
	return bool(C.URLSchemeTaskDidReceiveData(task, data, C.int(length)))
}

// urlSchemeTaskDidFinish completes the WKURLSchemeTask.
func urlSchemeTaskDidFinish(task unsafe.Pointer) bool {
	return bool(C.URLSchemeTaskDidFinish(task))
}
