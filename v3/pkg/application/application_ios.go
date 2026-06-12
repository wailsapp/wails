//go:build ios && !server

package application

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Foundation -framework UIKit -framework WebKit -framework UniformTypeIdentifiers -framework Network

#include <stdlib.h>
#include <string.h>
#include "application_ios.h"
#include "webview_window_ios.h"

*/
import "C"

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"unsafe"

	"encoding/json"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func iosConsoleLogf(level string, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	clevel := C.CString(level)
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(clevel))
	defer C.free(unsafe.Pointer(cmsg))
	C.ios_console_log(clevel, cmsg)
}

// iosDebugLogf is for the framework's internal diagnostics. It compiles to a
// no-op in production builds (see ios_logging_production.go).
func iosDebugLogf(format string, a ...interface{}) {
	if iosVerboseLogging {
		iosConsoleLogf("debug", format, a...)
	}
}

//export init_go
func init_go() {
	// Called from the iOS main function to initialize the Go runtime.
}

// iosLaunched is closed when UIApplicationDelegate's
// didFinishLaunchingWithOptions fires, signalling that UIKit is ready.
var (
	iosLaunched     = make(chan struct{})
	iosLaunchedOnce sync.Once
)

//export iosApplicationDidLaunch
func iosApplicationDidLaunch() {
	iosLaunchedOnce.Do(func() {
		close(iosLaunched)
	})
}

func (a *App) platformRun() {
	iosDebugLogf("[application_ios.go] platformRun: initialising")

	// Propagate the logging verbosity to the native layer.
	C.ios_set_verbose_logging(C.bool(iosVerboseLogging))

	C.ios_app_init()

	// Wait until the UIApplication delegate has finished launching.
	// UIApplicationMain runs on the main thread (started from main.m); the
	// delegate signals readiness via iosApplicationDidLaunch. The timeout is a
	// safety net only - it should never be hit in practice.
	select {
	case <-iosLaunched:
		iosDebugLogf("[application_ios.go] UIKit launch signal received")
	case <-time.After(10 * time.Second):
		iosConsoleLogf("warn", "[application_ios.go] Timed out waiting for UIKit launch signal; continuing")
	}

	// Emit the launch event from here rather than the delegate: by this point
	// setupCommonEvents has registered its listeners, so the event (and its
	// Common.ApplicationStarted mapping) cannot be dropped by startup races.
	applicationEvents <- newApplicationEvent(events.IOS.ApplicationDidFinishLaunching)

	// The WebView is created when the window runs (app.Window.NewWithOptions).
	// UIApplicationMain owns the main thread; keep the Go runtime alive here.
	select {}
}

func (a *App) platformQuit() {
	C.ios_app_quit()
}

func (a *App) isDarkMode() bool {
	return bool(C.ios_is_dark_mode())
}

func (a *App) isWindows() bool {
	return false
}

//export LogInfo
func LogInfo(source *C.char, message *C.char) {
	goSource := C.GoString(source)
	goMessage := C.GoString(message)

	if globalApplication != nil && globalApplication.Logger != nil {
		globalApplication.info("iOS log", "source", goSource, "message", goMessage)
	} else {
		iosDebugLogf("[iOS-%s] %s", goSource, goMessage)
	}
}

// Platform-specific app implementation for iOS
type iosApp struct {
	parent *App
}

