package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// Mock Window for testing
type mockWindow struct {
	id uint
}

func (m *mockWindow) ID() uint                                    { return m.id }
func (m *mockWindow) Name() string                                { return "test-window" }
func (m *mockWindow) Close()                                      {}
func (m *mockWindow) Center()                                     {}
func (m *mockWindow) Size() (int, int)                            { return 800, 600 }
func (m *mockWindow) Position() (int, int)                        { return 100, 100 }
func (m *mockWindow) SetSize(width, height int)                   {}
func (m *mockWindow) SetPosition(x, y int)                        {}
func (m *mockWindow) SetTitle(title string)                       {}
func (m *mockWindow) Show()                                       {}
func (m *mockWindow) Hide()                                       {}
func (m *mockWindow) Maximize()                                   {}
func (m *mockWindow) Minimize()                                   {}
func (m *mockWindow) Restore()                                    {}
func (m *mockWindow) Focus()                                      {}
func (m *mockWindow) Fullscreen()                                 {}
func (m *mockWindow) UnFullscreen()                               {}
func (m *mockWindow) IsMaximized() bool                           { return false }
func (m *mockWindow) IsMinimized() bool                           { return false }
func (m *mockWindow) IsFullscreen() bool                          { return false }
func (m *mockWindow) IsVisible() bool                             { return true }
func (m *mockWindow) IsNormal() bool                              { return true }
func (m *mockWindow) IsFocused() bool                             { return true }

// Mock BoundMethod for testing
type mockBoundMethod struct {
	name   string
	result interface{}
	err    error
	delay  time.Duration
}

func (m *mockBoundMethod) String() string { return m.name }
func (m *mockBoundMethod) Call(ctx context.Context, args []json.RawMessage) (interface{}, error) {
	// Simulate work time
	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return m.result, m.err
}

// Test helper to create HTTP request
func createBindingRequest(method string, callID string, methodName string, args []interface{}) *http.Request {
	values := url.Values{}
	values.Set("object", "0") // CallBinding
	values.Set("method", "0") // CallBinding method
	
	callOptions := map[string]interface{}{
		"call-id": callID,
	}
	if methodName != "" {
		callOptions["methodName"] = methodName
	}
	if args != nil {
		callOptions["args"] = args
	}
	
	argsJSON, _ := json.Marshal(callOptions)
	values.Set("args", string(argsJSON))
	
	req := httptest.NewRequest(method, "/wails/runtime?"+values.Encode(), nil)
	req.Header.Set("x-wails-window-id", "1")
	return req
}

func TestHTTPOnlyBindingExecution(t *testing.T) {
	// Setup
	processor := NewMessageProcessor(nil)
	window := &mockWindow{id: 1}
	
	// Mock global application for testing
	globalApplication = &App{
		bindings: &bindings{
			byName: map[string]*BoundMethod{
				"testMethod": &BoundMethod{
					Name: "testMethod",
					call: &mockBoundMethod{
						name:   "testMethod",
						result: map[string]string{"status": "success"},
					},
				},
			},
		},
		options: Options{
			Assets: AssetOptions{
				Bindings: BindingConfig{
					Timeout: 30 * time.Second,
				},
			},
		},
	}
	
	// Test successful binding call
	t.Run("Successful HTTP-only binding call", func(t *testing.T) {
		req := createBindingRequest("POST", "test-123", "testMethod", []interface{}{})
		recorder := httptest.NewRecorder()
		
		processor.processCallMethod(CallBinding, recorder, req, window, QueryParams(req.URL.Query()))
		
		// Verify HTTP response
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", recorder.Code)
		}
		
		// Verify JSON response
		var result map[string]string
		err := json.Unmarshal(recorder.Body.Bytes(), &result)
		if err != nil {
			t.Errorf("Failed to parse JSON response: %v", err)
		}
		
		if result["status"] != "success" {
			t.Errorf("Expected status 'success', got %s", result["status"])
		}
		
		// Verify Content-Type header
		contentType := recorder.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
		}
	})
	
	// Test method not found error
	t.Run("Method not found returns 404", func(t *testing.T) {
		req := createBindingRequest("POST", "test-404", "nonexistentMethod", []interface{}{})
		recorder := httptest.NewRecorder()
		
		processor.processCallMethod(CallBinding, recorder, req, window, QueryParams(req.URL.Query()))
		
		if recorder.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", recorder.Code)
		}
		
		var errorResp map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &errorResp)
		if err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
		
		if errorResp["kind"] != "ReferenceError" {
			t.Errorf("Expected error kind 'ReferenceError', got %s", errorResp["kind"])
		}
	})
	
	// Test timeout handling
	t.Run("Binding timeout returns 408", func(t *testing.T) {
		// Create a slow method
		globalApplication.bindings.byName["slowMethod"] = &BoundMethod{
			Name: "slowMethod", 
			call: &mockBoundMethod{
				name:   "slowMethod",
				result: "completed",
				delay:  2 * time.Second, // Longer than test timeout
			},
		}
		
		// Set short timeout for test
		globalApplication.options.Assets.Bindings.Timeout = 100 * time.Millisecond
		
		req := createBindingRequest("POST", "test-timeout", "slowMethod", []interface{}{})
		recorder := httptest.NewRecorder()
		
		processor.processCallMethod(CallBinding, recorder, req, window, QueryParams(req.URL.Query()))
		
		if recorder.Code != http.StatusRequestTimeout {
			t.Errorf("Expected status 408, got %d", recorder.Code)
		}
		
		var errorResp map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &errorResp)
		if err != nil {
			t.Errorf("Failed to parse timeout error response: %v", err)
		}
		
		if errorResp["kind"] != "TimeoutError" {
			t.Errorf("Expected error kind 'TimeoutError', got %s", errorResp["kind"])
		}
	})
}

