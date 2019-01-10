package wails

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dchest/htmlmin"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/websocket"
)

var headlessAssets = packr.NewBox("./assets/headless")
var defaultAssets = packr.NewBox("./assets/default")

// Window defines the main application window
// Default values in []
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
}

func (h *Headless) Initialise(appConfig *AppConfig, ipcManager *ipcManager, eventManager *eventManager) error {
	h.ipcManager = ipcManager
	h.appConfig = appConfig
	h.eventManager = eventManager
	h.log = newCustomLogger("Headless")
	return nil
}

func (h *Headless) evalJS(js string) error {
	if h.theConnection == nil {
		h.initialisationJS = append(h.initialisationJS, js)
	} else {
		h.sendMessage(h.theConnection, js)
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
	minifiedCSS = strings.Replace(minifiedCSS, "'", "\\'", -1)
	minifiedCSS = strings.Replace(minifiedCSS, "\n", " ", -1)
	inject := fmt.Sprintf("wails._.injectCSS('%s')", minifiedCSS)
	h.evalJS(inject)
}

func (h *Headless) rootHandler(w http.ResponseWriter, r *http.Request) {
	indexHTML := headlessAssets.String("index.html")
	fmt.Fprintf(w, "%s", indexHTML)
}

func (h *Headless) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	h.theConnection = conn
	h.log.Infof("Connection %p accepted.", h.theConnection)
	conn.SetCloseHandler(func(int, string) error {
		h.log.Infof("Connection %p dropped.", h.theConnection)
		h.theConnection = nil
		return nil
	})
	go h.start(conn)
}

func (h *Headless) sendMessage(conn *websocket.Conn, msg string) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		h.log.Error(err.Error())
	}
}

func (h *Headless) start(conn *websocket.Conn) {

	// set external.invoke
	h.log.Infof("Connected to frontend.")

	// If we are given an HTML fragment, load jquery
	// for the html() function
	if h.appConfig.isHTMLFragment {
		// Inject jquery
		jquery := defaultAssets.String("jquery.3.3.1.min.js")
		h.evalJS(jquery)
	}

	wailsRuntime := defaultAssets.String("wails.js")
	h.evalJS(wailsRuntime)

	// Inject the initial JS
	for _, js := range h.initialisationJS {
		h.sendMessage(conn, js)
	}

	// Inject bindings
	for _, binding := range h.bindingCache {
		h.evalJS(binding)
	}

	// Inject Framework
	if h.frameworkJS != "" {
		h.evalJS(h.frameworkJS)
	}
	if h.frameworkCSS != "" {
		h.injectCSS(h.frameworkCSS)
	}

	// If given an HMTL fragment, mount it on #app
	// Otherwise, replace the html tag
	var injectHTML string
	if h.appConfig.isHTMLFragment {
		injectHTML = fmt.Sprintf("$('#app').html('%s')", h.appConfig.HTML)
	} else {
		injectHTML = fmt.Sprintf("$('html').html('%s')", h.appConfig.HTML)
	}
	h.evalJS(injectHTML)

	// Inject user CSS
	if h.appConfig.CSS != "" {
		outputCSS := fmt.Sprintf("%.45s", h.appConfig.CSS)
		if len(outputCSS) > 45 {
			outputCSS += "..."
		}
		h.log.DebugFields("Inject User CSS", Fields{"css": outputCSS})
		h.injectCSS(h.appConfig.CSS)
	}

	// Inject all the CSS files that have been added
	for _, css := range h.cssCache {
		h.injectCSS(css)
	}

	// Inject all the JS files that have been added
	for _, js := range h.jsCache {
		h.evalJS(js)
	}

	// Inject user JS
	if h.appConfig.JS != "" {
		outputJS := fmt.Sprintf("%.45s", h.appConfig.JS)
		if len(outputJS) > 45 {
			outputJS += "..."
		}
		h.log.DebugFields("Inject User JS", Fields{"js": outputJS})
		h.evalJS(h.appConfig.JS)
	}

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
	http.HandleFunc("/ws", h.wsHandler)
	http.HandleFunc("/", h.rootHandler)

	h.log.Info("Started on port 34115")
	h.log.Info("Application running at http://localhost:34115")

	err := h.server.ListenAndServe()
	if err != nil {
		h.log.Fatal(err.Error())
	}
	return err
}

func (h *Headless) NewBinding(methodName string) error {
	objectCode := fmt.Sprintf("window.wails._.newBinding(`%s`);", methodName)
	h.bindingCache = append(h.bindingCache, objectCode)
	return nil
}

func (h *Headless) InjectFramework(js, css string) {
	h.frameworkJS = js
	h.frameworkCSS = css
}

func (h *Headless) SelectFile() string {
	h.log.Error("SelectFile() unsupported in headless mode")
	return ""
}
func (h *Headless) SelectDirectory() string {
	h.log.Error("SelectDirectory() unsupported in headless mode")
	return ""
}
func (h *Headless) SelectSaveFile() string {
	h.log.Error("SelectSaveFile() unsupported in headless mode")
	return ""
}
func (h *Headless) AddJSList(jsCache []string) {
	h.jsCache = jsCache
}
func (h *Headless) AddCSSList(cssCache []string) {
	h.cssCache = cssCache
}

// Callback sends a callback to the frontend
func (h *Headless) Callback(data string) error {
	callbackCMD := fmt.Sprintf("window.wails._.callback('%s');", data)
	return h.evalJS(callbackCMD)
}

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

	message := fmt.Sprintf("wails._.notify('%s','%s')", event.Name, data)
	return h.evalJS(message)
}

func (h *Headless) SetColour(colour string) error {
	h.log.WarnFields("SetColour ignored for headless more", Fields{"col": colour})
	return nil
}
func (h *Headless) Fullscreen() {
	h.log.Warn("Fullscreen() unsupported in headless mode")
}
func (h *Headless) UnFullscreen() {
	h.log.Warn("UnFullscreen() unsupported in headless mode")
}
func (h *Headless) SetTitle(title string) {
	h.log.WarnFields("SetTitle() unsupported in headless mode", Fields{"title": title})
}
func (h *Headless) Close() {
	h.log.Warn("Close() unsupported in headless mode")
}
