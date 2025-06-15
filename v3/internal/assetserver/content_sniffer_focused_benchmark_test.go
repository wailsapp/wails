package assetserver

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

// minimalResponseWriter is a minimal implementation for benchmarking
type minimalResponseWriter struct {
	headers http.Header
	buffer  bytes.Buffer
	status  int
}

func newMinimalResponseWriter() *minimalResponseWriter {
	return &minimalResponseWriter{
		headers: make(http.Header),
	}
}

func (m *minimalResponseWriter) Header() http.Header {
	return m.headers
}

func (m *minimalResponseWriter) Write(b []byte) (int, error) {
	return m.buffer.Write(b)
}

func (m *minimalResponseWriter) WriteHeader(statusCode int) {
	m.status = statusCode
}

// BenchmarkContentSnifferOnly focuses purely on sniffer allocations
func BenchmarkContentSnifferOnly(b *testing.B) {
	content := []byte("<!DOCTYPE html><html><head><title>Test</title></head><body>Hello World</body></html>")
	
	b.Run("Baseline-SnifferOnly", func(b *testing.B) {
		b.ReportAllocs()
		rw := newMinimalResponseWriter() // Create once, reuse
		
		for i := 0; i < b.N; i++ {
			// Only measure sniffer allocation
			sniffer := &contentTypeSniffer{
				rw:           rw,
				closeChannel: make(chan bool, 1),
			}
			sniffer.Write(content)
			sniffer.complete()
			
			// Reset the response writer for next iteration
			rw.buffer.Reset()
			rw.status = 0
		}
	})
	
	b.Run("Pooled-SnifferOnly", func(b *testing.B) {
		b.ReportAllocs()
		rw := newMinimalResponseWriter() // Create once, reuse
		
		for i := 0; i < b.N; i++ {
			// Only measure pooled sniffer
			sniffer := newContentTypeSniffer(rw)
			sniffer.Write(content)
			sniffer.complete()
			sniffer.returnToPool()
			
			// Reset the response writer for next iteration
			rw.buffer.Reset()
			rw.status = 0
		}
	})
}

// BenchmarkBufferAllocationOnly focuses on buffer allocations
func BenchmarkBufferAllocationOnly(b *testing.B) {
	b.Run("Baseline-BufferAlloc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			// Simulate old buffer allocation
			buf := make([]byte, 0, 512)
			buf = append(buf, []byte("test content")...)
			_ = buf
		}
	})
	
	b.Run("Pooled-BufferAlloc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			// Use pooled buffer
			buf := GetContentSnifferBuffer()
			buf = append(buf, []byte("test content")...)
			PutContentSnifferBuffer(buf)
		}
	})
}

// BenchmarkRealWorldHTTPHandler simulates actual HTTP handler usage
func BenchmarkRealWorldHTTPHandler(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "test")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "<!DOCTYPE html><html><body>Hello World</body></html>")
	})
	
	b.Run("Baseline-HTTPHandler", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			rw := newMinimalResponseWriter()
			
			// Old style wrapper
			wrapped := &contentTypeSniffer{
				rw:           rw,
				closeChannel: make(chan bool, 1),
			}
			
			req, _ := http.NewRequest("GET", "/test", nil)
			handler.ServeHTTP(wrapped, req)
			wrapped.complete()
		}
	})
	
	b.Run("Pooled-HTTPHandler", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			rw := newMinimalResponseWriter()
			
			// Pooled wrapper
			wrapped := newContentTypeSniffer(rw)
			
			req, _ := http.NewRequest("GET", "/test", nil)
			handler.ServeHTTP(wrapped, req)
			wrapped.complete()
			wrapped.returnToPool()
		}
	})
}

// BenchmarkContentSnifferMemoryImpact measures total memory impact
func BenchmarkContentSnifferMemoryImpact(b *testing.B) {
	// Various content sizes
	contents := [][]byte{
		bytes.Repeat([]byte("A"), 50),   // Small
		bytes.Repeat([]byte("B"), 200),  // Medium
		bytes.Repeat([]byte("C"), 512),  // Exact buffer
		bytes.Repeat([]byte("D"), 1024), // Large
	}
	
	b.Run("Baseline-MemoryImpact", func(b *testing.B) {
		b.ReportAllocs()
		rw := newMinimalResponseWriter()
		
		for i := 0; i < b.N; i++ {
			content := contents[i%len(contents)]
			
			sniffer := &contentTypeSniffer{
				rw:           rw,
				closeChannel: make(chan bool, 1),
			}
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(content)
			sniffer.complete()
			
			rw.buffer.Reset()
		}
	})
	
	b.Run("Pooled-MemoryImpact", func(b *testing.B) {
		b.ReportAllocs()
		rw := newMinimalResponseWriter()
		
		for i := 0; i < b.N; i++ {
			content := contents[i%len(contents)]
			
			sniffer := newContentTypeSniffer(rw)
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(content)
			sniffer.complete()
			sniffer.returnToPool()
			
			rw.buffer.Reset()
		}
	})
}