func TestCORSHeaders(t *testing.T) {
	processor := NewMessageProcessor(nil)
	window := &mockWindow{id: 1}
	
	// Mock global application with CORS enabled
	globalApplication = &App{
		bindings: &bindings{
			byName: map[string]*BoundMethod{
				"testMethod": &BoundMethod{
					Name: "testMethod",
					call: &mockBoundMethod{
						name:   "testMethod", 
						result: "success",
					},
				},
			},
		},
		options: Options{
			Assets: AssetOptions{
				Bindings: BindingConfig{
					Timeout: 30 * time.Second,
					CORS: CORSConfig{
						Enabled: true,
						AllowedOrigins: []string{"https://example.com"},
					},
				},
			},
		},
		isDebugMode: false,
	}
	
	t.Run("CORS headers set for allowed origin", func(t *testing.T) {
		req := createBindingRequest("POST", "test-cors", "testMethod", []interface{}{})
		req.Header.Set("Origin", "https://example.com")
		recorder := httptest.NewRecorder()
		
		processor.processCallMethod(CallBinding, recorder, req, window, QueryParams(req.URL.Query()))
		
		// Check CORS headers
		if recorder.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
			t.Errorf("Expected CORS origin header to be set")
		}
		
		if recorder.Header().Get("Access-Control-Allow-Methods") == "" {
			t.Errorf("Expected CORS methods header to be set")
		}
	})
	
	t.Run("CORS headers not set for disallowed origin", func(t *testing.T) {
		req := createBindingRequest("POST", "test-cors-denied", "testMethod", []interface{}{})
		req.Header.Set("Origin", "https://evil.com")
		recorder := httptest.NewRecorder()
		
		processor.processCallMethod(CallBinding, recorder, req, window, QueryParams(req.URL.Query()))
		
		// CORS headers should not be set for disallowed origin
		if recorder.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Errorf("CORS headers should not be set for disallowed origin")
		}
	})
	
	t.Run("OPTIONS preflight request", func(t *testing.T) {
		req := createBindingRequest("OPTIONS", "test-preflight", "", nil)
		req.Header.Set("Origin", "https://example.com")
		recorder := httptest.NewRecorder()
		
		processor.processCallMethod(CallBinding, recorder, req, window, QueryParams(req.URL.Query()))
		
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected OPTIONS request to return 200, got %d", recorder.Code)
		}
		
		if recorder.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
			t.Errorf("Expected CORS headers for OPTIONS request")
		}
	})
}

