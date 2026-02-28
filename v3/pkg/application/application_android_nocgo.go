//go:build android && !cgo && !server

package application

import (
	"fmt"
	"sync"
	"unsafe"

	"encoding/json"
)

var (
	// Global reference to the app for JNI callbacks
	globalApp     *App
	globalAppLock sync.RWMutex

	// JNI environment and class references
	javaVM       unsafe.Pointer
	bridgeObject unsafe.Pointer
)

func init() {
	androidLogf("info", "🤖 [application_android.go] init() called")
}

func androidLogf(level string, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	// For now, just use println - we'll connect to Android's Log.* later
	println(fmt.Sprintf("[Android/%s] %s", level, msg))
}

func (a *App) platformRun() {
	androidLogf("info", "🤖 [application_android.go] platformRun() called")

	// Store global reference for JNI callbacks
	globalAppLock.Lock()
	globalApp = a
	globalAppLock.Unlock()

	androidLogf("info", "🤖 [application_android.go] Waiting for Android lifecycle...")

	// Block forever - Android manages the app lifecycle via JNI callbacks
	select {}
}

func (a *App) platformQuit() {
	androidLogf("info", "🤖 [application_android.go] platformQuit() called")
	// Android will handle app termination
}

func (a *App) isDarkMode() bool {
	// TODO: Query Android for dark mode status
	return false
}

func (a *App) isWindows() bool {
	return false
}

func androidDeviceInfoViaJNI() map[string]interface{} {
	return map[string]interface{}{
		"platform": "android",
		"model":    "Unknown",
		"version":  "Unknown",
	}
}

// Platform-specific app implementation for Android
type androidApp struct {
	parent *App
}

func newPlatformApp(app *App) *androidApp {
	androidLogf("info", "🤖 [application_android.go] newPlatformApp() called")
	configureNativeTabs(app)
	return &androidApp{
		parent: app,
	}
}

func (a *androidApp) run() error {
	androidLogf("info", "🤖 [application_android.go] androidApp.run() called")

	// Wire common events
	a.setupCommonEvents()

	// Emit application started event
	a.parent.Event.Emit("ApplicationStarted")

	a.parent.platformRun()
	return nil
}

func (a *androidApp) destroy() {
	androidLogf("info", "🤖 [application_android.go] androidApp.destroy() called")
}

func (a *androidApp) setIcon(_ []byte) {
	// Android app icon is set through AndroidManifest.xml
}

func (a *androidApp) name() string {
	return a.parent.options.Name
}

func (a *androidApp) GetFlags(options Options) map[string]any {
	return nil
}

func (a *androidApp) getAccentColor() string {
	return ""
}

func (a *androidApp) getCurrentWindowID() uint {
	return 0
}

func (a *androidApp) hide() {
	// Android manages app visibility
}

func (a *androidApp) isDarkMode() bool {
	return a.parent.isDarkMode()
}

func (a *androidApp) on(_ uint) {
	// Android event handling
}

func (a *androidApp) setApplicationMenu(_ *Menu) {
	// Android doesn't have application menus in the same way
}

func (a *androidApp) show() {
	// Android manages app visibility
}

func (a *androidApp) showAboutDialog(_ string, _ string, _ []byte) {
	// TODO: Implement Android about dialog
}

func (a *androidApp) getPrimaryScreen() (*Screen, error) {
	screens, err := getScreens()
	if err != nil || len(screens) == 0 {
		return nil, err
	}
	return screens[0], nil
}

func (a *androidApp) getScreens() ([]*Screen, error) {
	return getScreens()
}

func (a *App) logPlatformInfo() {
	// Log Android platform info
	androidLogf("info", "Platform: Android")
}

func (a *App) platformEnvironment() map[string]any {
	return map[string]any{
		"platform": "android",
	}
}

func fatalHandler(errFunc func(error)) {
	// Android fatal handler
}

// Helper functions

func serveAssetForAndroid(app *App, path string) ([]byte, error) {
	// Normalize path
	if path == "" || path == "/" {
		path = "/index.html"
	}

	// TODO: Use the actual asset server to serve the file
	// For now, return a placeholder
	return nil, fmt.Errorf("asset serving not yet implemented: %s", path)
}

func handleMessageForAndroid(app *App, message string) string {
	// Parse the message
	var msg map[string]interface{}
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return fmt.Sprintf(`{"error":"%s"}`, err.Error())
	}

	// TODO: Route to appropriate handler based on message type
	// For now, return success
	return `{"success":true}`
}

func getMimeTypeForPath(path string) string {
	// Simple MIME type detection based on extension
	switch {
	case endsWith(path, ".html"), endsWith(path, ".htm"):
		return "text/html"
	case endsWith(path, ".js"), endsWith(path, ".mjs"):
		return "application/javascript"
	case endsWith(path, ".css"):
		return "text/css"
	case endsWith(path, ".json"):
		return "application/json"
	case endsWith(path, ".png"):
		return "image/png"
	case endsWith(path, ".jpg"), endsWith(path, ".jpeg"):
		return "image/jpeg"
	case endsWith(path, ".gif"):
		return "image/gif"
	case endsWith(path, ".svg"):
		return "image/svg+xml"
	case endsWith(path, ".ico"):
		return "image/x-icon"
	case endsWith(path, ".woff"):
		return "font/woff"
	case endsWith(path, ".woff2"):
		return "font/woff2"
	case endsWith(path, ".ttf"):
		return "font/ttf"
	default:
		return "application/octet-stream"
	}
}

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// executeJavaScript is a stub for non-cgo builds
func executeJavaScript(js string) {
	androidLogf("warn", "executeJavaScript called but cgo is not enabled")
}

func configureNativeTabs(app *App) {
	if app == nil {
		return
	}

	if len(app.options.Android.NativeTabsItems) > 0 || app.options.Android.EnableNativeTabs {
		androidLogf("warn", "native tabs are not available without cgo on Android")
	}
}
