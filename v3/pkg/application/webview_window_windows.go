//go:build windows

package application

import (
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/w32"
	"unsafe"
)

var showDevTools = func(window unsafe.Pointer) {}

type windowsWebviewWindow struct {
	windowImpl unsafe.Pointer
	parent     *WebviewWindow
	hwnd       w32.HWND

	// Fullscreen flags
	isCurrentlyFullscreen   bool
	previousWindowStyle     uint32
	previousWindowExStyle   uint32
	previousWindowPlacement w32.WINDOWPLACEMENT
}

func (w *windowsWebviewWindow) nativeWindowHandle() uintptr {
	return w.hwnd
}

func (w *windowsWebviewWindow) setTitle(title string) {
	w32.SetWindowText(w.hwnd, title)
}

func (w *windowsWebviewWindow) setSize(width, height int) {
	x, y := w.position()
	// TODO: Take scaling/DPI into consideration
	w32.MoveWindow(w.hwnd, x, y, width, height, true)
}

func (w *windowsWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	position := w32.HWND_NOTOPMOST
	if alwaysOnTop {
		position = w32.HWND_TOPMOST
	}
	w32.SetWindowPos(w.hwnd, position, 0, 0, 0, 0, uint(w32.SWP_NOMOVE|w32.SWP_NOSIZE))
}

func (w *windowsWebviewWindow) setURL(url string) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setResizable(resizable bool) {
	w.setStyle(resizable, w32.WS_THICKFRAME)
}

func (w *windowsWebviewWindow) setMinSize(width, height int) {
	w.parent.options.MinWidth = width
	w.parent.options.MinHeight = height
}

func (w *windowsWebviewWindow) setMaxSize(width, height int) {
	w.parent.options.MaxWidth = width
	w.parent.options.MaxHeight = height
}

func (w *windowsWebviewWindow) execJS(js string) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setBackgroundColour(color RGBA) {
	w32.SetBackgroundColour(w.hwnd, color.Red, color.Green, color.Blue)
}

func (w *windowsWebviewWindow) run() {
	globalApplication.dispatchOnMainThread(w._run)
}

func (w *windowsWebviewWindow) framelessWithDecorations() bool {
	return w.parent.options.Frameless && !w.parent.options.Windows.DisableFramelessWindowDecorations
}

func (w *windowsWebviewWindow) _run() {

	options := w.parent.options

	var exStyle uint
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
		w32.MustStringToUTF16Ptr(options.Title),
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

	// Register the window with the application
	windowsApp := globalApplication.impl.(*windowsApp)
	windowsApp.registerWindow(w)

	w.setResizable(!options.DisableResize)

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

	// Process the theme
	switch options.Windows.Theme {
	case SystemDefault:
		w.updateTheme(w32.IsCurrentlyDarkMode())
		w.parent.onApplicationEvent(events.Windows.SystemThemeChanged, func() {
			w.updateTheme(w32.IsCurrentlyDarkMode())
		})
	case Light:
		w.updateTheme(false)
	case Dark:
		w.updateTheme(true)
	}

	// Process StartState
	switch options.StartState {
	case WindowStateMaximised:
		if w.parent.Resizable() {
			w.maximise()
		}
	case WindowStateMinimised:
		w.minimise()
	case WindowStateFullscreen:
		w.fullscreen()
	}

	w.setForeground()

	if !options.Hidden {
		w.show()
		w.update()
	}
}

func (w *windowsWebviewWindow) center() {
	w32.CenterWindow(w.hwnd)
}

func (w *windowsWebviewWindow) disableSizeConstraints() {
	w.setMaxSize(0, 0)
	w.setMinSize(0, 0)
}

func (w *windowsWebviewWindow) enableSizeConstraints() {
	options := w.parent.options
	if options.MinWidth > 0 || options.MinHeight > 0 {
		w.setMinSize(options.MinWidth, options.MinHeight)
	}
	if options.MaxWidth > 0 || options.MaxHeight > 0 {
		w.setMaxSize(options.MaxWidth, options.MaxHeight)
	}
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
	info := w32.GetMonitorInfoForWindow(w.hwnd)
	workRect := info.RcWork
	w32.SetWindowPos(w.hwnd, w32.HWND_TOP, int(workRect.Left)+x, int(workRect.Top)+y, 0, 0, w32.SWP_NOSIZE)
}

// on is used to indicate that a particular event should be listened for
func (w *windowsWebviewWindow) on(eventID uint) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) minimise() {
	w32.ShowWindow(w.hwnd, w32.SW_MINIMIZE)
}

func (w *windowsWebviewWindow) unminimise() {
	w.restore()
}

func (w *windowsWebviewWindow) maximise() {
	w32.ShowWindow(w.hwnd, w32.SW_MAXIMIZE)
}

func (w *windowsWebviewWindow) unmaximise() {
	w.restore()
}

