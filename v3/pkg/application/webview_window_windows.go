//go:build windows

package application

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"

	"github.com/bep/debounce"
	"github.com/wailsapp/go-webview2/webviewloader"
	"github.com/wailsapp/wails/v3/internal/assetserver"
	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/internal/capabilities"
	"github.com/wailsapp/wails/v3/internal/runtime"

	"github.com/samber/lo"

	"github.com/wailsapp/go-webview2/pkg/edge"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

const (
	windowDidMoveDebounceMS = 200
)

var edgeMap = map[string]uintptr{
	"n-resize":  w32.HTTOP,
	"ne-resize": w32.HTTOPRIGHT,
	"e-resize":  w32.HTRIGHT,
	"se-resize": w32.HTBOTTOMRIGHT,
	"s-resize":  w32.HTBOTTOM,
	"sw-resize": w32.HTBOTTOMLEFT,
	"w-resize":  w32.HTLEFT,
	"nw-resize": w32.HTTOPLEFT,
}

type windowsWebviewWindow struct {
	windowImpl               unsafe.Pointer
	parent                   *WebviewWindow
	hwnd                     w32.HWND
	menu                     *Win32Menu
	currentlyOpenContextMenu *Win32Menu

	// Fullscreen flags
	isCurrentlyFullscreen   bool
	previousWindowStyle     uint32
	previousWindowExStyle   uint32
	previousWindowPlacement w32.WINDOWPLACEMENT

	// Webview
	chromium        *edge.Chromium
	hasStarted      bool
	resizeDebouncer func(func())

	// resizeBorder* is the width/height of the resize border in pixels.
	resizeBorderWidth  int32
	resizeBorderHeight int32
	focusingChromium   bool
	dropTarget         *w32.DropTarget
	onceDo             sync.Once

	// Window move debouncer
	moveDebouncer func(func())
}

func (w *windowsWebviewWindow) handleKeyEvent(_ string) {
	// Unused on windows
}

// getBorderSizes returns the extended border size for the window
func (w *windowsWebviewWindow) getBorderSizes() *LRTB {
	var result LRTB
	var frame w32.RECT
	w32.DwmGetWindowAttribute(w.hwnd, w32.DWMWA_EXTENDED_FRAME_BOUNDS, unsafe.Pointer(&frame), unsafe.Sizeof(frame))
	rect := w32.GetWindowRect(w.hwnd)
	result.Left = int(frame.Left - rect.Left)
	result.Top = int(frame.Top - rect.Top)
	result.Right = int(rect.Right - frame.Right)
	result.Bottom = int(rect.Bottom - frame.Bottom)
	return &result
}

func (w *windowsWebviewWindow) setPosition(x int, y int) {
	// Set the window's absolute position
	borderSize := w.getBorderSizes()
	w32.SetWindowPos(w.hwnd, 0, x-borderSize.Left, y-borderSize.Top, 0, 0, w32.SWP_NOSIZE|w32.SWP_NOZORDER)
}

func (w *windowsWebviewWindow) position() (int, int) {
	rect := w32.GetWindowRect(w.hwnd)
	borderSizes := w.getBorderSizes()
	x := int(rect.Left) + borderSizes.Left
	y := int(rect.Top) + borderSizes.Top
	left, right := w.scaleToDefaultDPI(x, y)
	return left, right
}

func (w *windowsWebviewWindow) setEnabled(enabled bool) {
	w32.EnableWindow(w.hwnd, enabled)
}

func (w *windowsWebviewWindow) print() error {
	w.execJS("window.print();")
	return nil
}

func (w *windowsWebviewWindow) startResize(border string) error {
	if !w32.ReleaseCapture() {
		return fmt.Errorf("unable to release mouse capture")
	}
	// Use PostMessage because we don't want to block the caller until resizing has been finished.
	w32.PostMessage(w.hwnd, w32.WM_NCLBUTTONDOWN, edgeMap[border], 0)
	return nil
}

func (w *windowsWebviewWindow) startDrag() error {
	if !w32.ReleaseCapture() {
		return fmt.Errorf("unable to release mouse capture")
	}
	// Use PostMessage because we don't want to block the caller until dragging has been finished.
	w32.PostMessage(w.hwnd, w32.WM_NCLBUTTONDOWN, w32.HTCAPTION, 0)
	return nil
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
	w.chromium.Resize()
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
	// Navigate to the given URL in the webview
	w.chromium.Navigate(url)
}

func (w *windowsWebviewWindow) setResizable(resizable bool) {
	w.setStyle(resizable, w32.WS_THICKFRAME)
	w.execJS(fmt.Sprintf("window._wails.setResizable(%v);", resizable))
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
	if w.chromium == nil {
		return
	}
	globalApplication.dispatchOnMainThread(func() {
		w.chromium.Eval(js)
	})
}

func (w *windowsWebviewWindow) setBackgroundColour(color RGBA) {
	w32.SetBackgroundColour(w.hwnd, color.Red, color.Green, color.Blue)
}

func (w *windowsWebviewWindow) framelessWithDecorations() bool {
	return w.parent.options.Frameless && !w.parent.options.Windows.DisableFramelessWindowDecorations
}

