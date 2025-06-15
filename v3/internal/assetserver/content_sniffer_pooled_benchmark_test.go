package assetserver

import (
	"bytes"
	"net/http"
	"testing"
)

// Benchmark pooled content type sniffer allocation
func BenchmarkContentTypeSnifferPooled(b *testing.B) {
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
			sniffer.returnToPool()
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
				sniffer.returnToPool()
			}
		})
	})
}

// Benchmark pooled memory allocation patterns
func BenchmarkContentSnifferMemoryPressurePooled(b *testing.B) {
	testData := bytes.Repeat([]byte("Test content for memory pressure analysis. "), 20)
	
	b.Run("Pooled-MemPressure", func(b *testing.B) {
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
					s.returnToPool()
				}
				sniffers = sniffers[:0]
			}
		}
		// Complete remaining sniffers
		for _, s := range sniffers {
			s.complete()
			s.returnToPool()
		}
	})
}

// Benchmark pooled realistic HTTP request processing
func BenchmarkHTTPRequestWorkflowPooled(b *testing.B) {
	// Different content types to test varied scenarios
	testContents := [][]byte{
		[]byte("<!DOCTYPE html><html><head><title>Test</title></head><body>Hello World</body></html>"), // HTML
		[]byte("{\"message\": \"Hello World\", \"timestamp\": 1234567890}"),                           // JSON
		[]byte("body { background-color: #fff; color: #000; }"),                                      // CSS
		bytes.Repeat([]byte("console.log('Hello World'); "), 20),                                     // JS
	}
	
	b.Run("Pooled-HTTPWorkflow", func(b *testing.B) {
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
			sniffer.returnToPool()
		}
	})
}

// Test pooled buffer growth scenarios (when content is written in chunks)
func BenchmarkChunkedContentWritingPooled(b *testing.B) {
	// Create test content that will be written in chunks
	fullContent := bytes.Repeat([]byte("This is chunked content data. "), 30)
	chunkSize := 100
	
	b.Run("Pooled-ChunkedWrite", func(b *testing.B) {
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
			sniffer.returnToPool()
		}
	})
}

// Benchmark pooled high-frequency asset serving (simulating real application usage)
func BenchmarkAssetServingSimulationPooled(b *testing.B) {
	// Simulate typical asset files
	assets := [][]byte{
		bytes.Repeat([]byte("/* CSS content */ .class { color: red; } "), 50),         // CSS file
		bytes.Repeat([]byte("// JS content\nfunction test() { return true; } "), 50),  // JS file  
		bytes.Repeat([]byte("<html><body>HTML content here</body></html> "), 20),      // HTML file
		bytes.Repeat([]byte("Binary content simulation data "), 100),                  // Binary file
	}
	
	b.Run("Pooled-AssetServing", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			asset := assets[i%len(assets)]
			
			rw := newMockResponseWriter()
			sniffer := newContentTypeSniffer(rw)
			
			// Simulate asset serving workflow
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(asset)
			sniffer.complete()
			sniffer.returnToPool()
		}
	})
}

// Comparative benchmark for baseline vs pooled
func BenchmarkContentSnifferComparison(b *testing.B) {
	testData := bytes.Repeat([]byte("Test content for comparison. "), 20)
	
	b.Run("Baseline", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			rw := newMockResponseWriter()
			// Create sniffer without pooling (simulating old behavior)
			sniffer := &contentTypeSniffer{
				rw:           rw,
				closeChannel: make(chan bool, 1),
			}
			
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(testData)
			sniffer.complete()
		}
	})
	
	b.Run("Pooled", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			rw := newMockResponseWriter()
			sniffer := newContentTypeSniffer(rw)
			
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(testData)
			sniffer.complete()
			sniffer.returnToPool()
		}
	})
}

// Test correctness of pooled content type sniffer behavior
func TestPooledContentTypeSnifferCorrectness(t *testing.T) {
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
			content:     []byte("{\"key\": \"value\", \"array\": [1, 2, 3], \"nested\": {\"a\": true}}"),
			expectedCT:  "text/plain; charset=utf-8", // Go's DetectContentType doesn't reliably detect JSON
		},
		{
			name:        "Plain Text",
			content:     []byte("This is plain text content"),
			expectedCT:  "text/plain; charset=utf-8",
		},
	}
	
	// Run multiple times to test pool reuse
	for i := 0; i < 10; i++ {
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
				
				// Return to pool for reuse
				sniffer.returnToPool()
			})
		}
	}
}

// Test pooled buffer reuse scenarios
func TestPooledSnifferBufferBehavior(t *testing.T) {
	// Test that buffer allocation and pooling works correctly
	for i := 0; i < 5; i++ {
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
		
		// Complete and return to pool
		sniffer.complete()
		sniffer.returnToPool()
	}
}

// Benchmark to demonstrate memory reduction
func BenchmarkMemoryReductionDemonstration(b *testing.B) {
	// Various content sizes to test different scenarios
	contents := [][]byte{
		bytes.Repeat([]byte("A"), 100),  // Small content
		bytes.Repeat([]byte("B"), 300),  // Medium content
		bytes.Repeat([]byte("C"), 500),  // Near buffer size
		bytes.Repeat([]byte("D"), 700),  // Larger than buffer
	}
	
	b.Run("BaselineMemory", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			content := contents[i%len(contents)]
			rw := newMockResponseWriter()
			
			// Simulate old behavior without pooling
			sniffer := &contentTypeSniffer{
				rw:           rw,
				closeChannel: make(chan bool, 1),
			}
			
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(content)
			sniffer.complete()
		}
	})
	
	b.Run("PooledMemory", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			content := contents[i%len(contents)]
			rw := newMockResponseWriter()
			sniffer := newContentTypeSniffer(rw)
			
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(content)
			sniffer.complete()
			sniffer.returnToPool()
		}
	})
}