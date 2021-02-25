package renderer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
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
	server *http.Server

	lock     sync.Mutex
	sessions map[string]*session
}

// Initialise the Bridge Renderer
func (h *Bridge) Initialise(appConfig interfaces.AppConfig, ipcManager interfaces.IPCManager, eventManager interfaces.EventManager) error {
	h.sessions = map[string]*session{}
	h.ipcManager = ipcManager
	h.appConfig = appConfig
	h.eventManager = eventManager
	ipcManager.BindRenderer(h)
	h.log = logger.NewCustomLogger("Bridge")
	return nil
}

func (h *Bridge) wsBridgeHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	h.log.Infof("Connection from frontend accepted [%s].", conn.RemoteAddr().String())
	h.startSession(conn)
}

func (h *Bridge) startSession(conn *websocket.Conn) {
	s := newSession(conn,
		h.bindingCache,
		h.ipcManager,
		logger.NewCustomLogger("BridgeSession"),
		h.eventManager)

	conn.SetCloseHandler(func(int, string) error {
		h.log.Infof("Connection dropped [%s].", s.Identifier())
		h.eventManager.Emit("wails:bridge:session:closed", s.Identifier())
		h.lock.Lock()
		defer h.lock.Unlock()
		delete(h.sessions, s.Identifier())
		return nil
	})

	h.lock.Lock()
	defer h.lock.Unlock()
	go s.start(len(h.sessions) == 0)
	h.sessions[s.Identifier()] = s
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
func (h *Bridge) SelectFile(title string, filter string) string {
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
func (h *Bridge) SelectSaveFile(title string, filter string) string {
	h.log.Warn("SelectSaveFile() unsupported in bridge mode")
	return ""
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
			h.log.Errorf("Cannot marshal JSON data in event: %s ", err.Error())
			return err
		}
	}

	// Double encode data to ensure everything is escaped correctly.
	data, err = json.Marshal(string(data))
	if err != nil {
		h.log.Errorf("Cannot marshal JSON data in event: %s ", err.Error())
		return err
	}

	message := "window.wails._.Notify('" + event.Name + "'," + string(data) + ")"
	dead := []*session{}
	for _, session := range h.sessions {
		err := session.evalJS(message, notifyMessage)
		if err != nil {
			h.log.Debugf("Failed to send message to %s - Removing listener : %v", session.Identifier(), err)
			h.log.Infof("Connection from [%v] unresponsive - dropping", session.Identifier())
			dead = append(dead, session)
		}
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	for _, session := range dead {
		delete(h.sessions, session.Identifier())
	}

	return nil
}

// SetColour is unsupported for Bridge but required
// for the Renderer interface
func (h *Bridge) SetColour(colour string) error {
	h.log.WarnFields("SetColour ignored for Bridge more", logger.Fields{"col": colour})
	return nil
}

// SetMinSize is unsupported for Bridge but required
// for the Renderer interface
func (h *Bridge) SetMinSize(width, height int) {
	h.log.Warn("SetMinSize() unsupported in bridge mode")
}

// SetMaxSize is unsupported for Bridge but required
// for the Renderer interface
func (h *Bridge) SetMaxSize(width, height int) {
	h.log.Warn("SetMaxSize() unsupported in bridge mode")
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
	for _, session := range h.sessions {
		session.Shutdown()
	}
	err := h.server.Close()
	if err != nil {
		h.log.Errorf(err.Error())
	}
}
