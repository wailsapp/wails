package webview

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
)

type contextRequest struct {
	Request
	ctx    context.Context
	closed atomic.Int32
}

func (r *contextRequest) Context() context.Context { return r.ctx }
func (r *contextRequest) Close() error             { r.closed.Add(1); return nil }

func TestFinalizerPreservesRequestCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	native := &contextRequest{ctx: ctx}
	request := newRequestFinalizer(native)
	defer request.Close()
	handlerContext := Context(request)
	if handlerContext.Err() != nil {
		t.Fatal("request cancelled before native abort")
	}
	cancel()
	if handlerContext.Err() != context.Canceled {
		t.Fatal("native abort lost through finalizer wrapper")
	}
	var workers sync.WaitGroup
	for i := 0; i < 16; i++ {
		workers.Add(1)
		go func() { defer workers.Done(); request.Close() }()
	}
	workers.Wait()
	if n := native.closed.Load(); n != 1 {
		t.Fatalf("native request closed %d times", n)
	}
}

func TestRequestWithoutNativeCancellation(t *testing.T) {
	// A legacy backend only implements Request, with no optional Context method.
	r := struct{ Request }{}
	if Context(r).Done() != nil {
		t.Fatal("backend without cancellation needs a background context")
	}
}
