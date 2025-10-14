package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Transport defines the interface for custom IPC transport implementations.
// Developers can provide their own transport (e.g., WebSocket, custom protocol)
// while retaining all Wails generated bindings and event communication.
//
// The transport is responsible for:
//   - Receiving runtime call requests from the frontend
//   - Processing them through Wails' MessageProcessor
//   - Sending responses back to the frontend
//
// Example use case: Implementing WebSocket-based transport instead of HTTP fetch.
type Transport interface {
	// Start initializes and starts the transport layer.
	// The provided handler should be called to process Wails runtime requests.
	// The context is the application context and will be cancelled on shutdown.
	Start(ctx context.Context, handler TransportHandler) error

	// Stop gracefully shuts down the transport.
	Stop() error
}

// AssetServerTransport is an optional interface that transports can implement
// to serve assets over HTTP, enabling browser-based deployments.
//
// When a transport implements this interface, Wails will call ServeAssets()
// after Start() to provide the asset server handler. The transport should
// integrate this handler into its HTTP server to serve HTML, CSS, JS, and
// other static assets alongside the IPC transport.
//
// This is useful for:
//   - Running Wails apps in a browser instead of a webview
//   - Exposing the app over a network
//   - Custom server configurations with both assets and IPC
type AssetServerTransport interface {
	Transport

	// ServeAssets configures the transport to serve assets.
	// The assetHandler is Wails' internal asset server that handles:
	//   - All static assets (HTML, CSS, JS, images, etc.)
	//   - /wails/runtime.js - The Wails runtime library
	//   - /wails/capabilities - Capability information
	//   - /wails/flags - Application flags
	//
	// The transport should integrate this handler into its HTTP server.
	// Typically this means mounting it at "/" and ensuring the IPC endpoint
	// (e.g., /wails/ws for WebSocket) is handled separately.
	//
	// This method is called after Start() completes successfully.
	ServeAssets(assetHandler http.Handler) error
}

// TransportHandler wraps the Wails message processor to handle runtime calls.
// Custom transports should invoke this handler to process bound method calls,
// events, dialogs, clipboard operations, etc.
type TransportHandler interface {
	// HandleRuntimeCall processes a runtime call request and returns the response.
	// This is the main entry point for processing all Wails runtime operations.
	//
	// Parameters:
	//   - ctx: Request context for cancellation/timeouts
	//   - req: The runtime call request
	//
	// Returns:
	//   - TransportResponse containing the result or error
	HandleRuntimeCall(ctx context.Context, req *TransportRequest) *TransportResponse
}

// TransportRequest represents a runtime call request from the frontend.
// This maps to the parameters sent by the frontend runtime in calls.ts and runtime.ts.
type TransportRequest struct {
	// Object identifies which Wails subsystem to call (Call=0, Clipboard=1, etc.)
	// See objectNames in runtime.ts
	Object int `json:"object"`

	// Method identifies which method within the object to call
	Method int `json:"method"`

	// Args contains the method arguments as a JSON string
	// For bound method calls (Object=0), this contains CallOptions with methodID/methodName and args
	Args string `json:"args,omitempty"`

	// WindowID identifies the source window (optional, sent via header x-wails-window-id)
	WindowID string `json:"windowId,omitempty"`

	// WindowName identifies the source window by name (optional, sent via header x-wails-window-name)
	WindowName string `json:"windowName,omitempty"`

	// ClientID identifies the frontend client (sent via header x-wails-client-id)
	ClientID string `json:"clientId,omitempty"`
}

// TransportResponse represents the response to a runtime call.
type TransportResponse struct {
	// StatusCode is the HTTP status code equivalent (200 for success, 422 for error, etc.)
	StatusCode int `json:"statusCode"`

	// ContentType indicates the response format ("application/json" or "text/plain")
	ContentType string `json:"contentType"`

	// Data contains the response body (can be []byte, string, or any encoded format)
	// The type depends on the TransportCodec being used
	Data interface{} `json:"data"`

	// Error contains error information if the call failed
	Error error `json:"error,omitempty"`
}

