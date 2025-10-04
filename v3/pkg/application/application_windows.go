//go:build windows

package application

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"github.com/wailsapp/go-webview2/webviewloader"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"

	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

var (
	wmTaskbarCreated = w32.RegisterWindowMessage(w32.MustStringToUTF16Ptr("TaskbarCreated"))
)

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
	isCurrentlyDarkMode bool
	currentWindowID     uint

	// Restart taskbar flag
	restartingTaskbar atomic.Bool
}

func (m *windowsApp) isDarkMode() bool {
	return w32.IsCurrentlyDarkMode()
}

func (m *windowsApp) getAccentColor() string {
	accentColor, err := w32.GetAccentColor()
	if err != nil {
		m.parent.error("failed to get accent color: %w", err)
		return "rgb(0,122,255)"
	}

	return accentColor
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

func (m *windowsApp) on(_ uint) {
}

func (m *windowsApp) setIcon(_ []byte) {
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
		menu = DefaultApplicationMenu()
	}
	menu.Update()

	m.parent.applicationMenu = menu
}

func (m *windowsApp) run() error {
	m.setupCommonEvents()
	for eventID := range m.parent.applicationEventListeners {
		m.on(eventID)
	}
	// EmitEvent application started event
	applicationEvents <- &ApplicationEvent{
		Id:  uint(events.Windows.ApplicationStarted),
		ctx: blankApplicationEventContext,
	}

	if len(os.Args) == 2 { // Case: program + 1 argument
		arg1 := os.Args[1]
		// Check if the argument is likely a URL from a custom protocol invocation
		if strings.Contains(arg1, "://") {
			m.parent.info("Application launched with argument, potentially a URL from custom protocol", "url", arg1)
			eventContext := newApplicationEventContext()
			eventContext.setURL(arg1)
			applicationEvents <- &ApplicationEvent{
				Id:  uint(events.Common.ApplicationLaunchedWithUrl),
				ctx: eventContext,
			}
		} else {
			// If not a URL-like string, check for file association
			if m.parent.options.FileAssociations != nil {
				ext := filepath.Ext(arg1)
				if slices.Contains(m.parent.options.FileAssociations, ext) {
					m.parent.info("Application launched with file via file association", "file", arg1)
					eventContext := newApplicationEventContext()
					eventContext.setOpenedWithFile(arg1)
					applicationEvents <- &ApplicationEvent{
						Id:  uint(events.Common.ApplicationOpenedWithFile),
						ctx: eventContext,
					}
				}
			}
		}
	} else if len(os.Args) > 2 {
		// Log if multiple arguments are passed, though typical protocol/file launch is a single arg.
		m.parent.info("Application launched with multiple arguments", "args", os.Args[1:])
	}

	_ = m.runMainLoop()

	return nil
}

