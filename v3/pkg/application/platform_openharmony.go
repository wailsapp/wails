//go:build openharmony && cgo && !server

package application

import (
	"encoding/json"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// Global state mirroring application_android.go: a global App reference for NAPI
// callbacks, the registered main function, and the app-ready signal.
var (
	globalApp     *App
	globalAppLock sync.RWMutex

	openharmonyMainFunc func()
	openharmonyMainLock sync.Mutex

	appReady     = make(chan struct{})
	appReadyOnce sync.Once
)

// RegisterOpenharmonyMain registers the main function to be called when the OHOS
// ability starts. Call it from init() in your main package (mirrors
// RegisterAndroidMain):
//
//	func init() {
//		application.RegisterOpenharmonyMain(main)
//	}
func RegisterOpenharmonyMain(mainFunc func()) {
	openharmonyMainLock.Lock()
	defer openharmonyMainLock.Unlock()
	openharmonyMainFunc = mainFunc
}

func signalAppReady() {
	appReadyOnce.Do(func() { close(appReady) })
}

// OpenharmonyReady reports whether the host can begin requesting local assets.
func OpenharmonyReady() bool {
	select {
	case <-appReady:
		return true
	default:
		return false
	}
}

// emitOpenharmonyApplicationEvent forwards a typed OHOS event onto the runtime
// event bus (mirrors emitAndroidApplicationEvent).
func emitOpenharmonyApplicationEvent(event events.ApplicationEventType) {
	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()
	if app == nil {
		return
	}
	applicationEvents <- newApplicationEvent(event)
}

func emitOpenharmonyApplicationEventWithData(event events.ApplicationEventType, jsonStr string) {
	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()
	if app == nil {
		return
	}
	evt := newApplicationEvent(event)
	if jsonStr != "" {
		evt.Context().setData(parseEventJSON(jsonStr))
	}
	applicationEvents <- evt
}

// emitNativeEventToJS forwards a named custom event (with JSON payload) to the
// frontend (mirrors android).
func emitNativeEventToJS(name string, jsonStr string) {
	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()
	if app == nil {
		return
	}
	app.Event.Emit(name, parseEventJSON(jsonStr))
}

// ---- NAPI entry points (mirrors the Java_com_wails_app_WailsBridge_* exports) --

// WailsOhNativeShutdown quits the app (nativeShutdown).
//
//export WailsOhNativeShutdown
func WailsOhNativeShutdown() {
	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()
	if app != nil {
		app.Quit()
	}
}

// WailsOhNativeOnStart / OnResume / OnPause / OnStop / OnLowMemory mirror the
// android lifecycle exports, emitting the typed openharmony: events.

//export WailsOhNativeOnStart
func WailsOhNativeOnStart() {
	emitOpenharmonyApplicationEvent(events.Openharmony.AbilityStarted)
}

//export WailsOhNativeOnResume
func WailsOhNativeOnResume() {
	emitOpenharmonyApplicationEvent(events.Openharmony.AbilityForeground)
}

//export WailsOhNativeOnPause
func WailsOhNativeOnPause() {
	emitOpenharmonyApplicationEvent(events.Openharmony.AbilityBackground)
}

//export WailsOhNativeOnStop
func WailsOhNativeOnStop() {
	emitOpenharmonyApplicationEvent(events.Openharmony.AbilityDestroyed)
}

//export WailsOhNativeOnLowMemory
func WailsOhNativeOnLowMemory() {
	emitOpenharmonyApplicationEvent(events.Openharmony.ApplicationLowMemory)
}

// WailsOhNativeEmitSystemEvent is the funnel the OHOS system-event receivers
// (battery, network, lock, theme) call to deliver a typed openharmony: event
// with its JSON payload (mirrors nativeEmitSystemEvent).
//
//export WailsOhNativeEmitSystemEvent
func WailsOhNativeEmitSystemEvent(name string, jsonStr string) {
	if eventType, ok := openharmonySystemEventTypes[name]; ok {
		emitOpenharmonyApplicationEventWithData(eventType, jsonStr)
	} else {
		openharmonyDebugLogf("[NAPI] nativeEmitSystemEvent: unknown event %q", name)
	}
}

// WailsOhNativeEmitEvent delivers an arbitrary custom event (with JSON payload)
// to the frontend (mirrors nativeEmitEvent).
//
//export WailsOhNativeEmitEvent
func WailsOhNativeEmitEvent(name string, jsonStr string) {
	emitNativeEventToJS(name, jsonStr)
}

// WailsOhNativeOnPageFinished forwards the typed WebViewPageFinished window
// event once the ArkWeb finishes loading (mirrors nativeOnPageFinished).
//
//export WailsOhNativeOnPageFinished
func WailsOhNativeOnPageFinished(url string) {
	openharmonyDebugLogf("[NAPI] page finished: %s", url)
	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()
	if app == nil {
		return
	}
	if windowID := openharmonyFirstWindowID(app); windowID != 0 {
		windowEvents <- &windowEvent{
			WindowID: windowID,
			EventID:  uint(events.Openharmony.WebViewPageFinished),
		}
	}
}

// WailsOhNativeServeAsset serves an app asset for the ArkWeb (mirrors
// nativeServeAsset). Returns the response body bytes ("" + false on error).
//
//export WailsOhNativeServeAsset
func WailsOhNativeServeAsset(path, method, _ string) ([]byte, bool) {
	if !OpenharmonyReady() {
		return nil, false
	}
	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()
	if app == nil || app.assets == nil {
		return nil, false
	}
	data, err := serveAssetForOpenharmony(app, method, path)
	if err != nil {
		openharmonyDebugLogf("[NAPI] error serving asset %s: %v", path, err)
		return nil, false
	}
	return data, true
}

// WailsOhNativeGetAssetMimeType returns the MIME type for a path (mirrors
// nativeGetAssetMimeType).
//
//export WailsOhNativeGetAssetMimeType
func WailsOhNativeGetAssetMimeType(path string) string {
	return getMimeTypeForPath(path)
}

// WailsOhNativeHandleMessage routes a window message into the app (mirrors
// nativeHandleMessage).
//
//export WailsOhNativeHandleMessage
func WailsOhNativeHandleMessage(message string) {
	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()
	if app == nil {
		return
	}
	handleMessageForOpenharmony(app, message)
}

// WailsOhNativeHandleRuntimeCall handles a runtime RPC and returns the JSON
// envelope (mirrors nativeHandleRuntimeCall).
//
//export WailsOhNativeHandleRuntimeCall
func WailsOhNativeHandleRuntimeCall(payload string) string {
	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()
	if app == nil {
		return `{"ok":false,"error":"App not initialized"}`
	}
	return handleRuntimeCallForOpenharmony(app, payload)
}

// WailsOhNativeDialogCallback resolves a pending message dialog (mirrors
// nativeDialogCallback).
//
//export WailsOhNativeDialogCallback
func WailsOhNativeDialogCallback(callbackID uint, buttonIndex int) {
	openharmonyDialogCallback(callbackID, buttonIndex)
}

// WailsOhNativeFilePickerResult delivers one selected file path (mirrors
// nativeFilePickerResult).
//
//export WailsOhNativeFilePickerResult
func WailsOhNativeFilePickerResult(callbackID uint, path string) {
	openharmonyFilePickerResult(callbackID, path)
}

// WailsOhNativeFilePickerDone closes a finished file picker (mirrors
// nativeFilePickerDone).
//
//export WailsOhNativeFilePickerDone
func WailsOhNativeFilePickerDone(callbackID uint) {
	openharmonyFilePickerDone(callbackID)
}

// WailsOhNativeMainThreadCallback runs a queued main-thread function (mirrors
// nativeMainThreadCallback).
//
//export WailsOhNativeMainThreadCallback
func WailsOhNativeMainThreadCallback(callbackID uint) {
	openharmonyMainThreadCallback(callbackID)
}

func parseEventJSON(jsonStr string) map[string]any {
	if jsonStr == "" {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		return nil
	}
	return m
}

// platformRun mirrors androidApp: store the global App reference, signal asset
// readiness, lay out screens, emit the launch event, then park (the OHOS
// lifecycle is driven by the ArkTS ability, not a Go event loop).
func (a *App) platformRun() {
	openharmonyDebugLogf("[application_openharmony.go] platformRun: initialising")

	globalAppLock.Lock()
	globalApp = a
	globalAppLock.Unlock()

	signalAppReady()

	if screens, err := getScreens(); err == nil && len(screens) > 0 {
		if err := a.Screen.LayoutScreens(screens); err != nil {
			openharmonyDebugLogf("[application_openharmony.go] LayoutScreens failed: %v", err)
		}
	}

	applicationEvents <- newApplicationEvent(events.Openharmony.AbilityCreated)

	select {}
}

// platformQuit is a no-op: OHOS terminates the app when the last ability exits.
func (a *App) platformQuit() {}