// httpTransport is the default HTTP-based transport using fetch() from the frontend.
// This is the standard Wails transport that routes through the asset server.
type httpTransport struct {
	// The asset server already handles the HTTP transport via the middleware in application.go
	// This is a placeholder to represent the default behavior.
}

func (h *httpTransport) Start(ctx context.Context, handler TransportHandler) error {
	// The default HTTP transport is handled by the asset server middleware
	// in application.go (lines 100-101), so nothing to start here.
	return nil
}

func (h *httpTransport) Stop() error {
	// Nothing to stop for HTTP transport
	return nil
}

// transportHandler implements TransportHandler by wrapping MessageProcessor
type transportHandler struct {
	messageProcessor *MessageProcessor
}

// HandleRuntimeCall processes a transport request through the MessageProcessor
func (t *transportHandler) HandleRuntimeCall(ctx context.Context, req *TransportRequest) *TransportResponse {
	// For binding calls (Object=0, Method=0), we need to intercept the async callback
	// and return it synchronously for non-webview transports
	if req.Object == 0 && req.Method == 0 {
		return t.handleBindingCall(ctx, req)
	}

	// For other runtime calls, use standard HTTP request/response
	return t.handleStandardCall(ctx, req)
}

// handleStandardCall processes non-binding runtime calls (clipboard, dialogs, etc.)
func (t *transportHandler) handleStandardCall(ctx context.Context, req *TransportRequest) *TransportResponse {
	r, err := http.NewRequestWithContext(ctx, "GET", "/wails/runtime", nil)
	if err != nil {
		return &TransportResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	// Set query parameters
	q := r.URL.Query()
	q.Set("object", fmt.Sprintf("%d", req.Object))
	q.Set("method", fmt.Sprintf("%d", req.Method))
	if req.Args != "" {
		q.Set("args", req.Args)
	}
	r.URL.RawQuery = q.Encode()

	// Set headers
	if req.WindowID != "" {
		r.Header.Set(webViewRequestHeaderWindowId, req.WindowID)
	}
	if req.WindowName != "" {
		r.Header.Set(webViewRequestHeaderWindowName, req.WindowName)
	}
	if req.ClientID != "" {
		r.Header.Set("x-wails-client-id", req.ClientID)
	}

	// Create response recorder
	rw := &transportResponseWriter{
		header:     make(http.Header),
		statusCode: http.StatusOK,
	}

	// Process through MessageProcessor
	t.messageProcessor.HandleRuntimeCallWithIDs(rw, r)

	return &TransportResponse{
		StatusCode:  rw.statusCode,
		ContentType: rw.header.Get("Content-Type"),
		Data:        rw.body,
		Error:       rw.err,
	}
}

// handleBindingCall processes binding calls synchronously by intercepting the callback
func (t *transportHandler) handleBindingCall(ctx context.Context, req *TransportRequest) *TransportResponse {
	// Parse the call-id from args
	var argsMap map[string]interface{}
	if err := json.Unmarshal([]byte(req.Args), &argsMap); err != nil {
		return &TransportResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Error:      fmt.Errorf("failed to parse binding call args: %w", err),
		}
	}

	callID, ok := argsMap["call-id"].(string)
	if !ok {
		return &TransportResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Error:      fmt.Errorf("missing call-id in binding call"),
		}
	}

	// Create a channel to receive the result
	resultChan := make(chan string, 1)
	errorChan := make(chan string, 1)

	// Create HTTP request
	r, err := http.NewRequestWithContext(ctx, "GET", "/wails/runtime", nil)
	if err != nil {
		return &TransportResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	q := r.URL.Query()
	q.Set("object", "0")
	q.Set("method", "0")
	q.Set("args", req.Args)
	r.URL.RawQuery = q.Encode()

	if req.ClientID != "" {
		r.Header.Set("x-wails-client-id", req.ClientID)
	}

	// Process through MessageProcessor with our callback window
	rw := &transportResponseWriter{
		header:     make(http.Header),
		statusCode: http.StatusOK,
	}

	// Temporarily replace the window in the request context
	// Actually, we need to call processCallMethod directly with our window
	// But that's private. Instead, let's use the standard flow and intercept at window level

	// Get the target window
	windows := globalApplication.Window.GetAll()
	if len(windows) == 0 {
		return &TransportResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Errorf("no windows available"),
		}
	}

	// Store original window callback handlers and replace them temporarily
	originalWindow := windows[0]
	callbackInterceptor := &windowCallbackInterceptor{
		Window:     originalWindow,
		callID:     callID,
		resultChan: resultChan,
		errorChan:  errorChan,
	}

	// Temporarily replace the window (this is hacky but necessary for sync responses)
	globalApplication.windowsLock.Lock()
	var windowID uint
	for id, win := range globalApplication.windows {
		if win == originalWindow {
			windowID = id
			globalApplication.windows[id] = callbackInterceptor
			break
		}
	}
	globalApplication.windowsLock.Unlock()

	defer func() {
		globalApplication.windowsLock.Lock()
		globalApplication.windows[windowID] = originalWindow
		globalApplication.windowsLock.Unlock()
	}()

	// Process the call
	t.messageProcessor.HandleRuntimeCallWithIDs(rw, r)

	// Wait for the result or error
	select {
	case result := <-resultChan:
		// Check if result is JSON
		isJSON := len(result) > 0 && (result[0] == '{' || result[0] == '[' || result[0] == '"')
		contentType := "text/plain"
		if isJSON {
			contentType = "application/json"
		}
		return &TransportResponse{
			StatusCode:  http.StatusOK,
			ContentType: contentType,
			Data:        []byte(result),
		}
	case errMsg := <-errorChan:
		return &TransportResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Data:       []byte(errMsg),
		}
	case <-ctx.Done():
		return &TransportResponse{
			StatusCode: http.StatusRequestTimeout,
			Error:      ctx.Err(),
		}
	}
}

