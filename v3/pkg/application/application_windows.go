//go:build windows

package application

import (
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/w32"
	"syscall"
	"unsafe"
)

var windowClassName = lo.Must(syscall.UTF16PtrFromString("WailsWebviewWindow"))

type windowsApp struct {
	parent *App

	instance w32.HINSTANCE
}

func (m *windowsApp) dispatchOnMainThread(id uint) {
	//TODO implement me
	panic("implement me")
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
}

func (m *windowsApp) wndProc(hwnd w32.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case w32.WM_SIZE, w32.WM_PAINT:
		return 0
	case w32.WM_CLOSE:
		w32.PostQuitMessage(0)
		return 0
	}
	return w32.DefWindowProc(hwnd, msg, wParam, lParam)
}

func (m *windowsApp) runMainLoop() int {
	msg := (*w32.MSG)(unsafe.Pointer(w32.GlobalAlloc(0, uint32(unsafe.Sizeof(w32.MSG{})))))
	defer w32.GlobalFree(w32.HGLOBAL(unsafe.Pointer(m)))

	for w32.GetMessage(msg, 0, 0, 0) != 0 {
		w32.TranslateMessage(msg)
		w32.DispatchMessage(msg)
	}

	return int(msg.WParam)
}

func newPlatformApp(app *App) *windowsApp {
	result := &windowsApp{
		parent:   app,
		instance: w32.GetModuleHandle(""),
	}

	result.init()

	return result
}
