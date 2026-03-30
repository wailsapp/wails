package webcontentsview

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// ErrAutomationNotSupported reports that the current platform implementation does not expose the automation bridge.
var ErrAutomationNotSupported = errors.New("webcontentsview automation is not supported on this platform")

var automationViews sync.Map

type automationObserver func(automationTargetEvent)

type automationNativeCapable interface {
	automationEnsureReady() error
	automationNativeCapabilities() automationNativeCapabilities
	automationEvaluate(expression string, world automationExecutionWorld, awaitPromise bool) (automationRemoteObject, error)
	automationInvoke(method string, params json.RawMessage) (any, error)
	automationCreatePDF() (string, error)
	automationSetInspectable(enabled bool) error
}

type automationTarget interface {
	targetID() string
	targetInfo() TargetInfo
	targetCapabilities() automationNativeCapabilities
	addAutomationObserver(automationObserver) uint64
	removeAutomationObserver(uint64)
	ensureAutomationReady() error
	navigate(string) error
	captureScreenshot() (string, error)
	printToPDF() (string, error)
	evaluate(expression string, world automationExecutionWorld, awaitPromise bool) (automationRemoteObject, error)
	invoke(method string, params json.RawMessage) (any, error)
	setInspectable(bool) error
	bufferedConsoleMessages() []AutomationConsoleMessage
}

type automationState struct {
	view         *WebContentsView
	targetID     string
	observers    map[uint64]automationObserver
	nextObserver atomic.Uint64
	mu           sync.RWMutex
	url          string
	title        string
	attached     bool
	loading      bool
	inspectable  bool
	console      []AutomationConsoleMessage
}

func (v *WebContentsView) ensureAutomationState() *automationState {
	if v.automation != nil {
		return v.automation
	}

	state := newAutomationState(v)
	v.automation = state
	registerAutomationView(v)
	return state
}

func newAutomationState(view *WebContentsView) *automationState {
	return &automationState{
		view:      view,
		targetID:  "wv-" + itoa(int(view.id)),
		observers: make(map[uint64]automationObserver),
		url:       view.options.URL,
	}
}

func registerAutomationView(view *WebContentsView) {
	automationViews.Store(view.id, view)
}

func lookupAutomationView(id uint) *WebContentsView {
	if value, ok := automationViews.Load(id); ok {
		return value.(*WebContentsView)
	}
	return nil
}

func (v *WebContentsView) targetID() string {
	return v.ensureAutomationState().targetID
}

func (v *WebContentsView) targetCapabilities() automationNativeCapabilities {
	if impl, ok := v.impl.(automationNativeCapable); ok {
		return impl.automationNativeCapabilities()
	}
	return defaultAutomationCapabilities()
}

func (v *WebContentsView) targetInfo() TargetInfo {
	state := v.ensureAutomationState()
	state.mu.RLock()
	defer state.mu.RUnlock()

	return TargetInfo{
		TargetID:          state.targetID,
		Type:              "webcontentsview",
		Name:              v.options.Name,
		URL:               state.url,
		Title:             state.title,
		Loading:           state.loading,
		Attached:          state.attached,
		InspectionEnabled: state.inspectable,
		Platform:          automationPlatform(),
	}
}

func (v *WebContentsView) addAutomationObserver(observer automationObserver) uint64 {
	state := v.ensureAutomationState()
	id := state.nextObserver.Add(1)
	state.mu.Lock()
	state.observers[id] = observer
	state.mu.Unlock()
	return id
}

func (v *WebContentsView) removeAutomationObserver(id uint64) {
	state := v.ensureAutomationState()
	state.mu.Lock()
	delete(state.observers, id)
	state.mu.Unlock()
}

func (v *WebContentsView) ensureAutomationReady() error {
	if impl, ok := v.impl.(automationNativeCapable); ok {
		return impl.automationEnsureReady()
	}
	return ErrAutomationNotSupported
}

func (v *WebContentsView) navigate(url string) error {
	v.SetURL(url)
	return nil
}