// transportResponseWriter captures HTTP responses for transport processing
type transportResponseWriter struct {
	header     http.Header
	body       []byte
	statusCode int
	err        error
}

func (w *transportResponseWriter) Header() http.Header {
	return w.header
}

func (w *transportResponseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	return len(data), nil
}

func (w *transportResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// windowCallbackInterceptor wraps a Window and intercepts CallResponse/CallError
// to provide synchronous responses for custom transports
type windowCallbackInterceptor struct {
	Window     // Embed the Window interface to inherit all methods
	callID     string
	resultChan chan string
	errorChan  chan string
}

// CallResponse intercepts the async callback and sends it to the result channel
func (w *windowCallbackInterceptor) CallResponse(callID string, result string) {
	if callID == w.callID {
		w.resultChan <- result
	} else {
		// Forward to original window if different call ID
		w.Window.CallResponse(callID, result)
	}
}

// CallError intercepts the async error callback and sends it to the error channel
func (w *windowCallbackInterceptor) CallError(callID string, message string, isJSON bool) {
	if callID == w.callID {
		w.errorChan <- message
	} else {
		// Forward to original window if different call ID
		w.Window.CallError(callID, message, isJSON)
	}
}

// NewHTTPTransport returns the default HTTP-based transport.
// This is used when no custom transport is specified.
func NewHTTPTransport() Transport {
	return &httpTransport{}
}

// TransportRequestFromHTTP converts an HTTP request to a TransportRequest.
// This is used internally by the asset server middleware.
func TransportRequestFromHTTP(r *http.Request) *TransportRequest {
	object := 0
	method := 0
	args := ""

	if objStr := r.URL.Query().Get("object"); objStr != "" {
		if o, ok := parseIntSafe(objStr); ok {
			object = o
		}
	}
	if methStr := r.URL.Query().Get("method"); methStr != "" {
		if m, ok := parseIntSafe(methStr); ok {
			method = m
		}
	}
	if argsStr := r.URL.Query().Get("args"); argsStr != "" {
		args = argsStr
	}

	return &TransportRequest{
		Object:     object,
		Method:     method,
		Args:       args,
		WindowID:   r.Header.Get(webViewRequestHeaderWindowId),
		WindowName: r.Header.Get(webViewRequestHeaderWindowName),
		ClientID:   r.Header.Get("x-wails-client-id"),
	}
}

// parseIntSafe safely parses an integer string without panicking
func parseIntSafe(s string) (int, bool) {
	result, err := strconv.Atoi(s)
	return result, err == nil
}
