//go:build darwin && !ios && purego

package webview

// CGO-free (purego) implementation of the WKURLSchemeTask response bridge.
//
// The cgo sibling (responsewriter_darwin.go) messages an `id<WKURLSchemeTask>`
// with -didReceiveResponse:/-didReceiveData:/-didFinish:. We reproduce that here
// through the Objective-C runtime via github.com/ebitengine/purego/objc, reusing
// the helper layer declared in request_darwin_purego.go (class/sel/nsString/
// taskID/withAutoreleasePool).
//
// LIMITATION vs cgo: the cgo bridge wraps each call in @try/@catch to swallow the
// NSException "This task has already been stopped" thrown by WebKit when the task
// was cancelled, returning false so Write/Finish surface errRequestStopped. An
// Objective-C @try/@catch cannot be expressed in pure Go/purego (it needs the
// zero-cost exception personality routine emitted by clang), so the helpers below
// always report success. If a stopped task is written to, the underlying
// NSException will propagate rather than being converted to errRequestStopped;
// catching it would require a small native (cgo/assembly) shim.

import (
	"net/http"
	"unsafe"

	"github.com/ebitengine/purego/objc"
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

	var content unsafe.Pointer
	var contentLen int
	if len(buf) > 0 {
		content = unsafe.Pointer(&buf[0])
		contentLen = len(buf)
	}

	if !urlSchemeTaskDidReceiveData(rw.r.task, content, contentLen) {
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

	// Flatten to a single value per key, matching the cgo path (which JSON
	// encoded a map[string]string and rebuilt an NSDictionary from it).
	header := map[string]string{}
	for k := range rw.Header() {
		header[k] = rw.Header().Get(k)
	}

	urlSchemeTaskDidReceiveResponse(rw.r.task, code, header)
}

func (rw *responseWriter) Finish() error {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusNotImplemented)
	}

	if rw.finished {
		return nil
	}
	rw.finished = true

	urlSchemeTaskDidFinish(rw.r.task)

	return nil
}

func (rw *responseWriter) Code() int {
	return rw.code
}

// ---------------------------------------------------------------------------
// Objective-C bridge (replaces the cgo C helpers)
//
// Each returns a bool mirroring the cgo contract (true == delivered). Because we
// cannot catch the "task already stopped" NSException in pure Go, the success
// path always returns true; see the file header note.
// ---------------------------------------------------------------------------

// urlSchemeTaskDidReceiveData wraps Go bytes in an NSData and calls
// -[task didReceiveData:].
func urlSchemeTaskDidReceiveData(task unsafe.Pointer, data unsafe.Pointer, length int) bool {
	t := taskID(task)
	if t == 0 {
		return false
	}

	withAutoreleasePool(func() {
		nsdata := class("NSData").Send(sel("dataWithBytes:length:"), data, uint(length))
		t.Send(sel("didReceiveData:"), nsdata)
	})
	return true
}

// urlSchemeTaskDidFinish calls -[task didFinish].
func urlSchemeTaskDidFinish(task unsafe.Pointer) bool {
	t := taskID(task)
	if t == 0 {
		return false
	}

	withAutoreleasePool(func() {
		t.Send(sel("didFinish"))
	})
	return true
}

// urlSchemeTaskDidReceiveResponse builds an NSHTTPURLResponse from the status
// code + headers and calls -[task didReceiveResponse:].
func urlSchemeTaskDidReceiveResponse(task unsafe.Pointer, statusCode int, headers map[string]string) bool {
	t := taskID(task)
	if t == 0 {
		return false
	}

	withAutoreleasePool(func() {
		// The response URL comes from the originating request, like the cgo code.
		var url objc.ID
		if req := t.Send(sel("request")); req != 0 {
			url = req.Send(sel("URL"))
		}

		headerFields := class("NSMutableDictionary").Send(sel("dictionary"))
		for k, v := range headers {
			headerFields.Send(sel("setObject:forKey:"), nsString(v), nsString(k))
		}

		// [[[NSHTTPURLResponse alloc] initWithURL:statusCode:HTTPVersion:headerFields:] autorelease]
		response := class("NSHTTPURLResponse").Send(sel("alloc")).Send(
			sel("initWithURL:statusCode:HTTPVersion:headerFields:"),
			url, statusCode, nsString("HTTP/1.1"), headerFields,
		)
		response.Send(sel("autorelease"))

		t.Send(sel("didReceiveResponse:"), response)
	})
	return true
}
