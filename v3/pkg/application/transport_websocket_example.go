package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketTransport is an example implementation of a WebSocket-based transport.
// This demonstrates how to create a custom transport that can replace the default
// HTTP fetch-based IPC while retaining all Wails bindings and event communication.
//
// This implementation is provided as an example and is not production-ready.
// You may need to add error handling, reconnection logic, authentication, etc.
type WebSocketTransport struct {
	addr     string
	server   *http.Server
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
	mu       sync.RWMutex
	handler  TransportHandler
}

// WebSocketMessage represents a message sent over the WebSocket transport
type WebSocketMessage struct {
	ID      string           `json:"id"`      // Unique message ID for request/response matching
	Type    string           `json:"type"`    // "request" or "response"
	Request *TransportRequest `json:"request,omitempty"`
	Response *TransportResponse `json:"response,omitempty"`
}

// NewWebSocketTransport creates a new WebSocket transport listening on the specified address.
// Example: NewWebSocketTransport(":9998")
func NewWebSocketTransport(addr string) *WebSocketTransport {
	return &WebSocketTransport{
		addr: addr,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// In production, you should validate the origin
				return true
			},
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

// Start initializes and starts the WebSocket server
func (w *WebSocketTransport) Start(ctx context.Context, handler TransportHandler) error {
	w.handler = handler

	// Create HTTP server but don't start it yet
	// We'll set up the handler in ServeAssets() if it's called
	w.server = &http.Server{
		Addr: w.addr,
	}

	// Handle context cancellation
	go func() {
		<-ctx.Done()
		w.Stop()
	}()

	return nil
}

// ServeAssets configures the transport to serve assets alongside WebSocket IPC.
// This implements the AssetServerTransport interface for browser-based deployments.
func (w *WebSocketTransport) ServeAssets(assetHandler http.Handler) error {
	mux := http.NewServeMux()

	// Mount WebSocket endpoint for IPC
	mux.HandleFunc("/wails/ws", w.handleWebSocket)

	// Mount asset server for all other requests
	mux.Handle("/", assetHandler)

	// Set the handler and start the server
	w.server.Handler = mux

	// Start server in background
	go func() {
		log.Printf("WebSocket transport serving assets and IPC on %s", w.addr)
		log.Printf("  - Assets: http://%s/", w.addr)
		log.Printf("  - WebSocket IPC: ws://%s/wails/ws", w.addr)
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("WebSocket server error: %v", err)
		}
	}()

	return nil
}

// Stop gracefully shuts down the WebSocket server
func (w *WebSocketTransport) Stop() error {
	if w.server == nil {
		return nil
	}

	// Close all client connections
	w.mu.Lock()
	for conn := range w.clients {
		conn.Close()
	}
	w.clients = make(map[*websocket.Conn]bool)
	w.mu.Unlock()

	return w.server.Shutdown(context.Background())
}

// handleWebSocket handles WebSocket connections
func (w *WebSocketTransport) handleWebSocket(rw http.ResponseWriter, r *http.Request) {
	conn, err := w.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	w.mu.Lock()
	w.clients[conn] = true
	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		delete(w.clients, conn)
		w.mu.Unlock()
		conn.Close()
	}()

	// Read messages from client
	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		// Process request
		if msg.Type == "request" && msg.Request != nil {
			go w.handleRequest(conn, msg.ID, msg.Request)
		}
	}
}

// handleRequest processes a runtime call request and sends the response
func (w *WebSocketTransport) handleRequest(conn *websocket.Conn, msgID string, req *TransportRequest) {
	log.Printf("[WebSocket] Received request: msgID=%s, object=%d, method=%d, args=%s", msgID, req.Object, req.Method, req.Args)

	// Call the Wails runtime handler
	response := w.handler.HandleRuntimeCall(context.Background(), req)

	log.Printf("[WebSocket] Response: statusCode=%d, contentType=%s, dataLen=%d", response.StatusCode, response.ContentType, len(response.Data))

	// Send response back to client
	responseMsg := WebSocketMessage{
		ID:       msgID,
		Type:     "response",
		Response: response,
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	if err := conn.WriteJSON(responseMsg); err != nil {
		log.Printf("[WebSocket] Failed to send response: %v", err)
	} else {
		log.Printf("[WebSocket] Successfully sent response for msgID=%s", msgID)
	}
}

// BroadcastEvent sends an event to all connected clients
// This can be used for server-pushed events
func (w *WebSocketTransport) BroadcastEvent(event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := WebSocketMessage{
		Type: "event",
		Response: &TransportResponse{
			StatusCode:  200,
			ContentType: "application/json",
			Data:        data,
		},
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	for conn := range w.clients {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Failed to broadcast event: %v", err)
		}
	}

	return nil
}
