//go:build windows

package application

import (
	"errors"
	"fmt"
	"github.com/samber/lo"
	"strconv"
	"unicode/utf16"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/w32"
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
	rect := w32.GetWindowRect(w.hwnd)
	width, height = w.scaleWithWindowDPI(width, height)
	w32.MoveWindow(w.hwnd, int(rect.Left), int(rect.Top), width, height, true)
}

func (w *windowsWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	w32.SetWindowPos(w.hwnd,
		lo.Ternary(alwaysOnTop, w32.HWND_TOPMOST, w32.HWND_NOTOPMOST),
		0,
		0,
		0,
		0,
		uint(w32.SWP_NOMOVE|w32.SWP_NOSIZE))
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

func (w *windowsWebviewWindow) framelessWithDecorations() bool {
	return w.parent.options.Frameless && !w.parent.options.Windows.DisableFramelessWindowDecorations
}

func (w *windowsWebviewWindow) run() {

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

	if options.Frameless {
		// Inform the application of the frame change this is needed to trigger the WM_NCCALCSIZE event.
		// => https://learn.microsoft.com/en-us/windows/win32/dwm/customframe#removing-the-standard-frame
		// This is normally done in WM_CREATE but we can't handle that there because that is emitted during CreateWindowEx
		// and at that time we can't yet register the window for calling our WndProc method.
		// This must be called after setResizable above!
		rcClient := w32.GetWindowRect(w.hwnd)
		w32.SetWindowPos(w.hwnd,
			0,
			int(rcClient.Left),
			int(rcClient.Top),
			int(rcClient.Right-rcClient.Left),
			int(rcClient.Bottom-rcClient.Top),
			w32.SWP_FRAMECHANGED)
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

	switch options.BackgroundType {
	case BackgroundTypeSolid:
		w.setBackgroundColour(options.BackgroundColour)
	case BackgroundTypeTransparent:
	case BackgroundTypeTranslucent:
		w.setBackdropType(options.Windows.BackdropType)
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

	// Process window mask
	if options.Windows.WindowMask != nil {
		w.setWindowMask(options.Windows.WindowMask)
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
	rect := w32.GetWindowRect(w.hwnd)
	width := int(rect.Right - rect.Left)
	height := int(rect.Bottom - rect.Top)
	width, height = w.scaleToDefaultDPI(width, height)
	return width, height
}

func (w *windowsWebviewWindow) setForeground() {
	w32.SetForegroundWindow(w.hwnd)
}

func (w *windowsWebviewWindow) update() {
	w32.UpdateWindow(w.hwnd)
}

func (w *windowsWebviewWindow) width() int {
	width, _ := w.size()
	return width
}

func (w *windowsWebviewWindow) height() int {
	_, height := w.size()
	return height
}

func (w *windowsWebviewWindow) position() (int, int) {
	rect := w32.GetWindowRect(w.hwnd)
	left, right := w.scaleToDefaultDPI(int(rect.Left), int(rect.Right))
	return left, right
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
	x, y = w.scaleWithWindowDPI(x, y)
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
	w.isCurrentlyFullscreen = true
	w32.SetWindowPos(w.hwnd, w32.HWND_TOP,
		int(monitorInfo.RcMonitor.Left),
		int(monitorInfo.RcMonitor.Top),
		int(monitorInfo.RcMonitor.Right-monitorInfo.RcMonitor.Left),
		int(monitorInfo.RcMonitor.Bottom-monitorInfo.RcMonitor.Top),
		w32.SWP_NOOWNERZORDER|w32.SWP_FRAMECHANGED)
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
	w.isCurrentlyFullscreen = false
	w32.SetWindowPos(w.hwnd, 0, 0, 0, 0, 0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOZORDER|w32.SWP_NOOWNERZORDER|w32.SWP_FRAMECHANGED)
	w.enableSizeConstraints()
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
	// => stffabi: This flag is essential since it indicates that we are in fullscreen mode even before the native properties
	//             reflect this, e.g. when needing to know if we are in fullscreen during a wndproc message.
	//             That's also why this flag is set before SetWindowPos in v2 in fullscreen/unfullscreen.
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

// Get the screen for the current window
func (w *windowsWebviewWindow) getScreen() (*Screen, error) {
	hMonitor := w32.MonitorFromWindow(w.hwnd, w32.MONITOR_DEFAULTTONEAREST)

	var mi w32.MONITORINFOEX
	mi.CbSize = uint32(unsafe.Sizeof(mi))
	w32.GetMonitorInfoEx(hMonitor, &mi)
	var thisScreen Screen
	thisScreen.X = int(mi.RcMonitor.Left)
	thisScreen.Y = int(mi.RcMonitor.Top)
	thisScreen.Size = Size{
		Width:  int(mi.RcMonitor.Right - mi.RcMonitor.Left),
		Height: int(mi.RcMonitor.Bottom - mi.RcMonitor.Top),
	}
	thisScreen.Bounds = Rect{
		X:      int(mi.RcMonitor.Left),
		Y:      int(mi.RcMonitor.Top),
		Width:  int(mi.RcMonitor.Right - mi.RcMonitor.Left),
		Height: int(mi.RcMonitor.Bottom - mi.RcMonitor.Top),
	}
	thisScreen.WorkArea = Rect{
		X:      int(mi.RcWork.Left),
		Y:      int(mi.RcWork.Top),
		Width:  int(mi.RcWork.Right - mi.RcWork.Left),
		Height: int(mi.RcWork.Bottom - mi.RcWork.Top),
	}
	thisScreen.ID = strconv.Itoa(int(hMonitor))
	thisScreen.Name = string(utf16.Decode(mi.SzDevice[:]))
	var xdpi, ydpi w32.UINT
	w32.GetDPIForMonitor(hMonitor, w32.MDT_EFFECTIVE_DPI, &xdpi, &ydpi)
	thisScreen.Scale = float32(xdpi) / 96.0
	thisScreen.IsPrimary = mi.DwFlags&w32.MONITORINFOF_PRIMARY != 0

	// TODO: Get screen rotation
	// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/ns-wingdi-devmodea

	//// get display settings for monitor
	//var dm w32.DEVMODE
	//dm.DmSize = uint16(unsafe.Sizeof(dm))
	//dm.DmDriverExtra = 0
	//w32.EnumDisplaySettingsEx(&mi.SzDevice[0], w32.ENUM_CURRENT_SETTINGS, &dm, 0)
	//
	//// check display settings for rotation
	//rotationAngle := dm.DmDi
	//if rotationAngle == DMDO_0 {
	//	printf("Monitor is not rotated\n")
	//} else if rotationAngle == DMDO_90 {
	//	printf("Monitor is rotated 90 degrees\n")
	//} else if rotationAngle == DMDO_180 {
	//	printf("Monitor is rotated 180 degrees\n")
	//} else if rotationAngle == DMDO_270 {
	//	printf("Monitor is rotated 270 degrees\n")
	//} else {
	//	printf("Monitor is rotated at an unknown angle\n")
	//}

	return &thisScreen, nil
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
		currentStyle = lo.Ternary(b, currentStyle|style, currentStyle&^style)
		w32.SetWindowLongPtr(w.hwnd, w32.GWL_STYLE, uintptr(currentStyle))
	}
}
func (w *windowsWebviewWindow) setExStyle(b bool, style int) {
	currentStyle := int(w32.GetWindowLongPtr(w.hwnd, w32.GWL_EXSTYLE))
	if currentStyle != 0 {
		currentStyle = lo.Ternary(b, currentStyle|style, currentStyle&^style)
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
		w32.DwmSetWindowAttribute(w.hwnd, w32.DwmwaSystemBackdropType, w32.PVOID(&backdropType), unsafe.Sizeof(backdropType))
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
	case w32.WM_DPICHANGED:
		newWindowSize := (*w32.RECT)(unsafe.Pointer(lparam))
		w32.SetWindowPos(w.hwnd,
			uintptr(0),
			int(newWindowSize.Left),
			int(newWindowSize.Top),
			int(newWindowSize.Right-newWindowSize.Left),
			int(newWindowSize.Bottom-newWindowSize.Top),
			w32.SWP_NOZORDER|w32.SWP_NOACTIVATE)

	}

	if w.parent.options.Windows.WindowMask != nil {
		switch msg {
		case w32.WM_NCHITTEST:
			return w32.HTCAPTION
		}
	}

	if options := w.parent.options; options.Frameless {
		switch msg {
		case w32.WM_ACTIVATE:
			// If we want to have a frameless window but with the default frame decorations, extend the DWM client area.
			// This Option is not affected by returning 0 in WM_NCCALCSIZE.
			// As a result we have hidden the titlebar but still have the default window frame styling.
			// See: https://docs.microsoft.com/en-us/windows/win32/api/dwmapi/nf-dwmapi-dwmextendframeintoclientarea#remarks
			if w.framelessWithDecorations() {
				w32.ExtendFrameIntoClientArea(w.hwnd, true)
			}
		case w32.WM_NCHITTEST:
			// Get the cursor position
			x := int32(w32.LOWORD(uint32(lparam)))
			y := int32(w32.HIWORD(uint32(lparam)))
			ptCursor := w32.POINT{X: x, Y: y}

			// Get the window rectangle
			rcWindow := w32.GetWindowRect(w.hwnd)

			// Determine if the cursor is in a resize area
			bOnResizeBorder := false
			resizeBorderWidth := int32(5) // change this to adjust the resize border width
			if ptCursor.X >= rcWindow.Right-resizeBorderWidth {
				bOnResizeBorder = true // right edge
			}
			if ptCursor.Y >= rcWindow.Bottom-resizeBorderWidth {
				bOnResizeBorder = true // bottom edge
			}
			if ptCursor.X <= rcWindow.Left+resizeBorderWidth {
				bOnResizeBorder = true // left edge
			}
			if ptCursor.Y <= rcWindow.Top+resizeBorderWidth {
				bOnResizeBorder = true // top edge
			}

			// Return the appropriate value
			if bOnResizeBorder {
				if ptCursor.X >= rcWindow.Right-resizeBorderWidth && ptCursor.Y >= rcWindow.Bottom-resizeBorderWidth {
					return w32.HTBOTTOMRIGHT
				} else if ptCursor.X <= rcWindow.Left+resizeBorderWidth && ptCursor.Y >= rcWindow.Bottom-resizeBorderWidth {
					return w32.HTBOTTOMLEFT
				} else if ptCursor.X >= rcWindow.Right-resizeBorderWidth && ptCursor.Y <= rcWindow.Top+resizeBorderWidth {
					return w32.HTTOPRIGHT
				} else if ptCursor.X <= rcWindow.Left+resizeBorderWidth && ptCursor.Y <= rcWindow.Top+resizeBorderWidth {
					return w32.HTTOPLEFT
				} else if ptCursor.X >= rcWindow.Right-resizeBorderWidth {
					return w32.HTRIGHT
				} else if ptCursor.Y >= rcWindow.Bottom-resizeBorderWidth {
					return w32.HTBOTTOM
				} else if ptCursor.X <= rcWindow.Left+resizeBorderWidth {
					return w32.HTLEFT
				} else if ptCursor.Y <= rcWindow.Top+resizeBorderWidth {
					return w32.HTTOP
				}
			}
			return w32.HTCLIENT

		case w32.WM_NCCALCSIZE:
			// Disable the standard frame by allowing the client area to take the full
			// window size.
			// See: https://docs.microsoft.com/en-us/windows/win32/winmsg/wm-nccalcsize#remarks
			// This hides the titlebar and also disables the resizing from user interaction because the standard frame is not
			// shown. We still need the WS_THICKFRAME style to enable resizing from the frontend.
			if wparam != 0 {
				rgrc := (*w32.RECT)(unsafe.Pointer(lparam))
				if w.isCurrentlyFullscreen {
					// In Full-Screen mode we don't need to adjust anything
					// It essential we have the flag here, that is set before SetWindowPos in fullscreen/unfullscreen
					// because the native size might not yet reflect we are in fullscreen during this event!
					// TODO: w.chromium.SetPadding(edge.Rect{})
				} else if w.isMaximised() {
					// If the window is maximized we must adjust the client area to the work area of the monitor. Otherwise
					// some content goes beyond the visible part of the monitor.
					// Make sure to use the provided RECT to get the monitor, because during maximizig there might be
					// a wrong monitor returned in multi screen mode when using MonitorFromWindow.
					// See: https://github.com/MicrosoftEdge/WebView2Feedback/issues/2549
					monitor := w32.MonitorFromRect(rgrc, w32.MONITOR_DEFAULTTONULL)

					var monitorInfo w32.MONITORINFO
					monitorInfo.CbSize = uint32(unsafe.Sizeof(monitorInfo))
					if monitor != 0 && w32.GetMonitorInfo(monitor, &monitorInfo) {
						*rgrc = monitorInfo.RcWork

						maxWidth := options.MaxWidth
						maxHeight := options.MaxHeight
						if maxWidth > 0 || maxHeight > 0 {
							var dpiX, dpiY uint
							w32.GetDPIForMonitor(monitor, w32.MDT_EFFECTIVE_DPI, &dpiX, &dpiY)

							maxWidth := int32(ScaleWithDPI(maxWidth, dpiX))
							if maxWidth > 0 && rgrc.Right-rgrc.Left > maxWidth {
								rgrc.Right = rgrc.Left + maxWidth
							}

							maxHeight := int32(ScaleWithDPI(maxHeight, dpiY))
							if maxHeight > 0 && rgrc.Bottom-rgrc.Top > maxHeight {
								rgrc.Bottom = rgrc.Top + maxHeight
							}
						}
					}
					// TODO: w.chromium.SetPadding(edge.Rect{})
				} else {
					// This is needed to workaround the resize flickering in frameless mode with WindowDecorations
					// See: https://stackoverflow.com/a/6558508
					// The workaround originally suggests to decrese the bottom 1px, but that seems to bring up a thin
					// white line on some Windows-Versions, due to DrawBackground using also this reduces ClientSize.
					// Increasing the bottom also worksaround the flickering but we would loose 1px of the WebView content
					// therefore let's pad the content with 1px at the bottom.
					rgrc.Bottom += 1
					// TODO: w.chromium.SetPadding(edge.Rect{Bottom: 1})
				}
				return 0
			}
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

func (w *windowsWebviewWindow) scaleToDefaultDPI(width, height int) (int, int) {
	dpix, dpiy := w.DPI()
	scaledWidth := ScaleToDefaultDPI(width, dpix)
	scaledHeight := ScaleToDefaultDPI(height, dpiy)

	return scaledWidth, scaledHeight
}

func (w *windowsWebviewWindow) setWindowMask(imageData []byte) {

	// Set the window to a WS_EX_LAYERED window
	newStyle := w32.GetWindowLong(w.hwnd, w32.GWL_EXSTYLE) | w32.WS_EX_LAYERED

	if w.isAlwaysOnTop() {
		newStyle |= w32.WS_EX_TOPMOST
	}
	// Save the current window style
	w.previousWindowExStyle = uint32(w32.GetWindowLong(w.hwnd, w32.GWL_EXSTYLE))

	w32.SetWindowLong(w.hwnd, w32.GWL_EXSTYLE, uint32(newStyle))

	data, err := pngToImage(imageData)
	if err != nil {
		panic(err)
	}

	bitmap, err := w32.CreateHBITMAPFromImage(data)
	hdc := w32.CreateCompatibleDC(0)
	defer w32.DeleteDC(hdc)

	oldBitmap := w32.SelectObject(hdc, bitmap)
	defer w32.SelectObject(hdc, oldBitmap)

	screenDC := w32.GetDC(0)
	defer w32.ReleaseDC(0, screenDC)

	size := w32.SIZE{CX: int32(data.Bounds().Dx()), CY: int32(data.Bounds().Dy())}
	ptSrc := w32.POINT{X: 0, Y: 0}
	ptDst := w32.POINT{X: int32(w.width()), Y: int32(w.height())}
	blend := w32.BLENDFUNCTION{
		BlendOp:             w32.AC_SRC_OVER,
		BlendFlags:          0,
		SourceConstantAlpha: 255,
		AlphaFormat:         w32.AC_SRC_ALPHA,
	}
	w32.UpdateLayeredWindow(w.hwnd, screenDC, &ptDst, &size, hdc, &ptSrc, 0, &blend, w32.ULW_ALPHA)
}

func (w *windowsWebviewWindow) isAlwaysOnTop() bool {
	return w32.GetWindowLong(w.hwnd, w32.GWL_EXSTYLE)&w32.WS_EX_TOPMOST != 0
}

func ScaleWithDPI(pixels int, dpi uint) int {
	return (pixels * int(dpi)) / 96
}

func ScaleToDefaultDPI(pixels int, dpi uint) int {
	return (pixels * 96) / int(dpi)
}

func NewIconFromResource(instance w32.HINSTANCE, resId uint16) (w32.HICON, error) {
	var err error
	var result w32.HICON
	if result = w32.LoadIconWithResourceID(instance, resId); result == 0 {
		err = errors.New(fmt.Sprintf("Cannot load icon from resource with id %v", resId))
	}
	return result, err
}
