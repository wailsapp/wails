package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// This file is the cross-platform surface of macOS window state restoration
// (NSWindowRestoration): marking a window restorable under an identifier,
// storing string data with its restorable state, recreating windows through
// WindowManager.OnRestore at the next launch, and saving or restoring the
// WKWebView interaction state (macOS 12+). The platform functions it calls
// live in webview_window_restoration_darwin.go, with documented stubs in
// webview_window_restoration_other.go.
//
// macOS restores windows when the application is relaunched after a crash,
// a forced quit or a reboot, and after a normal quit when the system setting
// "Close windows when quitting an application" is off (the
// NSQuitAlwaysKeepsWindows default). Only windows that are visible when the
// application terminates are saved.

var (
	// ErrMacRestorationUnsupported means window state restoration is
	// unavailable on the current platform.
	ErrMacRestorationUnsupported = errors.New("window state restoration is unavailable on this platform")
	// ErrMacInteractionStateUnsupported means the running macOS release has
	// no WKWebView interaction state (it needs macOS 12).
	ErrMacInteractionStateUnsupported = errors.New("the WebView interaction state requires macOS 12 or later")
	// ErrMacInteractionStateInvalid means WebKit rejected the data passed to
	// RestoreInteractionState.
	ErrMacInteractionStateInvalid = errors.New("the WebView interaction state data was rejected")
)

// RestorationState is the data saved with a restorable window and handed
// back to WindowManager.OnRestore at the next launch. Keep it to what is
// needed to recreate the window: identifiers, a document path, scroll or
// selection state. Binary values such as a WebView interaction state can
// be stored base64 encoded.
type RestorationState struct {
	Data map[string]string
}

// Get returns the value stored under key, or "" when it is absent.
func (s RestorationState) Get(key string) string {
	if s.Data == nil {
		return ""
	}
	return s.Data[key]
}

// encodeRestorationState serialises the data as a JSON object. A nil or
// empty map encodes as "{}".
func encodeRestorationState(state RestorationState) (string, error) {
	data := state.Data
	if data == nil {
		data = map[string]string{}
	}
	encoded, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("encoding restoration state: %w", err)
	}
	return string(encoded), nil
}

// decodeRestorationState parses the JSON produced by encodeRestorationState.
// An empty string decodes as an empty state.
func decodeRestorationState(encoded string) (RestorationState, error) {
	state := RestorationState{Data: map[string]string{}}
	if encoded == "" {
		return state, nil
	}
	if err := json.Unmarshal([]byte(encoded), &state.Data); err != nil {
		return RestorationState{Data: map[string]string{}}, fmt.Errorf("decoding restoration state: %w", err)
	}
	if state.Data == nil {
		state.Data = map[string]string{}
	}
	return state, nil
}

func copyRestorationData(data map[string]string) map[string]string {
	if data == nil {
		return nil
	}
	result := make(map[string]string, len(data))
	for key, value := range data {
		result[key] = value
	}
	return result
}

// Restore handler registered with WindowManager.OnRestore.

var (
	macWindowRestoreHandlerLock sync.RWMutex
	macWindowRestoreHandler     func(id string, state RestorationState) Window
)

func setMacWindowRestoreHandler(handler func(id string, state RestorationState) Window) {
	macWindowRestoreHandlerLock.Lock()
	macWindowRestoreHandler = handler
	macWindowRestoreHandlerLock.Unlock()
}

func getMacWindowRestoreHandler() func(id string, state RestorationState) Window {
	macWindowRestoreHandlerLock.RLock()
	defer macWindowRestoreHandlerLock.RUnlock()
	return macWindowRestoreHandler
}

// macWindowRestoreRequest is one restoreWindowWithIdentifier: call from
// AppKit. The native side keeps the completion handler under requestID.
type macWindowRestoreRequest struct {
	requestID uint64
	id        string
	encoded   string
}

var (
	macWindowRestoreRequests = make(chan macWindowRestoreRequest, 8)
	// macWindowRestoreRequested counts the restore requests AppKit issued at
	// this launch. They all arrive before ApplicationDidFinishLaunching, so
	// the count is complete by the time ApplicationStarted is delivered.
	macWindowRestoreRequested atomic.Int32
)

