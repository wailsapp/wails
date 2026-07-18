//go:build !production

package application

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

type trackingReadCloser struct {
	closed bool
}

func (*trackingReadCloser) Read([]byte) (int, error) { return 0, io.EOF }

func (body *trackingReadCloser) Close() error {
	body.closed = true
	return nil
}

func TestWaitForFrontendDevServerRetriesUntilReady(t *testing.T) {
	body := &trackingReadCloser{}
	attempts := 0
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		attempts++
		if attempts < 3 {
			return nil, errors.New("server is still starting")
		}
		return &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Body:       body,
			Header:     make(http.Header),
		}, nil
	})}

	retries := 0
	err := waitForFrontendDevServer(context.Background(), client, "http://localhost:9245", func() {
		retries++
	})
	if err != nil {
		t.Fatalf("waitForFrontendDevServer returned an error: %v", err)
	}
	if attempts != 3 {
		t.Fatalf("attempt count = %d, want 3", attempts)
	}
	if retries != 2 {
		t.Fatalf("retry count = %d, want 2", retries)
	}
	if !body.closed {
		t.Fatal("successful probe response body was not closed")
	}
}

func TestWaitForFrontendDevServerStopsWhenCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		cancel()
		return nil, errors.New("server is unavailable")
	})}

	err := waitForFrontendDevServer(ctx, client, "http://localhost:9245", nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("waitForFrontendDevServer error = %v, want context.Canceled", err)
	}
}

func TestWaitForFrontendDevServerRejectsInvalidURL(t *testing.T) {
	err := waitForFrontendDevServer(context.Background(), http.DefaultClient, "://not-a-url", nil)
	if err == nil || !strings.Contains(err.Error(), "invalid frontend dev server URL") {
		t.Fatalf("waitForFrontendDevServer error = %v, want invalid URL error", err)
	}
}
