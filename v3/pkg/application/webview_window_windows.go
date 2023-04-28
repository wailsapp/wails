//go:build windows

package application

import (
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v3/internal/w32"
	"syscall"
	"unsafe"

	"github.com/samber/lo"
)

var showDevTools = func(window unsafe.Pointer) {}

type windowsWebviewWindow struct {
	windowImpl unsafe.Pointer
	parent     *WebviewWindow
	hwnd       w32.HWND
}

func (w *windowsWebviewWindow) setTitle(title string) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setSize(width, height int) {
	x, y := w.position()
	w32.MoveWindow(w.hwnd, x, y, width, height, true)
}

func (w *windowsWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	globalApplication.dispatchOnMainThread(func() {
		position := w32.HWND_NOTOPMOST
		if alwaysOnTop {
			position = w32.HWND_TOPMOST
		}
		w32.SetWindowPos(w.hwnd, position, 0, 0, 0, 0, uint(w32.SWP_NOMOVE|w32.SWP_NOSIZE))
	})
}

func (w *windowsWebviewWindow) setURL(url string) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setResizable(resizable bool) {
	globalApplication.dispatchOnMainThread(func() {
		w.setStyle(resizable, w32.WS_THICKFRAME)
	})
}

func (w *windowsWebviewWindow) setMinSize(width, height int) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setMaxSize(width, height int) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) execJS(js string) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) restore() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setBackgroundColour(color RGBA) {
	w32.SetBackgroundColour(w.hwnd, color.Red, color.Green, color.Blue)
}

func (w *windowsWebviewWindow) run() {
	globalApplication.dispatchOnMainThread(w._run)
}

func (w *windowsWebviewWindow) _run() {
	var exStyle uint
	options := w.parent.options
	exStyle = w32.WS_EX_CONTROLPARENT | w32.WS_EX_APPWINDOW
	if options.BackgroundType != BackgroundTypeSolid {
		exStyle |= w32.WS_EX_NOREDIRECTIONBITMAP
	}
	if options.AlwaysOnTop {
		exStyle |= w32.WS_EX_TOPMOST
	}
	w.hwnd = w32.CreateWindowEx(
		exStyle,
		windowClassName,
		lo.Must(syscall.UTF16PtrFromString(options.Title)),
		w32.WS_OVERLAPPEDWINDOW,
		w32.CW_USEDEFAULT,
		w32.CW_USEDEFAULT,
		options.Width,
		options.Height,
		0,
		0,
		w32.GetModuleHandle(""),
		nil)

	if w.hwnd == 0 {
		panic("Unable to create window")
	}

	if options.DisableResize {
		w.setResizable(false)
	}

	// Icon
	if !options.Windows.DisableIcon {
		// App icon ID is 3
		icon, err := NewIconFromResource(w32.GetModuleHandle(""), uint16(3))
		if err == nil {
			w.setIcon(icon)
		}
	} else {
		w.disableIcon()
	}

	switch options.BackgroundType {
	case BackgroundTypeSolid:
		w.setBackgroundColour(options.BackgroundColour)
	case BackgroundTypeTransparent:
	case BackgroundTypeTranslucent:
		w.setBackdropType(options.Windows.BackdropType)
	}

	if !options.Hidden {
		w.show()
		w.update()
	}
	w.setForeground()
}

func (w *windowsWebviewWindow) center() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) size() (int, int) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setForeground() {
	w32.SetForegroundWindow(w.hwnd)
}

func (w *windowsWebviewWindow) update() {
	w32.UpdateWindow(w.hwnd)
}

func (w *windowsWebviewWindow) width() int {
	rect := w32.GetWindowRect(w.hwnd)
	return int(rect.Right - rect.Left)
}

func (w *windowsWebviewWindow) height() int {
	rect := w32.GetWindowRect(w.hwnd)
	return int(rect.Bottom - rect.Top)
}

func (w *windowsWebviewWindow) position() (int, int) {
	rect := w32.GetWindowRect(w.hwnd)
	return int(rect.Left), int(rect.Right)
}

