package appimage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
)

func TestToolDownloadCachesAndPreservesOnFailedRefresh(t *testing.T) {
	var calls atomic.Int32
	var fail atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if fail.Load() {
			http.Error(w, "unavailable", 503)
			return
		}
		_, _ = w.Write([]byte("tool contents"))
	}))
	defer server.Close()
	filename := filepath.Join(t.TempDir(), "cache", "tool")
	for range 2 {
		if err := ensureTool(context.Background(), server.Client(), server.URL, filename, false); err != nil {
			t.Fatal(err)
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("cache made %d requests", calls.Load())
	}
	fail.Store(true)
	if err := ensureTool(context.Background(), server.Client(), server.URL, filename, true); err == nil {
		t.Fatal("accepted HTTP error")
	}
	data, err := os.ReadFile(filename)
	if err != nil || string(data) != "tool contents" {
		t.Fatalf("cached tool lost: %q, %v", data, err)
	}
	matches, err := filepath.Glob(filepath.Join(filepath.Dir(filename), ".tool-*"))
	if err != nil || len(matches) != 0 {
		t.Fatalf("partial download leaked: %v, %v", matches, err)
	}
}
