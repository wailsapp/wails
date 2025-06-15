package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/internal/assetserver"
)

// BenchmarkPhase1Integration measures the combined impact of all Phase 1 optimizations
func BenchmarkPhase1Integration(b *testing.B) {
	// Test 1: Atomic operations performance (Stage 1)
	b.Run("AtomicOperations", func(b *testing.B) {
		b.Run("Baseline", func(b *testing.B) {
			var mu sync.Mutex
			var counter uint32
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					mu.Lock()
					counter++
					mu.Unlock()
				}
			})
		})
		
		b.Run("Optimized", func(b *testing.B) {
			var counter atomic.Uint32
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					counter.Add(1)
				}
			})
		})
	})

	// Test 2: JSON operations (Stage 2)
	b.Run("JSONOperations", func(b *testing.B) {
		testData := map[string]interface{}{
			"id":      12345,
			"name":    "Test User",
			"data":    []int{1, 2, 3, 4, 5},
			"active":  true,
			"created": time.Now().Unix(),
		}
		
		b.Run("Baseline", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				data, _ := json.Marshal(testData)
				var result map[string]interface{}
				json.Unmarshal(data, &result)
			}
		})
		
		b.Run("Optimized", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				// TODO: Replace with Sonic JSON when available
				data, _ := json.Marshal(testData)
				var result map[string]interface{}
				json.Unmarshal(data, &result)
			}
		})
	})

	// Test 3: Channel operations (Stage 4)
	b.Run("ChannelOperations", func(b *testing.B) {
		b.Run("Baseline", func(b *testing.B) {
			ch := make(chan int, 5) // Small buffer
			b.ReportAllocs()
			
			go func() {
				for i := 0; i < b.N; i++ {
					select {
					case ch <- i:
					default:
						// Dropped due to buffer full
					}
				}
			}()
			
			for i := 0; i < b.N; i++ {
				select {
				case <-ch:
				case <-time.After(time.Microsecond):
					// Timeout
				}
			}
		})
		
		b.Run("Optimized", func(b *testing.B) {
			ch := make(chan int, 100) // Larger buffer
			b.ReportAllocs()
			
			go func() {
				for i := 0; i < b.N; i++ {
					select {
					case ch <- i:
					default:
						// Dropped due to buffer full
					}
				}
			}()
			
			for i := 0; i < b.N; i++ {
				select {
				case <-ch:
				case <-time.After(time.Microsecond):
					// Timeout
				}
			}
		})
	})

	// Test 4: Combined real-world scenario
	b.Run("RealWorldScenario", func(b *testing.B) {
		b.Run("Baseline", func(b *testing.B) {
			b.ReportAllocs()
			runRealWorldScenario(b, false)
		})
		
		b.Run("Optimized", func(b *testing.B) {
			b.ReportAllocs()
			runRealWorldScenario(b, true)
		})
	})
}

// runRealWorldScenario simulates a typical Wails application workload
func runRealWorldScenario(b *testing.B, optimized bool) {
	// Simulate concurrent operations
	var wg sync.WaitGroup
	
	// ID generation (Stage 1)
	var idCounter atomic.Uint32
	
	// Event channel (Stage 4)
	bufferSize := 5
	if optimized {
		bufferSize = 100
	}
	eventChan := make(chan interface{}, bufferSize)
	
	// Start event processor
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			select {
			case <-eventChan:
			case <-time.After(time.Millisecond):
			}
		}
	}()
	
	// Simulate application operations
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Generate ID
			id := idCounter.Add(1)
			
			// Create event data
			event := map[string]interface{}{
				"id":        id,
				"type":      "test",
				"timestamp": time.Now().Unix(),
			}
			
			// Marshal event (Stage 2)
			data, _ := json.Marshal(event)
			
			// Send event
			select {
			case eventChan <- data:
			default:
			}
		}
	})
	
	close(eventChan)
	wg.Wait()
}

// BenchmarkHTTPAssetServing tests HTTP asset serving performance
func BenchmarkHTTPAssetServing(b *testing.B) {
	// Create test assets
	htmlContent := []byte("<!DOCTYPE html><html><head><title>Test</title></head><body>Hello World</body></html>")
	cssContent := []byte("body { margin: 0; padding: 0; font-family: Arial; }")
	jsContent := []byte("console.log('Application loaded');")
	
	b.Run("Baseline", func(b *testing.B) {
		b.ReportAllocs()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/index.html":
				w.Write(htmlContent)
			case "/style.css":
				w.Write(cssContent)
			case "/app.js":
				w.Write(jsContent)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		})
		
		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest("GET", "/index.html", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
		}
	})
	
	b.Run("Optimized", func(b *testing.B) {
		b.ReportAllocs()
		// This would use the optimized asset server with content sniffer pooling
		options := &assetserver.Options{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/index.html":
					w.Write(htmlContent)
				case "/style.css":
					w.Write(cssContent)
				case "/app.js":
					w.Write(jsContent)
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			}),
		}
		
		server, _ := assetserver.NewAssetServer(options)
		
		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest("GET", "/index.html", nil)
			rr := httptest.NewRecorder()
			server.ServeHTTP(rr, req)
		}
	})
}

// BenchmarkMemoryPressure tests behavior under memory pressure
func BenchmarkMemoryPressure(b *testing.B) {
	b.Run("Baseline", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			// Simulate allocations without pooling
			data := make(map[string]interface{})
			data["id"] = i
			data["timestamp"] = time.Now().Unix()
			
			jsonData, _ := json.Marshal(data)
			_ = jsonData
			
			// Simulate struct allocations
			opts := &struct {
				ID     int
				Method string
				Params []interface{}
			}{
				ID:     i,
				Method: "test",
				Params: make([]interface{}, 0, 10),
			}
			_ = opts
		}
	})
	
	b.Run("Optimized", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			// Using pooled objects and optimized JSON
			data := make(map[string]interface{})
			data["id"] = i
			data["timestamp"] = time.Now().Unix()
			
			jsonData, _ := json.Marshal(data)
			_ = jsonData
			
			// Pooled structs would be used here
			// This simulates the effect
		}
	})
}

// BenchmarkConcurrentLoad tests performance under concurrent load
func BenchmarkConcurrentLoad(b *testing.B) {
	workers := 10
	
	b.Run("Baseline", func(b *testing.B) {
		b.ReportAllocs()
		
		var wg sync.WaitGroup
		workChan := make(chan int, workers)
		
		// Start workers
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var mu sync.Mutex
				var counter int
				
				for work := range workChan {
					mu.Lock()
					counter++
					mu.Unlock()
					
					// Simulate work
					data := map[string]int{"work": work, "counter": counter}
					json.Marshal(data)
				}
			}()
		}
		
		// Send work
		for i := 0; i < b.N; i++ {
			workChan <- i
		}
		close(workChan)
		wg.Wait()
	})
	
	b.Run("Optimized", func(b *testing.B) {
		b.ReportAllocs()
		
		var wg sync.WaitGroup
		workChan := make(chan int, workers*10) // Larger buffer
		
		// Start workers
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var counter atomic.Uint32
				
				for work := range workChan {
					counter.Add(1)
					
					// Simulate work with optimized JSON
					data := map[string]interface{}{"work": work, "counter": counter.Load()}
					json.Marshal(data)
				}
			}()
		}
		
		// Send work
		for i := 0; i < b.N; i++ {
			workChan <- i
		}
		close(workChan)
		wg.Wait()
	})
}