func (w *windowsWebviewWindow) restore() {
	w32.ShowWindow(w.hwnd, w32.SW_RESTORE)
}

func (w *windowsWebviewWindow) fullscreen() {
	if w.isFullscreen() {
		return
	}
	if w.framelessWithDecorations() {
		w32.ExtendFrameIntoClientArea(w.hwnd, false)
	}
	w.disableSizeConstraints()
	w.previousWindowStyle = uint32(w32.GetWindowLongPtr(w.hwnd, w32.GWL_STYLE))
	w.previousWindowExStyle = uint32(w32.GetWindowLong(w.hwnd, w32.GWL_EXSTYLE))
	monitor := w32.MonitorFromWindow(w.hwnd, w32.MONITOR_DEFAULTTOPRIMARY)
	var monitorInfo w32.MONITORINFO
	monitorInfo.CbSize = uint32(unsafe.Sizeof(monitorInfo))
	if !w32.GetMonitorInfo(monitor, &monitorInfo) {
		return
	}
	if !w32.GetWindowPlacement(w.hwnd, &w.previousWindowPlacement) {
		return
	}
	// According to https://devblogs.microsoft.com/oldnewthing/20050505-04/?p=35703 one should use w32.WS_POPUP | w32.WS_VISIBLE
	w32.SetWindowLong(w.hwnd, w32.GWL_STYLE, w.previousWindowStyle & ^uint32(w32.WS_OVERLAPPEDWINDOW) | (w32.WS_POPUP|w32.WS_VISIBLE))
	w32.SetWindowLong(w.hwnd, w32.GWL_EXSTYLE, w.previousWindowExStyle & ^uint32(w32.WS_EX_DLGMODALFRAME))
	w32.SetWindowPos(w.hwnd, w32.HWND_TOP,
		int(monitorInfo.RcMonitor.Left),
		int(monitorInfo.RcMonitor.Top),
		int(monitorInfo.RcMonitor.Right-monitorInfo.RcMonitor.Left),
		int(monitorInfo.RcMonitor.Bottom-monitorInfo.RcMonitor.Top),
		w32.SWP_NOOWNERZORDER|w32.SWP_FRAMECHANGED)
	w.isCurrentlyFullscreen = true
}

func (w *windowsWebviewWindow) unfullscreen() {
	if !w.isFullscreen() {
		return
	}
	if w.framelessWithDecorations() {
		w32.ExtendFrameIntoClientArea(w.hwnd, true)
	}
	w32.SetWindowLong(w.hwnd, w32.GWL_STYLE, w.previousWindowStyle)
	w32.SetWindowLong(w.hwnd, w32.GWL_EXSTYLE, w.previousWindowExStyle)
	w32.SetWindowPlacement(w.hwnd, &w.previousWindowPlacement)
	w32.SetWindowPos(w.hwnd, 0, 0, 0, 0, 0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOZORDER|w32.SWP_NOOWNERZORDER|w32.SWP_FRAMECHANGED)
	w.enableSizeConstraints()
	w.isCurrentlyFullscreen = false
}

func (w *windowsWebviewWindow) isMinimised() bool {
	style := uint32(w32.GetWindowLong(w.hwnd, w32.GWL_STYLE))
	return style&w32.WS_MINIMIZE != 0
}

func (w *windowsWebviewWindow) isMaximised() bool {
	style := uint32(w32.GetWindowLong(w.hwnd, w32.GWL_STYLE))
	return style&w32.WS_MAXIMIZE != 0
}

func (w *windowsWebviewWindow) isFullscreen() bool {
	// TODO: Actually calculate this based on size of window against screen size
	return w.isCurrentlyFullscreen
}

func (w *windowsWebviewWindow) isNormal() bool {
	return !w.isMinimised() && !w.isMaximised() && !w.isFullscreen()
}

func (w *windowsWebviewWindow) isVisible() bool {
	style := uint32(w32.GetWindowLong(w.hwnd, w32.GWL_STYLE))
	return style&w32.WS_VISIBLE != 0
}

func (w *windowsWebviewWindow) setFullscreenButtonEnabled(_ bool) {
	// Unused in Windows
}

func (w *windowsWebviewWindow) show() {
	w32.ShowWindow(w.hwnd, w32.SW_SHOW)
}