func (w *windowsWebviewWindow) run() {

	options := w.parent.options

	w.chromium = edge.NewChromium()
	if globalApplication.options.ErrorHandler != nil {
		w.chromium.SetErrorCallback(globalApplication.options.ErrorHandler)
	}

	exStyle := w32.WS_EX_CONTROLPARENT
	if options.BackgroundType != BackgroundTypeSolid {
		exStyle |= w32.WS_EX_NOREDIRECTIONBITMAP
		if w.parent.options.IgnoreMouseEvents {
			exStyle |= w32.WS_EX_TRANSPARENT | w32.WS_EX_LAYERED
		}
	}
	if options.AlwaysOnTop {
		exStyle |= w32.WS_EX_TOPMOST
	}
	// If we're frameless, we need to add the WS_EX_TOOLWINDOW style to hide the window from the taskbar
	if options.Windows.HiddenOnTaskbar {
		exStyle |= w32.WS_EX_TOOLWINDOW
	} else {
		exStyle |= w32.WS_EX_APPWINDOW
	}

	if options.Windows.ExStyle != 0 {
		exStyle = options.Windows.ExStyle
	}

	// ToDo: X, Y should also be scaled, should it be always relative to the main monitor?
	var startX, _ = lo.Coalesce(options.X, w32.CW_USEDEFAULT)
	var startY, _ = lo.Coalesce(options.Y, w32.CW_USEDEFAULT)

	var appMenu w32.HMENU

	// Process Menu
	if !options.Windows.DisableMenu && !options.Frameless {
		theMenu := globalApplication.ApplicationMenu
		// Create the menu if we have one
		if w.parent.options.Windows.Menu != nil {
			theMenu = w.parent.options.Windows.Menu
		}
		if theMenu != nil {
			w.menu = NewApplicationMenu(w, theMenu)
			w.menu.parentWindow = w
			appMenu = w.menu.menu
		}
	}

	var parent w32.HWND

	var style uint = w32.WS_OVERLAPPEDWINDOW

	w.hwnd = w32.CreateWindowEx(
		uint(exStyle),
		w32.MustStringToUTF16Ptr(globalApplication.options.Windows.WndClass),
		w32.MustStringToUTF16Ptr(options.Title),
		style,
		startX,
		startY,
		w32.CW_USEDEFAULT,
		w32.CW_USEDEFAULT,
		parent,
		appMenu,
		w32.GetModuleHandle(""),
		nil)

	if w.hwnd == 0 {
		panic("Unable to create window")
	}

	w.setSize(options.Width, options.Height)

	w.setupChromium()

	// Initialise the window buttons
	w.setMinimiseButtonState(options.MinimiseButtonState)
	w.setMaximiseButtonState(options.MaximiseButtonState)
	w.setCloseButtonState(options.CloseButtonState)

	// Register the window with the application
	getNativeApplication().registerWindow(w)

	w.setResizable(!options.DisableResize)

	w.setIgnoreMouseEvents(options.IgnoreMouseEvents)

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
		if err != nil {
			// Try loading from the given icon
			if globalApplication.options.Icon != nil {
				icon, _ = w32.CreateLargeHIconFromImage(globalApplication.options.Icon)
			}
		}
		if icon != 0 {
			w.setIcon(icon)
		}
	} else {
		w.disableIcon()
	}

	// Process the theme
	switch options.Windows.Theme {
	case SystemDefault:
		w.updateTheme(w32.IsCurrentlyDarkMode())
		w.parent.onApplicationEvent(events.Windows.SystemThemeChanged, func(*Event) {
			w.updateTheme(w32.IsCurrentlyDarkMode())
		})
	case Light:
		w.updateTheme(false)
	case Dark:
		w.updateTheme(true)
	}

	switch options.BackgroundType {
	case BackgroundTypeSolid:
		var col = options.BackgroundColour
		w.setBackgroundColour(col)
		w.chromium.SetBackgroundColour(col.Red, col.Green, col.Blue, col.Alpha)
	case BackgroundTypeTransparent:
		w.chromium.SetBackgroundColour(0, 0, 0, 0)
	case BackgroundTypeTranslucent:
		w.chromium.SetBackgroundColour(0, 0, 0, 0)
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
	case WindowStateNormal:
	}

	// Process window mask
	if options.Windows.WindowMask != nil {
		w.setWindowMask(options.Windows.WindowMask)
	}

	if options.Windows.ResizeDebounceMS > 0 {
		w.resizeDebouncer = debounce.New(time.Duration(options.Windows.ResizeDebounceMS) * time.Millisecond)
	}

	if options.Centered {
		w.center()
	}

	if options.Frameless {
		// Trigger a resize to ensure the window is sized correctly
		w.chromium.Resize()
	}

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
	// Scaling appears to give invalid results...
	//width, height = w.scaleToDefaultDPI(width, height)
	return width, height
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

func (w *windowsWebviewWindow) relativePosition() (int, int) {
	// Get monitor for window
	monitor := w32.MonitorFromWindow(w.hwnd, w32.MONITOR_DEFAULTTONEAREST)
	var monitorInfo w32.MONITORINFO
	monitorInfo.CbSize = uint32(unsafe.Sizeof(monitorInfo))
	w32.GetMonitorInfo(monitor, &monitorInfo)

	// Get window rect
	rect := w32.GetWindowRect(w.hwnd)

	// Calculate relative position
	x := int(rect.Left) - int(monitorInfo.RcWork.Left)
	y := int(rect.Top) - int(monitorInfo.RcWork.Top)

	borderSize := w.getBorderSizes()
	x += borderSize.Left
	y += borderSize.Top

	return w.scaleToDefaultDPI(x, y)
}

func (w *windowsWebviewWindow) destroy() {
	w.parent.markAsDestroyed()
	if w.dropTarget != nil {
		w.dropTarget.Release()
	}
}

func (w *windowsWebviewWindow) reload() {
	w.execJS("window.location.reload();")
}

func (w *windowsWebviewWindow) forceReload() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) zoomReset() {
	w.setZoom(1.0)
}

func (w *windowsWebviewWindow) zoomIn() {
	// Increase the zoom level by 0.05
	currentZoom := w.getZoom()
	if currentZoom == -1 {
		return
	}
	w.setZoom(currentZoom + 0.05)
}

func (w *windowsWebviewWindow) zoomOut() {
	// Decrease the zoom level by 0.05
	currentZoom := w.getZoom()
	if currentZoom == -1 {
		return
	}
	if currentZoom > 1.05 {
		// Decrease the zoom level by 0.05
		w.setZoom(currentZoom - 0.05)
	} else {
		// Set the zoom level to 1.0
		w.setZoom(1.0)
	}
}

func (w *windowsWebviewWindow) getZoom() float64 {
	controller := w.chromium.GetController()
	factor, err := controller.GetZoomFactor()
	if err != nil {
		return -1
	}
	return factor
}