// WillRestoreWindows reports whether macOS asked the application to restore
// at least one window at this launch. It is complete once
// events.Common.ApplicationStarted is delivered; use it there to skip
// creating the default window when a restored one is on its way. It is
// always false off macOS.
func (wm *WindowManager) WillRestoreWindows() bool {
	return macWindowRestoreRequested.Load() > 0
}

// noteMacWindowRestoreRequest queues a restore request from AppKit.
func noteMacWindowRestoreRequest(request macWindowRestoreRequest) {
	macWindowRestoreRequested.Add(1)
	macWindowRestoreRequests <- request
}

func init() {
	registerChromeEventLoop(func(*App) {
		for {
			request := <-macWindowRestoreRequests
			go handleMacWindowRestore(request)
		}
	})
}

// handleMacWindowRestore runs the OnRestore handler for one window and
// completes the native request with the recreated window, or with an error
// when there is no handler, the handler declined, or the handler failed.
func handleMacWindowRestore(request macWindowRestoreRequest) {
	defer handlePanic()
	handler := getMacWindowRestoreHandler()
	if handler == nil {
		macWindowRestorationComplete(request.requestID, nil, "no WindowManager.OnRestore handler is registered")
		return
	}
	state, err := decodeRestorationState(request.encoded)
	if err != nil {
		macWindowRestorationComplete(request.requestID, nil, err.Error())
		return
	}
	window := runMacWindowRestoreHandler(handler, request.id, state)
	if window == nil {
		macWindowRestorationComplete(request.requestID, nil, "the OnRestore handler declined the window")
		return
	}
	if webview, ok := window.(*WebviewWindow); ok && webview != nil && webview.RestorationID() == "" {
		// Keep the window restorable at the next launch as well.
		webview.SetRestorationID(request.id)
	}
	nsWindow := macWindowRestorationAwaitNative(window)
	if nsWindow == nil {
		macWindowRestorationComplete(request.requestID, nil, "the restored window was not created")
		return
	}
	if webview, ok := window.(*WebviewWindow); ok && webview != nil {
		// Configure the identifier now rather than when the window first
		// appears, so AppKit sees a restorable window in its completion.
		webview.applyRestorationID(nsWindow, webview.RestorationID())
	}
	macWindowRestorationComplete(request.requestID, nsWindow, "")
}

// runMacWindowRestoreHandler calls the handler and turns a panic into a
// declined window so AppKit's completion handler is always called.
func runMacWindowRestoreHandler(handler func(id string, state RestorationState) Window, id string, state RestorationState) (window Window) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if globalApplication != nil {
				globalApplication.error("OnRestore handler for window %q panicked: %v", id, recovered)
			}
			window = nil
		}
	}()
	return handler(id, state)
}

// Per-window restoration state kept in Go: the identifier and data set
// before the native window exists, and the data to encode.

type macRestorationState struct {
	id   *string
	data map[string]string
}

var (
	macRestorationLock    sync.Mutex
	macRestorationPending = map[*WebviewWindow]*macRestorationState{}
)

// macRestorationStateFor returns the window's Go-side state, creating it on
// first use. The caller holds macRestorationLock. A new entry registers a
// WindowClosing listener that drops it again, so closed windows are not
// kept alive by the map.
func macRestorationStateFor(w *WebviewWindow) *macRestorationState {
	state := macRestorationPending[w]
	if state == nil {
		state = &macRestorationState{}
		macRestorationPending[w] = state
		if w.eventListeners != nil {
			w.OnWindowEvent(events.Common.WindowClosing, func(*WindowEvent) { forgetMacRestoration(w) })
		}
	}
	return state
}

// SetRestorationID makes the window restorable under id, or removes it from
// state restoration when id is "". Before the native window exists the
// value is applied when the window is first shown; MacWindow.RestorationID
// sets it in the options. It is a no-op off macOS.
func (w *WebviewWindow) SetRestorationID(id string) Window {
	if w == nil {
		return w
	}
	macRestorationLock.Lock()
	macRestorationStateFor(w).id = &id
	macRestorationLock.Unlock()
	if nsWindow := w.NativeWindow(); nsWindow != nil {
		w.applyRestorationID(nsWindow, id)
	}
	return w
}

