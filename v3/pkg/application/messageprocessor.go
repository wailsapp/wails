package application

import (
	"context"
	"errors"
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
	cancelCallRequesst = 10
)

type RuntimeRequest struct {
	Object            int         `json:"object"`
	Method            int         `json:"method"`
	Params            QueryParams `json:"params"`
	WebviewWindowName string      `json:"webviewWindowName,omitempty"`
	WebviewWindowId   uint32      `json:"webviewWindowId,omitempty"`
	ClientId          string      `json:"clientId,omitempty"`
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

func (m *MessageProcessor) HandleRuntimeCallWithIDs(ctx context.Context, req *RuntimeRequest) (any, error) {
	defer func() (any, error) {
		if handlePanic() {
			// TODO: should get error here and return it.
			return nil, errors.New("runtime panic detected!")
		}

		panic("todo!")
	}()
	targetWindow, nameOrID := m.getTargetWindow(req)

	var windowNotRequiredRequests = []int{callRequest, eventsRequest, applicationRequest, systemRequest}

	// Some operations (calls, events, application) don't require a window
	// This is useful for browser-based deployments with custom transports
	windowRequired := !slices.Contains(windowNotRequiredRequests, req.Object)
	if windowRequired && targetWindow == nil {
		return nil, errs.NewInvalidRuntimeCallErrorf("window '%s' not found", nameOrID)
	}

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
	case cancelCallRequesst:
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
	if req.WebviewWindowId == 0 {
		// No window specified - return the first available window
		// This is useful for custom transports that don't have automatic window context
		windows := globalApplication.Window.GetAll()
		if len(windows) > 0 {
			return windows[0], ""
		}
		return nil, ""
	}
	targetWindow, _ := globalApplication.Window.GetByID(uint(req.WebviewWindowId))
	if targetWindow == nil {
		m.Error("Window ID not found:", "id", req.WebviewWindowId)
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