func (w *windowsWebviewWindow) setZoom(zoom float64) {
	w.chromium.PutZoomFactor(zoom)
}

func (w *windowsWebviewWindow) close() {
	// Unregister the window with the application
	windowsApp := globalApplication.impl.(*windowsApp)
	windowsApp.unregisterWindow(w)
	w32.SendMessage(w.hwnd, w32.WM_CLOSE, 0, 0)
}

func (w *windowsWebviewWindow) zoom() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setHTML(html string) {
	// Render the given HTML in the webview window
	w.execJS(fmt.Sprintf("document.documentElement.innerHTML = %q;", html))
}

func (w *windowsWebviewWindow) setRelativePosition(x int, y int) {
	//x, y = w.scaleWithWindowDPI(x, y)
	info := w32.GetMonitorInfoForWindow(w.hwnd)
	workRect := info.RcWork
	borderSize := w.getBorderSizes()
	x -= borderSize.Left
	y -= borderSize.Top
	w32.SetWindowPos(w.hwnd, w32.HWND_TOP, int(workRect.Left)+x, int(workRect.Top)+y, 0, 0, w32.SWP_NOSIZE)
}

// on is used to indicate that a particular event should be listened for
func (w *windowsWebviewWindow) on(_ uint) {
	// We don't need to worry about this in Windows as we do not need
	// to optimise cgo calls
}

func (w *windowsWebviewWindow) minimise() {
	w32.ShowWindow(w.hwnd, w32.SW_MINIMIZE)
}

func (w *windowsWebviewWindow) unminimise() {
	w.restore()
}

func (w *windowsWebviewWindow) maximise() {
	w32.ShowWindow(w.hwnd, w32.SW_MAXIMIZE)
	w.chromium.Focus()
}

func (w *windowsWebviewWindow) unmaximise() {
	w.restore()
}

func (w *windowsWebviewWindow) restore() {
	w32.ShowWindow(w.hwnd, w32.SW_RESTORE)
	w.chromium.Focus()
}