func (w *windowsWebviewWindow) destroy() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) reload() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) forceReload() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) toggleDevTools() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) zoomReset() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) zoomIn() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) zoomOut() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) getZoom() float64 {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setZoom(zoom float64) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) close() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) zoom() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setHTML(html string) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setPosition(x int, y int) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) on(eventID uint) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) minimise() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) unminimise() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) maximise() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) unmaximise() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) fullscreen() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) unfullscreen() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) isMinimised() bool {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) isMaximised() bool {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) isFullscreen() bool {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) disableSizeConstraints() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setFullscreenButtonEnabled(enabled bool) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) show() {
	w32.ShowWindow(w.hwnd, w32.SW_SHOW)
}

func (w *windowsWebviewWindow) hide() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) getScreen() (*Screen, error) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setFrameless(b bool) {
	//TODO implement me
	panic("implement me")
}

func newWindowImpl(parent *WebviewWindow) *windowsWebviewWindow {
	result := &windowsWebviewWindow{
		parent: parent,
	}

	return result
}

func (w *windowsWebviewWindow) openContextMenu(menu *Menu, data *ContextMenuData) {
	// Create the menu
	thisMenu := newMenuImpl(menu)
	thisMenu.update()
	//C.windowShowMenu(w.nsWindow, thisMenu.nsMenu, C.int(data.X), C.int(data.Y))
}

func (w *windowsWebviewWindow) setStyle(b bool, style int) {
	currentStyle := int(w32.GetWindowLongPtr(w.hwnd, w32.GWL_STYLE))
	if currentStyle != 0 {
		if b {
			currentStyle |= style
		} else {
			currentStyle &^= style
		}
		w32.SetWindowLongPtr(w.hwnd, w32.GWL_STYLE, uintptr(currentStyle))
	}
}
func (w *windowsWebviewWindow) setExStyle(b bool, style int) {
	currentStyle := int(w32.GetWindowLongPtr(w.hwnd, w32.GWL_EXSTYLE))
	if currentStyle != 0 {
		if b {
			currentStyle |= style
		} else {
			currentStyle &^= style
		}
		w32.SetWindowLongPtr(w.hwnd, w32.GWL_EXSTYLE, uintptr(currentStyle))
	}
}

func (w *windowsWebviewWindow) setBackdropType(backdropType BackdropType) {
	if !w32.IsWindowsVersionAtLeast(10, 0, 22621) {
		var accent = w32.ACCENT_POLICY{
			AccentState: w32.ACCENT_ENABLE_BLURBEHIND,
		}
		var data w32.WINDOWCOMPOSITIONATTRIBDATA
		data.Attrib = w32.WCA_ACCENT_POLICY
		data.PvData = w32.PVOID(&accent)
		data.CbData = w32.SIZE_T(unsafe.Sizeof(accent))

		w32.SetWindowCompositionAttribute(w.hwnd, &data)
	} else {
		backdropValue := backdropType
		// We default to None, but in win32 None = 1 and Auto = 0
		// So we check if the value given was Auto and set it to 0
		if backdropType == Auto {
			backdropValue = None
		}
		w32.DwmSetWindowAttribute(w.hwnd, w32.DwmwaSystemBackdropType, w32.LPCVOID(&backdropValue), uint32(unsafe.Sizeof(backdropValue)))
	}
}

func (w *windowsWebviewWindow) setIcon(icon w32.HICON) {
	w32.SendMessage(w.hwnd, w32.BM_SETIMAGE, w32.IMAGE_ICON, uintptr(icon))
}

func (w *windowsWebviewWindow) disableIcon() {

	// TODO: If frameless, return

	exStyle := w32.GetWindowLong(w.hwnd, w32.GWL_EXSTYLE)
	w32.SetWindowLong(w.hwnd, w32.GWL_EXSTYLE, uint32(exStyle|w32.WS_EX_DLGMODALFRAME))
	w32.SetWindowPos(w.hwnd, 0, 0, 0, 0, 0,
		uint(
			w32.SWP_FRAMECHANGED|
				w32.SWP_NOMOVE|
				w32.SWP_NOSIZE|
				w32.SWP_NOZORDER),
	)
}

func NewIconFromResource(instance w32.HINSTANCE, resId uint16) (w32.HICON, error) {
	var err error
	var result w32.HICON
	if result = w32.LoadIconWithResourceID(instance, resId); result == 0 {
		err = errors.New(fmt.Sprintf("Cannot load icon from resource with id %v", resId))
	}
	return result, err
}
