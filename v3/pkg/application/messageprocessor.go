package application

import (
	"context"
	"log/slog"
	"slices"
	"strconv"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/errs"
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
	cancelCallRequest  = 10
)

var objectNames = map[int]string{
	callRequest:        "Call",
	clipboardRequest:   "Clipboard",
	applicationRequest: "Application",
	eventsRequest:      "Events",
	contextMenuRequest: "ContextMenu",
	dialogRequest:      "Dialog",
	windowRequest:      "Window",
	screensRequest:     "Screens",
	systemRequest:      "System",
	browserRequest:     "Browser",
	cancelCallRequest:  "CancellCall",
}

type RuntimeRequest struct {
	// Object identifies which Wails subsystem to call (Call=0, Clipboard=1, etc.)
	// See objectNames in runtime.ts
	Object int `json:"object"`

	// Method identifies which method within the object to call
	Method int `json:"method"`

	// Args contains the method arguments
	Args *Args `json:"args"`

	// WebviewWindowName identifies the source window by name (optional, sent via header x-wails-window-name)
	WebviewWindowName string `json:"webviewWindowName,omitempty"`

	// WebviewWindowID identifies the source window (optional, sent via header x-wails-window-id)
	WebviewWindowID uint32 `json:"webviewWindowId,omitempty"`

	// ClientID identifies the frontend client (sent via header x-wails-client-id)
	ClientID string `json:"clientId,omitempty"`
}

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

func (m *MessageProcessor) HandleRuntimeCallWithIDs(ctx context.Context, req *RuntimeRequest) (resp any, err error) {
	defer func() {
		if handlePanic() {
			// TODO: return panic error itself?
			err = errs.NewInvalidRuntimeCallErrorf("runtime panic detected!")
		}
	}()
	targetWindow, nameOrID := m.getTargetWindow(req)

	var windowNotRequiredRequests = []int{callRequest, eventsRequest, applicationRequest, systemRequest}

	// Some operations (calls, events, application) don't require a window
	// This is useful for browser-based deployments with custom transports
	windowRequired := !slices.Contains(windowNotRequiredRequests, req.Object)
	if windowRequired && targetWindow == nil {
		return nil, errs.NewInvalidRuntimeCallErrorf("window '%s' not found", nameOrID)
	}

	m.logRuntimeCall(req)

	switch req.Object {
	case windowRequest:
		return m.processWindowMethod(req, targetWindow)
	case clipboardRequest:
		return m.processClipboardMethod(req)
	case dialogRequest:
		return m.processDialogMethod(req, targetWindow)
	case eventsRequest:
		return m.processEventsMethod(req, targetWindow)
	case applicationRequest:
		return m.processApplicationMethod(req)
	case contextMenuRequest:
		return m.processContextMenuMethod(req, targetWindow)
	case screensRequest:
		return m.processScreensMethod(req)
	case callRequest:
		return m.processCallMethod(ctx, req, targetWindow)
	case systemRequest:
		return m.processSystemMethod(req)
	case browserRequest:
		return m.processBrowserMethod(req)
	case cancelCallRequest:
		return m.processCallCancelMethod(req)
	default:
		return nil, errs.NewInvalidRuntimeCallErrorf("unknown object %d", req.Object)
	}
}

func (m *MessageProcessor) getTargetWindow(req *RuntimeRequest) (Window, string) {
	if req.WebviewWindowName != "" {
		window, _ := globalApplication.Window.GetByName(req.WebviewWindowName)
		return window, req.WebviewWindowName
	}
	if req.WebviewWindowID == 0 {
		// No window specified - return the first available window
		// This is useful for custom transports that don't have automatic window context
		windows := globalApplication.Window.GetAll()
		if len(windows) > 0 {
			return windows[0], ""
		}
		return nil, ""
	}
	targetWindow, _ := globalApplication.Window.GetByID(uint(req.WebviewWindowID))
	if targetWindow == nil {
		m.Error("Window ID not found:", "id", req.WebviewWindowID)
		return nil, strconv.Itoa(int(windowID))
	}
	return targetWindow, strconv.Itoa(int(windowID))
}

func (m *MessageProcessor) Error(message string, args ...any) {
	m.logger.Error(message, args...)
}

func (m *MessageProcessor) Info(message string, args ...any) {
	m.logger.Info(message, args...)
}

func (m *MessageProcessor) logRuntimeCall(req *RuntimeRequest) {
	objectName := objectNames[req.Object]

	methodName := ""
	switch req.Object {
	case callRequest:
		return // logs done separately in call processor
	case clipboardRequest:
		methodName = clipboardMethods[req.Method]
	case applicationRequest:
		methodName = applicationMethodNames[req.Method]
	case eventsRequest:
		methodName = eventsMethodNames[req.Method]
	case contextMenuRequest:
		methodName = contextmenuMethodNames[req.Method]
	case dialogRequest:
		methodName = dialogMethodNames[req.Method]
	case windowRequest:
		methodName = windowMethodNames[req.Method]
	case screensRequest:
		methodName = screensMethodNames[req.Method]
	case systemRequest:
		methodName = systemMethodNames[req.Method]
	case browserRequest:
		methodName = browserMethodNames[req.Method]
	case cancelCallRequest:
		methodName = "Cancel"
	}

	m.Info("Runtime call:", "method", objectName+"."+methodName, "args", req.Args.String())
}
