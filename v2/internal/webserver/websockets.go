package webserver

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options/dialog"
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

func (wc *WebClient) SetTrayMenu(trayMenuJSON string) {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) UpdateTrayMenuLabel(trayMenuJSON string) {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) MessageDialog(dialogOptions *dialog.MessageDialog, callbackID string) {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) SetApplicationMenu(menuJSON string) {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) UpdateTrayMenu(trayMenuJSON string) {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) UpdateContextMenu(contextMenuJSON string) {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) OpenDialog(dialogOptions *dialog.OpenDialog, callbackID string) {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) SaveDialog(dialogOptions *dialog.SaveDialog, callbackID string) {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) WindowShow() {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) WindowHide() {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) WindowCenter() {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) WindowMaximise() {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) WindowUnmaximise() {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) WindowMinimise() {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) WindowUnminimise() {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) WindowPosition(x int, y int) {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) WindowSize(width int, height int) {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) DarkModeEnabled(callbackID string) {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) UpdateMenu(menu *menu.Menu) {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) UpdateTray(menu *menu.Menu) {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) UpdateTrayLabel(label string) {
	wc.logger.Info("Not implemented in server build")
}

func (wc *WebClient) UpdateTrayIcon(name string) {
	wc.logger.Info("Not implemented in server build")
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

	err := wc.conn.Close(ws.StatusNormalClosure, "Goodbye")
	if err != nil {
		w.logger.Error("Error encountered on socket: %v", err)
		return
	}
	w.logger.Debug("Connection closed: %v", wc.identifier)
}

// SendMessage converts the string to a []byte and passes it to
// the connection's Writer to send to the remote client
// The Writer itself prevents multiple users at the same time.
func (wc *WebClient) SendMessage(message string) {
	wc.logger.Debug("WebClient.SendMessage() - %s", message)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := wc.conn.Write(ctx, ws.MessageText, []byte(message))
	if err != nil {
		wc.logger.Error("Error encountered writing to webclient: %v", err)
	}
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
