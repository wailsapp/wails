//go:build windows

package application

import (
	"github.com/wailsapp/wails/v3/pkg/events"
	"syscall"
	"unsafe"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/w32"
)

var windowClassName = lo.Must(syscall.UTF16PtrFromString("WailsWebviewWindow"))

type windowsApp struct {
	parent *App

	instance w32.HINSTANCE

	windowMap map[w32.HWND]*windowsWebviewWindow

	mainThreadID         w32.HANDLE
	mainThreadWindowHWND w32.HWND

	// system theme
	isDarkMode bool
}

func (m *windowsApp) getPrimaryScreen() (*Screen, error) {
	//TODO implement me
	panic("implement me")
}

func (m *windowsApp) getScreens() ([]*Screen, error) {
	//TODO implement me
	panic("implement me")
}

func (m *windowsApp) hide() {
}

func (m *windowsApp) show() {
}

func (m *windowsApp) on(eventID uint) {
	//C.registerListener(C.uint(eventID))
}

func (m *windowsApp) setIcon(icon []byte) {
	//C.setApplicationIcon(unsafe.Pointer(&icon[0]), C.int(len(icon)))
}

func (m *windowsApp) name() string {
	//appName := C.getAppName()
	//defer C.free(unsafe.Pointer(appName))
	//return C.GoString(appName)
	return ""
}

func (m *windowsApp) getCurrentWindowID() uint {
	//return uint(C.getCurrentWindowID())
	return uint(0)
}

func (m *windowsApp) setApplicationMenu(menu *Menu) {
	if menu == nil {
		// Create a default menu for mac
		menu = defaultApplicationMenu()
	}
	menu.Update()

	// Convert impl to macosMenu object
	//m.applicationMenu = (menu.impl).(*macosMenu).nsMenu
	//C.setApplicationMenu(m.applicationMenu)
}

func (m *windowsApp) run() error {
	// Add a hook to the ApplicationDidFinishLaunching event
	//m.parent.On(events.Mac.ApplicationDidFinishLaunching, func() {
	//	C.setApplicationShouldTerminateAfterLastWindowClosed(C.bool(m.parent.options.Mac.ApplicationShouldTerminateAfterLastWindowClosed))
	//	C.setActivationPolicy(C.int(m.parent.options.Mac.ActivationPolicy))
	//	C.activateIgnoringOtherApps()
	//})
	// setup event listeners
	for eventID := range m.parent.applicationEventListeners {
		m.on(eventID)
	}

	_ = m.runMainLoop()

	//C.run()
	return nil
}

func (m *windowsApp) destroy() {
	//C.destroyApp()
}

func (m *windowsApp) init() {
	// Register the window class

	icon := w32.LoadIconWithResourceID(m.instance, w32.IDI_APPLICATION)

	var wc w32.WNDCLASSEX
	wc.Size = uint32(unsafe.Sizeof(wc))
	wc.Style = w32.CS_HREDRAW | w32.CS_VREDRAW
	wc.WndProc = syscall.NewCallback(m.wndProc)
	wc.Instance = m.instance
	wc.Background = w32.COLOR_BTNFACE + 1
	wc.Icon = icon
	wc.Cursor = w32.LoadCursorWithResourceID(0, w32.IDC_ARROW)
	wc.ClassName = windowClassName
	wc.MenuName = nil
	wc.IconSm = icon

	if ret := w32.RegisterClassEx(&wc); ret == 0 {
		panic(syscall.GetLastError())
	}

	m.isDarkMode = w32.IsCurrentlyDarkMode()
}

func (m *windowsApp) wndProc(hwnd w32.HWND, msg uint32, wParam, lParam uintptr) uintptr {

	// Handle the invoke callback
	if msg == wmInvokeCallback {
		m.invokeCallback(wParam, lParam)
		return 0
	}
	switch msg {
	case w32.WM_SETTINGCHANGE:
		settingChanged := w32.UTF16PtrToString((*uint16)(unsafe.Pointer(lParam)))
		if settingChanged == "ImmersiveColorSet" {
			isDarkMode := w32.IsCurrentlyDarkMode()
			if isDarkMode != m.isDarkMode {
				applicationEvents <- uint(events.Windows.SystemThemeChanged)
				m.isDarkMode = isDarkMode
			}
		}
		return 0
	case w32.WM_POWERBROADCAST:
		switch wParam {
		case w32.PBT_APMPOWERSTATUSCHANGE:
			applicationEvents <- uint(events.Windows.APMPowerStatusChange)
		case w32.PBT_APMSUSPEND:
			applicationEvents <- uint(events.Windows.APMSuspend)
		case w32.PBT_APMRESUMEAUTOMATIC:
			applicationEvents <- uint(events.Windows.APMResumeAutomatic)
		case w32.PBT_APMRESUMESUSPEND:
			applicationEvents <- uint(events.Windows.APMResumeSuspend)
		case w32.PBT_POWERSETTINGCHANGE:
			applicationEvents <- uint(events.Windows.APMPowerSettingChange)
		}
		return 0
	}

	if window, ok := m.windowMap[hwnd]; ok {
		return window.WndProc(msg, wParam, lParam)
	}

	// Dispatch the message to the appropriate window

	return w32.DefWindowProc(hwnd, msg, wParam, lParam)
}

func (m *windowsApp) registerWindow(result *windowsWebviewWindow) {
	m.windowMap[result.hwnd] = result
}

func (m *windowsApp) unregisterWindow(w *windowsWebviewWindow) {
	delete(m.windowMap, w.hwnd)

	// If this was the last window...
	if len(m.windowMap) == 0 {
		w32.PostQuitMessage(0)
	}
}

func newPlatformApp(app *App) *windowsApp {
	result := &windowsApp{
		parent:    app,
		instance:  w32.GetModuleHandle(""),
		windowMap: make(map[w32.HWND]*windowsWebviewWindow),
	}

	result.init()
	result.initMainLoop()

	return result
}
