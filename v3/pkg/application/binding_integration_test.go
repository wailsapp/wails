package application

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Integration test that tests the complete HTTP-only binding flow
func TestHTTPOnlyBindingIntegration(t *testing.T) {
	// Create a test application with bindings
	app := &App{
		options: Options{
			Assets: AssetOptions{
				Bindings: BindingConfig{
					Timeout: 10 * time.Second,
					CORS: CORSConfig{
						Enabled:        true,
						AllowedOrigins: []string{"https://example.com"},
					},
				},
			},
		},
		isDebugMode: false,
	}
	
	// Set up global application
	globalApplication = app
	
	// Create bindings
	bindings := &bindings{
		byName: make(map[string]*BoundMethod),
		byID:   make(map[uint32]*BoundMethod),
	}
	app.bindings = bindings
	
	// Add test methods
	bindings.byName["echo"] = &BoundMethod{
		Name: "echo",
		call: &mockBoundMethod{
			name: "echo",
			result: func(args []json.RawMessage) interface{} {
				if len(args) > 0 {
					var str string
					json.Unmarshal(args[0], &str)
					return map[string]string{"echo": str}
				}
				return map[string]string{"echo": "empty"}
			},
		},
	}
	
	bindings.byName["add"] = &BoundMethod{
		Name: "add",
		call: &mockBoundMethod{
			name: "add",
			result: func(args []json.RawMessage) interface{} {
				if len(args) >= 2 {
					var a, b float64
					json.Unmarshal(args[0], &a)
					json.Unmarshal(args[1], &b)
					return map[string]float64{"result": a + b}
				}
				return map[string]string{"error": "need two numbers"}
			},
		},
	}
	
	bindings.byName["slowMethod"] = &BoundMethod{
		Name: "slowMethod",
		call: &mockBoundMethod{
			name:  "slowMethod",
			delay: 2 * time.Second,
			result: "completed after delay",
		},
	}
	
	// Create HTTP server with MessageProcessor
	processor := NewMessageProcessor(nil)
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate the middleware that sets window headers
		if r.Header.Get("x-wails-window-id") == "" {
			r.Header.Set("x-wails-window-id", "1")
		}
		processor.ServeHTTP(w, r)
	})
	
	server := httptest.NewServer(handler)
	defer server.Close()
	
	window := &mockWindow{id: 1}
	
	t.Run("Complete binding call flow", func(t *testing.T) {
		// Simulate frontend calling binding
		callOptions := map[string]interface{}{
			"call-id":    "test-123",
			"methodName": "echo",
			"args":       []interface{}{"hello world"},
		}
		
		body, _ := json.Marshal(callOptions)
		
		req, err := http.NewRequest("POST", server.URL+"?object=0&method=0&args="+string(body), nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("x-wails-window-id", "1")
		req.Header.Set("x-wails-client-id", "test-client")
		
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
		}
		
		var result map[string]string
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatal(err)
		}
		
		if result["echo"] != "hello world" {
			t.Errorf("Expected echo 'hello world', got %s", result["echo"])
		}
	})
	
	t.Run("CORS headers in response", func(t *testing.T) {
		callOptions := map[string]interface{}{
			"call-id":    "test-cors",
			"methodName": "echo",
			"args":       []interface{}{"cors test"},
		}
		
		body, _ := json.Marshal(callOptions)
		
		req, err := http.NewRequest("POST", server.URL+"?object=0&method=0&args="+string(body), nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("x-wails-window-id", "1")
		req.Header.Set("Origin", "https://example.com")
		
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
		
		// Check CORS headers
		if resp.Header.Get("Access-Control-Allow-Origin") != "https://example.com" {
			t.Errorf("Expected CORS origin header, got %s", resp.Header.Get("Access-Control-Allow-Origin"))
		}
		
		if resp.Header.Get("Access-Control-Allow-Methods") == "" {
			t.Error("Expected CORS methods header to be set")
		}
	})
	
	t.Run("Method with parameters", func(t *testing.T) {
		callOptions := map[string]interface{}{
			"call-id":    "test-add",
			"methodName": "add",
			"args":       []interface{}{5.5, 4.5},
		}
		
		body, _ := json.Marshal(callOptions)
		
		req, err := http.NewRequest("POST", server.URL+"?object=0&method=0&args="+string(body), nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("x-wails-window-id", "1")
		
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
		}
		
		var result map[string]float64
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatal(err)
		}
		
		if result["result"] != 10.0 {
			t.Errorf("Expected result 10.0, got %f", result["result"])
		}
	})
	
	t.Run("Timeout handling", func(t *testing.T) {
		// Set a short timeout for this test
		app.options.Assets.Bindings.Timeout = 500 * time.Millisecond
		
		callOptions := map[string]interface{}{
			"call-id":    "test-timeout",
			"methodName": "slowMethod",
			"args":       []interface{}{},
		}
		
		body, _ := json.Marshal(callOptions)
		
		req, err := http.NewRequest("POST", server.URL+"?object=0&method=0&args="+string(body), nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("x-wails-window-id", "1")
		
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusRequestTimeout {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 408, got %d: %s", resp.StatusCode, string(body))
		}
		
		var errorResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		if err != nil {
			t.Fatal(err)
		}
		
		if errorResp["kind"] != "TimeoutError" {
			t.Errorf("Expected error kind 'TimeoutError', got %s", errorResp["kind"])
		}
		
		// Reset timeout
		app.options.Assets.Bindings.Timeout = 10 * time.Second
	})
	
	t.Run("Method not found error", func(t *testing.T) {
		callOptions := map[string]interface{}{
			"call-id":    "test-404",
			"methodName": "nonexistentMethod",
			"args":       []interface{}{},
		}
		
		body, _ := json.Marshal(callOptions)
		
		req, err := http.NewRequest("POST", server.URL+"?object=0&method=0&args="+string(body), nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("x-wails-window-id", "1")
		
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusNotFound {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 404, got %d: %s", resp.StatusCode, string(body))
		}
		
		var errorResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		if err != nil {
			t.Fatal(err)
		}
		
		if errorResp["kind"] != "ReferenceError" {
			t.Errorf("Expected error kind 'ReferenceError', got %s", errorResp["kind"])
		}
	})
	
	t.Run("OPTIONS preflight request", func(t *testing.T) {
		req, err := http.NewRequest("OPTIONS", server.URL+"?object=0&method=0", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 for OPTIONS, got %d", resp.StatusCode)
		}
		
		if resp.Header.Get("Access-Control-Allow-Origin") != "https://example.com" {
			t.Error("Expected CORS headers for OPTIONS request")
		}
	})
}

// Performance integration test
func TestHTTPOnlyBindingPerformance(t *testing.T) {
	app := &App{
		options: Options{
			Assets: AssetOptions{
				Bindings: BindingConfig{
					Timeout: 30 * time.Second,
				},
			},
		},
	}
	
	globalApplication = app
	
	// Create large response data
	largeData := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		largeData[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("large_value_%d_with_some_extra_content", i)
	}
	
	bindings := &bindings{
		byName: map[string]*BoundMethod{
			"largeResponse": &BoundMethod{
				Name: "largeResponse",
				call: &mockBoundMethod{
					name:   "largeResponse",
					result: largeData,
				},
			},
		},
	}
	app.bindings = bindings
	
	processor := NewMessageProcessor(nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-wails-window-id") == "" {
			r.Header.Set("x-wails-window-id", "1")
		}
		processor.ServeHTTP(w, r)
	})
	
	server := httptest.NewServer(handler)
	defer server.Close()
	
	t.Run("Large response performance", func(t *testing.T) {
		callOptions := map[string]interface{}{
			"call-id":    "perf-test",
			"methodName": "largeResponse",
			"args":       []interface{}{},
		}
		
		body, _ := json.Marshal(callOptions)
		
		start := time.Now()
		
		req, err := http.NewRequest("POST", server.URL+"?object=0&method=0&args="+string(body), nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("x-wails-window-id", "1")
		
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(bodyBytes))
		}
		
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatal(err)
		}
		
		duration := time.Since(start)
		
		if len(result) != 1000 {
			t.Errorf("Expected 1000 items in response, got %d", len(result))
		}
		
		// Performance check - should be reasonably fast
		if duration > 100*time.Millisecond {
			t.Logf("Large response took %v (consider optimizing if consistently slow)", duration)
		}
		
		t.Logf("Large response performance: %v for %d items", duration, len(result))
	})
	
	t.Run("Concurrent requests", func(t *testing.T) {
		concurrency := 10
		done := make(chan bool, concurrency)
		
		start := time.Now()
		
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				callOptions := map[string]interface{}{
					"call-id":    fmt.Sprintf("concurrent-%d", id),
					"methodName": "largeResponse",
					"args":       []interface{}{},
				}
				
				body, _ := json.Marshal(callOptions)
				
				req, err := http.NewRequest("POST", server.URL+"?object=0&method=0&args="+string(body), nil)
				if err != nil {
					t.Errorf("Failed to create request: %v", err)
					return
				}
				req.Header.Set("x-wails-window-id", "1")
				
				client := &http.Client{Timeout: 30 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					t.Errorf("Request failed: %v", err)
					return
				}
				defer resp.Body.Close()
				
				if resp.StatusCode != http.StatusOK {
					t.Errorf("Expected status 200, got %d", resp.StatusCode)
					return
				}
				
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
					return
				}
				
				if len(result) != 1000 {
					t.Errorf("Expected 1000 items, got %d", len(result))
				}
			}(i)
		}
		
		// Wait for all requests to complete
		for i := 0; i < concurrency; i++ {
			<-done
		}
		
		duration := time.Since(start)
		t.Logf("Concurrent requests (%d) completed in: %v", concurrency, duration)
		
		// Should handle concurrent requests reasonably well
		if duration > 1*time.Second {
			t.Logf("Concurrent requests took %v (consider optimizing)", duration)
		}
	})
}

// Enhanced mockBoundMethod that can handle dynamic results
type mockBoundMethodEnhanced struct {
	name   string
	result interface{}
	err    error
	delay  time.Duration
}

func (m *mockBoundMethodEnhanced) String() string { return m.name }
func (m *mockBoundMethodEnhanced) Call(ctx context.Context, args []json.RawMessage) (interface{}, error) {
	// Simulate work time
	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	
	// Handle function results
	if fn, ok := m.result.(func([]json.RawMessage) interface{}); ok {
		return fn(args), m.err
	}
	
	return m.result, m.err
}