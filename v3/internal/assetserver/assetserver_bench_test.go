//go:build bench

package assetserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

// resetMimeCache clears the mime cache for benchmark isolation
func resetMimeCache() {
	mimeCache = sync.Map{}
}

// BenchmarkGetMimetype measures MIME type detection performance
func BenchmarkGetMimetype(b *testing.B) {
	// Reset cache between runs
	resetMimeCache()

	b.Run("ByExtension/JS", func(b *testing.B) {
		data := []byte("function test() {}")
		for b.Loop() {
			_ = GetMimetype("script.js", data)
		}
	})

	resetMimeCache()
	b.Run("ByExtension/CSS", func(b *testing.B) {
		data := []byte(".class { color: red; }")
		for b.Loop() {
			_ = GetMimetype("style.css", data)
		}
	})

	resetMimeCache()
	b.Run("ByExtension/HTML", func(b *testing.B) {
		data := []byte("<!DOCTYPE html><html></html>")
		for b.Loop() {
			_ = GetMimetype("index.html", data)
		}
	})

	resetMimeCache()
	b.Run("ByExtension/JSON", func(b *testing.B) {
		data := []byte(`{"key": "value"}`)
		for b.Loop() {
			_ = GetMimetype("data.json", data)
		}
	})

	resetMimeCache()
	b.Run("Detection/Unknown", func(b *testing.B) {
		data := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}
		for b.Loop() {
			_ = GetMimetype("unknown.bin", data)
		}
	})

	resetMimeCache()
	b.Run("Detection/PNG", func(b *testing.B) {
		// PNG magic bytes
		data := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00}
		for b.Loop() {
			_ = GetMimetype("image.unknown", data)
		}
	})

	resetMimeCache()
	b.Run("CacheHit", func(b *testing.B) {
		data := []byte{0x00, 0x01, 0x02}
		// Prime the cache
		_ = GetMimetype("cached.bin", data)
		b.ResetTimer()
		for b.Loop() {
			_ = GetMimetype("cached.bin", data)
		}
	})
}

// BenchmarkGetMimetype_Concurrent tests concurrent MIME type lookups
func BenchmarkGetMimetype_Concurrent(b *testing.B) {
	resetMimeCache()
	data := []byte("function test() {}")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = GetMimetype("script.js", data)
		}
	})
}

// BenchmarkAssetServerServeHTTP measures request handling overhead
func BenchmarkAssetServerServeHTTP(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError + 1}))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<!DOCTYPE html><html><body>Hello</body></html>"))
	})

	server, err := NewAssetServer(&Options{
		Handler: handler,
		Logger:  logger,
	})
	if err != nil {
		b.Fatal(err)
	}

	b.Run("SimpleRequest", func(b *testing.B) {
		req := httptest.NewRequest("GET", "/index.html", nil)
		for b.Loop() {
			rr := httptest.NewRecorder()
			server.ServeHTTP(rr, req)
		}
	})

	b.Run("WithHeaders", func(b *testing.B) {
		req := httptest.NewRequest("GET", "/index.html", nil)
		req.Header.Set("x-wails-window-id", "1")
		req.Header.Set("x-wails-window-name", "main")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		for b.Loop() {
			rr := httptest.NewRecorder()
			server.ServeHTTP(rr, req)
		}
	})
}

// BenchmarkAssetServerServeHTTP_Concurrent tests concurrent request handling
func BenchmarkAssetServerServeHTTP_Concurrent(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError + 1}))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<!DOCTYPE html><html><body>Hello</body></html>"))
	})

	server, err := NewAssetServer(&Options{
		Handler: handler,
		Logger:  logger,
	})
	if err != nil {
		b.Fatal(err)
	}

	b.RunParallel(func(pb *testing.PB) {
		req := httptest.NewRequest("GET", "/index.html", nil)
		for pb.Next() {
			rr := httptest.NewRecorder()
			server.ServeHTTP(rr, req)
		}
	})
}

// BenchmarkContentTypeSniffer measures the content type sniffer overhead
func BenchmarkContentTypeSniffer(b *testing.B) {
	b.Run("SmallResponse", func(b *testing.B) {
		data := []byte("Hello, World!")
		for b.Loop() {
			rr := httptest.NewRecorder()
			sniffer := newContentTypeSniffer(rr)
			_, _ = sniffer.Write(data)
			_, _ = sniffer.complete()
		}
	})

	b.Run("HTMLResponse", func(b *testing.B) {
		data := []byte("<!DOCTYPE html><html><head><title>Test</title></head><body><h1>Hello</h1></body></html>")
		for b.Loop() {
			rr := httptest.NewRecorder()
			sniffer := newContentTypeSniffer(rr)
			_, _ = sniffer.Write(data)
			_, _ = sniffer.complete()
		}
	})

	b.Run("LargeResponse", func(b *testing.B) {
		data := make([]byte, 64*1024) // 64KB
		for i := range data {
			data[i] = byte(i % 256)
		}
		for b.Loop() {
			rr := httptest.NewRecorder()
			sniffer := newContentTypeSniffer(rr)
			_, _ = sniffer.Write(data)
			_, _ = sniffer.complete()
		}
	})
}

// BenchmarkServiceRouting measures service route matching performance
func BenchmarkServiceRouting(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError + 1}))

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server, err := NewAssetServer(&Options{
		Handler: dummyHandler,
		Logger:  logger,
	})
	if err != nil {
		b.Fatal(err)
	}

	// Attach multiple service routes
	for i := 0; i < 10; i++ {
		server.AttachServiceHandler(fmt.Sprintf("/api/v%d/", i), dummyHandler)
	}

	b.Run("FirstRoute", func(b *testing.B) {
		req := httptest.NewRequest("GET", "/api/v0/users", nil)
		for b.Loop() {
			rr := httptest.NewRecorder()
			server.ServeHTTP(rr, req)
		}
	})

	b.Run("LastRoute", func(b *testing.B) {
		req := httptest.NewRequest("GET", "/api/v9/users", nil)
		for b.Loop() {
			rr := httptest.NewRecorder()
			server.ServeHTTP(rr, req)
		}
	})

	b.Run("NoMatch", func(b *testing.B) {
		req := httptest.NewRequest("GET", "/static/app.js", nil)
		for b.Loop() {
			rr := httptest.NewRecorder()
			server.ServeHTTP(rr, req)
		}
	})
}