func (v *WebContentsView) captureScreenshot() (string, error) {
	return v.TakeSnapshot(), nil
}

func (v *WebContentsView) printToPDF() (string, error) {
	if impl, ok := v.impl.(automationNativeCapable); ok {
		if err := impl.automationEnsureReady(); err != nil {
			return "", err
		}
		return impl.automationCreatePDF()
	}
	return "", ErrAutomationNotSupported
}

func (v *WebContentsView) evaluate(expression string, world automationExecutionWorld, awaitPromise bool) (automationRemoteObject, error) {
	if impl, ok := v.impl.(automationNativeCapable); ok {
		if err := impl.automationEnsureReady(); err != nil {
			return automationRemoteObject{}, err
		}
		return impl.automationEvaluate(expression, world, awaitPromise)
	}
	return automationRemoteObject{}, ErrAutomationNotSupported
}

func (v *WebContentsView) invoke(method string, params json.RawMessage) (any, error) {
	if impl, ok := v.impl.(automationNativeCapable); ok {
		if err := impl.automationEnsureReady(); err != nil {
			return nil, err
		}
		return impl.automationInvoke(method, params)
	}
	return nil, ErrAutomationNotSupported
}

func (v *WebContentsView) setInspectable(enabled bool) error {
	if impl, ok := v.impl.(automationNativeCapable); ok {
		if err := impl.automationSetInspectable(enabled); err != nil {
			return err
		}
		state := v.ensureAutomationState()
		state.mu.Lock()
		state.inspectable = enabled
		state.mu.Unlock()
		v.emitAutomationEvent(automationTargetEvent{
			TargetID: v.targetID(),
			Method:   "Target.targetInfoChanged",
			Params: map[string]any{
				"targetInfo": v.targetInfo(),
			},
			Scope: automationEventScopeAll,
		})
		return nil
	}
	return ErrAutomationNotSupported
}

func (v *WebContentsView) bufferedConsoleMessages() []AutomationConsoleMessage {
	state := v.ensureAutomationState()
	state.mu.RLock()
	defer state.mu.RUnlock()

	result := make([]AutomationConsoleMessage, len(state.console))
	copy(result, state.console)
	return result
}

func (v *WebContentsView) updateAutomationURL(url string) {
	state := v.ensureAutomationState()
	state.mu.Lock()
	state.url = url
	state.mu.Unlock()
}

func (v *WebContentsView) updateAutomationAttached(attached bool) {
	state := v.ensureAutomationState()
	state.mu.Lock()
	state.attached = attached
	state.mu.Unlock()
	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Target.targetInfoChanged",
		Params: map[string]any{
			"targetInfo": v.targetInfo(),
		},
		Scope: automationEventScopeAll,
	})
}

func (v *WebContentsView) emitAutomationEvent(event automationTargetEvent) {
	state := v.ensureAutomationState()
	state.mu.RLock()
	observers := make([]automationObserver, 0, len(state.observers))
	for _, observer := range state.observers {
		observers = append(observers, observer)
	}
	state.mu.RUnlock()

	for _, observer := range observers {
		observer(event)
	}
}

func (v *WebContentsView) recordAutomationConsoleMessage(message AutomationConsoleMessage) {
	state := v.ensureAutomationState()
	state.mu.Lock()
	state.console = append(state.console, message)
	if len(state.console) > 200 {
		state.console = append([]AutomationConsoleMessage(nil), state.console[len(state.console)-200:]...)
	}
	state.mu.Unlock()

	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Console.messageAdded",
		Params: map[string]any{
			"targetId": v.targetID(),
			"message":  message,
		},
		Scope: automationEventScopeConsole,
	})
}

func (v *WebContentsView) recordAutomationException(exception AutomationException) {
	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Runtime.exceptionThrown",
		Params: map[string]any{
			"targetId":  v.targetID(),
			"exception": exception,
		},
		Scope: automationEventScopeConsole,
	})
}

