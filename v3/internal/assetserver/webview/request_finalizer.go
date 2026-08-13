package webview

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

var _ Request = &requestFinalizer{}

type requestFinalizer struct {
	Request
	context  context.Context
	cancel   context.CancelFunc
	nativeID uintptr
	closed   int32
}

var nativeRequestContexts = struct {
	sync.RWMutex
	active map[uintptr]*requestFinalizer
}{
	active: make(map[uintptr]*requestFinalizer),
}

// newRequestFinalizer returns a request with a runtime finalizer to make sure it will be closed from the finalizer
// if it has not been already closed.
// It also makes sure Close() of the wrapping request is only called once.
func newRequestFinalizer(r Request) Request {
	return newNativeRequestFinalizer(r, 0)
}

// newNativeRequestFinalizer additionally associates the request with a native
// scheme task that can report cancellation before the handler completes.
func newNativeRequestFinalizer(r Request, nativeID uintptr) Request {
	rf := &requestFinalizer{
		Request:  r,
		nativeID: nativeID,
	}
	if nativeID != 0 {
		rf.context, rf.cancel = context.WithCancel(context.Background())
		nativeRequestContexts.Lock()
		previous := nativeRequestContexts.active[nativeID]
		nativeRequestContexts.active[nativeID] = rf
		nativeRequestContexts.Unlock()

		// A retained native task should keep its address unique for its lifetime.
		// If a platform nevertheless reuses an address before cleanup, retire the
		// old context without letting its eventual Close remove the replacement.
		if previous != nil {
			previous.cancel()
		}
	}
	// Make sure to async release since it might block the finalizer goroutine for a longer period
	runtime.SetFinalizer(rf, func(obj *requestFinalizer) { rf.close(true) })
	return rf
}

// RequestContext returns the lifecycle context associated with a WebView request.
// Requests without a native lifecycle wrapper retain the historical background context.
func RequestContext(r Request) (context.Context, context.CancelFunc) {
	if contextual, ok := r.(interface{ Context() context.Context }); ok {
		if ctx := contextual.Context(); ctx != nil {
			return ctx, func() {}
		}
	}
	return context.WithCancel(context.Background())
}

// CancelRequest cancels the request associated with a native scheme task. It
// returns false when the task has already completed or was never registered.
func CancelRequest(nativeRequest unsafe.Pointer) bool {
	if nativeRequest == nil {
		return false
	}
	return cancelRequest(uintptr(nativeRequest))
}

func cancelRequest(nativeID uintptr) bool {
	nativeRequestContexts.RLock()
	rf := nativeRequestContexts.active[nativeID]
	nativeRequestContexts.RUnlock()
	if rf == nil {
		return false
	}
	rf.cancel()
	return true
}

func (r *requestFinalizer) Context() context.Context {
	return r.context
}

func (r *requestFinalizer) Close() error {
	return r.close(false)
}

func (r *requestFinalizer) close(asyncRelease bool) error {
	if atomic.CompareAndSwapInt32(&r.closed, 0, 1) {
		runtime.SetFinalizer(r, nil)
		if r.cancel != nil {
			r.cancel()
		}
		if r.nativeID != 0 {
			nativeRequestContexts.Lock()
			if nativeRequestContexts.active[r.nativeID] == r {
				delete(nativeRequestContexts.active, r.nativeID)
			}
			nativeRequestContexts.Unlock()
		}
		if asyncRelease {
			go r.Request.Close()
			return nil
		} else {
			return r.Request.Close()
		}
	}
	return nil
}
