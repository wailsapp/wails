package wails

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/dchest/htmlmin"
	"github.com/gorilla/websocket"
	"github.com/leaanthony/mewn"
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

// Headless is a backend that opens a local web server
// and renders the files over a websocket
type Headless struct {
	// Common
	log          *CustomLogger
	ipcManager   *ipcManager
	appConfig    *AppConfig
	eventManager *eventManager
	bindingCache []string
	frameworkJS  string
	frameworkCSS string
	jsCache      []string
	cssCache     []string

	// Headless specific
	initialisationJS []string
	server           *http.Server
	theConnection    *websocket.Conn

	// Mutex for writing to the socket
	lock sync.Mutex
}

// Initialise the Headless Renderer
func (h *Headless) Initialise(appConfig *AppConfig, ipcManager *ipcManager, eventManager *eventManager) error {
	h.ipcManager = ipcManager
	h.appConfig = appConfig
	h.eventManager = eventManager
	ipcManager.bindRenderer(h)
	h.log = newCustomLogger("Bridge")
	return nil
}

func (h *Headless) evalJS(js string, mtype messageType) error {

	message := mtype.toString() + js

	if h.theConnection == nil {
		h.initialisationJS = append(h.initialisationJS, message)
	} else {
		// Prepend message type to message
		h.sendMessage(h.theConnection, message)
	}

	return nil
}

func (h *Headless) injectCSS(css string) {
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
	inject := fmt.Sprintf("wails._.injectCSS('%s')", minifiedCSS)
	h.evalJS(inject, cssMessage)
}

func (h *Headless) wsBridgeHandler(w http.ResponseWriter, r *http.Request) {
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

func (h *Headless) sendMessage(conn *websocket.Conn, msg string) {

	h.lock.Lock()
	defer h.lock.Unlock()

	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		h.log.Error(err.Error())
	}
}

func (h *Headless) start(conn *websocket.Conn) {

	// set external.invoke
	h.log.Infof("Connected to frontend.")

	wailsRuntime := mewn.String("./wailsruntimeassets/default/wails.min.js")
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

		h.ipcManager.Dispatch(string(buffer))
	}
}

// Run the app in headless mode!
func (h *Headless) Run() error {
	h.server = &http.Server{Addr: ":34115"}
	http.HandleFunc("/bridge", h.wsBridgeHandler)

	h.log.Info("Bridge mode started.")
	h.log.Info("The frontend will connect automatically.")

	err := h.server.ListenAndServe()
	if err != nil {
		h.log.Fatal(err.Error())
	}
	return err
}

// NewBinding creates a new binding with the frontend
func (h *Headless) NewBinding(methodName string) error {
	h.bindingCache = append(h.bindingCache, methodName)
	return nil
}

// SelectFile is unsupported for Headless but required
// for the Renderer interface
func (h *Headless) SelectFile() string {
	h.log.Error("SelectFile() unsupported in bridge mode")
	return ""
}

// SelectDirectory is unsupported for Headless but required
// for the Renderer interface
func (h *Headless) SelectDirectory() string {
	h.log.Error("SelectDirectory() unsupported in bridge mode")
	return ""
}

// SelectSaveFile is unsupported for Headless but required
// for the Renderer interface
func (h *Headless) SelectSaveFile() string {
	h.log.Error("SelectSaveFile() unsupported in bridge mode")
	return ""
}

// AddJSList adds a slice of JS strings to the list of js
// files injected at startup
func (h *Headless) AddJSList(jsCache []string) {
	h.jsCache = jsCache
}

// AddCSSList adds a slice of CSS strings to the list of css
// files injected at startup
func (h *Headless) AddCSSList(cssCache []string) {
	h.cssCache = cssCache
}

// Callback sends a callback to the frontend
func (h *Headless) Callback(data string) error {
	return h.evalJS(data, callbackMessage)
}

// NotifyEvent notifies the frontend of an event
func (h *Headless) NotifyEvent(event *eventData) error {

	// Look out! Nils about!
	var err error
	if event == nil {
		err = fmt.Errorf("Sent nil event to renderer.webViewRenderer")
		logger.Error(err)
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

	message := fmt.Sprintf("window.wails._.notify('%s','%s')", event.Name, data)
	return h.evalJS(message, notifyMessage)
}

// SetColour is unsupported for Headless but required
// for the Renderer interface
func (h *Headless) SetColour(colour string) error {
	h.log.WarnFields("SetColour ignored for headless more", Fields{"col": colour})
	return nil
}

// Fullscreen is unsupported for Headless but required
// for the Renderer interface
func (h *Headless) Fullscreen() {
	h.log.Warn("Fullscreen() unsupported in bridge mode")
}

// UnFullscreen is unsupported for Headless but required
// for the Renderer interface
func (h *Headless) UnFullscreen() {
	h.log.Warn("UnFullscreen() unsupported in bridge mode")
}

// SetTitle is currently unsupported for Headless but required
// for the Renderer interface
func (h *Headless) SetTitle(title string) {
	h.log.WarnFields("SetTitle() unsupported in bridge mode", Fields{"title": title})
}

// Close is unsupported for Headless but required
// for the Renderer interface
func (h *Headless) Close() {
	h.log.Warn("Close() unsupported in bridge mode")
}
