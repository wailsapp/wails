package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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
	Start(ctx context.Context, messageProcessor *MessageProcessor) error

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

// EventTransport is an optional interface for transports that support
// server-to-client event delivery (e.g., WebSocket, Server-Sent Events).
//
// When a transport implements this interface, Wails automatically forwards
// all application events to the transport, eliminating the need for manual
// event wiring code.
//
// Example implementation:
//
//	type MyWebSocketTransport struct { ... }
//
//	func (t *MyWebSocketTransport) SendEvent(event *WailsEvent) error {
//	    // Broadcast event to all connected WebSocket clients
//	    return t.broadcast(event)
//	}
//
// The transport determines the delivery strategy:
//   - Broadcast to all clients (typical for WebSocket)
//   - Targeted delivery to specific windows/clients
//   - Buffering for disconnected clients
//   - Delivery guarantees and acknowledgments
type EventTransport interface {
	// SendEvent delivers an event to connected client(s).
	// Called automatically by Wails when any event is emitted via app.Events.Emit().
	//
	// The event contains:
	//   - Name: Event name (e.g., "greet:count", "user:login")
	//   - Data: Event payload (any JSON-serializable data)
	//   - Sender: Originating window name (optional)
	//
	// Returns an error if event delivery fails. Errors are logged but don't
	// prevent the event from being dispatched to other listeners.
	SendEvent(event *CustomEvent) error
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

	// Args contains the method arguments map[string]any
	// For bound method calls (Object=0), this contains CallOptions with methodID/methodName and args
	Args map[string]any `json:"args,omitempty"`

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

	// Data contains the response body (can be struct, string)
	Data any `json:"data"`
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
type legacyTransport struct {
	messageProcessor *MessageProcessor
	logger           *slog.Logger
}

func (t *legacyTransport) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("object") == "" {
		t.httpError(rw, "Invalid runtime call:", errors.New("missing object value"))
		return
	}

	object, err := strconv.Atoi(r.URL.Query().Get("object"))
	if err != nil {
		t.httpError(rw, "Invalid runtime call:", fmt.Errorf("error decoding object value: %w", err))
		return
	}
	method, err := strconv.Atoi(r.URL.Query().Get("method"))
	if err != nil {
		t.httpError(rw, "Invalid runtime call:", fmt.Errorf("error decoding method value: %w", err))
		return
	}
	params := QueryParams(r.URL.Query())

	windowIdStr := r.Header.Get(webViewRequestHeaderWindowId)
	windowId := 0
	if windowIdStr != "" {
		windowId, err = strconv.Atoi(windowIdStr)
		if err != nil {
			t.httpError(rw, "Invalid runtime call:", fmt.Errorf("error decoding windowId value: %w", err))
			return
		}
	}

	windowName := r.Header.Get(webViewRequestHeaderWindowName)
	clientId := r.Header.Get("x-wails-client-id")

	resp, err := t.messageProcessor.HandleRuntimeCallWithIDs(r.Context(), &RuntimeRequest{
		Object:            object,
		Method:            method,
		Params:            params,
		WebviewWindowId:   uint32(windowId),
		WebviewWindowName: windowName,
		ClientId:          clientId,
	})

	if err != nil {
		t.httpError(rw, "Failed to process runtime call:", err)
		return
	}

	if stringResp, ok := resp.(string); ok {
		t.text(rw, stringResp)
		return
	}

	t.json(rw, resp)
}

func (t *legacyTransport) text(rw http.ResponseWriter, data string) {
	_, err := rw.Write([]byte(data))
	if err != nil {
		t.error("Unable to write json payload. Please report this to the Wails team!", "error", err)
		return
	}
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
}

func (t *legacyTransport) json(rw http.ResponseWriter, data any) {
	rw.Header().Set("Content-Type", "application/json")
	// convert data to json
	var jsonPayload = []byte("{}")
	var err error
	if data != nil {
		jsonPayload, err = json.Marshal(data)
		if err != nil {
			t.error("Unable to convert data to JSON. Please report this to the Wails team!", "error", err)
			return
		}
	}
	_, err = rw.Write(jsonPayload)
	if err != nil {
		t.error("Unable to write json payload. Please report this to the Wails team!", "error", err)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (t *legacyTransport) httpError(rw http.ResponseWriter, message string, err error) {
	t.error(message, "error", err)
	rw.WriteHeader(http.StatusUnprocessableEntity)
	_, err = rw.Write([]byte(err.Error()))
	if err != nil {
		t.error("Unable to write error response:", "error", err)
	}
}

func (t *legacyTransport) error(message string, args ...any) {
	t.logger.Error(message, args...)
}
