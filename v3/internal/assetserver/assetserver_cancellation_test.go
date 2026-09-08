package assetserver

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
)

type cancellationRequest struct {
	ctx      context.Context
	response *cancellationResponse
}

func (r *cancellationRequest) Context() context.Context         { return r.ctx }
func (r *cancellationRequest) URL() (string, error)             { return "wails://localhost/slow", nil }
func (r *cancellationRequest) Method() (string, error)          { return http.MethodGet, nil }
func (r *cancellationRequest) Header() (http.Header, error)     { return http.Header{}, nil }
func (r *cancellationRequest) Body() (io.ReadCloser, error)     { return http.NoBody, nil }
func (r *cancellationRequest) Response() webview.ResponseWriter { return r.response }
func (r *cancellationRequest) Close() error                     { return nil }

type cancellationResponse struct{ *httptest.ResponseRecorder }

func (r *cancellationResponse) Finish() error { return nil }
func (r *cancellationResponse) Code() int     { return r.ResponseRecorder.Code }

func TestNativeAbortCancelsRunningHandler(t *testing.T) {
	ctx, abort := context.WithCancel(context.Background())
	defer abort()
	started, finished := make(chan struct{}), make(chan struct{})
	observed := make(chan error, 1)
	stop := make(chan struct{})
	defer close(stop)
	srv, err := NewAssetServer(&Options{Logger: slog.Default(), Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		select {
		case <-r.Context().Done():
			observed <- r.Context().Err()
		case <-stop:
		}
	})})
	if err != nil {
		t.Fatal(err)
	}
	request := &cancellationRequest{ctx: ctx, response: &cancellationResponse{httptest.NewRecorder()}}
	go func() { defer close(finished); srv.processWebViewRequest(request) }()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("handler never started")
	}
	abort()
	select {
	case err := <-observed:
		if err != context.Canceled {
			t.Fatalf("context error = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("native abort did not reach running handler")
	}
	<-finished
}
