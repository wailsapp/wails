package assetserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// A response shorter than the 512-byte sniffing prefix is held back until the
// handler returns. Flush has to release it, or a streaming handler that writes
// a short chunk and then waits delivers nothing while reporting success.
func TestContentTypeSnifferFlushReleasesShortPrefix(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := newContentTypeSniffer(rec)

	rw.WriteHeader(http.StatusOK) // no Content-Type set, so sniffing applies
	if _, err := io.WriteString(rw, "data: first\n\n"); err != nil {
		t.Fatalf("write: %v", err)
	}

	if got := rec.Body.Len(); got != 0 {
		t.Fatalf("prefix should still be buffered before Flush, got %d bytes", got)
	}

	rw.Flush()

	if got := rec.Body.String(); got != "data: first\n\n" {
		t.Errorf("Flush did not release the buffered prefix, body = %q", got)
	}
	if !rec.Flushed {
		t.Error("Flush did not reach the wrapped writer")
	}
	if ct := rec.Header().Get(HeaderContentType); ct == "" {
		t.Error("Content-Type was not resolved on flush")
	}
}

// After a flush the sniffer is done buffering, so later writes must pass
// straight through rather than being collected into a new prefix.
func TestContentTypeSnifferWritesPassThroughAfterFlush(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := newContentTypeSniffer(rec)

	rw.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(rw, "one")
	rw.Flush()

	if _, err := io.WriteString(rw, "two"); err != nil {
		t.Fatalf("write after flush: %v", err)
	}
	if got := rec.Body.String(); got != "onetwo" {
		t.Errorf("body = %q, want %q", got, "onetwo")
	}
}

// Flushing before anything is written must not panic or emit a header.
func TestContentTypeSnifferFlushBeforeWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := newContentTypeSniffer(rec)

	rw.Flush()

	if rec.Body.Len() != 0 {
		t.Errorf("expected no body, got %q", rec.Body.String())
	}
}

// An explicit Content-Type disables sniffing entirely; Flush must not disturb
// the header or duplicate the body.
func TestContentTypeSnifferFlushWithExplicitContentType(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := newContentTypeSniffer(rec)

	rw.Header().Set(HeaderContentType, "text/event-stream")
	rw.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(rw, "data: x\n\n")
	rw.Flush()

	if got := rec.Header().Get(HeaderContentType); got != "text/event-stream" {
		t.Errorf("Content-Type = %q, want text/event-stream", got)
	}
	if got := rec.Body.String(); got != "data: x\n\n" {
		t.Errorf("body = %q", got)
	}
}

// The existing behaviour for a full prefix must be unchanged: once 512 bytes
// have accumulated the sniffer completes on its own.
func TestContentTypeSnifferFullPrefixStillCompletesWithoutFlush(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := newContentTypeSniffer(rec)

	body := strings.Repeat("a", 600)
	rw.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(rw, body); err != nil {
		t.Fatalf("write: %v", err)
	}

	if got := rec.Body.String(); got != body {
		t.Errorf("body length = %d, want %d", len(got), len(body))
	}
	if ct := rec.Header().Get(HeaderContentType); ct == "" {
		t.Error("Content-Type was not sniffed once the prefix filled")
	}
}

// http.ResponseController must find the sniffer's Flush and report success.
func TestContentTypeSnifferSupportsResponseController(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := newContentTypeSniffer(rec)

	rw.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(rw, "chunk")

	if err := http.NewResponseController(rw).Flush(); err != nil {
		t.Fatalf("ResponseController.Flush: %v", err)
	}
	if got := rec.Body.String(); got != "chunk" {
		t.Errorf("body = %q, want %q", got, "chunk")
	}
}
