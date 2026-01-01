//go:build ios

package application

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Foundation -framework UIKit -framework WebKit

#include <stdlib.h>
#include <string.h>
#include "application_ios.h"
#include "webview_window_ios.h"

*/
import "C"

import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	json "github.com/goccy/go-json"

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

func init() {
	iosConsoleLogf("info", "🔵 [application_ios.go] START init()")
	// For iOS, we need to handle signals differently
	// Disable signal handling to avoid conflicts with iOS
	// DO NOT call runtime.LockOSThread() - it causes signal handling issues on iOS!
	iosConsoleLogf("info", "🔵 [application_ios.go] Skipping runtime.LockOSThread() on iOS")

	// Disable all signal handling on iOS
	// iOS apps run in a sandboxed environment where signal handling is restricted
	iosConsoleLogf("info", "🔵 [application_ios.go] END init()")
}

//export init_go
func init_go() {
	iosConsoleLogf("info", "🔵 [application_ios.go] init_go() called from iOS")
	// This is called from the iOS main function
	// to initialize the Go runtime
}

func (a *App) platformRun() {
	iosConsoleLogf("info", "🔵 [application_ios.go] START platformRun()")

	iosConsoleLogf("info", "🔵 [application_ios.go] platformRun called, initializing...")

	// Initialize what we need for the Go side
	iosConsoleLogf("info", "🔵 [application_ios.go] About to call C.ios_app_init()")
	C.ios_app_init()
	iosConsoleLogf("info", "🔵 [application_ios.go] C.ios_app_init() returned")

	// Wait a bit for the UI to be ready (UIApplicationMain is running in main thread)
	// The app delegate's didFinishLaunchingWithOptions will be called
	iosConsoleLogf("info", "🔵 [application_ios.go] Waiting for UI to be ready...")
	time.Sleep(2 * time.Second) // Give the app delegate time to initialize

	// The WebView will be created when the window runs (via app.Window.NewWithOptions in main.go)
	iosConsoleLogf("info", "🔵 [application_ios.go] WebView creation will be handled by window manager")

	// UIApplicationMain is running in the main thread (called from main.m)
	// We just need to keep the Go runtime alive
	iosConsoleLogf("info", "🔵 [application_ios.go] Blocking to keep Go runtime alive...")
	select {} // Block forever
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

	// Add iOS marker for HTML logger
	iosConsoleLogf("info", "[iOS-%s] %s", goSource, goMessage)

	if globalApplication != nil && globalApplication.Logger != nil {
		globalApplication.info("iOS log", "source", goSource, "message", goMessage)
	}
}

// Platform-specific app implementation for iOS
type iosApp struct {
	parent *App
}

func newPlatformApp(app *App) *iosApp {
	iosConsoleLogf("info", "🔵 [application_ios.go] START newPlatformApp()")
	// iOS initialization
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
	iosConsoleLogf("info", "🔵 [application_ios.go] Input accessory view %s", map[bool]string{true: "DISABLED", false: "ENABLED"}[disable])

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

	iosConsoleLogf("info", "🔵 [application_ios.go] END newPlatformApp() - iosApp created")
	return result
}

func (a *iosApp) run() error {
	iosConsoleLogf("info", "🔵 [application_ios.go] START iosApp.run()")

	// Initialize and create the WebView
	// UIApplicationMain is already running in the main thread (from main.m)
	// Wire common events (e.g. map ApplicationDidFinishLaunching → Common.ApplicationStarted)
	a.setupCommonEvents()
	iosConsoleLogf("info", "🔵 [application_ios.go] About to call parent.platformRun()")
	a.parent.platformRun()

	// platformRun blocks forever with select{}
	// If we get here, something went wrong
	iosConsoleLogf("error", "🔵 [application_ios.go] ERROR: platformRun() returned unexpectedly")
	return nil
}

