//go:build !production

package assetserver

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type retryTestTransport func(*http.Request) (*http.Response, error)

func (f retryTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestRetryTransportStopsAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	calls := 0
	transport := &retryTransport{
		base: retryTestTransport(func(*http.Request) (*http.Response, error) {
			calls++
			cancel()
			return nil, errors.New("connection refused")
		}),
		maxRetries: 3,
		delay:      time.Millisecond,
	}
	request := httptest.NewRequest(http.MethodGet, "http://localhost/asset.js", nil).WithContext(ctx)
	if _, err := transport.RoundTrip(request); !errors.Is(err, context.Canceled) {
		t.Errorf("RoundTrip error = %v, want context.Canceled", err)
	}
	if calls != 1 {
		t.Errorf("cancelled request made %d attempts, want 1", calls)
	}
	if _, err := transport.RoundTrip(request); !errors.Is(err, context.Canceled) || calls != 1 {
		t.Errorf("already-cancelled request: error %v, total attempts %d", err, calls)
	}
}

func TestRetryTransportCancellationInterruptsDelay(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	called := make(chan struct{}, 2)
	transport := &retryTransport{
		base: retryTestTransport(func(*http.Request) (*http.Response, error) {
			called <- struct{}{}
			return nil, errors.New("connection refused")
		}),
		maxRetries: 2,
		delay:      time.Hour,
	}
	result := make(chan error, 1)
	go func() {
		request := httptest.NewRequest(http.MethodGet, "http://localhost/asset.js", nil).WithContext(ctx)
		_, err := transport.RoundTrip(request)
		result <- err
	}()
	<-called
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("RoundTrip error = %v, want context.Canceled", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("cancellation did not interrupt the retry delay")
	}
}

func TestRetryTransportRetriesActiveRequest(t *testing.T) {
	calls := 0
	transport := &retryTransport{
		base: retryTestTransport(func(*http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return nil, errors.New("connection refused")
			}
			return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
		}),
		maxRetries: 3,
	}
	response, err := transport.RoundTrip(httptest.NewRequest(http.MethodGet, "http://localhost/asset.js", nil))
	if err != nil || response == nil || response.StatusCode != http.StatusOK || calls != 2 {
		t.Errorf("active request: response %v, error %v, attempts %d; want 200, nil, 2", response, err, calls)
	}
}

type retryTestBody struct {
	io.Reader
	closed bool
}

func (b *retryTestBody) Close() error { b.closed = true; return nil }

func TestRetryTransportClosesCancelledRequestBody(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	body := &retryTestBody{Reader: http.NoBody}
	request := httptest.NewRequest(http.MethodPost, "http://localhost/asset.js", body).WithContext(ctx)
	transport := &retryTransport{maxRetries: 3}
	if _, err := transport.RoundTrip(request); !errors.Is(err, context.Canceled) {
		t.Errorf("RoundTrip error = %v, want context.Canceled", err)
	}
	if !body.closed {
		t.Error("already-cancelled request body was not closed")
	}
}
