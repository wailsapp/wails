package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDownloadBytes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/redirect":
			http.Redirect(w, r, "/ok", http.StatusFound)
		case "/gzip":
			w.Header().Set("Content-Encoding", "gzip")
			writer := gzip.NewWriter(w)
			writer.Write([]byte("payload"))
			writer.Close()
		case "/unavailable":
			http.Error(w, "unavailable", http.StatusServiceUnavailable)
		case "/truncated":
			w.Header().Set("Content-Length", "100")
			w.Write([]byte("short"))
		default:
			w.Write([]byte("payload"))
		}
	}))
	defer server.Close()
	for _, path := range []string{"/ok", "/redirect", "/gzip", "/unavailable", "/truncated"} {
		t.Run(path, func(t *testing.T) {
			data, err := downloadBytes(server.URL + path)
			if path == "/unavailable" || path == "/truncated" {
				if err == nil {
					t.Fatal("failed download was accepted")
				}
				return
			}
			if err != nil || string(data) != "payload" {
				t.Fatalf("download = %q, %v", data, err)
			}
		})
	}
	if _, err := downloadBytes(":invalid"); err == nil {
		t.Fatal("invalid URL was accepted")
	}
}

type downloadRoundTrip func(*http.Request) (*http.Response, error)

func (f downloadRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type downloadBody struct {
	io.Reader
	closed bool
}

func (b *downloadBody) Close() error { b.closed = true; return nil }

func TestDownloadClosesResponseBody(t *testing.T) {
	previous := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = previous })
	for _, status := range []int{http.StatusOK, http.StatusServiceUnavailable} {
		body := &downloadBody{Reader: bytes.NewReader([]byte("payload"))}
		http.DefaultClient = &http.Client{Transport: downloadRoundTrip(func(r *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: status, Status: http.StatusText(status), Body: body, Header: make(http.Header)}, nil
		})}
		downloadBytes("https://example.invalid")
		if !body.closed {
			t.Errorf("HTTP %d response was not closed", status)
		}
	}
}

func TestDownloadIDLAndCache(t *testing.T) {
	t.Chdir(t.TempDir())
	var archive bytes.Buffer
	writer := zip.NewWriter(&archive)
	entry, err := writer.Create("WebView2.idl")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := entry.Write([]byte("fixture IDL")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	previous := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = previous })
	requests := 0
	http.DefaultClient = &http.Client{Transport: downloadRoundTrip(func(r *http.Request) (*http.Response, error) {
		requests++
		if r.URL.String() != "https://www.nuget.org/api/v2/package/Microsoft.Web.WebView2/1.2.3" {
			t.Errorf("unexpected URL %s", r.URL)
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(archive.Bytes())), Header: make(http.Header)}, nil
	})}
	for range 2 {
		data, err := DownloadIDL("1.2.3")
		if err != nil || string(data) != "fixture IDL" {
			t.Fatalf("IDL = %q, %v", data, err)
		}
	}
	data, err := os.ReadFile("WebView2.1.2.3.idl")
	if err != nil || string(data) != "fixture IDL" || requests != 1 {
		t.Fatalf("cache = %q, %v, requests=%d", data, err, requests)
	}
}

func BenchmarkDownloadBytes(b *testing.B) {
	payload := make([]byte, 4096)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(payload)
	}))
	defer server.Close()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		data, err := downloadBytes(server.URL)
		if err != nil || len(data) != len(payload) {
			b.Fatalf("download: len=%d err=%v", len(data), err)
		}
	}
}
