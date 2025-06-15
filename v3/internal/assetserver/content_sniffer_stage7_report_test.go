package assetserver

import (
	"bytes"
	"net/http"
	"testing"
)

// BenchmarkStage7Performance demonstrates Stage 7 improvements
func BenchmarkStage7Performance(b *testing.B) {
	// Test with various content sizes that represent typical web assets
	smallJS := bytes.Repeat([]byte("console.log('test'); "), 10)         // ~220 bytes
	mediumCSS := bytes.Repeat([]byte(".class { color: red; } "), 20)    // ~480 bytes
	largeHTML := bytes.Repeat([]byte("<div>Content here</div> "), 50)   // ~1200 bytes
	
	b.Run("SmallAsset-Baseline", func(b *testing.B) {
		b.ReportAllocs()
		rw := newMinimalResponseWriter()
		
		for i := 0; i < b.N; i++ {
			// Simulate old behavior
			sniffer := &contentTypeSniffer{
				rw:           rw,
				closeChannel: make(chan bool, 1),
			}
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(smallJS)
			sniffer.complete()
			rw.buffer.Reset()
		}
	})
	
	b.Run("SmallAsset-Pooled", func(b *testing.B) {
		b.ReportAllocs()
		rw := newMinimalResponseWriter()
		
		for i := 0; i < b.N; i++ {
			sniffer := newContentTypeSniffer(rw)
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(smallJS)
			sniffer.complete()
			sniffer.returnToPool()
			rw.buffer.Reset()
		}
	})
	
	b.Run("MediumAsset-Baseline", func(b *testing.B) {
		b.ReportAllocs()
		rw := newMinimalResponseWriter()
		
		for i := 0; i < b.N; i++ {
			sniffer := &contentTypeSniffer{
				rw:           rw,
				closeChannel: make(chan bool, 1),
			}
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(mediumCSS)
			sniffer.complete()
			rw.buffer.Reset()
		}
	})
	
	b.Run("MediumAsset-Pooled", func(b *testing.B) {
		b.ReportAllocs()
		rw := newMinimalResponseWriter()
		
		for i := 0; i < b.N; i++ {
			sniffer := newContentTypeSniffer(rw)
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(mediumCSS)
			sniffer.complete()
			sniffer.returnToPool()
			rw.buffer.Reset()
		}
	})
	
	b.Run("LargeAsset-Baseline", func(b *testing.B) {
		b.ReportAllocs()
		rw := newMinimalResponseWriter()
		
		for i := 0; i < b.N; i++ {
			sniffer := &contentTypeSniffer{
				rw:           rw,
				closeChannel: make(chan bool, 1),
			}
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(largeHTML)
			sniffer.complete()
			rw.buffer.Reset()
		}
	})
	
	b.Run("LargeAsset-Pooled", func(b *testing.B) {
		b.ReportAllocs()
		rw := newMinimalResponseWriter()
		
		for i := 0; i < b.N; i++ {
			sniffer := newContentTypeSniffer(rw)
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(largeHTML)
			sniffer.complete()
			sniffer.returnToPool()
			rw.buffer.Reset()
		}
	})
}

// BenchmarkHighConcurrencyServing tests performance under high concurrency
func BenchmarkHighConcurrencyServing(b *testing.B) {
	content := bytes.Repeat([]byte("/* CSS */ body { margin: 0; } "), 20)
	
	b.Run("Concurrency-Baseline", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			rw := newMinimalResponseWriter()
			for pb.Next() {
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
	})
	
	b.Run("Concurrency-Pooled", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			rw := newMinimalResponseWriter()
			for pb.Next() {
				sniffer := newContentTypeSniffer(rw)
				sniffer.WriteHeader(http.StatusOK)
				sniffer.Write(content)
				sniffer.complete()
				sniffer.returnToPool()
				rw.buffer.Reset()
			}
		})
	})
}

// BenchmarkRealWorldMixedContent simulates mixed content types
func BenchmarkRealWorldMixedContent(b *testing.B) {
	// Simulate typical web app asset distribution
	assets := []struct {
		content []byte
		name    string
	}{
		{bytes.Repeat([]byte("// JavaScript\nfunction app() {} "), 15), "app.js"},
		{bytes.Repeat([]byte("/* CSS */\n.btn { color: blue; } "), 10), "style.css"},
		{bytes.Repeat([]byte("<!DOCTYPE html><html> "), 20), "index.html"},
		{[]byte("{\"api\": \"v1\", \"status\": \"ok\"}"), "api.json"},
		{bytes.Repeat([]byte("PNG_DATA_HERE "), 100), "image.png"},
	}
	
	b.Run("MixedContent-Baseline", func(b *testing.B) {
		b.ReportAllocs()
		rw := newMinimalResponseWriter()
		
		for i := 0; i < b.N; i++ {
			asset := assets[i%len(assets)]
			sniffer := &contentTypeSniffer{
				rw:           rw,
				closeChannel: make(chan bool, 1),
			}
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(asset.content)
			sniffer.complete()
			rw.buffer.Reset()
		}
	})
	
	b.Run("MixedContent-Pooled", func(b *testing.B) {
		b.ReportAllocs()
		rw := newMinimalResponseWriter()
		
		for i := 0; i < b.N; i++ {
			asset := assets[i%len(assets)]
			sniffer := newContentTypeSniffer(rw)
			sniffer.WriteHeader(http.StatusOK)
			sniffer.Write(asset.content)
			sniffer.complete()
			sniffer.returnToPool()
			rw.buffer.Reset()
		}
	})
}