func (a *iosApp) destroy() {
	iosConsoleLogf("info", "🔵 [application_ios.go] iosApp.destroy() called")
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

func (a *iosApp) on(_ uint) {
	// iOS event handling
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

// ServeAssetRequest handles requests from the WebView
//
//export ServeAssetRequest
func ServeAssetRequest(windowID C.uint, urlSchemeTask unsafe.Pointer) {
	iosConsoleLogf("info", "[iOS-ServeAssetRequest] 🔵 Called with windowID=%d", windowID)

	// Route the request through the webviewRequests channel to use the asset server
	go func() {
		iosConsoleLogf("info", "[iOS-ServeAssetRequest] 🔵 Inside goroutine")

		// Use the webview package's NewRequest to wrap the task pointer
		req := webview.NewRequest(urlSchemeTask)
		url, _ := req.URL()

		// Log every single request with clear markers
		iosConsoleLogf("info", "===============================================")
		iosConsoleLogf("info", "[iOS-REQUEST] 🌐 RECEIVED REQUEST FOR: %s", url)
		iosConsoleLogf("info", "===============================================")

		// Special CSS logging with big markers
		if strings.Contains(url, ".css") || strings.Contains(url, "style") {
			iosConsoleLogf("warn", "🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨")
			iosConsoleLogf("warn", "[iOS-CSS] CSS FILE REQUESTED: %s", url)
			iosConsoleLogf("warn", "🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨")
		}

		// Log images separately
		if strings.Contains(url, ".png") || strings.Contains(url, ".jpg") || strings.Contains(url, ".svg") {
			iosConsoleLogf("info", "[iOS-IMAGE] 🇼 %s", url)
		}

		// Log JS files
		if strings.Contains(url, ".js") {
			iosConsoleLogf("info", "[iOS-JS] ⚙️ %s", url)
		}

		// Try to resolve the window name from the window ID so the AssetServer
		// receives both x-wails-window-id and x-wails-window-name headers.
		winName := ""
		if globalApplication != nil {
			if window, ok := globalApplication.Window.GetByID(uint(windowID)); ok && window != nil {
				winName = window.Name()
			} else {
				iosConsoleLogf("warn", "[iOS-ServeAssetRequest] 🟠 Could not resolve window name for id=%d", windowID)
			}
		}
		if winName != "" {
			iosConsoleLogf("info", "[iOS-ServeAssetRequest] ✅ Resolved window name: %s (id=%d)", winName, windowID)
		}

		request := &webViewAssetRequest{
			Request:    req,
			windowId:   uint(windowID),
			windowName: winName,
		}

		// Send through the channel to be handled by the asset server
		iosConsoleLogf("info", "[iOS-ServeAssetRequest] 🔵 Sending to webviewRequests channel")
		webviewRequests <- request
		iosConsoleLogf("info", "[iOS-ServeAssetRequest] 🔵 Request sent to channel successfully")
	}()
}

// HandleJSMessage handles messages from JavaScript
//
//export HandleJSMessage
func HandleJSMessage(windowID C.uint, message *C.char) {
	msg := C.GoString(message)

	// Try to parse as JSON first
	var msgData map[string]interface{}
	if err := json.Unmarshal([]byte(msg), &msgData); err == nil && msgData != nil {
		if name, ok := msgData["name"].(string); ok && name != "" {
			// Special handling for asset debug messages
			if name == "asset-debug" {
				if data, ok := msgData["data"].(map[string]interface{}); ok {
					iosConsoleLogf("info", "🔍 CLIENT ASSET DEBUG: %s %s - %s (status: %v)",
						data["type"], data["name"], data["src"], data["status"])
					if contentType, ok := data["contentType"].(map[string]interface{}); ok {
						iosConsoleLogf("info", "🔍 CLIENT CONTENT-TYPE: %s = %v", data["name"], contentType)
					}
					if code, ok := data["code"].(map[string]interface{}); ok {
						iosConsoleLogf("info", "🔍 CLIENT HTTP CODE: %s = %v", data["name"], code)
					}
					if errorMsg, ok := data["error"].(map[string]interface{}); ok {
						iosConsoleLogf("error", "🔍 CLIENT ERROR: %s = %v", data["name"], errorMsg)
					}
				}
				return // Don't send asset-debug messages to the main event system
			}

			if globalApplication != nil {
				globalApplication.info("HandleJSMessage received from client", "name", name)
			}
			windowMessageBuffer <- &windowMessage{
				windowId: uint(windowID),
				message:  name,
			}
			return
		}
		// Fallback for structured payloads without a "name" field
		if name, ok := msgData["message"].(string); ok && name != "" {
			if globalApplication != nil {
				globalApplication.info("HandleJSMessage received raw message field from client", "name", name)
			}
			windowMessageBuffer <- &windowMessage{
				windowId: uint(windowID),
				message:  name,
			}
			return
		}
	} else {
		if globalApplication != nil {
			globalApplication.error("[HandleJSMessage] Failed to parse JSON: %v", err)
		}
		iosConsoleLogf("warn", "🔍 RAW JS MESSAGE (unparsed JSON): %s", msg)
	}

	// If not JSON or JSON without name/message, treat the entire payload as a string event
	if msg != "" {
		if globalApplication != nil {
			globalApplication.info("HandleJSMessage received raw message from client", "message", msg)
		}
		windowMessageBuffer <- &windowMessage{
			windowId: uint(windowID),
			message:  msg,
		}
		return
	}

	iosConsoleLogf("warn", "[HandleJSMessage] Ignored empty JS message")
}

// Note: applicationEvents and windowEvents are already defined in events.go
// We'll use those existing channels

type iosWindowEvent struct {
	WindowID uint
	EventID  uint
}

//export processApplicationEvent
func processApplicationEvent(eventID C.uint, data unsafe.Pointer) {
	iosConsoleLogf("info", "🔵 [application_ios.go] processApplicationEvent called with eventID: %d", eventID)

	// Create and send the application event
	event := newApplicationEvent(events.ApplicationEventType(eventID))

	// Send to the applicationEvents channel for processing
	applicationEvents <- event

	iosConsoleLogf("info", "🔵 [application_ios.go] Application event sent to channel: %d", eventID)
}

//export processWindowEvent
func processWindowEvent(windowID C.uint, eventID C.uint) {
	// For now, just log the event
	iosConsoleLogf("info", "iOS: Window event received - Window: %d, Event: %d", windowID, eventID)
	windowEvents <- &windowEvent{
		WindowID: uint(windowID),
		EventID:  uint(eventID),
	}
}

//export hasListeners
func hasListeners(eventID C.uint) C.bool {
	// For now, return true to enable all events
	// TODO: Check actual listener registration
	return C.bool(true)
}
