//go:build darwin && !ios && purego

package webview

// CGO-free (purego) implementation of the WKURLSchemeTask request bridge.
//
// The cgo sibling (request_darwin.go) compiles an Objective-C snippet that
// messages an `id<WKURLSchemeTask>` to read the incoming NSURLRequest. Here we
// drive the Objective-C runtime directly through github.com/ebitengine/purego
// and its objc subpackage, so the package builds with CGO_ENABLED=0.
//
// This file also hosts the small Objective-C helper layer (framework loading,
// selector cache, NSString/NSData conversion, autorelease pool) shared with
// responsewriter_darwin_purego.go — both files are `package webview` under the
// same build tag, so the helpers are declared once here.
//
// Threading: the cgo version calls the WKURLSchemeTask methods directly on
// whatever goroutine drives the asset request and relies on an Objective-C
// @try/@catch to swallow the "This task has already been stopped" NSException
// raised when a send races WebKit cancelling the task. Pure Go cannot catch
// Objective-C exceptions, so the purego bridge instead confines the
// stopped-check + send atomically to the main thread (runOnMainSync below):
// WebKit delivers webView:stopURLSchemeTask: on the main thread, so a task can
// never transition to stopped between our check and our send.

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

// ---------------------------------------------------------------------------
// Framework loading
// ---------------------------------------------------------------------------

const (
	frameworkFoundation = "/System/Library/Frameworks/Foundation.framework/Foundation"
	frameworkWebKit     = "/System/Library/Frameworks/WebKit.framework/WebKit"
)

var loadFrameworksOnce sync.Once

