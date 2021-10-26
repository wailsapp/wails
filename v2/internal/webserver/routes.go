package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (w *WebServer) setUpRoutes() error {
	// Handle Index + assets
	http.HandleFunc("/", w.serveAssets)

	// Handle Bindings
	http.HandleFunc("/bindings.js", w.serveBindings)

	// Handle Wails
	http.HandleFunc("/wails.js", w.serveWails)

	// Handle websocket connection
	http.HandleFunc("/ws", w.websocketConnection)

	return nil
}

func (w *WebServer) serveAssets(resp http.ResponseWriter, req *http.Request) {
	fileserver := http.FileServer(w.assets)
	fileserver.ServeHTTP(resp, req)
}

func (w *WebServer) serveBindings(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("content-type", "application/javascript")
	bindings, err := w.bindings.ToJSON()
	if err != nil {
		w.logger.Error("Failed to convert bindings to JSON: %v", err)
		fmt.Fprintf(resp, "")
		return
	}
	w.logger.Debug("Sending bindings to webclient: %v", bindings)
	b, _ := json.Marshal(bindings)
	fmt.Fprintf(resp, fmt.Sprintf("window.wailsbindings=%s; window.SetBindings(window.wailsbindings);", b))
}

func (w *WebServer) serveWails(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("content-type", "application/javascript")
	if data, err := w.assets.String("/wails.js"); err == nil {
		fmt.Fprintf(resp, strings.Replace(data, ":8080", fmt.Sprintf(":%d", w.port), 1))
	}
}
