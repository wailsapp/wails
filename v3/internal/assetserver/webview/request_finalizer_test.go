package webview

import (
	"context"
	"io"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNativeRequestContextCancellationAndCleanup(t *testing.T) {
	const nativeID = uintptr(0x5963)
	raw := &finalizerTestRequest{}
	wrapped := newNativeRequestFinalizer(raw, nativeID)

	if !cancelRequest(nativeID) {
		t.Fatal("CancelRequest did not find the active native request")
	}
	select {
	case <-requestContext(wrapped).Done():
	case <-time.After(time.Second):
		t.Fatal("native cancellation did not cancel the request context")
	}

	if err := wrapped.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	if cancelRequest(nativeID) {
		t.Fatal("completed request remained in the native request registry")
	}
	if got := raw.closeCount.Load(); got != 1 {
		t.Fatalf("underlying request closed %d times, want 1", got)
	}
}

func TestNativeRequestPointerReuseKeepsReplacementRegistered(t *testing.T) {
	const nativeID = uintptr(0x5964)
	oldRequest := newNativeRequestFinalizer(&finalizerTestRequest{}, nativeID)
	newRequest := newNativeRequestFinalizer(&finalizerTestRequest{}, nativeID)

	select {
	case <-requestContext(oldRequest).Done():
	case <-time.After(time.Second):
		t.Fatal("replaced request context was not cancelled")
	}

	if err := oldRequest.Close(); err != nil {
		t.Fatalf("closing replaced request failed: %v", err)
	}
	if !cancelRequest(nativeID) {
		t.Fatal("closing replaced request removed the active replacement")
	}
	select {
	case <-requestContext(newRequest).Done():
	case <-time.After(time.Second):
		t.Fatal("replacement request context was not cancelled")
	}

	if err := newRequest.Close(); err != nil {
		t.Fatalf("closing replacement request failed: %v", err)
	}
	if cancelRequest(nativeID) {
		t.Fatal("replacement request leaked in the native request registry")
	}
}

func TestNativeRequestFinalizerCleansRegistration(t *testing.T) {
	const nativeID = uintptr(0x5966)
	raw := &finalizerTestRequest{}
	abandonNativeRequest(raw, nativeID)

	deadline := time.Now().Add(5 * time.Second)
	for raw.closeCount.Load() == 0 && time.Now().Before(deadline) {
		runtime.GC()
		runtime.Gosched()
		time.Sleep(time.Millisecond)
	}

	if got := raw.closeCount.Load(); got != 1 {
		t.Fatalf("underlying request closed %d times after finalization, want 1", got)
	}
	if cancelRequest(nativeID) {
		t.Fatal("finalized request remained in the native request registry")
	}
}

func TestNativeRequestConcurrentCancelAndClose(t *testing.T) {
	const nativeID = uintptr(0x5965)
	raw := &finalizerTestRequest{}
	wrapped := newNativeRequestFinalizer(raw, nativeID)

	var callers sync.WaitGroup
	for range 32 {
		callers.Add(2)
		go func() {
			defer callers.Done()
			cancelRequest(nativeID)
		}()
		go func() {
			defer callers.Done()
			_ = wrapped.Close()
		}()
	}
	callers.Wait()

	if cancelRequest(nativeID) {
		t.Fatal("concurrently closed request leaked in the native request registry")
	}
	if got := raw.closeCount.Load(); got != 1 {
		t.Fatalf("underlying request closed %d times, want 1", got)
	}
}

type finalizerTestRequest struct {
	closeCount atomic.Int32
}

func abandonNativeRequest(raw Request, nativeID uintptr) {
	wrapped := newNativeRequestFinalizer(raw, nativeID)
	runtime.KeepAlive(wrapped)
}

func requestContext(r Request) context.Context {
	ctx, _ := RequestContext(r)
	return ctx
}

func (*finalizerTestRequest) URL() (string, error)         { return "wails://localhost/test", nil }
func (*finalizerTestRequest) Method() (string, error)      { return http.MethodGet, nil }
func (*finalizerTestRequest) Header() (http.Header, error) { return http.Header{}, nil }
func (*finalizerTestRequest) Body() (io.ReadCloser, error) { return http.NoBody, nil }
func (*finalizerTestRequest) Response() ResponseWriter     { return &finalizerTestResponse{} }
func (r *finalizerTestRequest) Close() error               { r.closeCount.Add(1); return nil }

type finalizerTestResponse struct {
	header http.Header
}

func (r *finalizerTestResponse) Header() http.Header {
	if r.header == nil {
		r.header = make(http.Header)
	}
	return r.header
}
func (*finalizerTestResponse) Write(p []byte) (int, error) { return len(p), nil }
func (*finalizerTestResponse) WriteHeader(int)             {}
func (*finalizerTestResponse) Finish() error               { return nil }
func (*finalizerTestResponse) Code() int                   { return http.StatusOK }
