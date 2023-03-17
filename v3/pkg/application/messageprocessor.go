package application

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// TODO maybe we could use a new struct that has the targetWindow as an attribute so we could get rid of passing the targetWindow
// as parameter through every function call.

type MessageProcessor struct {
	pluginManager *PluginManager
}

func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{}
}

func (m *MessageProcessor) httpError(rw http.ResponseWriter, message string, args ...any) {
	m.Error(message, args...)
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte(fmt.Sprintf(message, args...)))
}

func (m *MessageProcessor) HandleRuntimeCall(rw http.ResponseWriter, r *http.Request) {
	// Read "method" from query string
	method := r.URL.Query().Get("method")
	if method == "" {
		m.httpError(rw, "No method specified")
		return
	}
	splitMethod := strings.Split(method, ".")
	if len(splitMethod) != 2 {
		m.httpError(rw, "Invalid method format")
		return
	}
	// Get the object
	object := splitMethod[0]
	// Get the method
	method = splitMethod[1]

	params := QueryParams(r.URL.Query())

	var windowID uint
	if hWindowID := r.Header.Get(webViewRequestHeaderWindowId); hWindowID != "" {
		// Get windowID out of the request header
		wID, err := strconv.ParseUint(hWindowID, 10, 64)
		if err != nil {
			m.Error("Window ID '%s' not parsable: %s", hWindowID, err)
			return
		}

		windowID = uint(wID)
	}

	if qWindowID := params.UInt("windowID"); qWindowID != nil {
		// Get windowID out of the query parameters if provided
		windowID = *qWindowID
	}

	targetWindow := globalApplication.getWindowForID(windowID)
	if targetWindow == nil {
		m.Error("Window ID %s not found", windowID)
		return
	}

	switch object {
	case "window":
		m.processWindowMethod(method, rw, r, targetWindow, params)
	case "clipboard":
		m.processClipboardMethod(method, rw, r, targetWindow, params)
	case "dialog":
		m.processDialogMethod(method, rw, r, targetWindow, params)
	case "events":
		m.processEventsMethod(method, rw, r, targetWindow, params)
	case "application":
		m.processApplicationMethod(method, rw, r, targetWindow, params)
	case "log":
		m.processLogMethod(method, rw, r, targetWindow, params)
	case "contextmenu":
		m.processContextMenuMethod(method, rw, r, targetWindow, params)
	case "screens":
		m.processScreensMethod(method, rw, r, targetWindow, params)
	case "call":
		m.processCallMethod(method, rw, r, targetWindow, params)
	default:
		m.httpError(rw, "Unknown runtime call: %s", object)
	}

}

func (m *MessageProcessor) Error(message string, args ...any) {
	fmt.Printf("[MessageProcessor] Error: "+message, args...)
}

func (m *MessageProcessor) Info(message string, args ...any) {
	fmt.Printf("[MessageProcessor] Info: "+message, args...)
}

func (m *MessageProcessor) json(rw http.ResponseWriter, data any) {
	rw.Header().Set("Content-Type", "application/json")
	// convert data to json
	var jsonPayload = []byte("{}")
	var err error
	if data != nil {
		jsonPayload, err = jsoniter.Marshal(data)
		if err != nil {
			m.Error("Unable to convert data to JSON. Please report this to the Wails team! Error: %s", err)
			return
		}
	}
	_, err = rw.Write(jsonPayload)
	if err != nil {
		m.Error("Unable to write json payload. Please report this to the Wails team! Error: %s", err)
		return
	}
	m.ok(rw)
}

func (m *MessageProcessor) text(rw http.ResponseWriter, data string) {
	_, err := rw.Write([]byte(data))
	if err != nil {
		m.Error("Unable to write json payload. Please report this to the Wails team! Error: %s", err)
		return
	}
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
}

func (m *MessageProcessor) ok(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusOK)
}
