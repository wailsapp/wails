package assetserver

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock response writer for testing
type mockResponseWriter struct {
	*httptest.ResponseRecorder
	headers http.Header
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		ResponseRecorder: httptest.NewRecorder(),
		headers:          make(http.Header),
	}
}

func (m *mockResponseWriter) Header() http.Header {
	return m.headers
}

// Benchmark baseline content type sniffer allocation
func BenchmarkContentTypeSnifferBaseline(b *testing.B) {
	testData := bytes.Repeat([]byte("Hello World! This is test content for MIME type detection. "), 10)
	
	b.Run("SingleSniffer", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			rw := newMockResponseWriter()
			sniffer := newContentTypeSniffer(rw)
			
			// Simulate typical content sniffing workflow
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(testData)
			sniffer.complete()
		}
	})

	b.Run("ContentionSniffer", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				rw := newMockResponseWriter()
				sniffer := newContentTypeSniffer(rw)
				
				// Simulate typical content sniffing workflow
				sniffer.WriteHeader(http.StatusOK)
				sniffer.Write(testData)
				sniffer.complete()
			}
		})
	})
}

// Benchmark memory allocation patterns
func BenchmarkContentSnifferMemoryPressure(b *testing.B) {
	testData := bytes.Repeat([]byte("Test content for memory pressure analysis. "), 20)
	
	b.Run("Baseline-MemPressure", func(b *testing.B) {
		b.ReportAllocs()
		var sniffers []*contentTypeSniffer
		for i := 0; i < b.N; i++ {
			rw := newMockResponseWriter()
			sniffer := newContentTypeSniffer(rw)
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(testData)
			sniffers = append(sniffers, sniffer)
			
			// Periodically clear to simulate GC pressure
			if len(sniffers) > 100 {
				for _, s := range sniffers {
					s.complete()
				}
				sniffers = sniffers[:0]
			}
		}
		// Complete remaining sniffers
		for _, s := range sniffers {
			s.complete()
		}
	})
}

// Benchmark realistic HTTP request processing
func BenchmarkHTTPRequestWorkflow(b *testing.B) {
	// Different content types to test varied scenarios
	testContents := [][]byte{
		[]byte("<!DOCTYPE html><html><head><title>Test</title></head><body>Hello World</body></html>"), // HTML
		[]byte("{\"message\": \"Hello World\", \"timestamp\": 1234567890}"),                           // JSON
		[]byte("body { background-color: #fff; color: #000; }"),                                      // CSS
		bytes.Repeat([]byte("console.log('Hello World'); "), 20),                                     // JS
	}
	
	b.Run("Baseline-HTTPWorkflow", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			content := testContents[i%len(testContents)]
			
			// Simulate complete HTTP request processing
			rw := newMockResponseWriter()
			sniffer := newContentTypeSniffer(rw)
			
			// Typical workflow: WriteHeader -> Write content -> Complete
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(content)
			sniffer.complete()
		}
	})
}

// Test buffer growth scenarios (when content is written in chunks)
func BenchmarkChunkedContentWriting(b *testing.B) {
	// Create test content that will be written in chunks
	fullContent := bytes.Repeat([]byte("This is chunked content data. "), 30)
	chunkSize := 100
	
	b.Run("Baseline-ChunkedWrite", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			rw := newMockResponseWriter()
			sniffer := newContentTypeSniffer(rw)
			sniffer.WriteHeader(http.StatusOK)
			
			// Write content in chunks to trigger buffer growth
			for offset := 0; offset < len(fullContent); offset += chunkSize {
				end := offset + chunkSize
				if end > len(fullContent) {
					end = len(fullContent)
				}
				sniffer.Write(fullContent[offset:end])
			}
			sniffer.complete()
		}
	})
}

// Benchmark high-frequency asset serving (simulating real application usage)
func BenchmarkAssetServingSimulation(b *testing.B) {
	// Simulate typical asset files
	assets := [][]byte{
		bytes.Repeat([]byte("/* CSS content */ .class { color: red; } "), 50),         // CSS file
		bytes.Repeat([]byte("// JS content\nfunction test() { return true; } "), 50),  // JS file  
		bytes.Repeat([]byte("<html><body>HTML content here</body></html> "), 20),      // HTML file
		bytes.Repeat([]byte("Binary content simulation data "), 100),                  // Binary file
	}
	
	b.Run("Baseline-AssetServing", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			asset := assets[i%len(assets)]
			
			rw := newMockResponseWriter()
			sniffer := newContentTypeSniffer(rw)
			
			// Simulate asset serving workflow
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(asset)
			sniffer.complete()
		}
	})
}

// Test correctness of content type sniffer behavior
func TestContentTypeSnifferCorrectness(t *testing.T) {
	testCases := []struct {
		name        string
		content     []byte
		expectedCT  string
	}{
		{
			name:        "HTML Content",
			content:     []byte("<!DOCTYPE html><html><head><title>Test</title></head>"),
			expectedCT:  "text/html; charset=utf-8",
		},
		{
			name:        "JSON Content", 
			content:     []byte("{\"key\": \"value\"}"),
			expectedCT:  "application/json",
		},
		{
			name:        "Plain Text",
			content:     []byte("This is plain text content"),
			expectedCT:  "text/plain; charset=utf-8",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rw := newMockResponseWriter()
			sniffer := newContentTypeSniffer(rw)
			
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(tc.content)
			sniffer.complete()
			
			contentType := rw.headers.Get("Content-Type")
			if contentType != tc.expectedCT {
				t.Errorf("Expected Content-Type %s, got %s", tc.expectedCT, contentType)
			}
		})
	}
}

// Test buffer reuse scenarios
func TestSnifferBufferBehavior(t *testing.T) {
	// Test that buffer allocation happens as expected
	rw := newMockResponseWriter()
	sniffer := newContentTypeSniffer(rw)
	
	// Initially prefix should be nil
	if sniffer.prefix != nil {
		t.Error("Expected initial prefix to be nil")
	}
	
	// Write small content (less than 512 bytes)
	smallContent := []byte("Small content for testing")
	sniffer.Write(smallContent)
	
	// After writing, prefix should be allocated with capacity 512
	if sniffer.prefix == nil {
		t.Error("Expected prefix to be allocated after first write")
	}
	
	if cap(sniffer.prefix) != 512 {
		t.Errorf("Expected prefix capacity to be 512, got %d", cap(sniffer.prefix))
	}
	
	if len(sniffer.prefix) != len(smallContent) {
		t.Errorf("Expected prefix length to be %d, got %d", len(smallContent), len(sniffer.prefix))
	}
}