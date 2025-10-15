//go:build windows

package application

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
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
	ignoreDPIChangeResizing  bool

	// Fullscreen flags
	isCurrentlyFullscreen   bool
	previousWindowStyle     uint32
	previousWindowExStyle   uint32
	previousWindowPlacement w32.WINDOWPLACEMENT

	// Webview
	chromium                   *edge.Chromium
	webviewNavigationCompleted bool

	// Window visibility management - robust fallback for issue #2861
	showRequested     bool        // Track if show() was called before navigation completed
	visibilityTimeout *time.Timer // Timeout to show window if navigation is delayed
	windowShown       bool        // Track if window container has been shown
	// Track whether content protection has been applied to the native window yet
	contentProtectionApplied bool

	// resizeBorder* is the width/height of the resize border in pixels.
	resizeBorderWidth  int32
	resizeBorderHeight int32
	focusingChromium   bool
	dropTarget         *w32.DropTarget
	onceDo             sync.Once

	// Window move debouncer
	moveDebouncer   func(func())
	resizeDebouncer func(func())

	// isMinimizing indicates whether the window is currently being minimized
	// Used to prevent unnecessary redraws during minimize/restore operations
	isMinimizing bool

	// menubarTheme is the theme for the menubar
	menubarTheme *w32.MenuBarTheme
}

func (w *windowsWebviewWindow) setMenu(menu *Menu) {
	menu.Update()
	w.menu = NewApplicationMenu(w, menu)
	w.menu.parentWindow = w
	w32.SetMenu(w.hwnd, w.menu.menu)

	// Set menu background if theme is active
	if w.menubarTheme != nil {
		globalApplication.debug("Applying menubar theme in setMenu", "window", w.parent.id)
		w.menubarTheme.SetMenuBackground(w.menu.menu)
		w32.DrawMenuBar(w.hwnd)
		// Force a repaint of the menu area
		w32.InvalidateRect(w.hwnd, nil, true)
	} else {
		globalApplication.debug("No menubar theme to apply in setMenu", "window", w.parent.id)
	}

	// Check if using translucent background with Mica - this makes menubars invisible
	if w.parent.options.BackgroundType == BackgroundTypeTranslucent &&
		(w.parent.options.Windows.BackdropType == Mica ||
			w.parent.options.Windows.BackdropType == Acrylic ||
			w.parent.options.Windows.BackdropType == Tabbed) {
		// Log warning about menubar visibility issue
		globalApplication.debug("Warning: Menubars may be invisible when using translucent backgrounds with Mica/Acrylic/Tabbed effects", "window", w.parent.id)
	}
}

func (w *windowsWebviewWindow) cut() {
	w.execJS("document.execCommand('cut')")
}

func (w *windowsWebviewWindow) paste() {
	w.execJS(`
		(async () => {
			try {
				// Try to read all available formats
				const clipboardItems = await navigator.clipboard.read();
				
				for (const clipboardItem of clipboardItems) {
					// Check for image types
					for (const type of clipboardItem.types) {
						if (type.startsWith('image/')) {
							const blob = await clipboardItem.getType(type);
							const url = URL.createObjectURL(blob);
							document.execCommand('insertHTML', false, '<img src="' + url + '">');
							return;
						}
					}
					
					// If no image found, try text
					if (clipboardItem.types.includes('text/plain')) {
						const text = await navigator.clipboard.readText();
						document.execCommand('insertText', false, text);
						return;
					}
				}
			} catch(err) {
				// Fallback to text-only paste if clipboard access fails
				try {
					const text = await navigator.clipboard.readText();
					document.execCommand('insertText', false, text);
				} catch(fallbackErr) {
					console.error('Failed to paste:', err, fallbackErr);
				}
			}
		})()
	`)
}

func (w *windowsWebviewWindow) copy() {
	w.execJS(`
		(async () => {
			try {
				const selection = window.getSelection();
				if (!selection.rangeCount) return;

				const range = selection.getRangeAt(0);
				const container = document.createElement('div');
				container.appendChild(range.cloneContents());

				// Check if we have any images in the selection
				const images = container.getElementsByTagName('img');
				if (images.length > 0) {
					// Handle image copy
					const img = images[0]; // Take the first image
					const response = await fetch(img.src);
					const blob = await response.blob();
					await navigator.clipboard.write([
						new ClipboardItem({
							[blob.type]: blob
						})
					]);
				} else {
					// Handle text copy
					const text = selection.toString();
					if (text) {
						await navigator.clipboard.writeText(text);
					}
				}
			} catch(err) {
				console.error('Failed to copy:', err);
			}
		})()
	`)
}

func (w *windowsWebviewWindow) selectAll() {
	w.execJS("document.execCommand('selectAll')")
}

func (w *windowsWebviewWindow) undo() {
	w.execJS("document.execCommand('undo')")
}

func (w *windowsWebviewWindow) redo() {
	w.execJS("document.execCommand('redo')")
}

func (w *windowsWebviewWindow) delete() {
	w.execJS("document.execCommand('delete')")
}

func (w *windowsWebviewWindow) handleKeyEvent(_ string) {
	// Unused on windows
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
		return errors.New("unable to release mouse capture")
	}
	// Use PostMessage because we don't want to block the caller until resizing has been finished.
	w32.PostMessage(w.hwnd, w32.WM_NCLBUTTONDOWN, edgeMap[border], 0)
	return nil
}

func (w *windowsWebviewWindow) startDrag() error {
	if !w32.ReleaseCapture() {
		return errors.New("unable to release mouse capture")
	}
	// Use PostMessage because we don't want to block the caller until dragging has been finished.
	w32.PostMessage(w.hwnd, w32.WM_NCLBUTTONDOWN, w32.HTCAPTION, 0)
	return nil
}

func (w *windowsWebviewWindow) nativeWindow() unsafe.Pointer {
	return unsafe.Pointer(w.hwnd)
}

func (w *windowsWebviewWindow) setTitle(title string) {
	w32.SetWindowText(w.hwnd, title)
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
	w.webviewNavigationCompleted = false
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
	switch w.parent.options.BackgroundType {
	case BackgroundTypeSolid:
		w32.SetBackgroundColour(w.hwnd, color.Red, color.Green, color.Blue)
		w.chromium.SetBackgroundColour(color.Red, color.Green, color.Blue, color.Alpha)
	case BackgroundTypeTransparent, BackgroundTypeTranslucent:
		w.chromium.SetBackgroundColour(0, 0, 0, 0)
	}
}

func (w *windowsWebviewWindow) framelessWithDecorations() bool {
	return w.parent.options.Frameless && !w.parent.options.Windows.DisableFramelessWindowDecorations
}

