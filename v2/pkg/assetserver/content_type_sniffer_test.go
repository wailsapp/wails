package assetserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContentTypeSnifferWasmFile(t *testing.T) {
	rec := httptest.NewRecorder()
	sniffer := &contentTypeSniffer{rw: rec, reqPath: "/app_bg.wasm"}

	wasmData := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}
	_, err := sniffer.Write(wasmData)
	if err != nil {
		t.Fatal(err)
	}

	ct := rec.Header().Get(HeaderContentType)
	if ct != "application/wasm" {
		t.Errorf("expected Content-Type 'application/wasm', got '%s'", ct)
	}
}

func TestContentTypeSnifferJSFile(t *testing.T) {
	rec := httptest.NewRecorder()
	sniffer := &contentTypeSniffer{rw: rec, reqPath: "/app.js"}

	_, err := sniffer.Write([]byte("console.log(1)"))
	if err != nil {
		t.Fatal(err)
	}

	ct := rec.Header().Get(HeaderContentType)
	if ct != "text/javascript; charset=utf-8" {
		t.Errorf("expected Content-Type 'text/javascript; charset=utf-8', got '%s'", ct)
	}
}

func TestContentTypeSnifferHTMLFile(t *testing.T) {
	rec := httptest.NewRecorder()
	sniffer := &contentTypeSniffer{rw: rec, reqPath: "/index.html"}

	_, err := sniffer.Write([]byte("<html></html>"))
	if err != nil {
		t.Fatal(err)
	}

	ct := rec.Header().Get(HeaderContentType)
	if ct != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type 'text/html; charset=utf-8', got '%s'", ct)
	}
}

func TestContentTypeSnifferPreservesExplicitContentType(t *testing.T) {
	rec := httptest.NewRecorder()
	sniffer := &contentTypeSniffer{rw: rec, reqPath: "/app.wasm"}

	sniffer.Header().Set(HeaderContentType, "custom/type")
	_, err := sniffer.Write([]byte("data"))
	if err != nil {
		t.Fatal(err)
	}

	ct := rec.Header().Get(HeaderContentType)
	if ct != "custom/type" {
		t.Errorf("expected Content-Type 'custom/type' to be preserved, got '%s'", ct)
	}
}

func TestContentTypeSnifferUnknownFileFallsBack(t *testing.T) {
	rec := httptest.NewRecorder()
	sniffer := &contentTypeSniffer{rw: rec, reqPath: "/unknown.dat"}

	_, err := sniffer.Write([]byte("some binary data"))
	if err != nil {
		t.Fatal(err)
	}

	ct := rec.Header().Get(HeaderContentType)
	if ct == "" {
		t.Error("expected a Content-Type to be set")
	}
}

func TestContentTypeSnifferNoReqPath(t *testing.T) {
	rec := httptest.NewRecorder()
	sniffer := &contentTypeSniffer{rw: rec}

	wasmData := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}
	_, err := sniffer.Write(wasmData)
	if err != nil {
		t.Fatal(err)
	}

	ct := rec.Header().Get(HeaderContentType)
	if ct == "" {
		t.Error("expected a Content-Type to be set")
	}
}

func TestGetMimeTypeByExt(t *testing.T) {
	tests := []struct {
		ext      string
		expected string
	}{
		{".wasm", "application/wasm"},
		{".js", "text/javascript; charset=utf-8"},
		{".css", "text/css; charset=utf-8"},
		{".html", "text/html; charset=utf-8"},
		{".json", "application/json"},
		{".svg", "image/svg+xml"},
		{".xyz", ""},
	}
	for _, tt := range tests {
		got := getMimeTypeByExt(tt.ext)
		if got != tt.expected {
			t.Errorf("getMimeTypeByExt(%q) = %q, want %q", tt.ext, got, tt.expected)
		}
	}
}

func TestContentTypeSnifferWriteHeaderOnce(t *testing.T) {
	rec := httptest.NewRecorder()
	sniffer := &contentTypeSniffer{rw: rec, reqPath: "/test.txt"}

	sniffer.WriteHeader(http.StatusOK)
	sniffer.WriteHeader(http.StatusBadRequest)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
