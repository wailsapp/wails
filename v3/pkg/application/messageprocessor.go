package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
)

// TODO maybe we could use a new struct that has the targetWindow as an attribute so we could get rid of passing the targetWindow
// as parameter through every function call.

const (
	callRequest        = 0
	clipboardRequest   = 1
	applicationRequest = 2
	eventsRequest      = 3
	contextMenuRequest = 4
	dialogRequest      = 5
	windowRequest      = 6
	screensRequest     = 7
	systemRequest      = 8
	browserRequest     = 9
	cancelCallRequesst = 10
)

type MessageProcessor struct {
	logger *slog.Logger

	runningCalls map[string]context.CancelFunc
	l            sync.Mutex
}

func NewMessageProcessor(logger *slog.Logger) *MessageProcessor {
	return &MessageProcessor{
		logger:       logger,
		runningCalls: map[string]context.CancelFunc{},
	}
}

func (m *MessageProcessor) httpError(rw http.ResponseWriter, message string, args ...any) {
	m.Error(message, args...)
	rw.WriteHeader(http.StatusBadRequest)
	_, err := rw.Write([]byte(fmt.Sprintf(message, args...)))
	if err != nil {
		m.Error("Unable to write error message: %s", err)
	}
}

func (m *MessageProcessor) getTargetWindow(r *http.Request) (Window, string) {
	windowName := r.Header.Get(webViewRequestHeaderWindowName)
	if windowName != "" {
		return globalApplication.GetWindowByName(windowName), windowName
	}
	windowID := r.Header.Get(webViewRequestHeaderWindowId)
	if windowID == "" {
		return nil, windowID
	}
	wID, err := strconv.ParseUint(windowID, 10, 64)
	if err != nil {
		m.Error("Window ID '%s' not parsable: %s", windowID, err)
		return nil, windowID
	}
	targetWindow := globalApplication.getWindowForID(uint(wID))
	if targetWindow == nil {
		m.Error("Window ID %d not found", wID)
		return nil, windowID
	}
	return targetWindow, windowID
}

func (m *MessageProcessor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	object := r.URL.Query().Get("object")
	if object == "" {
		m.httpError(rw, "Invalid runtime call")
		return
	}

	m.HandleRuntimeCallWithIDs(rw, r)
}

func (m *MessageProcessor) HandleRuntimeCallWithIDs(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if handlePanic() {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}()
	object, err := strconv.Atoi(r.URL.Query().Get("object"))
	if err != nil {
		m.httpError(rw, "Error decoding object value: "+err.Error())
		return
	}
	method, err := strconv.Atoi(r.URL.Query().Get("method"))
	if err != nil {
		m.httpError(rw, "Error decoding method value: "+err.Error())
		return
	}
	params := QueryParams(r.URL.Query())

	targetWindow, nameOrID := m.getTargetWindow(r)
	if targetWindow == nil {
		m.httpError(rw, fmt.Sprintf("Window '%s' not found", nameOrID))
		return
	}

	switch object {
	case windowRequest:
		m.processWindowMethod(method, rw, r, targetWindow, params)
	case clipboardRequest:
		m.processClipboardMethod(method, rw, r, targetWindow, params)
	case dialogRequest:
		m.processDialogMethod(method, rw, r, targetWindow, params)
	case eventsRequest:
		m.processEventsMethod(method, rw, r, targetWindow, params)
	case applicationRequest:
		m.processApplicationMethod(method, rw, r, targetWindow, params)
	case contextMenuRequest:
		m.processContextMenuMethod(method, rw, r, targetWindow, params)
	case screensRequest:
		m.processScreensMethod(method, rw, r, targetWindow, params)
	case callRequest:
		m.processCallMethod(method, rw, r, targetWindow, params)
	case systemRequest:
		m.processSystemMethod(method, rw, r, targetWindow, params)
	case browserRequest:
		m.processBrowserMethod(method, rw, r, targetWindow, params)
	case cancelCallRequesst:
		m.processCallCancelMethod(method, rw, r, targetWindow, params)
	default:
		m.httpError(rw, "Unknown runtime call: %d", object)
	}
}

func (m *MessageProcessor) Error(message string, args ...any) {
	m.logger.Error(message, args...)
}

func (m *MessageProcessor) Info(message string, args ...any) {
	m.logger.Info(message, args...)
}

func (m *MessageProcessor) json(rw http.ResponseWriter, data any) {
	rw.Header().Set("Content-Type", "application/json")
	// convert data to json
	var jsonPayload = []byte("{}")
	var err error
	if data != nil {
		jsonPayload, err = json.Marshal(data)
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
