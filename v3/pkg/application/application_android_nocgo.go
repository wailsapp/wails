//go:build android && !cgo && !server

package application

// This file keeps GOOS=android builds compiling without cgo (used by
// tooling such as `wails3 generate bindings`). A real Android app is always
// built with CGO_ENABLED=1 — see application_android.go for the JNI bridge.

import (
	"fmt"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/events"
)

var (
	globalApp     *App
	globalAppLock sync.RWMutex

	androidMainFunc func()
	androidMainLock sync.Mutex
)

func androidLogf(level string, format string, a ...interface{}) {
	println(fmt.Sprintf("[Android/%s] %s", level, fmt.Sprintf(format, a...)))
}

func androidDebugLogf(format string, a ...interface{}) {
	if androidVerboseLogging {
		androidLogf("debug", format, a...)
	}
}

// RegisterAndroidMain registers the main function to be called when the
// Android app starts. Call it from init() in your main package.
func RegisterAndroidMain(mainFunc func()) {
	androidMainLock.Lock()
	defer androidMainLock.Unlock()
	androidMainFunc = mainFunc
}

// Go-level bridge call API stubs (no JNI without cgo)

func androidBridgeString(method string) (string, bool) {
	return "", false
}

func androidBridgeVoidString(method string, arg string) {}

func androidBridgeVoidInt(method string, v int) {}

func androidBridgeVoidIntString(method string, id int, arg string) {}

func androidBridgeBool(method string) bool {
	return false
}

func executeJavaScript(js string) {
	androidLogf("warn", "executeJavaScript called but cgo is not enabled")
}

func (a *App) platformRun() {
	globalAppLock.Lock()
	globalApp = a
	globalAppLock.Unlock()

	applicationEvents <- newApplicationEvent(events.Android.ActivityCreated)

	// Block forever - Android manages the app lifecycle via JNI callbacks
	select {}
}

func (a *App) platformQuit() {
}

func (a *App) isDarkMode() bool {
	return false
}

func (a *App) isWindows() bool {
	return false
}

// Platform-specific app implementation for Android
type androidApp struct {
	parent *App
}

func newPlatformApp(app *App) *androidApp {
	return &androidApp{
		parent: app,
	}
}

func (a *androidApp) run() error {
	a.setupCommonEvents()
	a.parent.platformRun()
	return nil
}

func (a *androidApp) destroy() {
}

func (a *androidApp) setIcon(_ []byte) {
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
}

func (a *androidApp) isDarkMode() bool {
	return a.parent.isDarkMode()
}

func (a *androidApp) on(eventID uint) {
	registerAndroidListener(eventID)
}

func (a *androidApp) setApplicationMenu(_ *Menu) {
}

func (a *androidApp) show() {
}

func (a *androidApp) showAboutDialog(_ string, _ string, _ []byte) {
}

func (a *androidApp) getPrimaryScreen() (*Screen, error) {
	if a.parent.Screen.GetPrimary() == nil {
		screens, err := getScreens()
		if err != nil {
			return nil, err
		}
		if err := a.parent.Screen.LayoutScreens(screens); err != nil {
			return nil, err
		}
	}
	return a.parent.Screen.GetPrimary(), nil
}

func (a *androidApp) getScreens() ([]*Screen, error) {
	if len(a.parent.Screen.GetAll()) == 0 {
		screens, err := getScreens()
		if err != nil {
			return nil, err
		}
		if err := a.parent.Screen.LayoutScreens(screens); err != nil {
			return nil, err
		}
	}
	return a.parent.Screen.GetAll(), nil
}

func (a *App) logPlatformInfo() {
}

func (a *App) platformEnvironment() map[string]any {
	return map[string]any{
		"platform": "android",
	}
}

func fatalHandler(errFunc func(error)) {
}

var (
	androidEventListeners     = make(map[uint]bool)
	androidEventListenersLock sync.RWMutex
)

func registerAndroidListener(eventID uint) {
	androidEventListenersLock.Lock()
	defer androidEventListenersLock.Unlock()
	androidEventListeners[eventID] = true
}
