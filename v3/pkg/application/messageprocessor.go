package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"math"
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
	// Set CORS headers first
	m.setCORSHeaders(rw, r)
	
	// Handle preflight requests
	if r.Method == "OPTIONS" {
		rw.WriteHeader(http.StatusOK)
		return
	}
	
	object := r.URL.Query().Get("object")
	if object == "" {
		m.httpError(rw, "Invalid runtime call:", errors.New("missing object value"))
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
		m.httpError(rw, "Invalid runtime call:", fmt.Errorf("error decoding object value: %w", err))
		return
	}
	method, err := strconv.Atoi(r.URL.Query().Get("method"))
	if err != nil {
		m.httpError(rw, "Invalid runtime call:", fmt.Errorf("error decoding method value: %w", err))
		return
	}
	params := QueryParams(r.URL.Query())

	targetWindow, nameOrID := m.getTargetWindow(r)
	if targetWindow == nil {
		m.httpError(rw, "Invalid runtime call:", fmt.Errorf("window '%s' not found", nameOrID))
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
	
	config := globalApplication.getBindingConfig()
	
	// Use streaming encoder for better memory efficiency when enabled
	if config.EnableStreaming {
		encoder := json.NewEncoder(rw)
		err := encoder.Encode(data)
		if err != nil {
			m.Error("Unable to encode JSON response:", "error", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		return
	}
	
	// Traditional approach for compatibility
	var jsonPayload = []byte("{}")
	var err error
	if data != nil {
		jsonPayload, err = json.Marshal(data)
		if err != nil {
			m.Error("Unable to convert data to JSON. Please report this to the Wails team!", "error", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	
	_, err = rw.Write(jsonPayload)
	if err != nil {
		m.Error("Unable to write json payload. Please report this to the Wails team!", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
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

// setCORSHeaders sets CORS headers based on the application configuration
func (m *MessageProcessor) setCORSHeaders(rw http.ResponseWriter, r *http.Request) {
	config := globalApplication.getBindingConfig()
	if !config.CORS.Enabled {
		return
	}
	
	origin := r.Header.Get("Origin")
	if origin == "" {
		return // Not a cross-origin request
	}
	
	if m.isOriginAllowed(origin, config.CORS.AllowedOrigins) {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", strings.Join(config.CORS.AllowedMethods, ", "))
		rw.Header().Set("Access-Control-Allow-Headers", strings.Join(config.CORS.AllowedHeaders, ", "))
		rw.Header().Set("Access-Control-Max-Age", strconv.Itoa(int(config.CORS.MaxAge.Seconds())))
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

// isOriginAllowed checks if the given origin is allowed based on the configuration
func (m *MessageProcessor) isOriginAllowed(origin string, allowedOrigins []string) bool {
	// Development mode - allow all origins if no origins specified
	if globalApplication.isDebugMode && len(allowedOrigins) == 0 {
		return true
	}
	
	// Production mode - check whitelist
	for _, allowed := range allowedOrigins {
		if matched, _ := filepath.Match(allowed, origin); matched {
			return true
		}
	}
	
	return false
}

// handleBindingError returns appropriate HTTP status codes based on error type
func (m *MessageProcessor) handleBindingError(rw http.ResponseWriter, err error) {
	var callErr *CallError
	if errors.As(err, &callErr) {
		switch callErr.Kind {
		case ReferenceError:
			rw.WriteHeader(http.StatusNotFound)
		case TypeError:
			rw.WriteHeader(http.StatusBadRequest)
		case RuntimeError:
			rw.WriteHeader(http.StatusInternalServerError)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
		
		// Return structured error response
		errorResponse := map[string]interface{}{
			"error":     callErr.Message,
			"kind":      callErr.Kind,
			"cause":     callErr.Cause,
			"timestamp": callErr.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		}
		
		rw.Header().Set("Content-Type", "application/json")
		if jsonErr := json.NewEncoder(rw).Encode(errorResponse); jsonErr != nil {
			m.Error("Unable to encode error response:", "error", jsonErr)
		}
		return
	}
	
	// Handle context timeout
	if errors.Is(err, context.DeadlineExceeded) {
		rw.WriteHeader(http.StatusRequestTimeout)
		rw.Header().Set("Content-Type", "application/json")
		errorResponse := map[string]interface{}{
			"error":     "Binding execution timeout",
			"kind":      "TimeoutError",
			"timestamp": time.Now().Format("2006-01-02T15:04:05Z07:00"),
		}
		if jsonErr := json.NewEncoder(rw).Encode(errorResponse); jsonErr != nil {
			m.Error("Unable to encode timeout error response:", "error", jsonErr)
		}
		return
	}
	
	// Default error handling
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]interface{}{
		"error":     err.Error(),
		"kind":      "RuntimeError",
		"timestamp": time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}
	if jsonErr := json.NewEncoder(rw).Encode(errorResponse); jsonErr != nil {
		m.Error("Unable to encode error response:", "error", jsonErr)
	}
}