// newPlatformApp creates an iosApp for the provided App and applies iOS-specific
// configuration derived from app.options. It sets input accessory visibility,
// scrolling/bounce/indicator behavior, navigation gestures, link preview,
// media playback, inspector, user agent strings, app background color, and
// native tabs (marshaling items to JSON when enabled). The function invokes
// platform bindings to apply these settings and returns the configured *iosApp.
func newPlatformApp(app *App) *iosApp {
	result := &iosApp{
		parent: app,
	}
	// Configure input accessory visibility according to options
	// Default: false (show accessory) when not explicitly set to true
	disable := false
	if app != nil {
		disable = app.options.IOS.DisableInputAccessoryView
	}
	C.ios_set_disable_input_accessory(C.bool(disable))

	// Scrolling / Bounce / Indicators (defaults enabled; using Disable* flags)
	C.ios_set_disable_scroll(C.bool(app.options.IOS.DisableScroll))
	C.ios_set_disable_bounce(C.bool(app.options.IOS.DisableBounce))
	C.ios_set_disable_scroll_indicators(C.bool(app.options.IOS.DisableScrollIndicators))

	// Navigation gestures (Enable*)
	C.ios_set_enable_back_forward_gestures(C.bool(app.options.IOS.EnableBackForwardNavigationGestures))

	// Link preview (Disable*)
	C.ios_set_disable_link_preview(C.bool(app.options.IOS.DisableLinkPreview))

	// Media playback
	C.ios_set_enable_inline_media_playback(C.bool(app.options.IOS.EnableInlineMediaPlayback))
	C.ios_set_enable_autoplay_without_user_action(C.bool(app.options.IOS.EnableAutoplayWithoutUserAction))

	// Inspector (Disable*)
	C.ios_set_disable_inspectable(C.bool(app.options.IOS.DisableInspectable))

	// User agent strings
	if ua := strings.TrimSpace(app.options.IOS.UserAgent); ua != "" {
		cua := C.CString(ua)
		C.ios_set_user_agent(cua)
		C.free(unsafe.Pointer(cua))
	}
	if appName := strings.TrimSpace(app.options.IOS.ApplicationNameForUserAgent); appName != "" {
		cname := C.CString(appName)
		C.ios_set_app_name_for_user_agent(cname)
		C.free(unsafe.Pointer(cname))
	}
	// App-wide background colour for iOS window (pre-WebView)
	if app.options.IOS.AppBackgroundColourSet {
		rgba := app.options.IOS.BackgroundColour
		C.ios_set_app_background_color(
			C.uchar(rgba.Red), C.uchar(rgba.Green), C.uchar(rgba.Blue), C.uchar(rgba.Alpha), C.bool(true),
		)
	} else {
		// Ensure it's marked as not set to allow delegate to fallback to white
		C.ios_set_app_background_color(255, 255, 255, 255, C.bool(false))
	}
	// Native tabs option: only enable when explicitly requested
	if app.options.IOS.EnableNativeTabs {
		if len(app.options.IOS.NativeTabsItems) > 0 {
			if data, err := json.Marshal(app.options.IOS.NativeTabsItems); err == nil {
				cjson := C.CString(string(data))
				C.ios_native_tabs_set_items_json(cjson)
				C.free(unsafe.Pointer(cjson))
			} else if globalApplication != nil {
				globalApplication.error("Failed to marshal IOS.NativeTabsItems: %v", err)
			}
		}
		C.ios_native_tabs_set_enabled(C.bool(true))
	}

	return result
}

func (a *iosApp) run() error {
	// UIApplicationMain is already running in the main thread (from main.m).
	// Wire common events (e.g. map ApplicationDidFinishLaunching → Common.ApplicationStarted)
	a.setupCommonEvents()

	// Populate the ScreenManager so Screens.* runtime calls return data
	// (desktop platforms do this from their event loop; iOS has none).
	// getScreens() always yields at least a fallback screen, so this is safe.
	if screens, err := getScreens(); err == nil && len(screens) > 0 {
		if err := a.parent.Screen.LayoutScreens(screens); err != nil {
			iosConsoleLogf("error", "[application_ios.go] LayoutScreens failed: %v", err)
		}
	}

	// Start the native system-event monitors (battery, network, lock, theme,
	// app lifecycle, memory). They emit "system:*" custom events to JS.
	C.ios_start_system_event_monitors()

	a.parent.platformRun()

	// platformRun blocks forever with select{}
	// If we get here, something went wrong
	iosConsoleLogf("error", "[application_ios.go] ERROR: platformRun() returned unexpectedly")
	return nil
}

func (a *iosApp) destroy() {
	// Cleanup iOS resources
}

func (a *iosApp) setIcon(_ []byte) {
	// iOS app icon is set through Info.plist
}

func (a *iosApp) name() string {
	return a.parent.options.Name
}

func (a *iosApp) GetFlags(options Options) map[string]any {
	return nil
}

// dispatchOnMainThread is implemented in mainthread_ios.go

func (a *iosApp) getAccentColor() string {
	// iOS accent color
	return ""
}

func (a *iosApp) getCurrentWindowID() uint {
	// iOS current window ID
	return 0
}

func (a *iosApp) hide() {
	// iOS hide application - minimize to background
}

func (a *iosApp) isDarkMode() bool {
	return a.parent.isDarkMode()
}

// isOnMainThread is implemented in mainthread_ios.go

func (a *iosApp) on(eventID uint) {
	registerIOSListener(eventID)
}

func (a *iosApp) setApplicationMenu(_ *Menu) {
	// iOS doesn't have application menus
}

func (a *iosApp) show() {
	// iOS show application
}

func (a *iosApp) showAboutDialog(_ string, _ string, _ []byte) {
	// iOS about dialog
}

func (a *iosApp) getPrimaryScreen() (*Screen, error) {
	screens, err := getScreens()
	if err != nil || len(screens) == 0 {
		return nil, err
	}
	return screens[0], nil
}

func (a *iosApp) getScreens() ([]*Screen, error) {
	return getScreens()
}

func (a *App) logPlatformInfo() {
	// Log iOS platform info
}

func (a *App) platformEnvironment() map[string]any {
	return map[string]any{
		"platform": "ios",
	}
}

func fatalHandler(errFunc func(error)) {
	// iOS fatal handler
}

// ExecuteJavaScript runs JavaScript code in the WebView
func (a *App) ExecuteJavaScript(windowID uint, js string) {
	cjs := C.CString(js)
	defer C.free(unsafe.Pointer(cjs))
	C.ios_execute_javascript(C.uint(windowID), cjs)
}

// iosRuntimeReadyWindows tracks windows for which a synthetic
// "wails:runtime:ready" has been injected (see ServeAssetRequest).
var iosRuntimeReadyWindows sync.Map

