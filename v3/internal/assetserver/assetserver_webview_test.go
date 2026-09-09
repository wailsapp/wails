package assetserver

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
)

func TestWebViewRequestCancellationReachesHandler(t *testing.T) {
	started := make(chan struct{})
	contextErr := make(chan error, 1)

	srv, err := NewAssetServer(&Options{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			close(started)
			<-r.Context().Done()
			contextErr <- r.Context().Err()
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(bytes.Repeat([]byte("x"), 513))
		}),
		Logger: slog.Default(),
	})
	if err != nil {
		t.Fatalf("NewAssetServer failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	req := &contextWebViewRequest{
		ctx: ctx,
		rw:  &contextWebViewResponse{},
	}
	processed := make(chan struct{})
	go func() {
		srv.processWebViewRequest(req)
		close(processed)
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("handler did not start")
	}
	cancel()

	select {
	case err := <-contextErr:
		if err != context.Canceled {
			t.Fatalf("handler context error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("handler did not observe request cancellation")
	}
	select {
	case <-processed:
	case <-time.After(time.Second):
		t.Fatal("request processing did not finish after cancellation")
	}
	if got := req.closeCount.Load(); got != 1 {
		t.Fatalf("request closed %d times, want 1", got)
	}
	if got := req.rw.writeHeaderCount.Load(); got != 0 {
		t.Fatalf("response header written %d times after cancellation, want 0", got)
	}
	if got := req.rw.writeCount.Load(); got != 0 {
		t.Fatalf("response body written %d times after cancellation, want 0", got)
	}
	if got := req.rw.finishCount.Load(); got != 0 {
		t.Fatalf("response finished %d times after cancellation, want 0", got)
	}
}

type contextWebViewRequest struct {
	ctx        context.Context
	rw         *contextWebViewResponse
	closeCount atomic.Int32
}

func (r *contextWebViewRequest) Context() context.Context         { return r.ctx }
func (*contextWebViewRequest) URL() (string, error)               { return "wails://localhost/slow", nil }
func (*contextWebViewRequest) Method() (string, error)            { return http.MethodGet, nil }
func (*contextWebViewRequest) Header() (http.Header, error)       { return http.Header{}, nil }
func (*contextWebViewRequest) Body() (io.ReadCloser, error)       { return http.NoBody, nil }
func (r *contextWebViewRequest) Response() webview.ResponseWriter { return r.rw }
func (r *contextWebViewRequest) Close() error                     { r.closeCount.Add(1); return nil }

type contextWebViewResponse struct {
	header           http.Header
	body             bytes.Buffer
	code             int
	writeHeaderCount atomic.Int32
	writeCount       atomic.Int32
	finishCount      atomic.Int32
}

func (r *contextWebViewResponse) Header() http.Header {
	if r.header == nil {
		r.header = make(http.Header)
	}
	return r.header
}
func (r *contextWebViewResponse) Write(p []byte) (int, error) {
	r.writeCount.Add(1)
	return r.body.Write(p)
}
func (r *contextWebViewResponse) WriteHeader(code int) {
	r.writeHeaderCount.Add(1)
	r.code = code
}
func (r *contextWebViewResponse) Finish() error {
	r.finishCount.Add(1)
	return nil
}
func (r *contextWebViewResponse) Code() int { return r.code }