func (w *windowsWebviewWindow) fullscreen() {
	if w.isFullscreen() {
		return
	}
	if w.framelessWithDecorations() {
		err := w32.ExtendFrameIntoClientArea(w.hwnd, false)
		if err != nil {
			globalApplication.handleFatalError(err)
		}
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
	w.chromium.Focus()
}

func (w *windowsWebviewWindow) unfullscreen() {
	if !w.isFullscreen() {
		return
	}
	if w.framelessWithDecorations() {
		err := w32.ExtendFrameIntoClientArea(w.hwnd, true)
		if err != nil {
			globalApplication.handleFatalError(err)
		}
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

func (w *windowsWebviewWindow) isFocused() bool {
	// Returns true if the window is currently focused
	return w32.GetForegroundWindow() == w.hwnd
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

func (w *windowsWebviewWindow) focus() {
	w32.SetForegroundWindow(w.hwnd)
	w.focusingChromium = true
	w.chromium.Focus()
	w.focusingChromium = false
}

// printStyle takes a windows style and prints it in a human-readable format
// This is for debugging window style issues
func (w *windowsWebviewWindow) printStyle() {
	style := uint32(w32.GetWindowLong(w.hwnd, w32.GWL_STYLE))
	fmt.Printf("Style: ")
	if style&w32.WS_BORDER != 0 {
		fmt.Printf("WS_BORDER ")
	}
	if style&w32.WS_CAPTION != 0 {
		fmt.Printf("WS_CAPTION ")
	}
	if style&w32.WS_CHILD != 0 {
		fmt.Printf("WS_CHILD ")
	}
	if style&w32.WS_CLIPCHILDREN != 0 {
		fmt.Printf("WS_CLIPCHILDREN ")
	}
	if style&w32.WS_CLIPSIBLINGS != 0 {
		fmt.Printf("WS_CLIPSIBLINGS ")
	}
	if style&w32.WS_DISABLED != 0 {
		fmt.Printf("WS_DISABLED ")
	}
	if style&w32.WS_DLGFRAME != 0 {
		fmt.Printf("WS_DLGFRAME ")
	}
	if style&w32.WS_GROUP != 0 {
		fmt.Printf("WS_GROUP ")
	}
	if style&w32.WS_HSCROLL != 0 {
		fmt.Printf("WS_HSCROLL ")
	}
	if style&w32.WS_MAXIMIZE != 0 {
		fmt.Printf("WS_MAXIMIZE ")
	}
	if style&w32.WS_MAXIMIZEBOX != 0 {
		fmt.Printf("WS_MAXIMIZEBOX ")
	}
	if style&w32.WS_MINIMIZE != 0 {
		fmt.Printf("WS_MINIMIZE ")
	}
	if style&w32.WS_MINIMIZEBOX != 0 {
		fmt.Printf("WS_MINIMIZEBOX ")
	}
	if style&w32.WS_OVERLAPPED != 0 {
		fmt.Printf("WS_OVERLAPPED ")
	}
	if style&w32.WS_POPUP != 0 {
		fmt.Printf("WS_POPUP ")
	}
	if style&w32.WS_SYSMENU != 0 {
		fmt.Printf("WS_SYSMENU ")
	}
	if style&w32.WS_TABSTOP != 0 {
		fmt.Printf("WS_TABSTOP ")
	}
	if style&w32.WS_THICKFRAME != 0 {
		fmt.Printf("WS_THICKFRAME ")
	}
	if style&w32.WS_VISIBLE != 0 {
		fmt.Printf("WS_VISIBLE ")
	}
	if style&w32.WS_VSCROLL != 0 {
		fmt.Printf("WS_VSCROLL ")
	}
	fmt.Printf("\n")

	// Do the same for the extended style
	extendedStyle := uint32(w32.GetWindowLong(w.hwnd, w32.GWL_EXSTYLE))
	fmt.Printf("Extended Style: ")
	if extendedStyle&w32.WS_EX_ACCEPTFILES != 0 {
		fmt.Printf("WS_EX_ACCEPTFILES ")
	}
	if extendedStyle&w32.WS_EX_APPWINDOW != 0 {
		fmt.Printf("WS_EX_APPWINDOW ")
	}
	if extendedStyle&w32.WS_EX_CLIENTEDGE != 0 {
		fmt.Printf("WS_EX_CLIENTEDGE ")
	}
	if extendedStyle&w32.WS_EX_COMPOSITED != 0 {
		fmt.Printf("WS_EX_COMPOSITED ")
	}
	if extendedStyle&w32.WS_EX_CONTEXTHELP != 0 {
		fmt.Printf("WS_EX_CONTEXTHELP ")
	}
	if extendedStyle&w32.WS_EX_CONTROLPARENT != 0 {
		fmt.Printf("WS_EX_CONTROLPARENT ")
	}
	if extendedStyle&w32.WS_EX_DLGMODALFRAME != 0 {
		fmt.Printf("WS_EX_DLGMODALFRAME ")
	}
	if extendedStyle&w32.WS_EX_LAYERED != 0 {
		fmt.Printf("WS_EX_LAYERED ")
	}
	if extendedStyle&w32.WS_EX_LAYOUTRTL != 0 {
		fmt.Printf("WS_EX_LAYOUTRTL ")
	}
	if extendedStyle&w32.WS_EX_LEFT != 0 {
		fmt.Printf("WS_EX_LEFT ")
	}
	if extendedStyle&w32.WS_EX_LEFTSCROLLBAR != 0 {
		fmt.Printf("WS_EX_LEFTSCROLLBAR ")
	}
	if extendedStyle&w32.WS_EX_LTRREADING != 0 {
		fmt.Printf("WS_EX_LTRREADING ")
	}
	if extendedStyle&w32.WS_EX_MDICHILD != 0 {
		fmt.Printf("WS_EX_MDICHILD ")
	}
	if extendedStyle&w32.WS_EX_NOACTIVATE != 0 {
		fmt.Printf("WS_EX_NOACTIVATE ")
	}
	if extendedStyle&w32.WS_EX_NOINHERITLAYOUT != 0 {
		fmt.Printf("WS_EX_NOINHERITLAYOUT ")
	}
	if extendedStyle&w32.WS_EX_NOPARENTNOTIFY != 0 {
		fmt.Printf("WS_EX_NOPARENTNOTIFY ")
	}
	if extendedStyle&w32.WS_EX_NOREDIRECTIONBITMAP != 0 {
		fmt.Printf("WS_EX_NOREDIRECTIONBITMAP ")
	}
	if extendedStyle&w32.WS_EX_OVERLAPPEDWINDOW != 0 {
		fmt.Printf("WS_EX_OVERLAPPEDWINDOW ")
	}
	if extendedStyle&w32.WS_EX_PALETTEWINDOW != 0 {
		fmt.Printf("WS_EX_PALETTEWINDOW ")
	}
	if extendedStyle&w32.WS_EX_RIGHT != 0 {
		fmt.Printf("WS_EX_RIGHT ")
	}
	if extendedStyle&w32.WS_EX_RIGHTSCROLLBAR != 0 {
		fmt.Printf("WS_EX_RIGHTSCROLLBAR ")
	}
	if extendedStyle&w32.WS_EX_RTLREADING != 0 {
		fmt.Printf("WS_EX_RTLREADING ")
	}
	if extendedStyle&w32.WS_EX_STATICEDGE != 0 {
		fmt.Printf("WS_EX_STATICEDGE ")
	}
	if extendedStyle&w32.WS_EX_TOOLWINDOW != 0 {
		fmt.Printf("WS_EX_TOOLWINDOW ")
	}
	if extendedStyle&w32.WS_EX_TOPMOST != 0 {
		fmt.Printf("WS_EX_TOPMOST ")
	}
	if extendedStyle&w32.WS_EX_TRANSPARENT != 0 {
		fmt.Printf("WS_EX_TRANSPARENT ")
	}
	if extendedStyle&w32.WS_EX_WINDOWEDGE != 0 {
		fmt.Printf("WS_EX_WINDOWEDGE ")
	}
	fmt.Printf("\n")

}

func (w *windowsWebviewWindow) show() {
	w32.ShowWindow(w.hwnd, w32.SW_SHOW)
}

func (w *windowsWebviewWindow) hide() {
	w32.ShowWindow(w.hwnd, w32.SW_HIDE)
}

func getScreen(hwnd w32.HWND) (*Screen, error) {
	hMonitor := w32.MonitorFromWindow(hwnd, w32.MONITOR_DEFAULTTONEAREST)
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

	return &thisScreen, nil
}

// Get the screen for the current window
func (w *windowsWebviewWindow) getScreen() (*Screen, error) {
	return getScreen(w.hwnd)
}

func (w *windowsWebviewWindow) setFrameless(b bool) {
	// Remove or add the frame
	if b {
		w32.SetWindowLong(w.hwnd, w32.GWL_STYLE, w32.WS_VISIBLE|w32.WS_POPUP)
	} else {
		w32.SetWindowLong(w.hwnd, w32.GWL_STYLE, w32.WS_VISIBLE|w32.WS_OVERLAPPEDWINDOW)
	}
	w32.SetWindowPos(w.hwnd, 0, 0, 0, 0, 0, w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOZORDER|w32.SWP_FRAMECHANGED)
}

func newWindowImpl(parent *WebviewWindow) *windowsWebviewWindow {
	result := &windowsWebviewWindow{
		parent:             parent,
		resizeBorderWidth:  int32(w32.GetSystemMetrics(w32.SM_CXSIZEFRAME)),
		resizeBorderHeight: int32(w32.GetSystemMetrics(w32.SM_CYSIZEFRAME)),
	}

	return result
}

func (w *windowsWebviewWindow) openContextMenu(menu *Menu, _ *ContextMenuData) {
	// Create the menu
	thisMenu := NewPopupMenu(w.hwnd, menu)
	thisMenu.Update()
	w.currentlyOpenContextMenu = thisMenu
	thisMenu.ShowAtCursor()
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
	if !w32.SupportsBackdropTypes() {
		var accent = w32.ACCENT_POLICY{
			AccentState: w32.ACCENT_ENABLE_BLURBEHIND,
		}
		var data w32.WINDOWCOMPOSITIONATTRIBDATA
		data.Attrib = w32.WCA_ACCENT_POLICY
		data.PvData = w32.PVOID(&accent)
		data.CbData = unsafe.Sizeof(accent)

		w32.SetWindowCompositionAttribute(w.hwnd, &data)
	} else {
		w32.EnableTranslucency(w.hwnd, int32(backdropType))
	}
}

func (w *windowsWebviewWindow) setIcon(icon w32.HICON) {
	w32.SendMessage(w.hwnd, w32.WM_SETICON, w32.ICON_BIG, icon)
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
	case w32.WM_ACTIVATE:
		if int(wparam&0xffff) == w32.WA_INACTIVE {
			w.parent.emit(events.Windows.WindowInactive)
		}
		if wparam == w32.WA_ACTIVE {
			getNativeApplication().currentWindowID = w.parent.id
			w.parent.emit(events.Windows.WindowActive)
		}
		if wparam == w32.WA_CLICKACTIVE {
			getNativeApplication().currentWindowID = w.parent.id
			w.parent.emit(events.Windows.WindowClickActive)
		}
		// If we want to have a frameless window but with the default frame decorations, extend the DWM client area.
		// This Option is not affected by returning 0 in WM_NCCALCSIZE.
		// As a result we have hidden the titlebar but still have the default window frame styling.
		// See: https://docs.microsoft.com/en-us/windows/win32/api/dwmapi/nf-dwmapi-dwmextendframeintoclientarea#remarks
		if w.framelessWithDecorations() {
			err := w32.ExtendFrameIntoClientArea(w.hwnd, true)
			if err != nil {
				globalApplication.handleFatalError(err)
			}
		}
	case w32.WM_CLOSE:
		w.parent.emit(events.Windows.WindowClose)
		return 0
	case w32.WM_KILLFOCUS:
		if w.focusingChromium {
			return 0
		}
		w.parent.emit(events.Windows.WindowKillFocus)
	case w32.WM_ENTERSIZEMOVE:
		// This is needed to close open dropdowns when moving the window https://github.com/MicrosoftEdge/WebView2Feedback/issues/2290
		w32.SetFocus(w.hwnd)
	case w32.WM_SETFOCUS:
		w.focus()
		w.parent.emit(events.Windows.WindowSetFocus)
	case w32.WM_MOVE, w32.WM_MOVING:
		_ = w.chromium.NotifyParentWindowPositionChanged()
		if w.moveDebouncer == nil {
			w.moveDebouncer = debounce.New(time.Duration(windowDidMoveDebounceMS) * time.Millisecond)
		}
		w.moveDebouncer(func() {
			w.parent.emit(events.Windows.WindowDidMove)
		})
	// Check for keypress
	case w32.WM_KEYDOWN:
		w.processKeyBinding(uint(wparam))
	case w32.WM_SIZE:
		switch wparam {
		case w32.SIZE_MAXIMIZED:
			w.parent.emit(events.Windows.WindowMaximise)
		case w32.SIZE_RESTORED:
			w.parent.emit(events.Windows.WindowRestore)
		case w32.SIZE_MINIMIZED:
			w.parent.emit(events.Windows.WindowMinimise)
		}
		if w.parent.options.Frameless && wparam == w32.SIZE_MINIMIZED {
			// If the window is frameless, and we are minimizing, then we need to suppress the Resize on the
			// WebView2. If we don't do this, restoring does not work as expected and first restores with some wrong
			// size during the restore animation and only fully renders when the animation is done. This highly
			// depends on the content in the WebView, see https://github.com/wailsapp/wails/issues/1319
		} else if w.resizeDebouncer != nil {
			w.resizeDebouncer(func() {
				InvokeSync(func() {
					w.chromium.Resize()
				})
				w.parent.emit(events.Windows.WindowDidResize)
			})
		} else {
			w.chromium.Resize()
		}
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
		w.parent.emit(events.Common.WindowDPIChanged)
	}

	if w.parent.options.Windows.WindowMask != nil {
		switch msg {
		case w32.WM_NCHITTEST:
			if w.parent.options.Windows.WindowMaskDraggable {
				return w32.HTCAPTION
			}
			return w32.HTCLIENT
		}
	}

	if w.menu != nil || w.currentlyOpenContextMenu != nil {
		switch msg {
		case w32.WM_COMMAND:
			cmdMsgID := int(wparam & 0xffff)
			switch cmdMsgID {
			default:
				var processed bool
				if w.currentlyOpenContextMenu != nil {
					processed = w.currentlyOpenContextMenu.ProcessCommand(cmdMsgID)
					w.currentlyOpenContextMenu = nil

				}
				if !processed && w.menu != nil {
					processed = w.menu.ProcessCommand(cmdMsgID)
				}
			}
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
				err := w32.ExtendFrameIntoClientArea(w.hwnd, true)
				if err != nil {
					globalApplication.handleFatalError(err)
				}
			}

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
					w.chromium.SetPadding(edge.Rect{})
				} else if w.isMaximised() {
					// If the window is maximized we must adjust the client area to the work area of the monitor. Otherwise
					// some content goes beyond the visible part of the monitor.
					// Make sure to use the provided RECT to get the monitor, because during maximizig there might be
					// a wrong monitor returned in multiscreen mode when using MonitorFromWindow.
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
					w.chromium.SetPadding(edge.Rect{})
				} else {
					// This is needed to workaround the resize flickering in frameless mode with WindowDecorations
					// See: https://stackoverflow.com/a/6558508
					// The workaround originally suggests to decrese the bottom 1px, but that seems to bring up a thin
					// white line on some Windows-Versions, due to DrawBackground using also this reduces ClientSize.
					// Increasing the bottom also worksaround the flickering but we would loose 1px of the WebView content
					// therefore let's pad the content with 1px at the bottom.
					rgrc.Bottom += 1
					w.chromium.SetPadding(edge.Rect{Bottom: 1})
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

// processMessage is given a message sent from JS via the postMessage API
// We put it on the global window message buffer to be processed centrally
func (w *windowsWebviewWindow) processMessage(message string) {
	// We send all messages to the centralised window message buffer
	windowMessageBuffer <- &windowMessage{
		windowId: w.parent.id,
		message:  message,
	}
}

func (w *windowsWebviewWindow) processRequest(req *edge.ICoreWebView2WebResourceRequest, args *edge.ICoreWebView2WebResourceRequestedEventArgs) {

	// Setting the UserAgent on the CoreWebView2Settings clears the whole default UserAgent of the Edge browser, but
	// we want to just append our ApplicationIdentifier. So we adjust the UserAgent for every request.
	if reqHeaders, err := req.GetHeaders(); err == nil {
		useragent, _ := reqHeaders.GetHeader(assetserver.HeaderUserAgent)
		useragent = strings.Join([]string{useragent, assetserver.WailsUserAgentValue}, " ")
		err = reqHeaders.SetHeader(assetserver.HeaderUserAgent, useragent)
		if err != nil {
			globalApplication.fatal("Error setting UserAgent header: " + err.Error())
		}
		err = reqHeaders.SetHeader(webViewRequestHeaderWindowId, strconv.FormatUint(uint64(w.parent.id), 10))
		if err != nil {
			globalApplication.fatal("Error setting WindowId header: " + err.Error())
		}
		err = reqHeaders.Release()
		if err != nil {
			globalApplication.fatal("Error releasing headers: " + err.Error())
		}
	}

	if globalApplication.assets == nil {
		// We are using the devServer let the WebView2 handle the request with its default handler
		return
	}

	//Get the request
	uri, _ := req.GetUri()
	reqUri, err := url.ParseRequestURI(uri)
	if err != nil {
		globalApplication.error("Unable to parse request uri: uri='%s' error='%s'", uri, err)
		return
	}

	if reqUri.Scheme != "http" {
		// Let the WebView2 handle the request with its default handler
		return
	} else if !strings.HasPrefix(reqUri.Host, "wails.localhost") {
		// Let the WebView2 handle the request with its default handler
		return
	}

	webviewRequest, err := webview.NewRequest(
		w.chromium.Environment(),
		args,
		func(fn func()) {
			InvokeSync(fn)
		})
	if err != nil {
		globalApplication.error("%s: NewRequest failed: %s", uri, err)
		return
	}

	webviewRequests <- &webViewAssetRequest{
		Request:    webviewRequest,
		windowId:   w.parent.id,
		windowName: w.parent.options.Name,
	}
}

func (w *windowsWebviewWindow) setupChromium() {
	chromium := w.chromium
	debugMode := globalApplication.isDebugMode

	opts := w.parent.options.Windows

	webview2version, err := webviewloader.GetAvailableCoreWebView2BrowserVersionString(globalApplication.options.Windows.WebviewBrowserPath)
	if err != nil {
		globalApplication.error("Error getting WebView2 version: " + err.Error())
		return
	}
	globalApplication.capabilities = capabilities.NewCapabilities(webview2version)

	disableFeatues := []string{}
	if !opts.EnableFraudulentWebsiteWarnings {
		disableFeatues = append(disableFeatues, "msSmartScreenProtection")
	}
	if opts.WebviewGpuIsDisabled {
		chromium.AdditionalBrowserArgs = append(chromium.AdditionalBrowserArgs, "--disable-gpu")
	}

	if len(disableFeatues) > 0 {
		arg := fmt.Sprintf("--disable-features=%s", strings.Join(disableFeatues, ","))
		chromium.AdditionalBrowserArgs = append(chromium.AdditionalBrowserArgs, arg)
	}

	enableFeatures := []string{"msWebView2BrowserHitTransparent"}
	if len(enableFeatures) > 0 {
		arg := fmt.Sprintf("--enable-features=%s", strings.Join(enableFeatures, ","))
		chromium.AdditionalBrowserArgs = append(chromium.AdditionalBrowserArgs, arg)
	}

	chromium.DataPath = globalApplication.options.Windows.WebviewUserDataPath
	chromium.BrowserPath = globalApplication.options.Windows.WebviewBrowserPath

	if opts.Permissions != nil {
		for permission, state := range opts.Permissions {
			chromium.SetPermission(edge.CoreWebView2PermissionKind(permission),
				edge.CoreWebView2PermissionState(state))
		}
	}

	chromium.MessageCallback = w.processMessage
	chromium.MessageWithAdditionalObjectsCallback = w.processMessageWithAdditionalObjects
	chromium.WebResourceRequestedCallback = w.processRequest
	chromium.ContainsFullScreenElementChangedCallback = w.fullscreenChanged
	chromium.NavigationCompletedCallback = w.navigationCompleted
	chromium.AcceleratorKeyCallback = w.processKeyBinding

	chromium.Embed(w.hwnd)

	if chromium.HasCapability(edge.SwipeNavigation) {
		err := chromium.PutIsSwipeNavigationEnabled(opts.EnableSwipeGestures)
		if err != nil {
			globalApplication.fatal(err.Error())
		}
	}

	if chromium.HasCapability(edge.AllowExternalDrop) {
		err := chromium.AllowExternalDrag(false)
		if err != nil {
			globalApplication.fatal(err.Error())
		}
	}
	if w.parent.options.EnableDragAndDrop {
		w.dropTarget = w32.NewDropTarget()
		w.dropTarget.OnDrop = func(files []string) {
			w.parent.emit(events.Windows.WindowDragDrop)
			windowDragAndDropBuffer <- &dragAndDropMessage{
				windowId:  windowID,
				filenames: files,
			}
		}
		if opts.OnEnterEffect != 0 {
			w.dropTarget.OnEnterEffect = convertEffect(opts.OnEnterEffect)
		}
		if opts.OnOverEffect != 0 {
			w.dropTarget.OnOverEffect = convertEffect(opts.OnOverEffect)
		}
		w.dropTarget.OnEnter = func() {
			w.parent.emit(events.Windows.WindowDragEnter)
		}
		w.dropTarget.OnLeave = func() {
			w.parent.emit(events.Windows.WindowDragLeave)
		}
		w.dropTarget.OnOver = func() {
			w.parent.emit(events.Windows.WindowDragOver)
		}
		// Enumerate all the child windows for this window and register them as drop targets
		w32.EnumChildWindows(w.hwnd, func(hwnd w32.HWND, lparam w32.LPARAM) w32.LRESULT {
			// Check if the window class is "Chrome_RenderWidgetHostHWND"
			// If it is, then we register it as a drop target
			//windowName := w32.GetClassName(hwnd)
			//println(windowName)
			//if windowName == "Chrome_RenderWidgetHostHWND" {
			err := w32.RegisterDragDrop(hwnd, w.dropTarget)
			if err != nil && err != syscall.Errno(w32.DRAGDROP_E_ALREADYREGISTERED) {
				globalApplication.error("Error registering drag and drop: " + err.Error())
			}
			//}
			return 1
		})

	}

	// event mapping
	w.parent.On(events.Windows.WindowDidMove, func(e *WindowEvent) {
		w.parent.emit(events.Common.WindowDidMove)
	})
	w.parent.On(events.Windows.WindowDidResize, func(e *WindowEvent) {
		w.parent.emit(events.Common.WindowDidResize)
	})

	// We will get round to this
	//if chromium.HasCapability(edge.AllowExternalDrop) {
	//	err := chromium.AllowExternalDrag(w.parent.options.EnableDragAndDrop)
	//	if err != nil {
	//		globalApplication.fatal(err.Error())
	//	}
	//	if w.parent.options.EnableDragAndDrop {
	//		chromium.MessageWithAdditionalObjectsCallback = w.processMessageWithAdditionalObjects
	//	}
	//}

	chromium.Resize()
	settings, err := chromium.GetSettings()
	if err != nil {
		globalApplication.fatal(err.Error())
	}
	err = settings.PutAreDefaultContextMenusEnabled(debugMode || !w.parent.options.DefaultContextMenuDisabled)
	if err != nil {
		globalApplication.fatal(err.Error())
	}

	w.enableDevTools(settings)

	if w.parent.options.Zoom > 0.0 {
		chromium.PutZoomFactor(w.parent.options.Zoom)
	}
	err = settings.PutIsZoomControlEnabled(w.parent.options.ZoomControlEnabled)
	if err != nil {
		globalApplication.fatal(err.Error())
	}

	err = settings.PutIsStatusBarEnabled(false)
	if err != nil {
		globalApplication.fatal(err.Error())
	}
	err = settings.PutAreBrowserAcceleratorKeysEnabled(false)
	if err != nil {
		globalApplication.fatal(err.Error())
	}
	err = settings.PutIsSwipeNavigationEnabled(false)
	if err != nil {
		globalApplication.fatal(err.Error())
	}

	if debugMode && w.parent.options.OpenInspectorOnStartup {
		chromium.OpenDevToolsWindow()
	}

	// Set background colour
	w.setBackgroundColour(w.parent.options.BackgroundColour)

	chromium.SetGlobalPermission(edge.CoreWebView2PermissionStateAllow)
	chromium.AddWebResourceRequestedFilter("*", edge.COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)

	if w.parent.options.HTML != "" {
		var script string
		if w.parent.options.JS != "" {
			script = w.parent.options.JS
		}
		if w.parent.options.CSS != "" {
			script += fmt.Sprintf("; addEventListener(\"DOMContentLoaded\", (event) => { document.head.appendChild(document.createElement('style')).innerHTML=\"%s\"; });", strings.ReplaceAll(w.parent.options.CSS, `"`, `\"`))
		}
		chromium.Init(script)
		chromium.NavigateToString(w.parent.options.HTML)
	} else {
		startURL, err := assetserver.GetStartURL(w.parent.options.URL)
		if err != nil {
			globalApplication.fatal(err.Error())
		}
		chromium.Navigate(startURL)
	}

}

func (w *windowsWebviewWindow) fullscreenChanged(sender *edge.ICoreWebView2, _ *edge.ICoreWebView2ContainsFullScreenElementChangedEventArgs) {
	isFullscreen, err := sender.GetContainsFullScreenElement()
	if err != nil {
		globalApplication.fatal("Fatal error in callback fullscreenChanged: " + err.Error())
	}
	if isFullscreen {
		w.fullscreen()
	} else {
		w.unfullscreen()
	}
}

func convertEffect(effect DragEffect) w32.DWORD {
	switch effect {
	case DragEffectCopy:
		return w32.DROPEFFECT_COPY
	case DragEffectMove:
		return w32.DROPEFFECT_MOVE
	case DragEffectLink:
		return w32.DROPEFFECT_LINK
	default:
		return w32.DROPEFFECT_NONE
	}
}

func (w *windowsWebviewWindow) flash(enabled bool) {
	w32.FlashWindow(w.hwnd, enabled)
}

func (w *windowsWebviewWindow) navigationCompleted(sender *edge.ICoreWebView2, args *edge.ICoreWebView2NavigationCompletedEventArgs) {

	// Install the runtime core
	w.execJS(runtime.Core())

	// Emit DomReady Event
	windowEvents <- &windowEvent{EventID: uint(events.Windows.WebViewNavigationCompleted), WindowID: w.parent.id}

	if w.hasStarted {
		// NavigationCompleted is triggered for every Load. If an application uses reloads the Hide/Show will trigger
		// a flickering of the window with every reload. So we only do this once for the first NavigationCompleted.
		return
	}
	w.hasStarted = true

	wasFocused := w.isFocused()
	// Hack to make it visible: https://github.com/MicrosoftEdge/WebView2Feedback/issues/1077#issuecomment-825375026
	err := w.chromium.Hide()
	if err != nil {
		globalApplication.fatal(err.Error())
	}
	err = w.chromium.Show()
	if err != nil {
		globalApplication.fatal(err.Error())
	}
	if wasFocused {
		w.focus()
	}

	//f.mainWindow.hasBeenShown = true

}

func (w *windowsWebviewWindow) processKeyBinding(vkey uint) bool {

	globalApplication.debug("Processing key binding", "vkey", vkey)

	// Get the keyboard state and convert to an accelerator
	var keyState [256]byte
	if !w32.GetKeyboardState(keyState[:]) {
		globalApplication.error("Error getting keyboard state")
		return false
	}

	var acc accelerator
	// Check if CTRL is pressed
	if keyState[w32.VK_CONTROL]&0x80 != 0 {
		acc.Modifiers = append(acc.Modifiers, ControlKey)
	}
	// Check if ALT is pressed
	if keyState[w32.VK_MENU]&0x80 != 0 {
		acc.Modifiers = append(acc.Modifiers, OptionOrAltKey)
	}
	// Check if SHIFT is pressed
	if keyState[w32.VK_SHIFT]&0x80 != 0 {
		acc.Modifiers = append(acc.Modifiers, ShiftKey)
	}
	// Check if WIN is pressed
	if keyState[w32.VK_LWIN]&0x80 != 0 || keyState[w32.VK_RWIN]&0x80 != 0 {
		acc.Modifiers = append(acc.Modifiers, SuperKey)
	}

	if vkey != w32.VK_CONTROL && vkey != w32.VK_MENU && vkey != w32.VK_SHIFT && vkey != w32.VK_LWIN && vkey != w32.VK_RWIN {
		// Convert the vkey to a string
		accKey, ok := VirtualKeyCodes[vkey]
		if !ok {
			return false
		}
		acc.Key = accKey
	}

	accKey := acc.String()
	globalApplication.debug("Processing key binding", "vkey", vkey, "acc", accKey)

	// Process the key binding
	if w.parent.processKeyBinding(accKey) {
		return true
	}

	if accKey == "alt+f4" {
		w32.PostMessage(w.hwnd, w32.WM_CLOSE, 0, 0)
		return true
	}

	return false
}

func (w *windowsWebviewWindow) processMessageWithAdditionalObjects(message string, sender *edge.ICoreWebView2, args *edge.ICoreWebView2WebMessageReceivedEventArgs) {
	if strings.HasPrefix(message, "FilesDropped") {
		objs, err := args.GetAdditionalObjects()
		if err != nil {
			globalApplication.error(err.Error())
			return
		}

		defer func() {
			err = objs.Release()
			if err != nil {
				globalApplication.error("Error releasing objects: " + err.Error())
			}
		}()

		count, err := objs.GetCount()
		if err != nil {
			globalApplication.error(err.Error())
			return
		}

		var filenames []string
		for i := uint32(0); i < count; i++ {
			_file, err := objs.GetValueAtIndex(i)
			if err != nil {
				globalApplication.error("cannot get value at %d : %s", i, err.Error())
				return
			}

			file := (*edge.ICoreWebView2File)(unsafe.Pointer(_file))

			// TODO: Fix this
			defer file.Release()

			filepath, err := file.GetPath()
			if err != nil {
				globalApplication.error("cannot get path for object at %d : %s", i, err.Error())
				return
			}

			filenames = append(filenames, filepath)
		}

		addDragAndDropMessage(w.parent.id, filenames)
		return
	}
}

func (w *windowsWebviewWindow) setMaximiseButtonEnabled(enabled bool) {
	w.setStyle(enabled, w32.WS_MAXIMIZEBOX)
}

func (w *windowsWebviewWindow) setMinimiseButtonEnabled(enabled bool) {
	w.setStyle(enabled, w32.WS_MINIMIZEBOX)
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

func (w *windowsWebviewWindow) setMinimiseButtonState(state ButtonState) {
	switch state {
	case ButtonDisabled, ButtonHidden:
		w.setStyle(false, w32.WS_MINIMIZEBOX)
	case ButtonEnabled:
		w.setStyle(true, w32.WS_SYSMENU)
		w.setStyle(true, w32.WS_MINIMIZEBOX)

	}
}

func (w *windowsWebviewWindow) setMaximiseButtonState(state ButtonState) {
	switch state {
	case ButtonDisabled, ButtonHidden:
		w.setStyle(false, w32.WS_MAXIMIZEBOX)
	case ButtonEnabled:
		w.setStyle(true, w32.WS_SYSMENU)
		w.setStyle(true, w32.WS_MAXIMIZEBOX)
	}
}

func (w *windowsWebviewWindow) setCloseButtonState(state ButtonState) {
	switch state {
	case ButtonEnabled:
		w.setStyle(true, w32.WS_SYSMENU)
		_ = w32.EnableCloseButton(w.hwnd)
	case ButtonDisabled:
		w.setStyle(true, w32.WS_SYSMENU)
		_ = w32.DisableCloseButton(w.hwnd)
	case ButtonHidden:
		w.setStyle(false, w32.WS_SYSMENU)
	}
}

func (w *windowsWebviewWindow) setGWLStyle(style int) {
	w32.SetWindowLong(w.hwnd, w32.GWL_STYLE, uint32(style))
}

func (w *windowsWebviewWindow) isIgnoreMouseEvents() bool {
	exStyle := w32.GetWindowLong(w.hwnd, w32.GWL_EXSTYLE)
	return exStyle&w32.WS_EX_TRANSPARENT != 0
}

func (w *windowsWebviewWindow) setIgnoreMouseEvents(ignore bool) {
	exStyle := w32.GetWindowLong(w.hwnd, w32.GWL_EXSTYLE)
	if ignore {
		exStyle |= w32.WS_EX_LAYERED | w32.WS_EX_TRANSPARENT
	} else {
		exStyle &^= w32.WS_EX_TRANSPARENT
	}
	w32.SetWindowLong(w.hwnd, w32.GWL_EXSTYLE, uint32(exStyle))
}
