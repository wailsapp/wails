package webserver

import (
	"context"
	"net/http"
	"strings"

	"github.com/wailsapp/wails/v2/internal/logger"
	ws "nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// WebClient represents an individual web session
type WebClient struct {
	conn       *ws.Conn
	identifier string
	logger     *logger.Logger
	running    bool
}

// Quit terminates the webclient session
func (wc *WebClient) Quit() {
	wc.running = false
}

// NotifyEvent sends the event
func (wc *WebClient) NotifyEvent(message string) {
	wc.SendMessage("E" + message)
}

// CallResult sends the result of the Go function call back to the
// originator in the frontend
func (wc *WebClient) CallResult(message string) {
	wc.SendMessage("R" + message)
}

// SaveFileDialog is a noop in the webclient
func (wc *WebClient) SaveFileDialog(title string) string {
	return ""
}

// OpenFileDialog is a noop in the webclient
func (wc *WebClient) OpenFileDialog(title string) string {
	return ""
}

// OpenDirectoryDialog is a noop in the webclient
func (wc *WebClient) OpenDirectoryDialog(title string) string {
	return ""
}

// WindowSetTitle is a noop in the webclient
func (wc *WebClient) WindowSetTitle(title string) {}

// WindowFullscreen is a noop in the webclient
func (wc *WebClient) WindowFullscreen() {}

// WindowUnFullscreen is a noop in the webclient
func (wc *WebClient) WindowUnFullscreen() {}

// WindowSetColour is a noop in the webclient
func (wc *WebClient) WindowSetColour(colour int) {
}

// Run processes messages from the remote webclient
func (wc *WebClient) Run(w *WebServer) {
	dispatcher := w.dispatcher.RegisterClient(wc)
	defer w.dispatcher.RemoveClient(dispatcher)
	defer w.unregisterClient(wc.identifier)

	for wc.running {
		var v interface{}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if err := wsjson.Read(ctx, wc.conn, &v); err != nil {
			if ws.CloseStatus(err) == ws.StatusNormalClosure || ws.CloseStatus(err) == ws.StatusGoingAway {
				break
			}
			if ws.CloseStatus(err) != -1 {
				w.logger.Debug("Connection error: %s - %s", wc.identifier, err)
				break
			}

			if !strings.Contains(err.Error(), "status = Status") {
				w.logger.Debug("Error encountered on socket: %v", err)
				break
			}
		}
		dispatcher.DispatchMessage(v.(string))
	}

	wc.conn.Close(ws.StatusNormalClosure, "Goodbye")
	w.logger.Debug("Connection closed: %v", wc.identifier)
}

// SendMessage converts the string to a []byte and passes it to
// the connection's Writer to send to the remote client
// The Writer itself prevents multiple users at the same time.
func (wc *WebClient) SendMessage(message string) {
	wc.logger.Debug("WebClient.SendMessage() - %s", message)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wc.conn.Write(ctx, ws.MessageText, []byte(message))
}

// unregisterClient is called automatically by a WebClient session during termination
// so that it's registration can be removed
func (w *WebServer) unregisterClient(identifier string) {
	w.logger.Debug("Removing WebClient : %v", identifier)
	w.lock.Lock()
	delete(w.connections, identifier)
	w.lock.Unlock()
}

func (w *WebServer) websocketConnection(resp http.ResponseWriter, req *http.Request) {
	conn, err := ws.Accept(resp, req, nil)
	if err != nil {
		w.logger.Debug("Failed to upgrade websocket connection")
		return
	}
	wc := &WebClient{
		conn:       conn,
		identifier: req.RemoteAddr,
		logger:     w.logger,
		running:    true,
	}
	w.lock.Lock()
	w.connections[wc.identifier] = wc
	w.lock.Unlock()

	w.logger.Debug("Connection from: %v", wc.identifier)

	go wc.Run(w)

}