// loadFrameworks maps Foundation and WebKit into the process so objc_getClass
// can resolve NSString/NSData/NSHTTPURLResponse and the WebKit types. Classes
// only become resolvable once the framework that defines them is loaded.
func loadFrameworks() {
	loadFrameworksOnce.Do(func() {
		for _, fw := range []string{frameworkFoundation, frameworkWebKit} {
			// RTLD_GLOBAL so the class symbols become visible process-wide.
			_, _ = purego.Dlopen(fw, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		}
	})
}

// ---------------------------------------------------------------------------
// Main-thread confinement (libdispatch)
// ---------------------------------------------------------------------------

var (
	dispatchInitOnce  sync.Once
	dispatchMainQueue uintptr
	dispatchSyncFn    func(queue uintptr, block objc.Block)
)

func dispatchInit() {
	dispatchInitOnce.Do(func() {
		lib, err := purego.Dlopen("/usr/lib/libSystem.B.dylib", purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			panic("wails/purego: failed to load libSystem: " + err.Error())
		}
		mainQ, err := purego.Dlsym(lib, "_dispatch_main_q")
		if err != nil {
			panic("wails/purego: failed to resolve _dispatch_main_q: " + err.Error())
		}
		dispatchMainQueue = mainQ
		purego.RegisterLibFunc(&dispatchSyncFn, lib, "dispatch_sync")
	})
}

// runOnMainSync runs fn synchronously on the main thread. The WKURLSchemeTask
// stopped-check and message send must be one atomic main-thread unit (WebKit
// delivers stopURLSchemeTask: there), and dispatch_sync also keeps the caller's
// Go buffers alive for the duration of the send.
func runOnMainSync(fn func()) {
	if objc.Send[bool](class("NSThread"), sel("isMainThread")) {
		fn()
		return
	}
	dispatchInit()
	block := objc.NewBlock(func(objc.Block) { fn() })
	dispatchSyncFn(dispatchMainQueue, block)
	block.Release()
}

// ---------------------------------------------------------------------------
// Selector cache + class lookup
// ---------------------------------------------------------------------------

var (
	selMu    sync.RWMutex
	selCache = map[string]objc.SEL{}
)

// sel resolves (and caches) a selector by name. RegisterName takes the global
// Objective-C lock, so caching matters on the hot request/response path.
func sel(name string) objc.SEL {
	selMu.RLock()
	s, ok := selCache[name]
	selMu.RUnlock()
	if ok {
		return s
	}
	selMu.Lock()
	defer selMu.Unlock()
	if s, ok = selCache[name]; ok {
		return s
	}
	s = objc.RegisterName(name)
	selCache[name] = s
	return s
}

// class looks up a registered Objective-C class by name (loading the backing
// frameworks first) and returns it as an id so class methods can be sent to it.
func class(name string) objc.ID {
	loadFrameworks()
	return objc.ID(objc.GetClass(name))
}

// taskID reinterprets the raw WKURLSchemeTask pointer as an objc id.
func taskID(task unsafe.Pointer) objc.ID {
	return objc.ID(uintptr(task))
}

// ---------------------------------------------------------------------------
// String / data helpers
// ---------------------------------------------------------------------------

// nsString creates an autoreleased NSString from a Go string.
func nsString(s string) objc.ID {
	return class("NSString").Send(sel("stringWithUTF8String:"), s)
}

// nsStringToGo converts an NSString id back to a Go string via -UTF8String.
func nsStringToGo(s objc.ID) string {
	if s == 0 {
		return ""
	}
	cstr := objc.Send[uintptr](s, sel("UTF8String"))
	return goStringFromC(cstr)
}

// goStringFromC copies a NUL-terminated C string at the given address into a Go
// string without cgo.
func goStringFromC(p uintptr) string {
	if p == 0 {
		return ""
	}
	var n int
	for {
		if *(*byte)(unsafe.Pointer(p + uintptr(n))) == 0 {
			break
		}
		n++
	}
	if n == 0 {
		return ""
	}
	buf := make([]byte, n)
	copy(buf, unsafe.Slice((*byte)(unsafe.Pointer(p)), n))
	return string(buf)
}

// withAutoreleasePool runs fn inside a fresh NSAutoreleasePool and drains it
// afterwards, mirroring the @autoreleasepool blocks in the cgo bridge. The
// asset-request goroutine has no ambient pool, so without this the autoreleased
// NSData/NSString/NSHTTPURLResponse objects we create would leak.
func withAutoreleasePool(fn func()) {
	pool := class("NSAutoreleasePool").Send(sel("alloc")).Send(sel("init"))
	defer pool.Send(sel("drain"))
	fn()
}

// ---------------------------------------------------------------------------
// Request
// ---------------------------------------------------------------------------

// NewRequest creates a new WebViewRequest based on a pointer to an
// `id<WKURLSchemeTask>`.
func NewRequest(wkURLSchemeTask unsafe.Pointer) Request {
	// A freed task's heap address can be reused for a later task; drop any
	// stale stopped-registry entry so the new request doesn't inherit it.
	clearTaskStopped(wkURLSchemeTask)
	urlSchemeTaskRetain(wkURLSchemeTask)
	return newRequestFinalizer(&request{task: wkURLSchemeTask})
}

var _ Request = &request{}

type request struct {
	task unsafe.Pointer

	header http.Header
	body   io.ReadCloser
	rw     *responseWriter
}

// urlSchemeTaskRetain / urlSchemeTaskRelease balance the lifetime of the
// WKURLSchemeTask across the channel hop, exactly like the cgo -retain/-release.
func urlSchemeTaskRetain(task unsafe.Pointer) {
	taskID(task).Send(sel("retain"))
}

func urlSchemeTaskRelease(task unsafe.Pointer) {
	taskID(task).Send(sel("release"))
}

// nsurlRequest returns the task's NSURLRequest (`[task request]`), or 0.
func (r *request) nsurlRequest() objc.ID {
	t := taskID(r.task)
	if t == 0 {
		return 0
	}
	return t.Send(sel("request"))
}

func (r *request) URL() (string, error) {
	var url string
	withAutoreleasePool(func() {
		req := r.nsurlRequest()
		if req == 0 {
			return
		}
		url = nsStringToGo(req.Send(sel("URL")).Send(sel("absoluteString")))
	})
	return url, nil
}

func (r *request) Method() (string, error) {
	var method string
	withAutoreleasePool(func() {
		req := r.nsurlRequest()
		if req == 0 {
			return
		}
		method = nsStringToGo(req.Send(sel("HTTPMethod")))
	})
	return method, nil
}

func (r *request) Header() (http.Header, error) {
	if r.header != nil {
		return r.header, nil
	}

	header := http.Header{}
	withAutoreleasePool(func() {
		req := r.nsurlRequest()
		if req == 0 {
			return
		}
		// allHTTPHeaderFields is an NSDictionary<NSString*, NSString*> (single
		// value per key), so iterating it is equivalent to the cgo path that
		// JSON-serialised it into a map[string]string.
		dict := req.Send(sel("allHTTPHeaderFields"))
		if dict == 0 {
			return
		}
		keys := dict.Send(sel("allKeys"))
		count := objc.Send[uint](keys, sel("count"))
		for i := uint(0); i < count; i++ {
			key := keys.Send(sel("objectAtIndex:"), i)
			val := dict.Send(sel("objectForKey:"), key)
			header.Add(nsStringToGo(key), nsStringToGo(val))
		}
	})
	r.header = header
	return header, nil
}

func (r *request) Body() (io.ReadCloser, error) {
	if r.body != nil {
		return r.body, nil
	}

	withAutoreleasePool(func() {
		req := r.nsurlRequest()
		if req == 0 {
			return
		}

		// Prefer a materialised HTTPBody; fall back to the streaming body.
		if data := req.Send(sel("HTTPBody")); data != 0 {
			length := int(objc.Send[uint](data, sel("length")))
			if length > 0 {
				ptr := objc.Send[uintptr](data, sel("bytes"))
				// Copy out of the NSData before its pool drains.
				buf := make([]byte, length)
				copy(buf, unsafe.Slice((*byte)(unsafe.Pointer(ptr)), length))
				r.body = io.NopCloser(bytes.NewReader(buf))
			} else {
				r.body = http.NoBody
			}
			return
		}

		if stream := req.Send(sel("HTTPBodyStream")); stream != 0 {
			stream.Send(sel("open"))
			r.body = &requestBodyStreamReader{task: r.task}
		}
	})

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
	clearTaskStopped(r.task)
	urlSchemeTaskRelease(r.task)
	return err
}

// ---------------------------------------------------------------------------
// Streaming request body
// ---------------------------------------------------------------------------

// NSStreamStatus values (Foundation) used by the streaming body reader.
const (
	nsStreamStatusOpen  = 2
	nsStreamStatusAtEnd = 5
)

var _ io.ReadCloser = &requestBodyStreamReader{}

type requestBodyStreamReader struct {
	task   unsafe.Pointer
	closed bool
}

// stream returns the task's HTTPBodyStream (NSInputStream), or 0.
func (r *requestBodyStreamReader) stream() objc.ID {
	t := taskID(r.task)
	if t == 0 {
		return 0
	}
	req := t.Send(sel("request"))
	if req == 0 {
		return 0
	}
	return req.Send(sel("HTTPBodyStream"))
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

// readStream mirrors the numeric contract of the cgo helper
// URLSchemeTaskRequestBodyStreamRead: >0 bytes read, 0 at end / no bytes,
// -1 read error, -2 no stream, -3 stream not open.
func (r *requestBodyStreamReader) readStream(buf unsafe.Pointer, bufLen int) int {
	var res int
	withAutoreleasePool(func() {
		stream := r.stream()
		if stream == 0 {
			res = -2
			return
		}

		status := objc.Send[int](stream, sel("streamStatus"))
		hasBytes := objc.Send[bool](stream, sel("hasBytesAvailable"))
		if status == nsStreamStatusAtEnd || !hasBytes {
			res = 0
			return
		} else if status != nsStreamStatusOpen {
			res = -3
			return
		}

		// -[NSInputStream read:maxLength:] returns NSInteger.
		res = objc.Send[int](stream, sel("read:maxLength:"), buf, uint(bufLen))
	})
	return res
}

func (r *requestBodyStreamReader) Close() error {
	if r.closed {
		return nil
	}
	r.closed = true

	withAutoreleasePool(func() {
		if stream := r.stream(); stream != 0 {
			stream.Send(sel("close"))
		}
	})
	return nil
}
