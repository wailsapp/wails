package application

import (
	"fmt"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type MessageProcessor struct {
	window *WebviewWindow
}

func NewMessageProcessor(w *WebviewWindow) *MessageProcessor {
	return &MessageProcessor{
		window: w,
	}
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

	var targetWindow = m.window
	windowID := params.UInt("windowID")
	if windowID != nil {
		// Get window for ID
		targetWindow = globalApplication.getWindowForID(*windowID)
		if targetWindow == nil {
			m.Error("Window ID %s not found", *windowID)
			return
		}
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

func (m *MessageProcessor) ProcessMessage(message string) {
	m.Info("ProcessMessage from front end:", message)
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