func (w *windowsWebviewWindow) run() {

	options := w.parent.options

	// Initialize showRequested based on whether window should be hidden
	// Non-hidden windows should be shown by default
	w.showRequested = !options.Hidden

	w.chromium = edge.NewChromium()
	if globalApplication.options.ErrorHandler != nil {
		w.chromium.SetErrorCallback(globalApplication.options.ErrorHandler)
	}

	exStyle := w32.WS_EX_CONTROLPARENT
	if options.BackgroundType != BackgroundTypeSolid {
		if (options.Frameless && options.BackgroundType == BackgroundTypeTransparent) ||
			w.parent.options.IgnoreMouseEvents {
			// Always if transparent and frameless
			exStyle |= w32.WS_EX_TRANSPARENT | w32.WS_EX_LAYERED
		} else {
			// Only WS_EX_NOREDIRECTIONBITMAP if not (and not solid)
			exStyle |= w32.WS_EX_NOREDIRECTIONBITMAP
		}
	}
	if options.AlwaysOnTop {
		exStyle |= w32.WS_EX_TOPMOST
	}
	// If we're frameless, we need to add the WS_EX_TOOLWINDOW style to hide the window from the taskbar
	if options.Windows.HiddenOnTaskbar {
		//exStyle |= w32.WS_EX_TOOLWINDOW
		exStyle |= w32.WS_EX_NOACTIVATE
	} else {
		exStyle |= w32.WS_EX_APPWINDOW
	}

	if options.Windows.ExStyle != 0 {
		exStyle = options.Windows.ExStyle
	}

	bounds := Rect{
		X:      options.X,
		Y:      options.Y,
		Width:  options.Width,
		Height: options.Height,
	}
	initialScreen := ScreenNearestDipRect(bounds)
	physicalBounds := initialScreen.dipToPhysicalRect(bounds)

	// Default window position applied by the system
	// TODO: provide a way to set (0,0) as an initial position?
	if options.X == 0 && options.Y == 0 {
		physicalBounds.X = w32.CW_USEDEFAULT
		physicalBounds.Y = w32.CW_USEDEFAULT
	}

	var appMenu w32.HMENU

	// Process Menu
	if !options.Frameless {
		userMenu := w.parent.options.Windows.Menu
		if userMenu != nil {
			userMenu.Update()
			w.menu = NewApplicationMenu(w, userMenu)
			w.menu.parentWindow = w
			appMenu = w.menu.menu
		}
	}

	var parent w32.HWND

	var style uint = w32.WS_OVERLAPPEDWINDOW
	// If the window should be hidden initially, exclude WS_VISIBLE from the style
	// This prevents the white window flash reported in issue #4611
	if options.Hidden {
		style = style &^ uint(w32.WS_VISIBLE)
	}

	w.hwnd = w32.CreateWindowEx(
		uint(exStyle),
		w32.MustStringToUTF16Ptr(globalApplication.options.Windows.WndClass),
		w32.MustStringToUTF16Ptr(options.Title),
		style,
		physicalBounds.X,
		physicalBounds.Y,
		physicalBounds.Width,
		physicalBounds.Height,
		parent,
		appMenu,
		w32.GetModuleHandle(""),
		nil)

	if w.hwnd == 0 {
		globalApplication.fatal("unable to create window")
	}

	// Ensure correct window size in case the scale factor of current screen is different from the initial one.
	// This could happen when using the default window position and the window launches on a secondary monitor.
	currentScreen, _ := w.getScreen()
	if currentScreen.ScaleFactor != initialScreen.ScaleFactor {
		w.setSize(options.Width, options.Height)
	}

	w.setupChromium()

	if options.Windows.WindowDidMoveDebounceMS == 0 {
		options.Windows.WindowDidMoveDebounceMS = 50
	}
	w.moveDebouncer = debounce.New(
		time.Duration(options.Windows.WindowDidMoveDebounceMS) * time.Millisecond,
	)

	if options.Windows.ResizeDebounceMS > 0 {
		w.resizeDebouncer = debounce.New(
			time.Duration(options.Windows.ResizeDebounceMS) * time.Millisecond,
		)
	}

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
		isDark := w32.IsCurrentlyDarkMode()
		if isDark {
			w32.AllowDarkModeForWindow(w.hwnd, true)
		}
		w.updateTheme(isDark)
		// Don't initialize default dark theme here if custom theme might be set
		// The updateTheme call above will handle both default and custom themes
		w.parent.onApplicationEvent(events.Windows.SystemThemeChanged, func(*ApplicationEvent) {
			InvokeAsync(func() {
				w.updateTheme(w32.IsCurrentlyDarkMode())
			})
		})
	case Light:
		w.updateTheme(false)
	case Dark:
		w32.AllowDarkModeForWindow(w.hwnd, true)
		w.updateTheme(true)
		// Don't initialize default dark theme here if custom theme might be set
		// The updateTheme call above will handle custom themes
	}

	w.setBackgroundColour(options.BackgroundColour)
	if options.BackgroundType == BackgroundTypeTranslucent {
		w.setBackdropType(w.parent.options.Windows.BackdropType)
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

	if options.InitialPosition == WindowCentered {
		w.center()
	} else {
		w.setPosition(options.X, options.Y)
	}

	if options.Frameless {
		// Trigger a resize to ensure the window is sized correctly
		w.chromium.Resize()
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

func (w *windowsWebviewWindow) update() {
	w32.UpdateWindow(w.hwnd)
}

// getBorderSizes returns the extended border size for the window
func (w *windowsWebviewWindow) getBorderSizes() *LRTB {
	var result LRTB
	var frame w32.RECT
	w32.DwmGetWindowAttribute(
		w.hwnd,
		w32.DWMWA_EXTENDED_FRAME_BOUNDS,
		unsafe.Pointer(&frame),
		unsafe.Sizeof(frame),
	)
	rect := w32.GetWindowRect(w.hwnd)
	result.Left = int(frame.Left - rect.Left)
	result.Top = int(frame.Top - rect.Top)
	result.Right = int(rect.Right - frame.Right)
	result.Bottom = int(rect.Bottom - frame.Bottom)
	return &result
}

// convertWindowToWebviewCoordinates converts window-relative coordinates to webview-relative coordinates
func (w *windowsWebviewWindow) convertWindowToWebviewCoordinates(windowX, windowY int) (int, int) {
	// Get the client area of the window (this excludes borders, title bar, etc.)
	clientRect := w32.GetClientRect(w.hwnd)
	if clientRect == nil {
		// Fallback: return coordinates as-is if we can't get client rect
		globalApplication.debug("[DragDropDebug] convertWindowToWebviewCoordinates: Failed to get client rect, returning original coordinates", "windowX", windowX, "windowY", windowY)
		return windowX, windowY
	}

	// Get the window rect to calculate the offset
	windowRect := w32.GetWindowRect(w.hwnd)

	globalApplication.debug("[DragDropDebug] convertWindowToWebviewCoordinates: Input window coordinates", "windowX", windowX, "windowY", windowY)
	globalApplication.debug("[DragDropDebug] convertWindowToWebviewCoordinates: Window rect",
		"left", windowRect.Left, "top", windowRect.Top, "right", windowRect.Right, "bottom", windowRect.Bottom,
		"width", windowRect.Right-windowRect.Left, "height", windowRect.Bottom-windowRect.Top)
	globalApplication.debug("[DragDropDebug] convertWindowToWebviewCoordinates: Client rect",
		"left", clientRect.Left, "top", clientRect.Top, "right", clientRect.Right, "bottom", clientRect.Bottom,
		"width", clientRect.Right-clientRect.Left, "height", clientRect.Bottom-clientRect.Top)

	// Convert client (0,0) to screen coordinates to find where the client area starts
	var point w32.POINT
	point.X = 0
	point.Y = 0

	// Convert client (0,0) to screen coordinates
	clientX, clientY := w32.ClientToScreen(w.hwnd, int(point.X), int(point.Y))

	// The window coordinates from drag drop are relative to the window's top-left
	// But we need them relative to the client area's top-left
	// So we need to subtract the difference between window origin and client origin
	windowOriginX := int(windowRect.Left)
	windowOriginY := int(windowRect.Top)

	globalApplication.debug("[DragDropDebug] convertWindowToWebviewCoordinates: Client (0,0) in screen coordinates", "clientX", clientX, "clientY", clientY)
	globalApplication.debug("[DragDropDebug] convertWindowToWebviewCoordinates: Window origin in screen coordinates", "windowOriginX", windowOriginX, "windowOriginY", windowOriginY)

	// Calculate the offset from window origin to client origin
	offsetX := clientX - windowOriginX
	offsetY := clientY - windowOriginY

	globalApplication.debug("[DragDropDebug] convertWindowToWebviewCoordinates: Calculated offset", "offsetX", offsetX, "offsetY", offsetY)

	// Convert window-relative coordinates to webview-relative coordinates
	webviewX := windowX - offsetX
	webviewY := windowY - offsetY

	globalApplication.debug("[DragDropDebug] convertWindowToWebviewCoordinates: Final webview coordinates", "webviewX", webviewX, "webviewY", webviewY)

	return webviewX, webviewY
}

func (w *windowsWebviewWindow) physicalBounds() Rect {
	// var rect w32.RECT
	// // Get the extended frame bounds instead of the window rect to offset the invisible borders in Windows 10
	// w32.DwmGetWindowAttribute(w.hwnd, w32.DWMWA_EXTENDED_FRAME_BOUNDS, unsafe.Pointer(&rect), unsafe.Sizeof(rect))
	rect := w32.GetWindowRect(w.hwnd)
	return Rect{
		X:      int(rect.Left),
		Y:      int(rect.Top),
		Width:  int(rect.Right - rect.Left),
		Height: int(rect.Bottom - rect.Top),
	}
}

func (w *windowsWebviewWindow) setPhysicalBounds(physicalBounds Rect) {
	// // Offset invisible borders
	// borderSize := w.getBorderSizes()
	// physicalBounds.X -= borderSize.Left
	// physicalBounds.Y -= borderSize.Top
	// physicalBounds.Width += borderSize.Left + borderSize.Right
	// physicalBounds.Height += borderSize.Top + borderSize.Bottom

	// Set flag to ignore resizing the window with DPI change because we already calculated the correct size
	// for the target position, this prevents double resizing issue when the window is moved between screens
	previousFlag := w.ignoreDPIChangeResizing
	w.ignoreDPIChangeResizing = true
	w32.SetWindowPos(
		w.hwnd,
		0,
		physicalBounds.X,
		physicalBounds.Y,
		physicalBounds.Width,
		physicalBounds.Height,
		w32.SWP_NOZORDER|w32.SWP_NOACTIVATE,
	)
	w.ignoreDPIChangeResizing = previousFlag
}

// Get window dip bounds
func (w *windowsWebviewWindow) bounds() Rect {
	return PhysicalToDipRect(w.physicalBounds())
}

// Set window dip bounds
func (w *windowsWebviewWindow) setBounds(bounds Rect) {
	w.setPhysicalBounds(DipToPhysicalRect(bounds))
}

func (w *windowsWebviewWindow) size() (int, int) {
	bounds := w.bounds()
	return bounds.Width, bounds.Height
}

func (w *windowsWebviewWindow) width() int {
	return w.bounds().Width
}

func (w *windowsWebviewWindow) height() int {
	return w.bounds().Height
}

func (w *windowsWebviewWindow) setSize(width, height int) {
	bounds := w.bounds()
	bounds.Width = width
	bounds.Height = height

	w.setBounds(bounds)
}

func (w *windowsWebviewWindow) position() (int, int) {
	bounds := w.bounds()
	return bounds.X, bounds.Y
}

func (w *windowsWebviewWindow) setPosition(x int, y int) {
	bounds := w.bounds()
	bounds.X = x
	bounds.Y = y

	w.setBounds(bounds)
}

// Get window position relative to the screen WorkArea on which it is
func (w *windowsWebviewWindow) relativePosition() (int, int) {
	screen, _ := w.getScreen()
	pos := screen.absoluteToRelativeDipPoint(w.bounds().Origin())
	// Relative to WorkArea origin
	pos.X -= (screen.WorkArea.X - screen.X)
	pos.Y -= (screen.WorkArea.Y - screen.Y)
	return pos.X, pos.Y
}

// Set window position relative to the screen WorkArea on which it is
func (w *windowsWebviewWindow) setRelativePosition(x int, y int) {
	screen, _ := w.getScreen()
	pos := screen.relativeToAbsoluteDipPoint(Point{X: x, Y: y})
	// Relative to WorkArea origin
	pos.X += (screen.WorkArea.X - screen.X)
	pos.Y += (screen.WorkArea.Y - screen.Y)
	w.setPosition(pos.X, pos.Y)
}

func (w *windowsWebviewWindow) destroy() {
	w.parent.markAsDestroyed()
	if w.dropTarget != nil {
		w.dropTarget.Release()
	}
	// destroy the window
	w32.DestroyWindow(w.hwnd)
}

func (w *windowsWebviewWindow) reload() {
	w.execJS("window.location.reload();")
}

func (w *windowsWebviewWindow) forceReload() {
	// noop
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
	// Send WM_CLOSE message to trigger the same flow as clicking the X button
	w32.SendMessage(w.hwnd, w32.WM_CLOSE, 0, 0)
}

func (w *windowsWebviewWindow) zoom() {
	// Noop
}

func (w *windowsWebviewWindow) setHTML(html string) {
	// Render the given HTML in the webview window
	w.execJS(fmt.Sprintf("document.documentElement.innerHTML = %q;", html))
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
	w.parent.emit(events.Windows.WindowUnMaximise)
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
	w32.SetWindowLong(
		w.hwnd,
		w32.GWL_STYLE,
		w.previousWindowStyle & ^uint32(w32.WS_OVERLAPPEDWINDOW) | (w32.WS_POPUP|w32.WS_VISIBLE),
	)
	w32.SetWindowLong(
		w.hwnd,
		w32.GWL_EXSTYLE,
		w.previousWindowExStyle & ^uint32(w32.WS_EX_DLGMODALFRAME),
	)
	w.isCurrentlyFullscreen = true
	w32.SetWindowPos(w.hwnd, w32.HWND_TOP,
		int(monitorInfo.RcMonitor.Left),
		int(monitorInfo.RcMonitor.Top),
		int(monitorInfo.RcMonitor.Right-monitorInfo.RcMonitor.Left),
		int(monitorInfo.RcMonitor.Bottom-monitorInfo.RcMonitor.Top),
		w32.SWP_NOOWNERZORDER|w32.SWP_FRAMECHANGED)

	// Hide the menubar in fullscreen mode
	w32.SetMenu(w.hwnd, 0)

	w.chromium.Focus()
	w.parent.emit(events.Windows.WindowFullscreen)
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

	// Restore the menubar when exiting fullscreen
	if w.menu != nil {
		w32.SetMenu(w.hwnd, w.menu.menu)
	}

	w32.SetWindowPos(w.hwnd, 0, 0, 0, 0, 0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOZORDER|w32.SWP_NOOWNERZORDER|w32.SWP_FRAMECHANGED)
	w.enableSizeConstraints()
	w.parent.emit(events.Windows.WindowUnFullscreen)
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

func (w *windowsWebviewWindow) focus() {
	w32.SetForegroundWindow(w.hwnd)

	if w.isDisabled() {
		return
	}
	if w.isMinimised() {
		w.unminimise()
	}

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
	// Always show the window container immediately (decouple from WebView state)
	// This fixes issue #2861 where efficiency mode prevents window visibility
	w32.ShowWindow(w.hwnd, w32.SW_SHOW)
	w.windowShown = true
	w.showRequested = true
	w.updateContentProtection()

	// Show WebView if navigation has completed
	if w.webviewNavigationCompleted {
		w.chromium.Show()
		// Cancel timeout since we can show immediately
		if w.visibilityTimeout != nil {
			w.visibilityTimeout.Stop()
			w.visibilityTimeout = nil
		}
	} else {
		// Start timeout to show WebView if navigation is delayed (fallback for efficiency mode)
		if w.visibilityTimeout == nil {
			w.visibilityTimeout = time.AfterFunc(3*time.Second, func() {
				// Show WebView even if navigation hasn't completed
				// This prevents permanent invisibility in efficiency mode
				if !w.webviewNavigationCompleted && w.chromium != nil {
					w.chromium.Show()
				}
				w.visibilityTimeout = nil
			})
		}
	}
}

func (w *windowsWebviewWindow) hide() {
	w32.ShowWindow(w.hwnd, w32.SW_HIDE)
	w.windowShown = false
	w.showRequested = false

	// Cancel any pending visibility timeout
	if w.visibilityTimeout != nil {
		w.visibilityTimeout.Stop()
		w.visibilityTimeout = nil
	}
}

// Get the screen for the current window
func (w *windowsWebviewWindow) getScreen() (*Screen, error) {
	return getScreenForWindow(w)
}

func (w *windowsWebviewWindow) setFrameless(b bool) {
	// Remove or add the frame
	if b {
		w32.SetWindowLong(w.hwnd, w32.GWL_STYLE, w32.WS_VISIBLE|w32.WS_POPUP)
	} else {
		w32.SetWindowLong(w.hwnd, w32.GWL_STYLE, w32.WS_VISIBLE|w32.WS_OVERLAPPEDWINDOW)
	}
	w32.SetWindowPos(
		w.hwnd,
		0,
		0,
		0,
		0,
		0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOZORDER|w32.SWP_FRAMECHANGED,
	)
}

func newWindowImpl(parent *WebviewWindow) *windowsWebviewWindow {
	result := &windowsWebviewWindow{
		parent:             parent,
		resizeBorderWidth:  int32(w32.GetSystemMetrics(w32.SM_CXSIZEFRAME)),
		resizeBorderHeight: int32(w32.GetSystemMetrics(w32.SM_CYSIZEFRAME)),
		// Initialize visibility tracking fields
		showRequested:     false,
		visibilityTimeout: nil,
		windowShown:       false,
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
		w32.EnableTranslucency(w.hwnd, uint32(backdropType))
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

func (w *windowsWebviewWindow) processThemeColour(fn func(w32.HWND, uint32), value *uint32) {
	if value == nil {
		return
	}
	fn(w.hwnd, *value)
}

func (w *windowsWebviewWindow) isDisabled() bool {
	style := uint32(w32.GetWindowLong(w.hwnd, w32.GWL_STYLE))
	return style&w32.WS_DISABLED != 0
}

func (w *windowsWebviewWindow) updateTheme(isDarkMode bool) {

	if w32.IsCurrentlyHighContrastMode() {
		return
	}

	if !w32.SupportsThemes() {
		return
	}

	w32.SetTheme(w.hwnd, isDarkMode)

	// Clear any existing theme first
	if w.menubarTheme != nil && !isDarkMode {
		// Reset menu to default Windows theme when switching to light mode
		w.menubarTheme = nil
		if w.menu != nil {
			// Clear the menu background by setting it to default
			var mi w32.MENUINFO
			mi.CbSize = uint32(unsafe.Sizeof(mi))
			mi.FMask = w32.MIIM_BACKGROUND | w32.MIIM_APPLYTOSUBMENUS
			mi.HbrBack = 0 // NULL brush resets to default
			w32.SetMenuInfo(w.menu.menu, &mi)
		}
	}

	// Custom theme processing
	customTheme := w.parent.options.Windows.CustomTheme
	// Custom theme
	if w32.SupportsCustomThemes() {
		var userTheme *MenuBarTheme
		if isDarkMode {
			userTheme = customTheme.DarkModeMenuBar
		} else {
			userTheme = customTheme.LightModeMenuBar
		}

		if userTheme != nil {
			modeStr := "light"
			if isDarkMode {
				modeStr = "dark"
			}
			globalApplication.debug("Setting custom "+modeStr+" menubar theme", "window", w.parent.id)
			w.menubarTheme = &w32.MenuBarTheme{
				TitleBarBackground:     userTheme.Default.Background,
				TitleBarText:           userTheme.Default.Text,
				MenuBarBackground:      userTheme.Default.Background, // Use default background for menubar
				MenuHoverBackground:    userTheme.Hover.Background,
				MenuHoverText:          userTheme.Hover.Text,
				MenuSelectedBackground: userTheme.Selected.Background,
				MenuSelectedText:       userTheme.Selected.Text,
			}
			w.menubarTheme.Init()

			// If menu is already set, update it
			if w.menu != nil {
				w.menubarTheme.SetMenuBackground(w.menu.menu)
				w32.DrawMenuBar(w.hwnd)
				w32.InvalidateRect(w.hwnd, nil, true)
			}
		} else if userTheme == nil && isDarkMode {
			// Use default dark theme if no custom theme provided
			globalApplication.debug("Setting default dark menubar theme", "window", w.parent.id)
			w.menubarTheme = &w32.MenuBarTheme{
				TitleBarBackground:     w32.RGBptr(45, 45, 45),    // Dark titlebar
				TitleBarText:           w32.RGBptr(222, 222, 222), // Slightly muted white
				MenuBarBackground:      w32.RGBptr(33, 33, 33),    // Standard dark mode (#212121)
				MenuHoverBackground:    w32.RGBptr(48, 48, 48),    // Slightly lighter for hover (#303030)
				MenuHoverText:          w32.RGBptr(222, 222, 222), // Slightly muted white
				MenuSelectedBackground: w32.RGBptr(48, 48, 48),    // Same as hover
				MenuSelectedText:       w32.RGBptr(222, 222, 222), // Slightly muted white
			}
			w.menubarTheme.Init()

			// If menu is already set, update it
			if w.menu != nil {
				w.menubarTheme.SetMenuBackground(w.menu.menu)
				w32.DrawMenuBar(w.hwnd)
				w32.InvalidateRect(w.hwnd, nil, true)
			}
		} else if userTheme == nil && !isDarkMode && w.menu != nil {
			// No custom theme for light mode - ensure menu is reset to default
			globalApplication.debug("Resetting menu to default light theme", "window", w.parent.id)
			var mi w32.MENUINFO
			mi.CbSize = uint32(unsafe.Sizeof(mi))
			mi.FMask = w32.MIIM_BACKGROUND | w32.MIIM_APPLYTOSUBMENUS
			mi.HbrBack = 0 // NULL brush resets to default
			w32.SetMenuInfo(w.menu.menu, &mi)
			w32.DrawMenuBar(w.hwnd)
			w32.InvalidateRect(w.hwnd, nil, true)
		}
		// Define a map for theme selection
		themeMap := map[bool]map[bool]*WindowTheme{
			true: { // Window is active
				true:  customTheme.DarkModeActive,  // Dark mode
				false: customTheme.LightModeActive, // Light mode
			},
			false: { // Window is inactive
				true:  customTheme.DarkModeInactive,  // Dark mode
				false: customTheme.LightModeInactive, // Light mode
			},
		}

		// Select the appropriate theme
		theme := themeMap[w.isActive()][isDarkMode]

		// Apply theme colors
		if theme != nil {
			w.processThemeColour(w32.SetTitleBarColour, theme.TitleBarColour)
			w.processThemeColour(w32.SetTitleTextColour, theme.TitleTextColour)
			w.processThemeColour(w32.SetBorderColour, theme.BorderColour)
		}
	}
}

func (w *windowsWebviewWindow) isActive() bool {
	return w32.GetForegroundWindow() == w.hwnd
}

var resizePending int32

func (w *windowsWebviewWindow) WndProc(msg uint32, wparam, lparam uintptr) uintptr {

	// Use the original implementation that works perfectly for maximized
	processed, code := w32.MenuBarWndProc(w.hwnd, msg, wparam, lparam, w.menubarTheme)
	if processed {
		return code
	}

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

		if atomic.LoadUint32(&w.parent.unconditionallyClose) == 0 {
			// We were called by `Close()` or pressing the close button on the window
			w.parent.emit(events.Windows.WindowClosing)
			return 0
		}

		defer func() {
			windowsApp := globalApplication.impl.(*windowsApp)
			windowsApp.unregisterWindow(w)

		}()

		// Now do the actual close
		w.chromium.ShuttingDown()
		return w32.DefWindowProc(w.hwnd, w32.WM_CLOSE, 0, 0)

	case w32.WM_KILLFOCUS:
		if w.focusingChromium {
			return 0
		}
		w.parent.emit(events.Windows.WindowKillFocus)
	case w32.WM_ENTERSIZEMOVE:
		// This is needed to close open dropdowns when moving the window https://github.com/MicrosoftEdge/WebView2Feedback/issues/2290
		w32.SetFocus(w.hwnd)
		if int(w32.GetKeyState(w32.VK_LBUTTON))&(0x8000) != 0 {
			// Left mouse button is down - window is being moved
			w.parent.emit(events.Windows.WindowStartMove)
		} else {
			// Window is being resized
			w.parent.emit(events.Windows.WindowStartResize)
		}
	case w32.WM_EXITSIZEMOVE:
		if int(w32.GetKeyState(w32.VK_LBUTTON))&0x8000 != 0 {
			w.parent.emit(events.Windows.WindowEndMove)
		} else {
			w.parent.emit(events.Windows.WindowEndResize)
		}
	case w32.WM_SETFOCUS:
		w.focus()
		w.parent.emit(events.Windows.WindowSetFocus)
	case w32.WM_MOVE, w32.WM_MOVING:
		_ = w.chromium.NotifyParentWindowPositionChanged()
		w.moveDebouncer(func() {
			w.parent.emit(events.Windows.WindowDidMove)
		})
	case w32.WM_SHOWWINDOW:
		if wparam == 1 {
			w.parent.emit(events.Windows.WindowShow)
			w.updateContentProtection()
		} else {
			w.parent.emit(events.Windows.WindowHide)
		}
	case w32.WM_WINDOWPOSCHANGED:
		windowPos := (*w32.WINDOWPOS)(unsafe.Pointer(lparam))
		if windowPos.Flags&w32.SWP_NOZORDER == 0 {
			w.parent.emit(events.Windows.WindowZOrderChanged)
		}
	case w32.WM_PAINT:
		w.parent.emit(events.Windows.WindowPaint)
	case w32.WM_ERASEBKGND:
		w.parent.emit(events.Windows.WindowBackgroundErase)
		return 1 // Let WebView2 handle background erasing
	// WM_UAHDRAWMENUITEM is handled by MenuBarWndProc at the top of this function
	// Check for keypress
	case w32.WM_SYSCOMMAND:
		switch wparam {
		case w32.SC_KEYMENU:
			if lparam == 0 {
				// F10 or plain Alt key
				if w.processKeyBinding(w32.VK_F10) {
					return 0
				}
			} else {
				// Alt + key combination
				// The character code is in the low word of lparam
				char := byte(lparam & 0xFF)
				// Convert ASCII to virtual key code if needed
				vkey := w32.VkKeyScan(uint16(char))
				if w.processKeyBinding(uint(vkey)) {
					return 0
				}
			}
		}
	case w32.WM_SYSKEYDOWN:
		globalApplication.info("w32.WM_SYSKEYDOWN: %v", uint(wparam))
		w.parent.emit(events.Windows.WindowKeyDown)
		if w.processKeyBinding(uint(wparam)) {
			return 0
		}
	case w32.WM_SYSKEYUP:
		w.parent.emit(events.Windows.WindowKeyUp)
	case w32.WM_KEYDOWN:
		w.parent.emit(events.Windows.WindowKeyDown)
		w.processKeyBinding(uint(wparam))
	case w32.WM_KEYUP:
		w.parent.emit(events.Windows.WindowKeyUp)
	case w32.WM_SIZE:
		switch wparam {
		case w32.SIZE_MAXIMIZED:
			if w.isMinimizing {
				w.parent.emit(events.Windows.WindowUnMinimise)
			}
			w.isMinimizing = false
			w.parent.emit(events.Windows.WindowMaximise)
			// Force complete redraw when maximized
			if w.menu != nil && w.menubarTheme != nil {
				// Invalidate the entire window to force complete redraw
				w32.RedrawWindow(w.hwnd, nil, 0, w32.RDW_FRAME|w32.RDW_INVALIDATE|w32.RDW_UPDATENOW)
			}
		case w32.SIZE_RESTORED:
			if w.isMinimizing {
				w.parent.emit(events.Windows.WindowUnMinimise)
			}
			w.isMinimizing = false
			w.parent.emit(events.Windows.WindowRestore)
		case w32.SIZE_MINIMIZED:
			w.isMinimizing = true
			w.parent.emit(events.Windows.WindowMinimise)
		}

		doResize := func() {
			// Get the new size from lparam
			width := int32(lparam & 0xFFFF)
			height := int32((lparam >> 16) & 0xFFFF)
			bounds := &edge.Rect{
				Left:   0,
				Top:    0,
				Right:  width,
				Bottom: height,
			}
			InvokeSync(func() {
				time.Sleep(1 * time.Nanosecond)
				w.chromium.ResizeWithBounds(bounds)
				atomic.StoreInt32(&resizePending, 0)
				w.parent.emit(events.Windows.WindowDidResize)
			})
		}

		if w.parent.options.Frameless && wparam == w32.SIZE_MINIMIZED {
			// If the window is frameless, and we are minimizing, then we need to suppress the Resize on the
			// WebView2. If we don't do this, restoring does not work as expected and first restores with some wrong
			// size during the restore animation and only fully renders when the animation is done. This highly
			// depends on the content in the WebView, see https://github.com/MicrosoftEdge/WebView2Feedback/issues/2549
		} else if w.resizeDebouncer != nil {
			w.resizeDebouncer(doResize)
		} else {
			if atomic.CompareAndSwapInt32(&resizePending, 0, 1) {
				doResize()
			}
		}
		return 0

	case w32.WM_GETMINMAXINFO:
		mmi := (*w32.MINMAXINFO)(unsafe.Pointer(lparam))
		hasConstraints := false
		options := w.parent.options
		// Using ScreenManager to get the closest screen and scale according to its DPI is problematic
		// here because in multi-monitor setup, when dragging the window between monitors with the mouse
		// on the side with the higher DPI, the DPI change point is offset beyond the mid point, causing
		// wrong scaling and unwanted resizing when using the monitor DPI. To avoid this issue, we use
		// scaleWithWindowDPI() instead which retrieves the correct DPI with GetDpiForWindow().
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
		if !w.ignoreDPIChangeResizing {
			newWindowRect := (*w32.RECT)(unsafe.Pointer(lparam))
			w32.SetWindowPos(w.hwnd,
				uintptr(0),
				int(newWindowRect.Left),
				int(newWindowRect.Top),
				int(newWindowRect.Right-newWindowRect.Left),
				int(newWindowRect.Bottom-newWindowRect.Top),
				w32.SWP_NOZORDER|w32.SWP_NOACTIVATE)
		}
		w.parent.emit(events.Windows.WindowDPIChanged)
	}

	if w.parent.options.Windows.WindowMask != nil {
		switch msg {
		case w32.WM_NCHITTEST:
			if w.parent.options.Windows.WindowMaskDraggable {
				return w32.HTCAPTION
			}
			w.parent.emit(events.Windows.WindowNonClientHit)
			return w32.HTCLIENT
		case w32.WM_NCLBUTTONDOWN:
			w.parent.emit(events.Windows.WindowNonClientMouseDown)
		case w32.WM_NCLBUTTONUP:
			w.parent.emit(events.Windows.WindowNonClientMouseUp)
		case w32.WM_NCMOUSEMOVE:
			w.parent.emit(events.Windows.WindowNonClientMouseMove)
		case w32.WM_NCMOUSELEAVE:
			w.parent.emit(events.Windows.WindowNonClientMouseLeave)
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
					w.setPadding(edge.Rect{})
				} else if w.isMaximised() {
					// If the window is maximized we must adjust the client area to the work area of the monitor. Otherwise
					// some content goes beyond the visible part of the monitor.
					// Make sure to use the provided RECT to get the monitor, because during maximizig there might be
					// a wrong monitor returned in multiscreen mode when using MonitorFromWindow.
					// See: https://github.com/MicrosoftEdge/WebView2Feedback/issues/2549
					screen := ScreenNearestPhysicalRect(Rect{
						X:      int(rgrc.Left),
						Y:      int(rgrc.Top),
						Width:  int(rgrc.Right - rgrc.Left),
						Height: int(rgrc.Bottom - rgrc.Top),
					})

					rect := screen.PhysicalWorkArea

					maxWidth := options.MaxWidth
					maxHeight := options.MaxHeight

					if maxWidth > 0 {
						maxWidth = screen.scale(maxWidth, false)
						if rect.Width > maxWidth {
							rect.Width = maxWidth
						}
					}

					if maxHeight > 0 {
						maxHeight = screen.scale(maxHeight, false)
						if rect.Height > maxHeight {
							rect.Height = maxHeight
						}
					}

					*rgrc = w32.RECT{
						Left:   int32(rect.X),
						Top:    int32(rect.Y),
						Right:  int32(rect.X + rect.Width),
						Bottom: int32(rect.Y + rect.Height),
					}
					w.setPadding(edge.Rect{})
				} else {
					// This is needed to work around the resize flickering in frameless mode with WindowDecorations
					// See: https://stackoverflow.com/a/6558508
					// The workaround from the SO answer suggests to reduce the bottom of the window by 1px.
					// However, this would result in losing 1px of the WebView content.
					// Increasing the bottom also worksaround the flickering, but we would lose 1px of the WebView content
					// therefore let's pad the content with 1px at the bottom.
					rgrc.Bottom += 1
					w.setPadding(edge.Rect{Bottom: 1})
				}
				return 0
			}
		}
	}
	return w32.DefWindowProc(w.hwnd, msg, wparam, lparam)
}

func (w *windowsWebviewWindow) DPI() (w32.UINT, w32.UINT) {
	if w32.HasGetDpiForWindowFunc() {
		// GetDpiForWindow is supported beginning with Windows 10, 1607 and is the most accurate
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
		globalApplication.fatal("fatal error in callback setWindowMask: %w", err)
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
func (w *windowsWebviewWindow) processMessage(message string, sender *edge.ICoreWebView2, args *edge.ICoreWebView2WebMessageReceivedEventArgs) {
	// We send all messages to the centralised window message buffer
	windowMessageBuffer <- &windowMessage{
		windowId: w.parent.id,
		message:  message,
	}
}

func (w *windowsWebviewWindow) processRequest(
	req *edge.ICoreWebView2WebResourceRequest,
	args *edge.ICoreWebView2WebResourceRequestedEventArgs,
) {

	// Setting the UserAgent on the CoreWebView2Settings clears the whole default UserAgent of the Edge browser, but
	// we want to just append our ApplicationIdentifier. So we adjust the UserAgent for every request.
	if reqHeaders, err := req.GetHeaders(); err == nil {
		useragent, _ := reqHeaders.GetHeader(assetserver.HeaderUserAgent)
		useragent = strings.Join([]string{useragent, assetserver.WailsUserAgentValue}, " ")
		err = reqHeaders.SetHeader(assetserver.HeaderUserAgent, useragent)
		if err != nil {
			globalApplication.fatal("error setting UserAgent header: %w", err)
		}
		err = reqHeaders.SetHeader(
			webViewRequestHeaderWindowId,
			strconv.FormatUint(uint64(w.parent.id), 10),
		)
		if err != nil {
			globalApplication.fatal("error setting WindowId header: %w", err)
		}
		err = reqHeaders.Release()
		if err != nil {
			globalApplication.fatal("error releasing headers: %w", err)
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
		globalApplication.error("unable to parse request uri: uri='%s' error='%w'", uri, err)
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
		globalApplication.error("%s: NewRequest failed: %w", uri, err)
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

	webview2version, err := webviewloader.GetAvailableCoreWebView2BrowserVersionString(
		globalApplication.options.Windows.WebviewBrowserPath,
	)
	if err != nil {
		globalApplication.error("error getting WebView2 version: %w", err)
		return
	}
	globalApplication.capabilities = capabilities.NewCapabilities(webview2version)

	// We disable this by default. Can be overridden with the `EnableFraudulentWebsiteWarnings` option
	opts.DisabledFeatures = append(opts.DisabledFeatures, "msSmartScreenProtection")

	if len(opts.DisabledFeatures) > 0 {
		opts.DisabledFeatures = lo.Uniq(opts.DisabledFeatures)
		arg := fmt.Sprintf("--disable-features=%s", strings.Join(opts.DisabledFeatures, ","))
		chromium.AdditionalBrowserArgs = append(chromium.AdditionalBrowserArgs, arg)
	}

	if len(opts.EnabledFeatures) > 0 {
		opts.EnabledFeatures = lo.Uniq(opts.EnabledFeatures)
		arg := fmt.Sprintf("--enable-features=%s", strings.Join(opts.EnabledFeatures, ","))
		chromium.AdditionalBrowserArgs = append(chromium.AdditionalBrowserArgs, arg)
	}

	if len(opts.AdditionalLaunchArgs) > 0 {
		chromium.AdditionalBrowserArgs = append(chromium.AdditionalBrowserArgs, opts.AdditionalLaunchArgs...)
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

	// Prevent efficiency mode by keeping WebView2 visible (fixes issue #2861)
	// Microsoft recommendation: keep IsVisible = true to avoid efficiency mode
	// See: https://github.com/MicrosoftEdge/WebView2Feedback/discussions/4021
	// TODO: Re-enable when PutIsVisible method is available in go-webview2 package
	// err := chromium.PutIsVisible(true)
	// if err != nil {
	//	globalApplication.error("Failed to set WebView2 visibility for efficiency mode prevention: %v", err)
	// }

	if chromium.HasCapability(edge.SwipeNavigation) {
		err := chromium.PutIsSwipeNavigationEnabled(opts.EnableSwipeGestures)
		if err != nil {
			globalApplication.handleFatalError(err)
		}
	}

	if w.parent.options.EnableDragAndDrop {
		if chromium.HasCapability(edge.AllowExternalDrop) {
			err := chromium.AllowExternalDrag(false)
			if err != nil {
				globalApplication.handleFatalError(err)
			}
		}
		w.dropTarget = w32.NewDropTarget()
		w.dropTarget.OnDrop = func(files []string, x int, y int) {
			w.parent.emit(events.Windows.WindowDragDrop)
			globalApplication.debug("[DragDropDebug] Windows DropTarget OnDrop: Raw screen coordinates", "x", x, "y", y)

			// Convert screen coordinates to window-relative coordinates first
			// Windows DropTarget gives us screen coordinates, but we need window-relative coordinates
			windowRect := w32.GetWindowRect(w.hwnd)
			windowRelativeX := x - int(windowRect.Left)
			windowRelativeY := y - int(windowRect.Top)

			globalApplication.debug("[DragDropDebug] Windows DropTarget OnDrop: After screen-to-window conversion", "windowRelativeX", windowRelativeX, "windowRelativeY", windowRelativeY)

			// Convert window-relative coordinates to webview-relative coordinates
			webviewX, webviewY := w.convertWindowToWebviewCoordinates(windowRelativeX, windowRelativeY)
			globalApplication.debug("[DragDropDebug] Windows DropTarget OnDrop: Final webview coordinates", "webviewX", webviewX, "webviewY", webviewY)
			w.parent.InitiateFrontendDropProcessing(files, webviewX, webviewY)
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
			if err != nil && !errors.Is(err, syscall.Errno(w32.DRAGDROP_E_ALREADYREGISTERED)) {
				globalApplication.error("error registering drag and drop: %w", err)
			}
			//}
			return 1
		})

	}

	err = chromium.PutIsGeneralAutofillEnabled(opts.GeneralAutofillEnabled)
	if err != nil {
		if errors.Is(err, edge.UnsupportedCapabilityError) {
			globalApplication.warning("unsupported capability: GeneralAutofillEnabled")
		} else {
			globalApplication.handleFatalError(err)
		}
	}

	err = chromium.PutIsPasswordAutosaveEnabled(opts.PasswordAutosaveEnabled)
	if err != nil {
		if errors.Is(err, edge.UnsupportedCapabilityError) {
			globalApplication.warning("unsupported capability: PasswordAutosaveEnabled")
		} else {
			globalApplication.handleFatalError(err)
		}
	}

	chromium.Resize()
	settings, err := chromium.GetSettings()
	if err != nil {
		globalApplication.handleFatalError(err)
	}
	if settings == nil {
		globalApplication.fatal("error getting settings")
	}
	err = settings.PutAreDefaultContextMenusEnabled(
		debugMode || !w.parent.options.DefaultContextMenuDisabled,
	)
	if err != nil {
		globalApplication.handleFatalError(err)
	}

	w.enableDevTools(settings)

	if w.parent.options.Zoom > 0.0 {
		chromium.PutZoomFactor(w.parent.options.Zoom)
	}
	err = settings.PutIsZoomControlEnabled(w.parent.options.ZoomControlEnabled)
	if err != nil {
		globalApplication.handleFatalError(err)
	}

	err = settings.PutIsStatusBarEnabled(false)
	if err != nil {
		globalApplication.handleFatalError(err)
	}
	err = settings.PutAreBrowserAcceleratorKeysEnabled(false)
	if err != nil {
		globalApplication.handleFatalError(err)
	}
	err = settings.PutIsSwipeNavigationEnabled(false)
	if err != nil {
		globalApplication.handleFatalError(err)
	}

	if debugMode && w.parent.options.OpenInspectorOnStartup {
		chromium.OpenDevToolsWindow()
	}

	// Set background colour
	w.setBackgroundColour(w.parent.options.BackgroundColour)
	chromium.SetBackgroundColour(
		w.parent.options.BackgroundColour.Red,
		w.parent.options.BackgroundColour.Green,
		w.parent.options.BackgroundColour.Blue,
		w.parent.options.BackgroundColour.Alpha,
	)

	chromium.SetGlobalPermission(edge.CoreWebView2PermissionStateAllow)
	chromium.AddWebResourceRequestedFilter("*", edge.COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)

	if w.parent.options.HTML != "" {
		var script string
		if w.parent.options.JS != "" {
			script = w.parent.options.JS
		}
		if w.parent.options.CSS != "" {
			script += fmt.Sprintf(
				"; addEventListener(\"DOMContentLoaded\", (event) => { document.head.appendChild(document.createElement('style')).innerHTML=\"%s\"; });",
				strings.ReplaceAll(w.parent.options.CSS, `"`, `\"`),
			)
		}
		if script != "" {
			chromium.Init(script)
		}
		chromium.NavigateToString(w.parent.options.HTML)
	} else {
		startURL, err := assetserver.GetStartURL(w.parent.options.URL)
		if err != nil {
			globalApplication.handleFatalError(err)
		}
		w.webviewNavigationCompleted = false
		chromium.Navigate(startURL)
	}

}

func (w *windowsWebviewWindow) fullscreenChanged(
	sender *edge.ICoreWebView2,
	_ *edge.ICoreWebView2ContainsFullScreenElementChangedEventArgs,
) {
	isFullscreen, err := sender.GetContainsFullScreenElement()
	if err != nil {
		globalApplication.fatal("fatal error in callback fullscreenChanged: %w", err)
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

func (w *windowsWebviewWindow) navigationCompleted(
	sender *edge.ICoreWebView2,
	args *edge.ICoreWebView2NavigationCompletedEventArgs,
) {

	// Install the runtime core
	w.execJS(runtime.Core())

	// EmitEvent DomReady ApplicationEvent
	windowEvents <- &windowEvent{EventID: uint(events.Windows.WebViewNavigationCompleted), WindowID: w.parent.id}

	if w.webviewNavigationCompleted {
		// NavigationCompleted is triggered for every Load. If an application uses reloads the Hide/Show will trigger
		// a flickering of the window with every reload. So we only do this once for the first NavigationCompleted.
		return
	}
	w.webviewNavigationCompleted = true

	// Cancel any pending visibility timeout since navigation completed
	if w.visibilityTimeout != nil {
		w.visibilityTimeout.Stop()
		w.visibilityTimeout = nil
	}

	wasFocused := w.isFocused()
	// Hack to make it visible: https://github.com/MicrosoftEdge/WebView2Feedback/issues/1077#issuecomment-825375026
	err := w.chromium.Hide()
	if err != nil {
		globalApplication.handleFatalError(err)
	}
	err = w.chromium.Show()
	if err != nil {
		globalApplication.handleFatalError(err)
	}
	if wasFocused {
		w.focus()
	}

	// Only call parent.Show() if not hidden and show was requested but window wasn't shown yet
	// The new robust show() method handles window visibility independently
	if !w.parent.options.Hidden {
		if w.showRequested && !w.windowShown {
			w.parent.Show()
		}
		w.update()
	}
}

func (w *windowsWebviewWindow) processKeyBinding(vkey uint) bool {

	globalApplication.debug("Processing key binding", "vkey", vkey)

	// Get the keyboard state and convert to an accelerator
	var keyState [256]byte
	if !w32.GetKeyboardState(keyState[:]) {
		globalApplication.error("error getting keyboard state")
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

	if vkey != w32.VK_CONTROL && vkey != w32.VK_MENU && vkey != w32.VK_SHIFT &&
		vkey != w32.VK_LWIN &&
		vkey != w32.VK_RWIN {
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

func (w *windowsWebviewWindow) processMessageWithAdditionalObjects(
	message string,
	sender *edge.ICoreWebView2,
	args *edge.ICoreWebView2WebMessageReceivedEventArgs,
) {
	if strings.HasPrefix(message, "FilesDropped") {
		objs, err := args.GetAdditionalObjects()
		if err != nil {
			globalApplication.handleError(err)
			return
		}

		defer func() {
			err = objs.Release()
			if err != nil {
				globalApplication.error("error releasing objects: %w", err)
			}
		}()

		count, err := objs.GetCount()
		if err != nil {
			globalApplication.error("cannot get count: %w", err)
			return
		}

		var filenames []string
		for i := uint32(0); i < count; i++ {
			_file, err := objs.GetValueAtIndex(i)
			if err != nil {
				globalApplication.error("cannot get value at %d: %w", i, err)
				return
			}

			file := (*edge.ICoreWebView2File)(unsafe.Pointer(_file))

			// TODO: Fix this
			defer file.Release()

			filepath, err := file.GetPath()
			if err != nil {
				globalApplication.error("cannot get path for object at %d: %w", i, err)
				return
			}

			filenames = append(filenames, filepath)
		}

		// Extract X/Y coordinates from message - format should be "FilesDropped:x:y"
		var x, y int
		parts := strings.Split(message, ":")
		if len(parts) >= 3 {
			if parsedX, err := strconv.Atoi(parts[1]); err == nil {
				x = parsedX
			}
			if parsedY, err := strconv.Atoi(parts[2]); err == nil {
				y = parsedY
			}
		}

		globalApplication.debug("[DragDropDebug] processMessageWithAdditionalObjects: Raw WebView2 coordinates", "x", x, "y", y)

		// Convert webview-relative coordinates to window-relative coordinates, then to webview-relative coordinates
		// Note: The coordinates from WebView2 are already webview-relative, but let's log them for debugging
		webviewX, webviewY := x, y

		globalApplication.debug("[DragDropDebug] processMessageWithAdditionalObjects: Using coordinates as-is (already webview-relative)", "webviewX", webviewX, "webviewY", webviewY)

		w.parent.InitiateFrontendDropProcessing(filenames, webviewX, webviewY)
		return
	}
}

func (w *windowsWebviewWindow) setMaximiseButtonEnabled(enabled bool) {
	w.setStyle(enabled, w32.WS_MAXIMIZEBOX)
}

func (w *windowsWebviewWindow) setMinimiseButtonEnabled(enabled bool) {
	w.setStyle(enabled, w32.WS_MINIMIZEBOX)
}

func (w *windowsWebviewWindow) toggleMenuBar() {
	if w.menu != nil {
		if w32.GetMenu(w.hwnd) == 0 {
			w32.SetMenu(w.hwnd, w.menu.menu)
		} else {
			w32.SetMenu(w.hwnd, 0)
		}

		// Get the bounds of the client area
		//bounds := w32.GetClientRect(w.hwnd)

		// Resize the webview
		w.chromium.Resize()

		// Update size of webview
		w.update()
		// Restore focus to the webview after toggling menu
		w.focus()
	}
}

func (w *windowsWebviewWindow) enableRedraw() {
	w32.SendMessage(w.hwnd, w32.WM_SETREDRAW, 1, 0)
	w32.RedrawWindow(
		w.hwnd,
		nil,
		0,
		w32.RDW_ERASE|w32.RDW_FRAME|w32.RDW_INVALIDATE|w32.RDW_ALLCHILDREN,
	)
}

func (w *windowsWebviewWindow) disableRedraw() {
	w32.SendMessage(w.hwnd, w32.WM_SETREDRAW, 0, 0)
}

func (w *windowsWebviewWindow) disableRedrawWithCallback(callback func()) {
	w.disableRedraw()
	callback()
	w.enableRedraw()

}

func NewIconFromResource(instance w32.HINSTANCE, resId uint16) (w32.HICON, error) {
	var err error
	var result w32.HICON
	if result = w32.LoadIconWithResourceID(instance, resId); result == 0 {
		err = fmt.Errorf("cannot load icon from resource with id %v", resId)
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

func (w *windowsWebviewWindow) setPadding(padding edge.Rect) {
	// Skip SetPadding if window is being minimized to prevent flickering
	if w.isMinimizing {
		return
	}
	w.chromium.SetPadding(padding)
}

func (w *windowsWebviewWindow) showMenuBar() {
	if w.menu != nil {
		w32.SetMenu(w.hwnd, w.menu.menu)
	}
}

func (w *windowsWebviewWindow) hideMenuBar() {
	if w.menu != nil {
		w32.SetMenu(w.hwnd, 0)
	}
}

func (w *windowsWebviewWindow) snapAssist() {
	// Simulate Win+Z key combination to trigger Snap Assist
	// Press Windows key
	w32.KeybdEvent(byte(w32.VK_LWIN), 0, 0, 0)
	// Press Z key
	w32.KeybdEvent(byte('Z'), 0, 0, 0)
	// Release Z key
	w32.KeybdEvent(byte('Z'), 0, w32.KEYEVENTF_KEYUP, 0)
	// Release Windows key
	w32.KeybdEvent(byte(w32.VK_LWIN), 0, w32.KEYEVENTF_KEYUP, 0)
}

func (w *windowsWebviewWindow) setContentProtection(enabled bool) {
	// Ensure the option reflects the requested state for future show() calls
	w.parent.options.ContentProtectionEnabled = enabled
	w.updateContentProtection()
}

func (w *windowsWebviewWindow) updateContentProtection() {
	if w.hwnd == 0 {
		return
	}

	if !w.isVisible() {
		// Defer updates until the window is visible to avoid affinity glitches.
		return
	}

	desired := w.parent.options.ContentProtectionEnabled

	if desired {
		if w.applyDisplayAffinity(w32.WDA_EXCLUDEFROMCAPTURE) {
			w.contentProtectionApplied = true
		}
		return
	}

	if w.applyDisplayAffinity(w32.WDA_NONE) {
		w.contentProtectionApplied = false
	}
}

func (w *windowsWebviewWindow) applyDisplayAffinity(affinity uint32) bool {
	if ok := w32.SetWindowDisplayAffinity(w.hwnd, affinity); !ok {
		// Note: wrapper already falls back to WDA_MONITOR on older Windows.
		globalApplication.warning("SetWindowDisplayAffinity failed: window=%v, affinity=%v", w.parent.id, affinity)
		return false
	}
	return true
}
