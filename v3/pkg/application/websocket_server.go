//go:build server

package application

import (
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/wailsapp/wails/v3/internal/mailbox"
)

// globalBroadcaster holds a reference to the WebSocket broadcaster for server mode.
// Used by HTTP handler to look up BrowserWindows from runtime client IDs.
var globalBroadcaster *WebSocketBroadcaster

// GetBrowserWindow returns the BrowserWindow for a given runtime clientId.
// Only available in server mode. Returns nil if not found.
func GetBrowserWindow(clientId string) *BrowserWindow {
	if globalBroadcaster == nil {
		return nil
	}
	return globalBroadcaster.GetBrowserWindow(clientId)
}

// clientInfo holds information about a connected WebSocket client.
type clientInfo struct {
	conn   *websocket.Conn
	window *BrowserWindow
	events *mailbox.Mailbox[*CustomEvent]
}

// WebSocketBroadcaster manages WebSocket connections and broadcasts events to all connected clients.
// It implements WailsEventListener to receive events from the application.
type WebSocketBroadcaster struct {
	clients map[*websocket.Conn]*clientInfo
	windows map[string]*BrowserWindow // maps runtime clientId (nanoid) to BrowserWindow
	mu      sync.RWMutex
	app     *App
	nextID  atomic.Uint64
}

// NewWebSocketBroadcaster creates a new WebSocket broadcaster.
func NewWebSocketBroadcaster(app *App) *WebSocketBroadcaster {
	return &WebSocketBroadcaster{
		clients: make(map[*websocket.Conn]*clientInfo),
		windows: make(map[string]*BrowserWindow),
		app:     app,
	}
}

// GetBrowserWindow returns the BrowserWindow for a given runtime clientId.
// Returns nil if not found.
func (b *WebSocketBroadcaster) GetBrowserWindow(clientId string) *BrowserWindow {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.windows[clientId]
}

// ServeHTTP handles WebSocket upgrade requests.
func (b *WebSocketBroadcaster) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// Allow connections from any origin in server mode
		InsecureSkipVerify: true,
	})
	if err != nil {
		b.app.error("WebSocket accept error: %v", err)
		return
	}

	// Get runtime's clientId from query parameter
	runtimeClientID := r.URL.Query().Get("clientId")

	// Create BrowserWindow for this connection
	browserWindow := NewBrowserWindow(uint(b.nextID.Add(1)), runtimeClientID)

	// Store mapping if runtime clientId was provided
	if runtimeClientID != "" {
		b.mu.Lock()
		b.windows[runtimeClientID] = browserWindow
		b.mu.Unlock()
	}

	b.register(conn, browserWindow)
	defer b.unregister(conn, runtimeClientID)

	// Keep connection alive - read loop for ping/pong and detecting disconnects
	// Events from client to server are sent via HTTP, not WebSocket
	for {
		_, _, err := conn.Read(r.Context())
		if err != nil {
			break
		}
	}
}

// register adds a client connection with its BrowserWindow.
func (b *WebSocketBroadcaster) register(conn *websocket.Conn, window *BrowserWindow) {
	client := &clientInfo{
		conn:   conn,
		window: window,
	}
	client.events = mailbox.New(func(event *CustomEvent) {
		if err := wsjson.Write(b.app.ctx, client.conn, event); err != nil {
			b.app.debug("WebSocket write error", "error", err, "client", client.window.Name())
		}
	})

	b.mu.Lock()
	b.clients[conn] = client
	clientCount := len(b.clients)
	b.mu.Unlock()
	b.app.info("WebSocket client connected", "id", window.Name(), "clients", clientCount)
}

// unregister removes a client connection and its BrowserWindow.
func (b *WebSocketBroadcaster) unregister(conn *websocket.Conn, runtimeClientID string) {
	client, clientCount := b.removeClient(conn, runtimeClientID)
	conn.Close(websocket.StatusNormalClosure, "")
	if client != nil {
		b.app.info("WebSocket client disconnected", "id", client.window.Name(), "clients", clientCount)
	}
}

func (b *WebSocketBroadcaster) removeClient(conn *websocket.Conn, runtimeClientID string) (*clientInfo, int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	client := b.clients[conn]
	delete(b.clients, conn)
	if runtimeClientID != "" && client != nil && b.windows[runtimeClientID] == client.window {
		delete(b.windows, runtimeClientID)
	}

	return client, len(b.clients)
}

// DispatchWailsEvent implements WailsEventListener interface.
// It broadcasts the event to all connected WebSocket clients.
func (b *WebSocketBroadcaster) DispatchWailsEvent(event *CustomEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, client := range b.clients {
		client.events.Send(event)
	}
}