// ServeAssetRequest handles requests from the WebView
//
//export ServeAssetRequest
func ServeAssetRequest(windowID C.uint, urlSchemeTask unsafe.Pointer) {
	// Route the request through the webviewRequests channel to use the asset server
	go func() {
		// Use the webview package's NewRequest to wrap the task pointer
		req := webview.NewRequest(urlSchemeTask)
		url, _ := req.URL()

		iosDebugLogf("[iOS-request] %s (window %d)", url, windowID)

		// The JavaScript runtime announces itself with a
		// "wails:runtime:ready" postMessage, but a message posted during the
		// initial page load can be dropped by WebKit before the bridge is
		// fully up. A call to /wails/runtime proves the runtime is mounted,
		// so treat the first one as an implicit ready signal. processMessage
		// handles duplicate ready messages gracefully.
		if strings.Contains(url, "/wails/runtime") {
			if _, alreadyReady := iosRuntimeReadyWindows.LoadOrStore(uint(windowID), true); !alreadyReady {
				windowMessageBuffer <- &windowMessage{
					windowId: uint(windowID),
					message:  "wails:runtime:ready",
				}
			}
		}

		// Try to resolve the window name from the window ID so the AssetServer
		// receives both x-wails-window-id and x-wails-window-name headers.
		winName := ""
		if globalApplication != nil {
			if window, ok := globalApplication.Window.GetByID(uint(windowID)); ok && window != nil {
				winName = window.Name()
			}
		}

		request := &webViewAssetRequest{
			Request:    req,
			windowId:   uint(windowID),
			windowName: winName,
		}

		// Send through the channel to be handled by the asset server
		webviewRequests <- request
	}()
}

// HandleJSMessage handles messages from JavaScript
//
//export HandleJSMessage
func HandleJSMessage(windowID C.uint, message *C.char) {
	msg := C.GoString(message)
	if msg == "" {
		return
	}

	iosDebugLogf("[iOS-message] window %d: %s", windowID, msg)

	// Structured payloads carry the message in a "name" or "message" field;
	// plain strings (e.g. "wails:runtime:ready") are forwarded as-is.
	var msgData map[string]interface{}
	if err := json.Unmarshal([]byte(msg), &msgData); err == nil && msgData != nil {
		if name, ok := msgData["name"].(string); ok && name != "" {
			msg = name
		} else if name, ok := msgData["message"].(string); ok && name != "" {
			msg = name
		}
	}

	windowMessageBuffer <- &windowMessage{
		windowId: uint(windowID),
		message:  msg,
	}
}

// Note: applicationEvents and windowEvents are already defined in events.go
// We'll use those existing channels

type iosWindowEvent struct {
	WindowID uint
	EventID  uint
}

//export processApplicationEvent
func processApplicationEvent(eventID C.uint, data unsafe.Pointer) {
	iosDebugLogf("[application_ios.go] application event: %d", eventID)

	// Create and send the application event
	event := newApplicationEvent(events.ApplicationEventType(eventID))

	// Send to the applicationEvents channel for processing
	applicationEvents <- event
}

//export processWindowEvent
func processWindowEvent(windowID C.uint, eventID C.uint) {
	iosDebugLogf("[application_ios.go] window event: window %d, event %d", windowID, eventID)
	windowEvents <- &windowEvent{
		WindowID: uint(windowID),
		EventID:  uint(eventID),
	}
}

// emitMobileSystemEvent emits a custom "system:*" event to JS from a
// native-triggered callback (battery/network/lock/theme/lifecycle). jsonStr
// (which may be empty) is decoded into the event payload; if it is not valid
// JSON it is passed through as a plain string.
func emitMobileSystemEvent(name, jsonStr string) {
	if globalApplication == nil || name == "" {
		return
	}
	var data any
	if jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			data = jsonStr
		}
	}
	globalApplication.Event.Emit(name, data)
}

// emitSystemEvent is the C entry point the iOS system-event monitors call (see
// ios_system_events.m): name is e.g. "system:battery", json is the payload.
//
//export emitSystemEvent
func emitSystemEvent(name *C.char, jsonData *C.char) {
	emitMobileSystemEvent(C.GoString(name), C.GoString(jsonData))
}

// iosEventListeners records which native event IDs have at least one Go-side
// listener. Registration happens via iosApp.on / iosWebviewWindow.on, which
// the cross-platform layer invokes whenever a listener is added. Listeners
// are never unregistered natively (same behaviour as macOS).
var (
	iosEventListeners     = make(map[uint]bool)
	iosEventListenersLock sync.RWMutex
)

func registerIOSListener(eventID uint) {
	iosEventListenersLock.Lock()
	defer iosEventListenersLock.Unlock()
	iosEventListeners[eventID] = true
}

//export hasListeners
func hasListeners(eventID C.uint) C.bool {
	iosEventListenersLock.RLock()
	defer iosEventListenersLock.RUnlock()
	return C.bool(iosEventListeners[uint(eventID)])
}
