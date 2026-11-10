package application

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestHTTPTransportCleanupLifecycle(t *testing.T) {
	transport := NewHTTPTransport()
	t.Cleanup(func() { _ = transport.Stop() })

	for cycle := 0; cycle < 100; cycle++ {
		if err := transport.Start(context.Background(), nil); err != nil {
			t.Fatalf("cycle %d: Start: %v", cycle, err)
		}

		transport.cleanupMu.Lock()
		done := transport.cleanupDone
		transport.cleanupMu.Unlock()
		if done == nil {
			t.Fatalf("cycle %d: cleanup goroutine was not started", cycle)
		}
		if err := transport.Start(context.Background(), nil); err != nil {
			t.Fatalf("cycle %d: repeated Start: %v", cycle, err)
		}
		transport.cleanupMu.Lock()
		repeatedDone := transport.cleanupDone
		transport.cleanupMu.Unlock()
		if repeatedDone != done {
			t.Fatalf("cycle %d: repeated Start launched another cleanup goroutine", cycle)
		}

		// Stop is allowed to be called defensively from several shutdown paths.
		// It must close the channel only once and the call that owns the cleanup
		// goroutine must wait for it to finish before returning.
		var stops sync.WaitGroup
		stops.Add(4)
		for range 4 {
			go func() {
				defer stops.Done()
				if err := transport.Stop(); err != nil {
					t.Errorf("cycle %d: Stop: %v", cycle, err)
				}
			}()
		}
		stops.Wait()

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatalf("cycle %d: cleanup goroutine did not stop", cycle)
		}
	}
}

// Supply JSON followed by a streamed padding body so the fixture itself does not
// allocate tens of megabytes. Missing object/method keeps accepted requests from
// needing a live application; oversized requests must fail before JSON parsing.
func TestHTTPTransportRuntimeBodyLimit(t *testing.T) {
	for _, size := range []int64{maxRuntimeBodyBytes - 1, maxRuntimeBodyBytes, maxRuntimeBodyBytes + 1} {
		for _, knownLength := range []bool{false, true} {
			t.Run(fmt.Sprintf("size=%d/knownLength=%t", size, knownLength), func(t *testing.T) {
				padding := &runtimeBodyPadding{remaining: size - 2}
				body := io.MultiReader(strings.NewReader("{}"), padding)
				req := httptest.NewRequest(http.MethodPost, "/wails/runtime", body)
				if knownLength {
					req.ContentLength = size
				}
				rec := httptest.NewRecorder()
				transport := NewHTTPTransport()
				transport.handleRuntimeRequest(rec, req)
				want := http.StatusUnprocessableEntity
				if size > maxRuntimeBodyBytes {
					want = http.StatusRequestEntityTooLarge
				}
				if rec.Code != want {
					t.Fatalf("status = %d, want %d: %s", rec.Code, want, rec.Body.String())
				}
			})
		}
	}
}

type runtimeBodyPadding struct{ remaining int64 }

func (p *runtimeBodyPadding) Read(buf []byte) (int, error) {
	if p.remaining == 0 {
		return 0, io.EOF
	}
	n := int(min(int64(len(buf)), p.remaining))
	for i := range buf[:n] {
		buf[i] = ' '
	}
	p.remaining -= int64(n)
	return n, nil
}

func TestHTTPTransportRuntimeBodyStopsReadingAtLimit(t *testing.T) {
	padding := &runtimeBodyPadding{remaining: 2 * maxRuntimeBodyBytes}
	req := httptest.NewRequest(http.MethodPost, "/wails/runtime", padding)
	rec := httptest.NewRecorder()
	NewHTTPTransport().handleRuntimeRequest(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
	if padding.remaining != maxRuntimeBodyBytes-1 {
		t.Fatalf("read beyond limit probe: %d bytes remain", padding.remaining)
	}
}

// A chunk the client sends again after a failed attempt replaces the first copy; an upload of exactly the limit must
// not be refused because the retry was counted twice.
func TestHTTPTransportChunkRetryDoesNotCountTwice(t *testing.T) {
	const chunks = maxAssembledBytes / maxChunkBodyBytes
	transport := NewHTTPTransport()
	send := func(index int) *httptest.ResponseRecorder {
		body := bytes.Repeat([]byte(" "), maxChunkBodyBytes)
		if index == 0 {
			copy(body, "{}")
		}
		req := httptest.NewRequest(http.MethodPost, "/wails/runtime", bytes.NewReader(body))
		req.Header.Set("x-wails-client-id", "test")
		req.Header.Set(chunkIDHeader, "retry")
		req.Header.Set(chunkIndexHeader, strconv.Itoa(index))
		req.Header.Set(chunkTotalHeader, strconv.Itoa(chunks))
		rec := httptest.NewRecorder()
		transport.handleRuntimeRequest(rec, req)
		return rec
	}
	if rec := send(0); rec.Code != http.StatusOK {
		t.Fatalf("chunk 0: %d", rec.Code)
	}
	if rec := send(0); rec.Code != http.StatusOK {
		t.Fatalf("retried chunk 0: %d: %s", rec.Code, rec.Body.String())
	}
	var last *httptest.ResponseRecorder
	for i := 1; i < chunks; i++ {
		last = send(i)
	}
	// The assembled body is "{}" plus spaces, which parses but names no method, so it is answered 422; a refusal
	// for size is also a 422, told apart by its message.
	if strings.Contains(last.Body.String(), "too large") {
		t.Fatalf("final chunk refused for size: %s", last.Body.String())
	}
	if last.Code != http.StatusUnprocessableEntity {
		t.Fatalf("final chunk: %d, want the assembled body processed: %s", last.Code, last.Body.String())
	}
}