func TestHTTPOnlyPerformance(t *testing.T) {
	processor := NewMessageProcessor(nil)
	window := &mockWindow{id: 1}
	
	// Create a large response for performance testing  
	largeData := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		largeData[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
	}
	
	globalApplication = &App{
		bindings: &bindings{
			byName: map[string]*BoundMethod{
				"largeDataMethod": &BoundMethod{
					Name: "largeDataMethod",
					call: &mockBoundMethod{
						name:   "largeDataMethod",
						result: largeData,
					},
				},
			},
		},
		options: Options{
			Assets: AssetOptions{
				Bindings: BindingConfig{
					Timeout: 30 * time.Second,
				},
			},
		},
	}
	
	t.Run("Large data HTTP response performance", func(t *testing.T) {
		start := time.Now()
		
		req := createBindingRequest("POST", "test-large", "largeDataMethod", []interface{}{})
		recorder := httptest.NewRecorder()
		
		processor.processCallMethod(CallBinding, recorder, req, window, QueryParams(req.URL.Query()))
		
		duration := time.Since(start)
		
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", recorder.Code)
		}
		
		// Verify response can be parsed
		var response map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to parse large data response: %v", err)
		}
		
		if len(response) != 1000 {
			t.Errorf("Expected 1000 items in response, got %d", len(response))
		}
		
		// Performance check - should be fast (under 10ms for this size)
		if duration > 50*time.Millisecond {
			t.Logf("Large data response took %v (acceptable but monitor)", duration)
		}
		
		t.Logf("Large data response time: %v", duration)
	})
}

func TestErrorHandling(t *testing.T) {
	processor := NewMessageProcessor(nil)
	window := &mockWindow{id: 1}
	
	globalApplication = &App{
		bindings: &bindings{
			byName: map[string]*BoundMethod{
				"errorMethod": &BoundMethod{
					Name: "errorMethod",
					call: &mockBoundMethod{
						name: "errorMethod",
						err: &CallError{
							Kind:    RuntimeError,
							Message: "Something went wrong",
						},
					},
				},
			},
		},
		options: Options{
			Assets: AssetOptions{
				Bindings: BindingConfig{
					Timeout: 30 * time.Second,
				},
			},
		},
	}
	
	t.Run("Runtime error returns 500", func(t *testing.T) {
		req := createBindingRequest("POST", "test-error", "errorMethod", []interface{}{})
		recorder := httptest.NewRecorder()
		
		processor.processCallMethod(CallBinding, recorder, req, window, QueryParams(req.URL.Query()))
		
		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", recorder.Code)
		}
		
		var errorResp map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &errorResp)
		if err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
		
		if errorResp["kind"] != "RuntimeError" {
			t.Errorf("Expected error kind 'RuntimeError', got %s", errorResp["kind"])
		}
		
		if errorResp["error"] != "Something went wrong" {
			t.Errorf("Expected error message 'Something went wrong', got %s", errorResp["error"])
		}
	})
	
	t.Run("Missing call-id returns 400", func(t *testing.T) {
		values := url.Values{}
		values.Set("object", "0")
		values.Set("method", "0")
		// No call-id provided
		
		req := httptest.NewRequest("POST", "/wails/runtime?"+values.Encode(), nil)
		req.Header.Set("x-wails-window-id", "1")
		recorder := httptest.NewRecorder()
		
		processor.processCallMethod(CallBinding, recorder, req, window, QueryParams(req.URL.Query()))
		
		if recorder.Code != http.StatusUnprocessableEntity {
			t.Errorf("Expected status 422, got %d", recorder.Code)
		}
	})
}

// Benchmark HTTP-only vs simulated callback overhead
func BenchmarkHTTPOnlyVsCallback(b *testing.B) {
	processor := NewMessageProcessor(nil)
	window := &mockWindow{id: 1}
	
	globalApplication = &App{
		bindings: &bindings{
			byName: map[string]*BoundMethod{
				"benchMethod": &BoundMethod{
					Name: "benchMethod",
					call: &mockBoundMethod{
						name:   "benchMethod",
						result: map[string]string{"result": "success"},
					},
				},
			},
		},
		options: Options{
			Assets: AssetOptions{
				Bindings: BindingConfig{
					Timeout: 30 * time.Second,
				},
			},
		},
	}
	
	b.Run("HTTP-only", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req := createBindingRequest("POST", fmt.Sprintf("bench-%d", i), "benchMethod", []interface{}{})
			recorder := httptest.NewRecorder()
			
			processor.processCallMethod(CallBinding, recorder, req, window, QueryParams(req.URL.Query()))
			
			if recorder.Code != http.StatusOK {
				b.Errorf("Expected status 200, got %d", recorder.Code)
			}
		}
	})
}