func (v *WebContentsView) recordAutomationDOMContentLoaded(url string, timestamp int64) {
	state := v.ensureAutomationState()
	state.mu.Lock()
	state.url = url
	state.loading = false
	state.mu.Unlock()

	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Page.domContentEventFired",
		Params: map[string]any{
			"targetId":  v.targetID(),
			"url":       url,
			"timestamp": timestamp,
		},
		Scope: automationEventScopeAttached,
	})
	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Target.targetInfoChanged",
		Params: map[string]any{
			"targetInfo": v.targetInfo(),
		},
		Scope: automationEventScopeAll,
	})
}

func (v *WebContentsView) recordAutomationNavigationStarted(url string, timestamp int64) {
	state := v.ensureAutomationState()
	state.mu.Lock()
	state.url = url
	state.loading = true
	state.mu.Unlock()

	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Page.frameStartedLoading",
		Params: map[string]any{
			"targetId":  v.targetID(),
			"url":       url,
			"timestamp": timestamp,
		},
		Scope: automationEventScopeAttached,
	})
	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Target.targetInfoChanged",
		Params: map[string]any{
			"targetInfo": v.targetInfo(),
		},
		Scope: automationEventScopeAll,
	})
}

func (v *WebContentsView) recordAutomationNavigationCommitted(url, title string, timestamp int64) {
	state := v.ensureAutomationState()
	state.mu.Lock()
	state.url = url
	state.title = title
	state.loading = true
	state.mu.Unlock()

	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Page.frameNavigated",
		Params: map[string]any{
			"targetId":  v.targetID(),
			"url":       url,
			"title":     title,
			"timestamp": timestamp,
		},
		Scope: automationEventScopeAttached,
	})
	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Target.targetInfoChanged",
		Params: map[string]any{
			"targetInfo": v.targetInfo(),
		},
		Scope: automationEventScopeAll,
	})
}

func (v *WebContentsView) recordAutomationLoadFinished(url, title string, timestamp int64) {
	state := v.ensureAutomationState()
	state.mu.Lock()
	state.url = url
	state.title = title
	state.loading = false
	state.mu.Unlock()

	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Page.loadEventFired",
		Params: map[string]any{
			"targetId":  v.targetID(),
			"url":       url,
			"title":     title,
			"timestamp": timestamp,
		},
		Scope: automationEventScopeAttached,
	})
	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Target.targetInfoChanged",
		Params: map[string]any{
			"targetInfo": v.targetInfo(),
		},
		Scope: automationEventScopeAll,
	})
}

func (v *WebContentsView) recordAutomationPageState(url, title string) {
	state := v.ensureAutomationState()
	state.mu.Lock()
	state.url = url
	state.title = title
	state.mu.Unlock()
	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Target.targetInfoChanged",
		Params: map[string]any{
			"targetInfo": v.targetInfo(),
		},
		Scope: automationEventScopeAll,
	})
}

func (v *WebContentsView) recordAutomationWindowOpenRequested(url string, timestamp int64) {
	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Page.windowOpenRequested",
		Params: map[string]any{
			"targetId":  v.targetID(),
			"url":       url,
			"timestamp": timestamp,
		},
		Scope: automationEventScopeAttached,
	})
}

func (v *WebContentsView) recordAutomationProcessTerminated(timestamp int64) {
	state := v.ensureAutomationState()
	state.mu.Lock()
	state.loading = false
	state.mu.Unlock()
	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Page.webContentProcessTerminated",
		Params: map[string]any{
			"targetId":  v.targetID(),
			"timestamp": timestamp,
		},
		Scope: automationEventScopeAttached,
	})
	v.emitAutomationEvent(automationTargetEvent{
		TargetID: v.targetID(),
		Method:   "Target.targetInfoChanged",
		Params: map[string]any{
			"targetInfo": v.targetInfo(),
		},
		Scope: automationEventScopeAll,
	})
}

func (v *WebContentsView) syncAutomationState() {
	v.updateAutomationURL(v.GetURL())
}

func itoa(value int) string {
	return strconv.Itoa(value)
}

func (v *WebContentsView) updateBounds(bounds application.Rect) {
	v.options.Bounds = bounds
}