func (w *windowsWebviewWindow) hide() {
	w32.ShowWindow(w.hwnd, w32.SW_HIDE)
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
		// We default to None, but in w32 None = 1 and Auto = 0
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

func (w *windowsWebviewWindow) updateTheme(isDarkMode bool) {

	if w32.IsCurrentlyHighContrastMode() {
		return
	}

	if !w32.SupportsThemes() {
		return
	}

	w32.SetTheme(w.hwnd, isDarkMode)

	// Custom theme processing
	customTheme := w.parent.options.Windows.CustomTheme
	// Custom theme
	if w32.SupportsCustomThemes() && customTheme != nil {
		if w.isActive() {
			if isDarkMode {
				w32.SetTitleBarColour(w.hwnd, customTheme.DarkModeTitleBar)
				w32.SetTitleTextColour(w.hwnd, customTheme.DarkModeTitleText)
				w32.SetBorderColour(w.hwnd, customTheme.DarkModeBorder)
			} else {
				w32.SetTitleBarColour(w.hwnd, customTheme.LightModeTitleBar)
				w32.SetTitleTextColour(w.hwnd, customTheme.LightModeTitleText)
				w32.SetBorderColour(w.hwnd, customTheme.LightModeBorder)
			}
		} else {
			if isDarkMode {
				w32.SetTitleBarColour(w.hwnd, customTheme.DarkModeTitleBarInactive)
				w32.SetTitleTextColour(w.hwnd, customTheme.DarkModeTitleTextInactive)
				w32.SetBorderColour(w.hwnd, customTheme.DarkModeBorderInactive)
			} else {
				w32.SetTitleBarColour(w.hwnd, customTheme.LightModeTitleBarInactive)
				w32.SetTitleTextColour(w.hwnd, customTheme.LightModeTitleTextInactive)
				w32.SetBorderColour(w.hwnd, customTheme.LightModeBorderInactive)
			}
		}
	}
}

func (w *windowsWebviewWindow) isActive() bool {
	return w32.GetForegroundWindow() == w.hwnd
}

func (w *windowsWebviewWindow) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_SIZE:
		return 0
	case w32.WM_CLOSE:
		w32.PostMessage(w.hwnd, w32.WM_QUIT, 0, 0)
		// Unregister the window with the application
		windowsApp := globalApplication.impl.(*windowsApp)
		windowsApp.unregisterWindow(w)
		return 0
	case w32.WM_GETMINMAXINFO:
		mmi := (*w32.MINMAXINFO)(unsafe.Pointer(lparam))
		hasConstraints := false
		options := w.parent.options
		if options.MinWidth > 0 || options.MinHeight > 0 {
			hasConstraints = true

			width, height := w.scaleWithWindowDPI(options.MinWidth, options.MinHeight)
			if width > 0 {
				mmi.PtMinTrackSize.X = int32(width)
			}
			if height > 0 {
				mmi.PtMinTrackSize.Y = int32(height)
			}
		}
		if options.MaxWidth > 0 || options.MaxHeight > 0 {
			hasConstraints = true

			width, height := w.scaleWithWindowDPI(options.MaxWidth, options.MaxHeight)
			if width > 0 {
				mmi.PtMaxTrackSize.X = int32(width)
			}
			if height > 0 {
				mmi.PtMaxTrackSize.Y = int32(height)
			}
		}
		if hasConstraints {
			return 0
		}
	}
	return w32.DefWindowProc(w.hwnd, msg, wparam, lparam)
}

func (w *windowsWebviewWindow) DPI() (w32.UINT, w32.UINT) {
	if w32.HasGetDpiForWindowFunc() {
		// GetDpiForWindow is supported beginning with Windows 10, 1607 and is the most accureate
		// one, especially it is consistent with the WM_DPICHANGED event.
		dpi := w32.GetDpiForWindow(w.hwnd)
		return dpi, dpi
	}

	if w32.HasGetDPIForMonitorFunc() {
		// GetDpiForWindow is supported beginning with Windows 8.1
		monitor := w32.MonitorFromWindow(w.hwnd, w32.MONITOR_DEFAULTTONEAREST)
		if monitor == 0 {
			return 0, 0
		}
		var dpiX, dpiY w32.UINT
		w32.GetDPIForMonitor(monitor, w32.MDT_EFFECTIVE_DPI, &dpiX, &dpiY)
		return dpiX, dpiY
	}

	// If none of the above is supported fallback to the System DPI.
	screen := w32.GetDC(0)
	x := w32.GetDeviceCaps(screen, w32.LOGPIXELSX)
	y := w32.GetDeviceCaps(screen, w32.LOGPIXELSY)
	w32.ReleaseDC(0, screen)
	return w32.UINT(x), w32.UINT(y)
}

func (w *windowsWebviewWindow) scaleWithWindowDPI(width, height int) (int, int) {
	dpix, dpiy := w.DPI()
	scaledWidth := ScaleWithDPI(width, dpix)
	scaledHeight := ScaleWithDPI(height, dpiy)

	return scaledWidth, scaledHeight
}

func ScaleWithDPI(pixels int, dpi uint) int {
	return (pixels * int(dpi)) / 96
}

func NewIconFromResource(instance w32.HINSTANCE, resId uint16) (w32.HICON, error) {
	var err error
	var result w32.HICON
	if result = w32.LoadIconWithResourceID(instance, resId); result == 0 {
		err = errors.New(fmt.Sprintf("Cannot load icon from resource with id %v", resId))
	}
	return result, err
}
