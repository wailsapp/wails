package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v3/pkg/application"
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
	clients  map[*websocket.Conn]chan *WebSocketMessage
	mu       sync.RWMutex
	handler  *application.MessageProcessor
}

// WebSocketTransportOption is a functional option for configuring WebSocketTransport
type WebSocketTransportOption func(*WebSocketTransport)

// wsResponse represents the response to a runtime call.
type wsResponse struct {
	// StatusCode is the HTTP status code equivalent (200 for success, 422 for error, etc.)
	StatusCode int `json:"statusCode"`

	// Data contains the response body (can be struct, string)
	Data any `json:"data"`
}

// WebSocketMessage represents a message sent over the WebSocket transport
type WebSocketMessage struct {
	ID       string                      `json:"id"`   // Unique message ID for request/response matching
	Type     string                      `json:"type"` // "request" or "response"
	Request  *application.RuntimeRequest `json:"request,omitempty"`
	Response *wsResponse                 `json:"response,omitempty"`
	Event    *application.CustomEvent    `json:"event,omitempty"`
}

// NewWebSocketTransport creates a new WebSocket transport listening on the specified address.
// Example: NewWebSocketTransport(":9099")
func NewWebSocketTransport(addr string, opts ...WebSocketTransportOption) *WebSocketTransport {
	t := &WebSocketTransport{
		addr: addr,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// In production, you should validate the origin
				return true
			},
		},
		clients: make(map[*websocket.Conn]chan *WebSocketMessage),
	}

	// Apply options
	for _, opt := range opts {
		opt(t)
	}

	return t
}

// Start initializes and starts the WebSocket server
func (w *WebSocketTransport) Start(ctx context.Context, handler *application.MessageProcessor) error {
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

	w.mu.Lock()
	for conn := range w.clients {
		conn.Close()
	}
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
	messageChan := make(chan *WebSocketMessage, 100)
	w.clients[conn] = messageChan
	w.mu.Unlock()

	ctx, cancel := context.WithCancel(r.Context())

	defer func() {
		w.mu.Lock()
		cancel()
		close(messageChan)
		if _, ok := w.clients[conn]; ok {
			delete(w.clients, conn)
		}
		w.mu.Unlock()
		conn.Close()
	}()

	// write responses in one place, as concurrent writeJSON is not allowed
	go func() {
		for {
			select {
			case msg, ok := <-messageChan:
				if !ok {
					return
				}

				w.mu.RLock()
				if err := conn.WriteJSON(msg); err != nil {
					log.Printf("[WebSocket] Failed to send message: %v", err)
				} else {
					if msg.Type == "response" {
						log.Printf("[WebSocket] Successfully sent response for msgID=%s", msg.ID)
					}
				}
				w.mu.RUnlock()
			}
		}
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
			if msg.Request.Args == nil {
				msg.Request.Args = &application.Args{}
			}
			go w.handleRequest(ctx, messageChan, msg.ID, msg.Request)
		}
	}
}

// handleRequest processes a runtime call request and sends the response
func (w *WebSocketTransport) handleRequest(ctx context.Context, messageChan chan *WebSocketMessage, msgID string, req *application.RuntimeRequest) {
	log.Printf("[WebSocket] Received request: msgID=%s, object=%d, method=%d, args=%s", msgID, req.Object, req.Method, req.Args.String())

	// Call the Wails runtime handler
	response, err := w.handler.HandleRuntimeCallWithIDs(ctx, req)

	w.sendResponse(ctx, messageChan, msgID, response, err)
}

// sendResponse sends a response message to the client
func (w *WebSocketTransport) sendResponse(ctx context.Context, messageChan chan *WebSocketMessage, msgID string, resp any, err error) {
	response := &wsResponse{
		StatusCode: 200,
		Data:       resp,
	}
	if err != nil {
		response.StatusCode = 422
		response.Data = err.Error()
	}

	responseMsg := &WebSocketMessage{
		ID:       msgID,
		Type:     "response",
		Response: response,
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	select {
	case <-ctx.Done():
		log.Println("[WebSocket] Context cancelled before sending response.")
	default:
		messageChan <- responseMsg
	}
}

// BroadcastEvent sends an event to all connected clients
// This can be used for server-pushed events
func (w *WebSocketTransport) DispatchWailsEvent(event *application.CustomEvent) {
	msg := &WebSocketMessage{
		Type:  "event",
		Event: event,
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	for _, channel := range w.clients {
		channel <- msg
	}
}
