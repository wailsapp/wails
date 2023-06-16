//go:build windows

package application

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"strconv"
	"sync"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/w32"

	"github.com/samber/lo"
)

var windowClassName = lo.Must(syscall.UTF16PtrFromString("WailsWebviewWindow"))

type windowsApp struct {
	parent *App

	windowClass w32.WNDCLASSEX
	instance    w32.HINSTANCE

	windowMap     map[w32.HWND]*windowsWebviewWindow
	windowMapLock sync.RWMutex

	systrayMap     map[w32.HMENU]*windowsSystemTray
	systrayMapLock sync.RWMutex

	mainThreadID         w32.HANDLE
	mainThreadWindowHWND w32.HWND

	// Windows hidden by application.Hide()
	hiddenWindows []*windowsWebviewWindow
	focusedWindow w32.HWND

	// system theme
	isDarkMode      bool
	currentWindowID uint
}

func (m *windowsApp) isOnMainThread() bool {
	return m.mainThreadID == w32.GetCurrentThreadId()
}

func (m *windowsApp) GetFlags(options Options) map[string]any {
	if options.Flags == nil {
		options.Flags = make(map[string]any)
	}
	options.Flags["system"] = map[string]any{
		"resizeHandleWidth":  w32.GetSystemMetrics(w32.SM_CXSIZEFRAME),
		"resizeHandleHeight": w32.GetSystemMetrics(w32.SM_CYSIZEFRAME),
	}
	return options.Flags
}

func (m *windowsApp) getWindowForHWND(hwnd w32.HWND) *windowsWebviewWindow {
	m.windowMapLock.RLock()
	defer m.windowMapLock.RUnlock()
	return m.windowMap[hwnd]
}

func getNativeApplication() *windowsApp {
	return globalApplication.impl.(*windowsApp)
}

func (m *windowsApp) getPrimaryScreen() (*Screen, error) {
	screens, err := m.getScreens()
	if err != nil {
		return nil, err
	}
	for _, screen := range screens {
		if screen.IsPrimary {
			return screen, nil
		}
	}
	return nil, fmt.Errorf("no primary screen found")
}

func (m *windowsApp) getScreens() ([]*Screen, error) {
	allScreens, err := w32.GetAllScreens()
	if err != nil {
		return nil, err
	}
	// Convert result to []*Screen
	screens := make([]*Screen, len(allScreens))
	for id, screen := range allScreens {
		x := int(screen.MONITORINFOEX.RcMonitor.Left)
		y := int(screen.MONITORINFOEX.RcMonitor.Top)
		right := int(screen.MONITORINFOEX.RcMonitor.Right)
		bottom := int(screen.MONITORINFOEX.RcMonitor.Bottom)
		width := right - x
		height := bottom - y
		screens[id] = &Screen{
			ID:     strconv.Itoa(id),
			Name:   windows.UTF16ToString(screen.MONITORINFOEX.SzDevice[:]),
			X:      x,
			Y:      y,
			Size:   Size{Width: width, Height: height},
			Bounds: Rect{X: x, Y: y, Width: width, Height: height},
			WorkArea: Rect{
				X:      int(screen.MONITORINFOEX.RcWork.Left),
				Y:      int(screen.MONITORINFOEX.RcWork.Top),
				Width:  int(screen.MONITORINFOEX.RcWork.Right - screen.MONITORINFOEX.RcWork.Left),
				Height: int(screen.MONITORINFOEX.RcWork.Bottom - screen.MONITORINFOEX.RcWork.Top),
			},
			IsPrimary: screen.IsPrimary,
			Scale:     screen.Scale,
			Rotation:  0,
		}
	}
	return screens, nil
}

func (m *windowsApp) hide() {
	// Get the current focussed window
	m.focusedWindow = w32.GetForegroundWindow()

	// Iterate over all windows and hide them if they aren't already hidden
	for _, window := range m.windowMap {
		if window.isVisible() {
			// Add to hidden windows
			m.hiddenWindows = append(m.hiddenWindows, window)
			window.hide()
		}
	}
	// Switch focus to the next application
	hwndNext := w32.GetWindow(m.mainThreadWindowHWND, w32.GW_HWNDNEXT)
	w32.SetForegroundWindow(hwndNext)
}

func (m *windowsApp) show() {
	// Iterate over all windows and show them if they were previously hidden
	for _, window := range m.hiddenWindows {
		window.show()
	}
	// Show the foreground window
	w32.SetForegroundWindow(m.focusedWindow)
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
	return m.currentWindowID
}

func (m *windowsApp) setApplicationMenu(menu *Menu) {
	if menu == nil {
		// Create a default menu for windows
		menu = defaultApplicationMenu()
	}
	menu.Update()

	m.parent.ApplicationMenu = menu
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

	m.windowClass.Size = uint32(unsafe.Sizeof(m.windowClass))
	m.windowClass.Style = w32.CS_HREDRAW | w32.CS_VREDRAW
	m.windowClass.WndProc = syscall.NewCallback(m.wndProc)
	m.windowClass.Instance = m.instance
	m.windowClass.Background = w32.COLOR_BTNFACE + 1
	m.windowClass.Icon = icon
	m.windowClass.Cursor = w32.LoadCursorWithResourceID(0, w32.IDC_ARROW)
	m.windowClass.ClassName = windowClassName
	m.windowClass.MenuName = nil
	m.windowClass.IconSm = icon

	if ret := w32.RegisterClassEx(&m.windowClass); ret == 0 {
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

	// If the WndProcInterceptor is set in options, pass the message on
	if m.parent.options.Windows.WndProcInterceptor != nil {
		returnValue, shouldReturn := m.parent.options.Windows.WndProcInterceptor(hwnd, msg, wParam, lParam)
		if shouldReturn {
			return returnValue
		}
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

	m.systrayMapLock.Lock()
	systray, ok := m.systrayMap[hwnd]
	m.systrayMapLock.Unlock()
	if ok {
		return systray.wndProc(msg, wParam, lParam)
	}

	// Dispatch the message to the appropriate window

	return w32.DefWindowProc(hwnd, msg, wParam, lParam)
}

func (m *windowsApp) registerWindow(result *windowsWebviewWindow) {
	m.windowMapLock.Lock()
	m.windowMap[result.hwnd] = result
	m.windowMapLock.Unlock()
}

func (m *windowsApp) registerSystemTray(result *windowsSystemTray) {
	m.systrayMapLock.Lock()
	defer m.systrayMapLock.Unlock()
	m.systrayMap[result.hwnd] = result
}

func (m *windowsApp) unregisterSystemTray(result *windowsSystemTray) {
	m.systrayMapLock.Lock()
	defer m.systrayMapLock.Unlock()
	delete(m.systrayMap, result.hwnd)
}

func (m *windowsApp) unregisterWindow(w *windowsWebviewWindow) {
	m.windowMapLock.Lock()
	delete(m.windowMap, w.hwnd)
	m.windowMapLock.Unlock()

	// If this was the last window...
	if len(m.windowMap) == 0 && !m.parent.options.Windows.DisableQuitOnLastWindowClosed {
		w32.PostQuitMessage(0)
	}
}

func newPlatformApp(app *App) *windowsApp {
	err := w32.SetProcessDPIAware()
	if err != nil {
		println("Fatal error in application initialisation: ", err.Error())
		os.Exit(1)
	}

	result := &windowsApp{
		parent:     app,
		instance:   w32.GetModuleHandle(""),
		windowMap:  make(map[w32.HWND]*windowsWebviewWindow),
		systrayMap: make(map[w32.HWND]*windowsSystemTray),
	}

	result.init()
	result.initMainLoop()

	return result
}