func (m *windowsApp) destroy() {
	if !globalApplication.shouldQuit() {
		return
	}
	globalApplication.cleanup()
	// destroy the main thread window
	w32.DestroyWindow(m.mainThreadWindowHWND)
	// Post a quit message to the main thread
	w32.PostQuitMessage(0)
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
	m.windowClass.ClassName = w32.MustStringToUTF16Ptr(m.parent.options.Windows.WndClass)
	m.windowClass.MenuName = nil
	m.windowClass.IconSm = icon

	if ret := w32.RegisterClassEx(&m.windowClass); ret == 0 {
		panic(syscall.GetLastError())
	}
	m.isCurrentlyDarkMode = w32.IsCurrentlyDarkMode()
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

	// Handle the main thread window
	// Quit the application if requested
	// Reprocess and cache screens when display settings change
	if hwnd == m.mainThreadWindowHWND {
		if msg == w32.WM_ENDSESSION || msg == w32.WM_DESTROY || msg == w32.WM_CLOSE {
			globalApplication.Quit()
		}
		if msg == w32.WM_DISPLAYCHANGE || (msg == w32.WM_SETTINGCHANGE && wParam == w32.SPI_SETWORKAREA) {
			err := m.processAndCacheScreens()
			if err != nil {
				m.parent.handleError(err)
			}
		}
	}

	switch msg {
	case wmTaskbarCreated:
		if m.restartingTaskbar.Load() {
			break
		}
		m.restartingTaskbar.Store(true)
		m.reshowSystrays()
		go func() {
			// 1 second debounce
			time.Sleep(1000)
			m.restartingTaskbar.Store(false)
		}()
	case w32.WM_SETTINGCHANGE:
		settingChanged := w32.UTF16PtrToString((*uint16)(unsafe.Pointer(lParam)))
		if settingChanged == "ImmersiveColorSet" {
			isDarkMode := w32.IsCurrentlyDarkMode()
			if isDarkMode != m.isCurrentlyDarkMode {
				eventContext := newApplicationEventContext()
				eventContext.setIsDarkMode(isDarkMode)
				applicationEvents <- &ApplicationEvent{
					Id:  uint(events.Windows.SystemThemeChanged),
					ctx: eventContext,
				}
				m.isCurrentlyDarkMode = isDarkMode
			}
		}
		return 0
	case w32.WM_POWERBROADCAST:
		switch wParam {
		case w32.PBT_APMPOWERSTATUSCHANGE:
			applicationEvents <- newApplicationEvent(events.Windows.APMPowerStatusChange)
		case w32.PBT_APMSUSPEND:
			applicationEvents <- newApplicationEvent(events.Windows.APMSuspend)
		case w32.PBT_APMRESUMEAUTOMATIC:
			applicationEvents <- newApplicationEvent(events.Windows.APMResumeAutomatic)
		case w32.PBT_APMRESUMESUSPEND:
			applicationEvents <- newApplicationEvent(events.Windows.APMResumeSuspend)
		case w32.PBT_POWERSETTINGCHANGE:
			applicationEvents <- newApplicationEvent(events.Windows.APMPowerSettingChange)
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

func (m *windowsApp) reshowSystrays() {
	m.systrayMapLock.Lock()
	defer m.systrayMapLock.Unlock()
	for _, systray := range m.systrayMap {
		systray.reshow()
	}
}

func setupDPIAwareness() error {
	// https://learn.microsoft.com/en-us/windows/win32/hidpi/setting-the-default-dpi-awareness-for-a-process
	// https://learn.microsoft.com/en-us/windows/win32/hidpi/high-dpi-desktop-application-development-on-windows

	if w32.HasSetProcessDpiAwarenessContextFunc() {
		// This is most recent version with the best results
		// supported beginning with Windows 10, version 1703
		return w32.SetProcessDpiAwarenessContext(w32.DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2)
	}

	if w32.HasSetProcessDpiAwarenessFunc() {
		// Supported beginning with Windows 8.1
		return w32.SetProcessDpiAwareness(w32.PROCESS_PER_MONITOR_DPI_AWARE)
	}

	if w32.HasSetProcessDPIAwareFunc() {
		// If none of the above is supported, fallback to SetProcessDPIAware
		// which is supported beginning with Windows Vista
		return w32.SetProcessDPIAware()
	}

	return errors.New("no DPI awareness method supported")
}

func (m *windowsApp) setStartAtLogin(enabled bool) error {
	// Get the current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Resolve any symbolic links to get the real path
	realPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	// Validate that the executable exists
	if _, err := os.Stat(realPath); os.IsNotExist(err) {
		return fmt.Errorf("executable does not exist at path: %s", realPath)
	}

	// Get the registry key name (derive from executable name)
	exeName := filepath.Base(realPath)
	keyName := strings.TrimSuffix(exeName, filepath.Ext(exeName)) // Remove extension

	// Validate key name - avoid problematic characters
	if strings.ContainsAny(keyName, "\\/:*?\"<>|") {
		return fmt.Errorf("invalid executable name for registry key: %s", keyName)
	}

	if enabled {
		return m.addToStartup(keyName, realPath)
	}
	return m.removeFromStartup(keyName)
}

func (m *windowsApp) startsAtLogin() (bool, error) {
	// Get the current executable path
	exePath, err := os.Executable()
	if err != nil {
		return false, fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the registry key name
	exeName := filepath.Base(exePath)
	keyName := strings.TrimSuffix(exeName, filepath.Ext(exeName))

	return m.isInStartup(keyName)
}

func (m *windowsApp) addToStartup(keyName, exePath string) error {
	key, err := w32.RegOpenKeyEx(
		w32.HKEY_CURRENT_USER,
		w32.StringToUTF16Ptr(`Software\Microsoft\Windows\CurrentVersion\Run`),
		0,
		w32.KEY_SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer w32.RegCloseKey(key)

	// Set the registry value
	err = w32.RegSetValueEx(
		key,
		w32.StringToUTF16Ptr(keyName),
		0,
		w32.REG_SZ,
		(*byte)(unsafe.Pointer(w32.StringToUTF16Ptr(exePath))),
		uint32((len(exePath)+1)*2), // UTF-16 byte length
	)
	if err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	return nil
}

func (m *windowsApp) removeFromStartup(keyName string) error {
	key, err := w32.RegOpenKeyEx(
		w32.HKEY_CURRENT_USER,
		w32.StringToUTF16Ptr(`Software\Microsoft\Windows\CurrentVersion\Run`),
		0,
		w32.KEY_SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer w32.RegCloseKey(key)

	// Delete the registry value
	err = w32.RegDeleteValue(key, w32.StringToUTF16Ptr(keyName))
	if err != nil && err != w32.ERROR_FILE_NOT_FOUND {
		return fmt.Errorf("failed to delete registry value: %w", err)
	}

	return nil
}

func (m *windowsApp) isInStartup(keyName string) (bool, error) {
	key, err := w32.RegOpenKeyEx(
		w32.HKEY_CURRENT_USER,
		w32.StringToUTF16Ptr(`Software\Microsoft\Windows\CurrentVersion\Run`),
		0,
		w32.KEY_QUERY_VALUE,
	)
	if err != nil {
		return false, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer w32.RegCloseKey(key)

	// Query the registry value
	var valueType uint32
	var dataSize uint32
	err = w32.RegQueryValueEx(
		key,
		w32.StringToUTF16Ptr(keyName),
		nil,
		&valueType,
		nil,
		&dataSize,
	)

	if err == w32.ERROR_FILE_NOT_FOUND {
		return false, nil // Key doesn't exist, not in startup
	}
	if err != nil {
		return false, fmt.Errorf("failed to query registry value: %w", err)
	}

	return true, nil
}

func newPlatformApp(app *App) *windowsApp {

	err := setupDPIAwareness()
	if err != nil {
		app.handleError(err)
	}

	result := &windowsApp{
		parent:     app,
		instance:   w32.GetModuleHandle(""),
		windowMap:  make(map[w32.HWND]*windowsWebviewWindow),
		systrayMap: make(map[w32.HWND]*windowsSystemTray),
	}

	err = result.processAndCacheScreens()
	if err != nil {
		app.handleFatalError(err)
	}

	result.init()
	result.initMainLoop()

	return result
}

func (a *App) logPlatformInfo() {
	var args []any
	args = append(args, "Go-WebView2Loader", webviewloader.UsingGoWebview2Loader)
	webviewVersion, err := webviewloader.GetAvailableCoreWebView2BrowserVersionString(
		a.options.Windows.WebviewBrowserPath,
	)
	if err != nil {
		args = append(args, "WebView2", "Error: "+err.Error())
	} else {
		args = append(args, "WebView2", webviewVersion)
	}

	osInfo, _ := operatingsystem.Info()
	args = append(args, osInfo.AsLogSlice()...)

	a.info("Platform Info:", args...)
}

func (a *App) platformEnvironment() map[string]any {
	result := map[string]any{}
	webviewVersion, _ := webviewloader.GetAvailableCoreWebView2BrowserVersionString(
		a.options.Windows.WebviewBrowserPath,
	)
	result["Go-WebView2Loader"] = webviewloader.UsingGoWebview2Loader
	result["WebView2"] = webviewVersion
	return result
}

func fatalHandler(errFunc func(error)) {
	w32.Fatal = errFunc
	return
}