// RestorationID returns the identifier set with SetRestorationID or
// MacWindow.RestorationID, or "" when the window is not restorable.
func (w *WebviewWindow) RestorationID() string {
	if w == nil {
		return ""
	}
	macRestorationLock.Lock()
	defer macRestorationLock.Unlock()
	if state := macRestorationPending[w]; state != nil && state.id != nil {
		return *state.id
	}
	return w.options.Mac.RestorationID
}

// SetRestorationData replaces the string data saved with the window's
// restorable state and delivered to WindowManager.OnRestore at the next
// launch. Call it whenever the state changes; macOS re-encodes the window
// state shortly afterwards and at termination. Values are copied. It is a
// no-op off macOS.
func (w *WebviewWindow) SetRestorationData(data map[string]string) Window {
	if w == nil {
		return w
	}
	copied := copyRestorationData(data)
	macRestorationLock.Lock()
	macRestorationStateFor(w).data = copied
	macRestorationLock.Unlock()
	if nsWindow := w.NativeWindow(); nsWindow != nil {
		encoded, err := encodeRestorationState(RestorationState{Data: copied})
		if err != nil {
			w.Error("SetRestorationData: %v", err)
			return w
		}
		macWindowRestorationSetData(nsWindow, encoded)
	}
	return w
}

// RestorationData returns a copy of the data set with SetRestorationData.
func (w *WebviewWindow) RestorationData() map[string]string {
	if w == nil {
		return nil
	}
	macRestorationLock.Lock()
	defer macRestorationLock.Unlock()
	if state := macRestorationPending[w]; state != nil {
		return copyRestorationData(state.data)
	}
	return nil
}

// applyRestorationID configures the native window and pushes the data
// stored so far. It runs the native calls through InvokeSync.
func (w *WebviewWindow) applyRestorationID(nsWindow unsafe.Pointer, id string) {
	macWindowRestorationConfigure(nsWindow, id)
	if id == "" {
		return
	}
	macRestorationLock.Lock()
	var data map[string]string
	if state := macRestorationPending[w]; state != nil {
		data = copyRestorationData(state.data)
	}
	macRestorationLock.Unlock()
	if data == nil {
		return
	}
	encoded, err := encodeRestorationState(RestorationState{Data: data})
	if err != nil {
		w.Error("SetRestorationData: %v", err)
		return
	}
	macWindowRestorationSetData(nsWindow, encoded)
}

// forgetMacRestoration drops the Go-side state for a window that is
// closing; see macRestorationStateFor.
func forgetMacRestoration(w *WebviewWindow) {
	macRestorationLock.Lock()
	delete(macRestorationPending, w)
	macRestorationLock.Unlock()
}

// InteractionState returns the WKWebView interaction state (macOS 12+): an
// opaque blob holding the back-forward list and scroll positions that
// RestoreInteractionState can apply to a WebView showing the same content,
// typically stored base64 encoded in the RestorationState. The native window
// must exist. It returns ErrMacInteractionStateUnsupported before macOS 12
// and ErrMacOnly off macOS.
func (w *WebviewWindow) InteractionState() ([]byte, error) {
	if !macWindowRestorationSupported() {
		return nil, ErrMacOnly
	}
	if w == nil || w.impl == nil || w.isDestroyed() {
		return nil, ErrMacWindowNotCreated
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		return nil, ErrMacWindowNotCreated
	}
	return macWindowRestorationInteractionState(nsWindow)
}

// RestoreInteractionState applies a blob returned by InteractionState to the
// window's WebView (macOS 12+). WebKit navigates to the recorded page, so
// call it once the window exists and before or instead of the initial
// navigation settling. See InteractionState for the errors.
func (w *WebviewWindow) RestoreInteractionState(state []byte) error {
	if !macWindowRestorationSupported() {
		return ErrMacOnly
	}
	if w == nil || w.impl == nil || w.isDestroyed() {
		return ErrMacWindowNotCreated
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		return ErrMacWindowNotCreated
	}
	return macWindowRestorationSetInteractionState(nsWindow, state)
}
