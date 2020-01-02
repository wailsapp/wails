package renderer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/dchest/htmlmin"
	"github.com/gorilla/websocket"
	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails/lib/interfaces"
	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/lib/messages"
)

type messageType int

const (
	jsMessage messageType = iota
	cssMessage
	htmlMessage
	notifyMessage
	bindingMessage
	callbackMessage
	wailsRuntimeMessage
)

func (m messageType) toString() string {
	return [...]string{"j", "s", "h", "n", "b", "c", "w"}[m]
}

// Bridge is a backend that opens a local web server
// and renders the files over a websocket
type Bridge struct {
	// Common
	log          *logger.CustomLogger
	ipcManager   interfaces.IPCManager
	appConfig    interfaces.AppConfig
	eventManager interfaces.EventManager
	bindingCache []string

	// Bridge specific
	initialisationJS []string
	server           *http.Server
	theConnection    *websocket.Conn

	// Mutex for writing to the socket
	lock sync.Mutex
}

// Initialise the Bridge Renderer
func (h *Bridge) Initialise(appConfig interfaces.AppConfig, ipcManager interfaces.IPCManager, eventManager interfaces.EventManager) error {
	h.ipcManager = ipcManager
	h.appConfig = appConfig
	h.eventManager = eventManager
	ipcManager.BindRenderer(h)
	h.log = logger.NewCustomLogger("Bridge")
	return nil
}

func (h *Bridge) evalJS(js string, mtype messageType) error {

	message := mtype.toString() + js

	if h.theConnection == nil {
		h.initialisationJS = append(h.initialisationJS, message)
	} else {
		// Prepend message type to message
		h.sendMessage(h.theConnection, message)
	}

	return nil
}

// EnableConsole not needed for bridge!
func (h *Bridge) EnableConsole() {
}

func (h *Bridge) injectCSS(css string) {
	// Minify css to overcome issues in the browser with carriage returns
	minified, err := htmlmin.Minify([]byte(css), &htmlmin.Options{
		MinifyStyles: true,
	})
	if err != nil {
		h.log.Fatal("Unable to minify CSS: " + css)
	}
	minifiedCSS := string(minified)
	minifiedCSS = strings.Replace(minifiedCSS, "\\", "\\\\", -1)
	minifiedCSS = strings.Replace(minifiedCSS, "'", "\\'", -1)
	minifiedCSS = strings.Replace(minifiedCSS, "\n", " ", -1)
	inject := fmt.Sprintf("wails._.InjectCSS('%s')", minifiedCSS)
	h.evalJS(inject, cssMessage)
}

func (h *Bridge) wsBridgeHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	h.theConnection = conn
	h.log.Infof("Connection from frontend accepted [%p].", h.theConnection)
	conn.SetCloseHandler(func(int, string) error {
		h.log.Infof("Connection dropped [%p].", h.theConnection)
		h.theConnection = nil
		return nil
	})
	go h.start(conn)
}

func (h *Bridge) sendMessage(conn *websocket.Conn, msg string) {

	h.lock.Lock()
	defer h.lock.Unlock()

	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		h.log.Error(err.Error())
	}
}

func (h *Bridge) start(conn *websocket.Conn) {

	// set external.invoke
	h.log.Infof("Connected to frontend.")

	wailsRuntime := mewn.String("../../runtime/assets/wails.js")
	h.evalJS(wailsRuntime, wailsRuntimeMessage)

	// Inject bindings
	for _, binding := range h.bindingCache {
		h.evalJS(binding, bindingMessage)
	}

	// Emit that everything is loaded and ready
	h.eventManager.Emit("wails:ready")

	for {
		messageType, buffer, err := conn.ReadMessage()
		if messageType == -1 {
			return
		}
		if err != nil {
			h.log.Errorf("Error reading message: ", err)
			continue
		}

		h.log.Debugf("Got message: %#v\n", string(buffer))

		h.ipcManager.Dispatch(string(buffer), h.Callback)
	}
}

// Run the app in Bridge mode!
func (h *Bridge) Run() error {
	h.server = &http.Server{Addr: ":34115"}
	http.HandleFunc("/bridge", h.wsBridgeHandler)

	h.log.Info("Bridge mode started.")
	h.log.Info("The frontend will connect automatically.")

	err := h.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		h.log.Fatal(err.Error())
	}
	return err
}

// NewBinding creates a new binding with the frontend
func (h *Bridge) NewBinding(methodName string) error {
	h.bindingCache = append(h.bindingCache, methodName)
	return nil
}

// SelectFile is unsupported for Bridge but required
// for the Renderer interface
func (h *Bridge) SelectFile() string {
	h.log.Warn("SelectFile() unsupported in bridge mode")
	return ""
}

// SelectDirectory is unsupported for Bridge but required
// for the Renderer interface
func (h *Bridge) SelectDirectory() string {
	h.log.Warn("SelectDirectory() unsupported in bridge mode")
	return ""
}

// SelectSaveFile is unsupported for Bridge but required
// for the Renderer interface
func (h *Bridge) SelectSaveFile() string {
	h.log.Warn("SelectSaveFile() unsupported in bridge mode")
	return ""
}

// Callback sends a callback to the frontend
func (h *Bridge) Callback(data string) error {
	return h.evalJS(data, callbackMessage)
}

// NotifyEvent notifies the frontend of an event
func (h *Bridge) NotifyEvent(event *messages.EventData) error {

	// Look out! Nils about!
	var err error
	if event == nil {
		err = fmt.Errorf("Sent nil event to renderer.webViewRenderer")
		h.log.Error(err.Error())
		return err
	}

	// Default data is a blank array
	data := []byte("[]")

	// Process event data
	if event.Data != nil {
		// Marshall the data
		data, err = json.Marshal(event.Data)
		if err != nil {
			h.log.Errorf("Cannot unmarshall JSON data in event: %s ", err.Error())
			return err
		}
	}

	message := fmt.Sprintf("window.wails._.Notify('%s','%s')", event.Name, data)
	return h.evalJS(message, notifyMessage)
}

// SetColour is unsupported for Bridge but required
// for the Renderer interface
func (h *Bridge) SetColour(colour string) error {
	h.log.WarnFields("SetColour ignored for Bridge more", logger.Fields{"col": colour})
	return nil
}

// Fullscreen is unsupported for Bridge but required
// for the Renderer interface
func (h *Bridge) Fullscreen() {
	h.log.Warn("Fullscreen() unsupported in bridge mode")
}

// UnFullscreen is unsupported for Bridge but required
// for the Renderer interface
func (h *Bridge) UnFullscreen() {
	h.log.Warn("UnFullscreen() unsupported in bridge mode")
}

// SetTitle is currently unsupported for Bridge but required
// for the Renderer interface
func (h *Bridge) SetTitle(title string) {
	h.log.WarnFields("SetTitle() unsupported in bridge mode", logger.Fields{"title": title})
}

// Close is unsupported for Bridge but required
// for the Renderer interface
func (h *Bridge) Close() {
	h.log.Debug("Shutting down")
	err := h.server.Close()
	if err != nil {
		h.log.Errorf(err.Error())
	}
}
