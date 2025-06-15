package assetserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkSnifferAllocationPatterns tests specific allocation scenarios
func BenchmarkSnifferAllocationPatterns(b *testing.B) {
	// Test 1: Object allocation only (no content)
	b.Run("ObjectAllocation-Baseline", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			rw := httptest.NewRecorder()
			sniffer := &contentTypeSniffer{
				rw:           rw,
				closeChannel: make(chan bool, 1),
			}
			_ = sniffer
		}
	})
	
	b.Run("ObjectAllocation-Pooled", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			rw := httptest.NewRecorder()
			sniffer := newContentTypeSniffer(rw)
			sniffer.returnToPool()
		}
	})
	
	// Test 2: Small content write (triggers buffer allocation)
	smallContent := []byte("Small test content")
	
	b.Run("SmallContent-Baseline", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			rw := httptest.NewRecorder()
			sniffer := &contentTypeSniffer{
				rw:           rw,
				closeChannel: make(chan bool, 1),
			}
			sniffer.Write(smallContent)
			sniffer.complete()
		}
	})
	
	b.Run("SmallContent-Pooled", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			rw := httptest.NewRecorder()
			sniffer := newContentTypeSniffer(rw)
			sniffer.Write(smallContent)
			sniffer.complete()
			sniffer.returnToPool()
		}
	})
	
	// Test 3: Exact 512 byte content (buffer size boundary)
	exactContent := make([]byte, 512)
	for i := range exactContent {
		exactContent[i] = byte(i % 256)
	}
	
	b.Run("ExactBuffer-Baseline", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			rw := httptest.NewRecorder()
			sniffer := &contentTypeSniffer{
				rw:           rw,
				closeChannel: make(chan bool, 1),
			}
			sniffer.Write(exactContent)
			sniffer.complete()
		}
	})
	
	b.Run("ExactBuffer-Pooled", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			rw := httptest.NewRecorder()
			sniffer := newContentTypeSniffer(rw)
			sniffer.Write(exactContent)
			sniffer.complete()
			sniffer.returnToPool()
		}
	})
}

// BenchmarkHighFrequencyServing simulates rapid-fire requests
func BenchmarkHighFrequencyServing(b *testing.B) {
	content := []byte("<!DOCTYPE html><html><body>Test Page</body></html>")
	
	b.Run("Baseline", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				rw := httptest.NewRecorder()
				sniffer := &contentTypeSniffer{
					rw:           rw,
					closeChannel: make(chan bool, 1),
				}
				sniffer.WriteHeader(http.StatusOK)
				sniffer.Write(content)
				sniffer.complete()
			}
		})
	})
	
	b.Run("Pooled", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				rw := httptest.NewRecorder()
				sniffer := newContentTypeSniffer(rw)
				sniffer.WriteHeader(http.StatusOK)
				sniffer.Write(content)
				sniffer.complete()
				sniffer.returnToPool()
			}
		})
	})
}

// TestPoolEffectiveness verifies the pool is actually being used
func TestPoolEffectiveness(t *testing.T) {
	// Pre-warm the pool
	var sniffers []*contentTypeSniffer
	for i := 0; i < 10; i++ {
		rw := httptest.NewRecorder()
		s := newContentTypeSniffer(rw)
		sniffers = append(sniffers, s)
	}
	
	// Return all to pool
	for _, s := range sniffers {
		s.returnToPool()
	}
	
	// Now get new sniffers and verify they're reused
	reusedCount := 0
	for i := 0; i < 10; i++ {
		rw := httptest.NewRecorder()
		s := newContentTypeSniffer(rw)
		
		// Check if this is a reused sniffer by seeing if it's one we saw before
		for _, original := range sniffers {
			if s == original {
				reusedCount++
				break
			}
		}
		s.returnToPool()
	}
	
	if reusedCount == 0 {
		t.Error("Pool doesn't appear to be reusing sniffers")
	}
	
	t.Logf("Reused %d/%d sniffers from pool", reusedCount, 10)
}