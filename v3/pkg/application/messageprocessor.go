package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
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

func (m *MessageProcessor) httpError(rw http.ResponseWriter, message string, err error) {
	m.Error(message, "error", err)
	rw.WriteHeader(http.StatusUnprocessableEntity)
	_, err = rw.Write([]byte(err.Error()))
	if err != nil {
		m.Error("Unable to write error response:", "error", err)
	}
}

func (m *MessageProcessor) getTargetWindow(r *http.Request) (Window, string) {
	windowName := r.Header.Get(webViewRequestHeaderWindowName)
	if windowName != "" {
		window, _ := globalApplication.Window.GetByName(windowName)
		return window, windowName
	}
	windowID := r.Header.Get(webViewRequestHeaderWindowId)
	if windowID == "" {
		return nil, windowID
	}
	wID, err := strconv.ParseUint(windowID, 10, 64)
	if err != nil {
		m.Error("Window ID not parsable:", "id", windowID, "error", err)
		return nil, windowID
	}
	// Check if wID is within the valid range for uint
	if wID > math.MaxUint32 {
		m.Error("Window ID out of range for uint:", "id", wID)
		return nil, windowID
	}
	targetWindow, _ := globalApplication.Window.GetByID(uint(wID))
	if targetWindow == nil {
		m.Error("Window ID not found:", "id", wID)
		return nil, windowID
	}
	return targetWindow, windowID
}

func (m *MessageProcessor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, err := m.decodeBody(r.Body)
	if err != nil {
		m.httpError(rw, "Invalid runtime call:", err)
		return
	}
	if string(body.Object) == "" {
		m.httpError(rw, "Invalid runtime call:", errors.New("missing object values"))
		return
	}

	m.HandleRuntimeCallWithIDs(rw, r, body)
}

func (m *MessageProcessor) HandleRuntimeCallWithIDs(rw http.ResponseWriter, r *http.Request, body *runtimeRequest) {
	defer func() {
		if handlePanic() {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}()

	if body == nil {
		m.httpError(rw, "Invalid runtime call:", errors.New("request body is nil"))
		return
	}

	object := body.Object
	method := body.Method

	targetWindow, nameOrID := m.getTargetWindow(r)
	if targetWindow == nil {
		m.httpError(rw, "Invalid runtime call:", fmt.Errorf("window '%s' not found", nameOrID))
		return
	}

	switch object {
	case windowRequest:
		m.processWindowMethod(method, rw, r, targetWindow, *body)
	case clipboardRequest:
		m.processClipboardMethod(method, rw, r, targetWindow, *body)
	case dialogRequest:
		m.processDialogMethod(method, rw, r, targetWindow, *body)
	case eventsRequest:
		m.processEventsMethod(method, rw, r, targetWindow, *body)
	case applicationRequest:
		m.processApplicationMethod(method, rw, r, targetWindow, *body)
	case contextMenuRequest:
		m.processContextMenuMethod(method, rw, r, targetWindow, *body)
	case screensRequest:
		m.processScreensMethod(method, rw, r, targetWindow, *body)
	case callRequest:
		m.processCallMethod(method, rw, r, targetWindow, *body)
	case systemRequest:
		m.processSystemMethod(method, rw, r, targetWindow, *body)
	case browserRequest:
		m.processBrowserMethod(method, rw, r, targetWindow, *body)
	case cancelCallRequesst:
		m.processCallCancelMethod(method, rw, r, targetWindow, *body)
	default:
		m.httpError(rw, "Invalid runtime call:", fmt.Errorf("unknown object %d", object))
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
			m.Error("Unable to convert data to JSON. Please report this to the Wails team!", "error", err)
			return
		}
	}
	_, err = rw.Write(jsonPayload)
	if err != nil {
		m.Error("Unable to write json payload. Please report this to the Wails team!", "error", err)
		return
	}
	m.ok(rw)
}

func (m *MessageProcessor) text(rw http.ResponseWriter, data string) {
	_, err := rw.Write([]byte(data))
	if err != nil {
		m.Error("Unable to write json payload. Please report this to the Wails team!", "error", err)
		return
	}
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
}

func (m *MessageProcessor) ok(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusOK)
}
