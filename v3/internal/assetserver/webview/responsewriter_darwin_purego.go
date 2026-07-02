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
// The cgo bridge wraps each call in @try/@catch to swallow the NSException
// "This task has already been stopped" thrown by WebKit when the task was
// cancelled. An Objective-C @try/@catch cannot be expressed in pure Go/purego,
// so instead of catching the exception we AVOID it: the window backend calls
// MarkTaskStopped from its webView:stopURLSchemeTask: handler (delivered on
// the main thread), and each helper below performs its stopped-check and its
// message send as ONE synchronous main-thread unit (runOnMainSync). Because
// stop can only be delivered on the same thread, a task can never transition
// to stopped between the check and the send — the check-then-act race that a
// bare registry lookup would leave open. Write/Finish surface
// errRequestStopped, matching the cgo behaviour.

import (
	"net/http"
	"sync"
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

// stoppedTasks tracks WKURLSchemeTasks that WebKit has cancelled (via
// webView:stopURLSchemeTask:). Messaging a stopped task throws an NSException
// that cannot be caught in pure Go, so instead we skip the call and report the
// request as stopped. MarkTaskStopped is called from the window backend's
// stopURLSchemeTask: handler.
var stoppedTasks sync.Map // unsafe.Pointer -> struct{}

// MarkTaskStopped records that the given WKURLSchemeTask has been cancelled.
func MarkTaskStopped(task unsafe.Pointer) { stoppedTasks.Store(task, struct{}{}) }

func taskStopped(task unsafe.Pointer) bool {
	_, ok := stoppedTasks.Load(task)
	return ok
}

func clearTaskStopped(task unsafe.Pointer) { stoppedTasks.Delete(task) }

// urlSchemeTaskDidReceiveData wraps Go bytes in an NSData and calls
// -[task didReceiveData:]. Check + send run atomically on the main thread.
func urlSchemeTaskDidReceiveData(task unsafe.Pointer, data unsafe.Pointer, length int) bool {
	t := taskID(task)
	if t == 0 {
		return false
	}

	delivered := false
	runOnMainSync(func() {
		if taskStopped(task) {
			return
		}
		withAutoreleasePool(func() {
			nsdata := class("NSData").Send(sel("dataWithBytes:length:"), data, uint(length))
			t.Send(sel("didReceiveData:"), nsdata)
		})
		delivered = true
	})
	return delivered
}

// urlSchemeTaskDidFinish calls -[task didFinish]. Check + send run atomically
// on the main thread.
func urlSchemeTaskDidFinish(task unsafe.Pointer) bool {
	t := taskID(task)
	if t == 0 {
		return false
	}

	delivered := false
	runOnMainSync(func() {
		if taskStopped(task) {
			return
		}
		withAutoreleasePool(func() {
			t.Send(sel("didFinish"))
		})
		delivered = true
	})
	return delivered
}

// urlSchemeTaskDidReceiveResponse builds an NSHTTPURLResponse from the status
// code + headers and calls -[task didReceiveResponse:]. Check + send run
// atomically on the main thread.
func urlSchemeTaskDidReceiveResponse(task unsafe.Pointer, statusCode int, headers map[string]string) bool {
	t := taskID(task)
	if t == 0 {
		return false
	}

	delivered := false
	runOnMainSync(func() {
		if taskStopped(task) {
			return
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
		delivered = true
	})
	return delivered
}
