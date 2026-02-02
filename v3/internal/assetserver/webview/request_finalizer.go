package webview

import (
	"runtime"
	"sync/atomic"
)

var _ Request = &requestFinalizer{}

type requestFinalizer struct {
	Request
	closed int32
}

// newRequestFinalizer returns a request with a runtime finalizer to make sure it will be closed from the finalizer
// if it has not been already closed.
// It also makes sure Close() of the wrapping request is only called once.
func newRequestFinalizer(r Request) Request {
	rf := &requestFinalizer{Request: r}
	// Make sure to async release since it might block the finalizer goroutine for a longer period
	runtime.SetFinalizer(rf, func(obj *requestFinalizer) { rf.close(true) })
	return rf
}

func (r *requestFinalizer) Close() error {
	return r.close(false)
}

func (r *requestFinalizer) close(asyncRelease bool) error {
	if atomic.CompareAndSwapInt32(&r.closed, 0, 1) {
		runtime.SetFinalizer(r, nil)
		if asyncRelease {
			go r.Request.Close()
			return nil
		} else {
			return r.Request.Close()
		}
	}
	return nil
}
