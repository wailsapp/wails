//go:build windows

package application

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/assetserver"
	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/internal/capabilities"
	"github.com/wailsapp/wails/v3/internal/debounce"
	"github.com/wailsapp/wails/v3/internal/runtime"
	"github.com/wailsapp/wails/v3/internal/sliceutil"
	"github.com/wailsapp/wails/webview2/webviewloader"

	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/w32"
	"github.com/wailsapp/wails/webview2/pkg/edge"
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
	onceDo             sync.Once

	// Window move debouncer
	moveDebouncer func(func())

	// isMinimizing indicates whether the window is currently being minimized
	// Used to prevent unnecessary redraws during minimize/restore operations
	isMinimizing bool

	// lastSizeWParam is the wParam from the most-recent WM_SIZE message.
	// Gate the dark-menubar force-repaint on a state transition so it does not fire
	// on every SIZE_RESTORED during live drag-resize. WM_ENTERSIZEMOVE/WM_EXITSIZEMOVE
	// cannot be used for this because keyboard snap (Win+Left) bypasses those messages.
	lastSizeWParam uintptr

	// lastKnownDPI is the window's DPI the last time it was in a non-minimised
	// state. It is used on the un-minimise path to decide whether the WebView2
	// rasterization scale actually needs resyncing: when the DPI is unchanged
	// (the common case — same monitor, fixed DPI) we must avoid making any COM
	// call into the controller, which can be fatal if WebView2 suspended or its
	// render/GPU process died while minimised (#5605). It is read with the Win32
	// GetDpiForWindow, never via COM, and only updated while not minimised so it
	// never captures the parked minimised position's DPI (while minimised the
	// window is repositioned off the monitor it will restore to).
	lastKnownDPI w32.UINT

	// menubarTheme is the theme for the menubar
	menubarTheme *w32.MenuBarTheme

	// Modal window tracking
	parentHWND w32.HWND // Parent window HWND when this window is a modal

	// WebView2 process-failure recovery state (webview_recovery_windows.go).
	// All fields are touched on the main thread only, so no lock is needed.
	// processFailedLogAt throttles the ProcessFailed stack trace per kind
	// (RENDER_PROCESS_UNRESPONSIVE re-fires every few seconds while hung; a
	// new kind is always logged immediately). webviewRebuildTimes rate-caps
	// controller rebuilds; webviewHealthProbeFailures counts consecutive
	// watchdog probe failures.
	processFailedLogAt         map[string]time.Time
	lastRebuildSuppressedLogAt time.Time
	webviewRebuildInProgress   bool
	webviewRebuildTimes        []time.Time
	webviewHealthProbeFailures int

	// DPI-flap breaker state (noteDPITransitionAndDetectFlap, #5701).
	// Main-thread only (WM_DPICHANGED handling). lastAppliedDPI is the DPI of
	// the last WM_DPICHANGED actually processed (dedup key — unlike
	// lastKnownDPI it is NOT refreshed by WM_SIZE, so it cannot race ahead of
	// the transition stream). inSizeMove tracks the modal move/size loop
	// (WM_ENTERSIZEMOVE..WM_EXITSIZEMOVE); the flap resolver must not
	// reposition a window the user is actively dragging.
	lastDPITransitionAt   time.Time
	lastDPITransitionFrom uint32
	dpiFlapReversals      int
	dpiFlapSuppressUntil  time.Time
	dpiFlapStormStartAt   time.Time
	lastDPIFlapSettledAt  time.Time
	dpiFlapResumeCount    int
	lastAppliedDPI        uint32
	inSizeMove            bool
	lastStraddleResolveAt time.Time
	// lastScalePutAt stamps every app-initiated rasterization-scale put (the
	// verify ladder's corrective puts and, in app-owner mode, the per-flip
	// puts). decideScaleReconcile rate-limits corrective puts against it —
	// put churn faster than 1/s means something else is rewriting the scale,
	// and fighting it per-event is the field-fatal pattern (#5701). Main
	// thread only.
	lastScalePutAt time.Time
	// lastProgrammaticPlacementAt is stamped by setPhysicalBounds (every
	// app-driven SetBounds/SetPosition/SetSize funnels through it). The
	// parked-reversal fast path must not trip inside dpiPlacementGrace of it:
	// a multi-monitor placement legitimately produces back-to-back DPI
	// reversals, and the resolver "resolving" one moves the window the app
	// just placed — with a mid-placement rect, onto the wrong monitor
	// (v200.0.22 alarm field trace, #5701).
	lastProgrammaticPlacementAt time.Time
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

func (w *windowsWebviewWindow) attachModal(modalWindow *WebviewWindow) {
	if modalWindow == nil || modalWindow.impl == nil || modalWindow.isDestroyed() {
		return
	}

	// Get the modal window's Windows implementation
	modalWindowsImpl, ok := modalWindow.impl.(*windowsWebviewWindow)
	if !ok {
		return
	}

	parentHWND := w.hwnd
	modalHWND := modalWindowsImpl.hwnd

	// Set parent-child relationship using GWLP_HWNDPARENT
	// This ensures the modal stays above parent and moves with it
	w32.SetWindowLongPtr(modalHWND, w32.GWLP_HWNDPARENT, uintptr(parentHWND))

	// Track the parent HWND in the modal window for cleanup
	modalWindowsImpl.parentHWND = parentHWND

	// Disable the parent window to block interaction (Microsoft's recommended approach)
	// This follows Windows modal dialog best practices
	w32.EnableWindow(parentHWND, false)

	// Ensure modal window is shown and brought to front
	w32.ShowWindow(modalHWND, w32.SW_SHOW)
	w32.SetForegroundWindow(modalHWND)
	w32.BringWindowToTop(modalHWND)
}

func (w *windowsWebviewWindow) nativeWindow() unsafe.Pointer {
	return unsafe.Pointer(w.hwnd)
}

func (w *windowsWebviewWindow) setTitle(title string) {
	w32.SetWindowText(w.hwnd, title)
}

func (w *windowsWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	var hwndInsertAfter uintptr
	if alwaysOnTop {
		hwndInsertAfter = w32.HWND_TOPMOST
	} else {
		hwndInsertAfter = w32.HWND_NOTOPMOST
	}
	w32.SetWindowPos(w.hwnd,
		hwndInsertAfter,
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
		if w.parent.options.IgnoreMouseEvents {
			// WS_EX_TRANSPARENT makes WM_NCHITTEST return HTTRANSPARENT, so the window
			// passes mouse input through to whatever is behind it. Only apply it when
			// the caller explicitly opts in via IgnoreMouseEvents — applying it to every
			// frameless + transparent window causes clicks to fall through to the desktop
			// in any area the child WebView2 HWND does not currently cover (issue #4871).
			exStyle |= w32.WS_EX_TRANSPARENT | w32.WS_EX_LAYERED
		} else {
			// Transparent/translucent composition via DirectComposition.
			exStyle |= w32.WS_EX_NOREDIRECTIONBITMAP
		}
	}
	if options.AlwaysOnTop {
		exStyle |= w32.WS_EX_TOPMOST
	}
	// WS_EX_TOOLWINDOW hides the window from the taskbar without blocking keyboard focus.
	// WS_EX_NOACTIVATE (previously used here) prevents the window from being activated,
	// which blocks keyboard focus and causes click-through issues after Win+D or Alt+Tab.
	if options.Windows.HiddenOnTaskbar {
		exStyle |= w32.WS_EX_TOOLWINDOW
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
			// Explicit window menu takes priority
			userMenu.Update()
			w.menu = NewApplicationMenu(w, userMenu)
			w.menu.parentWindow = w
			appMenu = w.menu.menu
		} else if options.UseApplicationMenu && globalApplication.applicationMenu != nil {
			// Use the global application menu if opted in
			globalApplication.applicationMenu.Update()
			w.menu = NewApplicationMenu(w, globalApplication.applicationMenu)
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

	// Seed the last-known DPI so the first un-minimise has a baseline to compare
	// against (it is otherwise refreshed on every non-minimised WM_SIZE). See
	// resyncWebviewDPIAfterUnminimiseIfDPIChanged (#5605). lastAppliedDPI seeds
	// the WM_DPICHANGED dedup (#5701) with the starting DPI.
	if dpi, _ := w.DPI(); dpi != 0 {
		w.lastKnownDPI = dpi
		w.lastAppliedDPI = uint32(dpi)
	}

	// Ensure correct window size in case the scale factor of current screen is different from the initial one.
	// This could happen when using the default window position and the window launches on a secondary monitor.
	currentScreen, _ := w.getScreen()
	if currentScreen.ScaleFactor != initialScreen.ScaleFactor {
		w.setSize(options.Width, options.Height)
	}

	w.setupChromium()

	// Backstop for WebView2 wedges whose ProcessFailed event is missed or
	// never fires: a low-frequency controller health probe that rebuilds the
	// controller after consecutive failures (#5701).
	w.startWebviewHealthWatchdog()

	if options.Windows.WindowDidMoveDebounceMS == 0 {
		options.Windows.WindowDidMoveDebounceMS = 50
	}
	w.moveDebouncer = debounce.New(
		time.Duration(options.Windows.WindowDidMoveDebounceMS) * time.Millisecond,
	)

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

	if options.Screen != nil {
		if options.InitialPosition == WindowCentered {
			w.centerOnScreen(options.Screen)
		} else {
			workArea := options.Screen.WorkArea
			w.setPosition(workArea.X+options.X, workArea.Y+options.Y)
		}
	} else if options.InitialPosition == WindowCentered {
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
	webviewPhysicalX := windowX - offsetX
	webviewPhysicalY := windowY - offsetY

	globalApplication.debug("[DragDropDebug] convertWindowToWebviewCoordinates: Webview coordinates before DPI Scaling", "webviewPhysicalX", webviewPhysicalX, "webviewPhysicalY", webviewPhysicalY)

	// Get DPI for this window
	dpi := w32.GetDpiForWindow(w.hwnd)
	if dpi == 0 {
		globalApplication.debug("[DragDropDebug] convertWindowToWebviewCoordinates: Failed to get dpi, returning physical coordinates", "webviewPhysicalX", webviewPhysicalX, "webviewPhysicalY", webviewPhysicalY)
		return webviewPhysicalX, webviewPhysicalY
	}

	// Convert to scale factor: 96 DPI == 1.0 (100%)
	scaleFactor := float64(dpi) / 96.0
	globalApplication.debug("[DragDropDebug] convertWindowToWebviewCoordinates: DPI info", "dpi", dpi, "scaleFactor", scaleFactor)

	// Convert physical pixels -> logical/CSS pixels by dividing by the scale factor
	// Use rounding to avoid truncation artefacts
	webviewLogicalX := int(math.Round(float64(webviewPhysicalX) / scaleFactor))
	webviewLogicalY := int(math.Round(float64(webviewPhysicalY) / scaleFactor))
	globalApplication.debug("[DragDropDebug] convertWindowToWebviewCoordinates: Final webview coordinates (logical/CSS pixels)",
		"webviewLogicalX", webviewLogicalX, "webviewLogicalY", webviewLogicalY)

	return webviewLogicalX, webviewLogicalY
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
	w.lastProgrammaticPlacementAt = time.Now()
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

	// For WS_EX_LAYERED windows (frameless+transparent or IgnoreMouseEvents), the hit-test
	// region is not updated by SetWindowPos alone. Calling SetLayeredWindowAttributes refreshes
	// the layered region so that the full new window area responds to mouse events.
	if exStyle := w32.GetWindowLong(w.hwnd, w32.GWL_EXSTYLE); exStyle&w32.WS_EX_LAYERED != 0 {
		w32.SetLayeredWindowAttributes(w.hwnd, 0, 255, w32.LWA_ALPHA)
	}
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

func (w *windowsWebviewWindow) centerOnScreen(screen *Screen) {
	workArea := screen.WorkArea
	width, height := w.size()
	x := workArea.X + (workArea.Width-width)/2
	y := workArea.Y + (workArea.Height-height)/2
	w.setPosition(x, y)
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
	// Re-enable parent window if this was a modal window
	if w.parentHWND != 0 {
		w32.EnableWindow(w.parentHWND, true)
		w.parentHWND = 0
	}

	w.parent.markAsDestroyed()
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
	if zoom < 1.0 {
		zoom = 1.0
	}
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
	if w.chromium.GetController() != nil {
		w.chromium.Focus()
	}
}

func (w *windowsWebviewWindow) unmaximise() {
	w.restore()
	w.parent.emit(events.Windows.WindowUnMaximise)
}

func (w *windowsWebviewWindow) restore() {
	w32.ShowWindow(w.hwnd, w32.SW_RESTORE)
	if w.chromium.GetController() != nil {
		w.chromium.Focus()
	}
	w.enforceMinSizeConstraints()
}

func (w *windowsWebviewWindow) enforceMinSizeConstraints() {
	options := w.parent.options
	if options.MinWidth <= 0 && options.MinHeight <= 0 {
		return
	}
	b := w.bounds()
	changed := false
	if options.MinWidth > 0 && b.Width < options.MinWidth {
		b.Width = options.MinWidth
		changed = true
	}
	if options.MinHeight > 0 && b.Height < options.MinHeight {
		b.Height = options.MinHeight
		changed = true
	}
	if changed {
		w.setBounds(b)
	}
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
		w.previousWindowExStyle & ^uint32(w32.WS_EX_DLGMODALFRAME|w32.WS_EX_TRANSPARENT),
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

	// The SetWindowPos above re-puts controller bounds via its synchronous
	// WM_SIZE, but that put silently no-ops when the controller is not ready
	// yet (a freshly created alarm window fullscreened while Embed is still
	// pumping), and a Window-to-Visual internal visual resize can lag the put.
	// Field report (#5701, v200.0.21): all-monitors fullscreen alarm content
	// parked bottom-right with correct text size on the secondary monitor —
	// the signature of content centered on a stale, larger canvas. Re-assert
	// the bounds once, and leave a breadcrumb capturing which value (bounds
	// vs scale) was stale if it ever recurs.
	if controller := w.chromium.GetController(); controller != nil {
		var bw, bh int32
		if b, err := controller.GetBounds(); err == nil && b != nil {
			bw, bh = b.Right-b.Left, b.Bottom-b.Top
		}
		var cw, ch int32
		if rect := w32.GetClientRect(w.hwnd); rect != nil {
			cw, ch = rect.Right-rect.Left, rect.Bottom-rect.Top
		}
		globalApplication.warning("fullscreen: window %d client %dx%d px, controller bounds %dx%d px, raster=%.2f — re-asserting bounds (#5701)",
			w.parent.id, cw, ch, bw, bh, w.currentWebviewRasterizationScale())
		w.chromium.Resize()
		w.chromium.Focus()
	}
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

	// Guard against calling Focus when the WebView2 controller is not yet
	// initialized or has already been torn down (e.g. during dev hot-reload).
	// go-webview2's Focus() calls os.Exit(1) on any MoveFocus error, so we
	// must not call it when the controller is in a nil/invalid state.
	if w.chromium.GetController() == nil {
		return
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

	// Show WebView if navigation has completed. IsReady also guards the
	// recovery window: during a controller rebuild (#5701) a queued show()
	// can be dispatched by the rebuild's nested Embed message pump while the
	// fresh chromium has no controller yet, and chromium.Show() dereferences
	// the controller unguarded — the rebuild makes the webview visible itself.
	if w.webviewNavigationCompleted {
		if w.chromium.IsReady() {
			w.chromium.Show()
		}
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
				if !w.webviewNavigationCompleted && w.chromium != nil && w.chromium.IsReady() {
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
	if w.isFullscreen() {
		// Avoid disrupting fullscreen (WS_POPUP) styling; frame trimming will be handled on exit.
		return
	}
	// Keep the full overlapped-window style in both states and let the
	// WM_NCCALCSIZE handler — keyed on options.Frameless, which the caller
	// has already updated — trim the frame, exactly like Frameless: true at
	// window creation. The previous implementation switched frameless
	// windows to a bare WS_POPUP style, which silently loses the styles DWM
	// animations are keyed on: minimise/restore/maximise transitions, Aero
	// snap and the resize borders all stopped working after
	// SetFrameless(true), unlike creation-time frameless windows (#5541).
	//
	// Preserve the live state/button bits rather than overwriting GWL_STYLE
	// wholesale: a blanket WS_VISIBLE|WS_OVERLAPPEDWINDOW would clear the
	// WS_MAXIMIZE/WS_MINIMIZE state bits and re-enable the minimise/maximise/
	// close buttons even when they were disabled via options or the runtime
	// setMinimiseButtonState/setMaximiseButtonState/setCloseButtonState calls.
	const preserve = w32.WS_VISIBLE | w32.WS_DISABLED | w32.WS_MAXIMIZE | w32.WS_MINIMIZE |
		w32.WS_MINIMIZEBOX | w32.WS_MAXIMIZEBOX | w32.WS_SYSMENU | w32.WS_THICKFRAME
	current := uint(w32.GetWindowLongPtr(w.hwnd, w32.GWL_STYLE))
	style := (uint(w32.WS_OVERLAPPEDWINDOW) &^ preserve) | (current & preserve)
	w32.SetWindowLongPtr(w.hwnd, w32.GWL_STYLE, uintptr(style))
	// Inform the application of the frame change to trigger WM_NCCALCSIZE.
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
	// Destroy previous context menu if it exists to prevent memory leak
	if w.currentlyOpenContextMenu != nil {
		w.currentlyOpenContextMenu.Destroy()
	}
	// Create the menu from current Go-side menu state
	thisMenu := NewPopupMenu(w.hwnd, menu)
	thisMenu.Update()
	w.currentlyOpenContextMenu = thisMenu
	thisMenu.ShowAtCursor()
}

func (w *windowsWebviewWindow) setStyle(b bool, style int) {
	currentStyle := int(w32.GetWindowLongPtr(w.hwnd, w32.GWL_STYLE))
	if currentStyle != 0 {
		if b {
			currentStyle = currentStyle | style
		} else {
			currentStyle = currentStyle &^ style
		}
		w32.SetWindowLongPtr(w.hwnd, w32.GWL_STYLE, uintptr(currentStyle))
	}
}
func (w *windowsWebviewWindow) setExStyle(b bool, style int) {
	currentStyle := int(w32.GetWindowLongPtr(w.hwnd, w32.GWL_EXSTYLE))
	if currentStyle != 0 {
		if b {
			currentStyle = currentStyle | style
		} else {
			currentStyle = currentStyle &^ style
		}
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

func (w *windowsWebviewWindow) WndProc(msg uint32, wparam, lparam uintptr) uintptr {

	// Use the original implementation that works perfectly for maximized
	processed, code := w32.MenuBarWndProc(w.hwnd, msg, wparam, lparam, w.menubarTheme)
	if processed {
		return code
	}

	if msg == w32.WM_NCHITTEST && w.isCurrentlyFullscreen {
		return w32.HTCLIENT
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
			// Re-enable parent window if this was a modal window
			if w.parentHWND != 0 {
				w32.EnableWindow(w.parentHWND, true)
				w.parentHWND = 0
			}

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
		w.inSizeMove = true
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
		w.inSizeMove = false
		if !w.dpiFlapSuppressUntil.IsZero() {
			// A drag ended while a DPI-flap storm was open. Once released, a
			// window off the boundary produces no further transitions, so the
			// frozen rasterization scale can snap back fast instead of waiting
			// out the full quiet window (#5701, v200.0.7 field report).
			// The breadcrumb captures the exact moment the user parks and looks
			// at the window: raster vs dpi/96 here says whether they were left
			// staring at wrongly-sized content until the settle.
			dpi, _ := w.DPI()
			raster := 0.0
			if !w.isMinimizing {
				raster = w.currentWebviewRasterizationScale()
			}
			globalApplication.warning("DPI flap: drag released mid-storm on window %d (dpi %d → target %.2f, raster=%.2f) — fast settle armed (#5701)",
				w.parent.id, dpi, w.targetRasterizationScale(uint32(dpi)), raster)
			w.scheduleDPIFlapSettle(dpiFlapReleaseSettle + 50*time.Millisecond)
		}
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
		// Paint the background with the configured colour so that areas not yet
		// covered by WebView2 during a resize show the correct colour instead of white.
		if w.parent.options.BackgroundType == BackgroundTypeSolid {
			col := w.parent.options.BackgroundColour
			hdc := w32.HDC(wparam)
			rc := w32.GetClientRect(w.hwnd)
			colorRef := w32.COLORREF(uint32(col.Red) | uint32(col.Green)<<8 | uint32(col.Blue)<<16)
			hbrush := w32.CreateSolidBrush(colorRef)
			w32.FillRect(hdc, rc, hbrush)
			w32.DeleteObject(w32.HGDIOBJ(hbrush))
		}
		return 1
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
		globalApplication.debug("w32.WM_SYSKEYDOWN", "wparam", uint(wparam))
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
				// A maximised window leaving the minimised state arrives
				// here (not at SIZE_RESTORED), and needs the same DPI
				// resync as the restore path below (#5544). Clear
				// isMinimizing FIRST: the window has left the minimised
				// state, and the verify ladder's #5605 gate reads it.
				w.isMinimizing = false
				w.resyncWebviewDPIAfterUnminimiseIfDPIChanged()
				w.parent.emit(events.Windows.WindowUnMinimise)
			}
			w.isMinimizing = false
			w.parent.emit(events.Windows.WindowMaximise)
			if w.menu != nil && w.menubarTheme != nil {
				w32.RedrawWindow(w.hwnd, nil, 0, w32.RDW_FRAME|w32.RDW_INVALIDATE)
			}
		case w32.SIZE_RESTORED:
			restoredFromMaximised := w.lastSizeWParam == w32.SIZE_MAXIMIZED
			if restoredFromMaximised {
				// Native caption drag from a maximised frameless window bypasses UnMaximise(),
				// so we need to restore the saved constraints here before the next resize.
				w.parent.restoreSavedSizeConstraintOptions()
			}
			if w.isMinimizing {
				// While minimised the window is repositioned off the monitor it
				// will restore to (Windows parks it off-screen / at the primary
				// monitor's corner), which on mixed-DPI systems can re-associate
				// it with another monitor's DPI. Nothing corrects WebView2's
				// rasterization scale on restore, so window.devicePixelRatio
				// keeps the wrong monitor's value until a manual resize (#5544).
				// Clear isMinimizing FIRST: the window has left the minimised
				// state, and the verify ladder's #5605 gate reads it.
				w.isMinimizing = false
				w.resyncWebviewDPIAfterUnminimiseIfDPIChanged()
				w.parent.emit(events.Windows.WindowUnMinimise)
			}
			w.isMinimizing = false
			w.parent.emit(events.Windows.WindowRestore)
			// Repaint the dark menubar only when leaving a maximized/snapped state.
			// SIZE_RESTORED fires on every WM_SIZE during live drag-resize; gating on the
			// previous state avoids per-frame invalidations and the associated flicker.
			// WM_ENTERSIZEMOVE/WM_EXITSIZEMOVE cannot guard this because keyboard snap
			// (Win+Left) bypasses those messages entirely.
			if restoredFromMaximised && w.menu != nil && w.menubarTheme != nil {
				w32.RedrawWindow(w.hwnd, nil, 0, w32.RDW_FRAME|w32.RDW_INVALIDATE)
			}
		case w32.SIZE_MINIMIZED:
			w.isMinimizing = true
			w.parent.emit(events.Windows.WindowMinimise)
		}
		w.lastSizeWParam = wparam

		// Keep the last-known DPI current while the window is in a normal
		// (non-minimised) state. The un-minimise transitions above read this to
		// decide whether a WebView2 rasterization resync is actually needed; it
		// must never be sampled while minimised, where the window is repositioned
		// off its restore monitor and GetDpiForWindow can report a different
		// monitor's DPI (#5605, #5544).
		if wparam != w32.SIZE_MINIMIZED {
			if dpi, _ := w.DPI(); dpi != 0 {
				w.lastKnownDPI = dpi
			}
		}

		if !(w.parent.options.Frameless && wparam == w32.SIZE_MINIMIZED) {
			// If the window is frameless and minimizing, suppress the WebView2 resize.
			// Without this, restoring has wrong dimensions during the animation.
			// See https://github.com/MicrosoftEdge/WebView2Feedback/issues/2549
			width := int32(lparam & 0xFFFF)
			height := int32((lparam >> 16) & 0xFFFF)
			bounds := &edge.Rect{Left: 0, Top: 0, Right: width, Bottom: height}
			InvokeSync(func() {
				time.Sleep(1 * time.Nanosecond)
				w.chromium.ResizeWithBounds(bounds)
				w.parent.emit(events.Windows.WindowDidResize)
			})
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
		newDPI := uint32(wparam & 0xFFFF)
		suggested := (*w32.RECT)(unsafe.Pointer(lparam))

		// Redundant-transition dedup (tao/winit/Chromium all carry the same
		// guard): Windows can deliver a WM_DPICHANGED for the DPI the window
		// was already switched to (a "dpi 120 -> 120" was observed in the
		// field). Processing it would only add resize/COM churn during exactly
		// the storms that crash the GPU process (#5701).
		if newDPI != 0 && newDPI == w.lastAppliedDPI {
			globalApplication.warning("DPI transition dedup: window %d already at dpi %d — skipped (#5701)", w.parent.id, newDPI)
			break
		}

		// DPI-flap breaker (#5701): a window resting near a mixed-DPI boundary
		// can OSCILLATE — each applied suggested rect pushes the window's
		// majority back across the boundary, Windows fires the opposite
		// WM_DPICHANGED, and the resize/scale storm assert-crashes the WebView2
		// GPU process (exit code 0x80000003) until Chromium intentionally kills
		// the browser process. Detect the reversal pattern and — unless the
		// user is mid-drag — move the window fully onto one monitor so the
		// loop's root condition ends; the recovery ladder heals any GPU
		// deaths the residual churn still causes.
		stormActive, resolveStraddle := w.noteDPITransitionAndDetectFlap(uint32(w.lastKnownDPI), newDPI)
		// Diagnostic breadcrumb (#5701): mixed-DPI monitor crossings are the
		// trigger for WebView2 GPU-process deaths in the field, so record every
		// transition — a "WebView2 process failed" shortly after one of these
		// confirms the correlation. Warning level so log bridges that forward
		// only Warn+ still ship it; DPI transitions are rare enough not to spam.
		// raster= is the controller's LIVE rasterization scale at flip time: it
		// shows whether the browser's native scale detection has already caught
		// up with the monitor the window is entering (raster == newDPI/96) or is
		// still on the old monitor's scale — the losing side of that race is
		// what the user sees as wrongly-sized content (#5677). Property read
		// only; skipped while minimised (#5605).
		rasterScale := 0.0
		if !w.isMinimizing {
			rasterScale = w.currentWebviewRasterizationScale()
		}
		globalApplication.warning("DPI transition: window %d dpi %d -> %d (suggested %dx%d px, raster=%.2f, minimised=%v, inSizeMove=%v, stormActive=%v) (#5701)",
			w.parent.id, w.lastKnownDPI, newDPI,
			suggested.Right-suggested.Left, suggested.Bottom-suggested.Top, rasterScale, w.isMinimizing, w.inSizeMove, stormActive)
		resolved := false
		if resolveStraddle && !w.inSizeMove && !w.isMinimizing {
			w.resolveDPIFlapStraddle(suggested)
			// The resolver placed the window fully onto the target monitor
			// (or is rate-limited/converged from doing so a moment ago) —
			// applying the ORIGINAL suggested rect below would put the
			// straddle right back and re-feed the loop.
			resolved = true
		}
		if !w.isMinimizing {
			w.lastAppliedDPI = newDPI
		}
		// While minimised the window is repositioned off its restore monitor; a
		// DPI change in that state delivers a suggested rect scaled for whatever
		// monitor the parked position maps to. Applying it resizes the
		// window's restore bookkeeping (e.g. a maximised 1920x1080 window
		// restored at 3072x1728 after crossing a 200%→125% boundary while
		// minimised, #5544). Skip the resize; a fresh WM_DPICHANGED with a
		// correct rect arrives if the DPI really differs on restore.
		// Geometry stays OS-applied even during a storm (Steps 15/16, #5701):
		// the suggested rect is applied on every transition (except when the
		// resolver just placed the window itself), so window sizes can never
		// go stale and snap/maximise sizes are never fought. Content stays
		// visible and tracks the correct size via native monitor-scale
		// detection; our manual resync is deferred to the settle, and the
		// parked-straddle resolver ends the rect feedback loop at its root.
		if !w.ignoreDPIChangeResizing && !w.isMinimizing && !resolved {
			newWindowRect := (*w32.RECT)(unsafe.Pointer(lparam))
			flags := w32.SWP_NOZORDER | w32.SWP_NOACTIVATE
			// For frameless windows, include SWP_FRAMECHANGED to trigger WM_NCCALCSIZE
			// and recalculate hit-test regions for proper mouse interaction after DPI change.
			// See: https://github.com/wailsapp/wails/issues/4691
			if w.parent.options.Frameless {
				flags |= w32.SWP_FRAMECHANGED
			}
			w32.SetWindowPos(w.hwnd,
				uintptr(0),
				int(newWindowRect.Left),
				int(newWindowRect.Top),
				int(newWindowRect.Right-newWindowRect.Left),
				int(newWindowRect.Bottom-newWindowRect.Top),
				uint(flags))
			// For frameless windows with decorations, re-extend the frame into client area
			// to ensure proper window frame styling after DPI change.
			if w.framelessWithDecorations() {
				if err := w32.ExtendFrameIntoClientArea(w.hwnd, true); err != nil {
					globalApplication.handleFatalError(err)
				}
			}
			// Refresh the layered hit-test region after DPI resize, same as setPhysicalBounds.
			if exStyle := w32.GetWindowLong(w.hwnd, w32.GWL_EXSTYLE); exStyle&w32.WS_EX_LAYERED != 0 {
				w32.SetLayeredWindowAttributes(w.hwnd, 0, 255, w32.LWA_ALPHA)
			}
		}
		// No manual scale correction here: with native monitor-scale detection
		// on, the browser owns the rasterization scale, and the
		// RasterizationScaleChanged handler re-puts the bounds so the content
		// re-lays-out at whatever scale the browser commits (the #5677
		// scale-without-relayout gap). The manual clean-path resync that used
		// to live here was the crash-adjacent actor in the field: every
		// remaining GPU-process exit in the v200.0.18 AND v200.0.22 sessions
		// sat 0.6-2s after one of its scale-puts, while the event-driven path
		// alone closed the wrong-size race in ≤57ms (#5701). Historical note —
		// under WINDOWED hosting, per-flip correction in either ordering
		// (v200.0.16/17) tripped Chromium's browser-process kill within ~21s;
		// Window-to-Visual hosting removed that cost class, and the manual
		// puts turned out to be redundant anyway (every settle logged
		// put=false: native detection had always corrected first).
		// While minimised skip the bookkeeping: the window sits off its
		// restore monitor and GetDpiForWindow can report the wrong monitor's
		// DPI (#5605); a genuine DPI difference is caught on restore by
		// resyncWebviewDPIAfterUnminimiseIfDPIChanged.
		if !w.isMinimizing {
			// Track the new DPI so the un-minimise comparison stays accurate.
			if dpi, _ := w.DPI(); dpi != 0 {
				w.lastKnownDPI = dpi
			}
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

const (
	// dpiFlapReversalWindow bounds how far apart two opposing WM_DPICHANGED
	// transitions may be and still count as oscillation. Field data (#5701)
	// shows the feedback loop running anywhere from ~350 ms to ~2.5 s per
	// transition depending on where the window sits relative to the boundary —
	// the first breaker shipped with 600 ms and missed the slower storms
	// entirely (every breadcrumb logged flapSuppressed=false while the GPU
	// process was assert-crashing). A human CAN reverse a drag within 3 s, so
	// occasional false positives are accepted: the cost is a 3 s resize pause
	// and a settle resync, versus a GPU-process crash cascade for a miss.
	dpiFlapReversalWindow = 3 * time.Second

	// dpiFlapSuppression is the resize-suppression period once oscillation is
	// detected. Transitions arriving while suppressed EXTEND the suppression —
	// the storm is evidently still alive — so the total pause is storm-scoped,
	// not fixed.
	dpiFlapSuppression = 3 * time.Second

	// dpiParkedReversalWindow bounds the PARKED fast path. A window that is
	// not in the modal move/size loop receives opposing WM_DPICHANGED
	// transitions only when it rests straddling a mixed-DPI boundary and the
	// applied rects are flipping its majority monitor — a user cannot
	// ping-pong a parked window's DPI except via two full keyboard monitor
	// moves, and for those the straddle resolver is behaviourally equivalent
	// to the default handling (it applies the transition's own suggested size,
	// placed fully on the target monitor). Field data (#5701, v200.0.2): the
	// parked oscillation runs as slowly as ~6 s per transition, far under the
	// 3-reversals-in-3s in-drag threshold, and every GPU-process crash in that
	// run sat next to these unsuppressed slow reversals. One parked reversal
	// is therefore enough to trip.
	dpiParkedReversalWindow = 10 * time.Second

	// dpiFlapQuietSettle is how long the transition stream must stay quiet
	// before the settle resync runs and suppression ends. Suppression freezes
	// the rasterization scale, and that freeze is user-visible — content
	// renders too small on a higher-DPI monitor (or too large on a lower-DPI
	// one) until the settle resync (v200.0.6 field report). The storm cadence
	// tops out around ~1 s per transition, so 1.5 s of quiet means the storm
	// is over; chasing the end of the storm instead of the full suppression
	// window cuts the visible wrong-scale tail from ~3.5-6.5 s to ~1.6 s.
	dpiFlapQuietSettle = 1500 * time.Millisecond

	// dpiFlapReleaseSettle lives in dpiflap_settle.go: decideScaleReconcile
	// uses it as the post-flip quiet gate, and that helper must compile on
	// every OS for the mac/linux CI regression tests.

	// dpiFlapResumeWindow: a reversal this soon after a settle is the SAME
	// storm resuming (v200.0.7 field trace: fresh transitions 25 ms and
	// 138 ms after two settles), so re-trip on the FIRST reversal even
	// in-drag — demanding 3 fresh reversals hands the resumed storm 2-3
	// unsuppressed transitions of full resize+resync churn, exactly the
	// exposure the breaker exists to remove.
	dpiFlapResumeWindow = 10 * time.Second

	// dpiPlacementGrace: DPI reversals this soon after an app-driven
	// setPhysicalBounds are placement side-effects, not a straddle
	// oscillation — a multi-monitor placement legitimately ping-pongs the
	// window's DPI while it lands. The parked-reversal fast path and the
	// resolver stand down inside the grace: the v200.0.22 field trace shows
	// the resolver "resolving" a freshly placed fullscreen alarm window onto
	// the wrong monitor with a mid-placement 3456x1944 rect (#5701). 2 s
	// comfortably covers the placement's transition burst (~10 ms spacing in
	// the trace) while leaving genuinely parked straddles — which oscillate
	// for as long as they straddle — to trip right after it expires.
	dpiPlacementGrace = 2 * time.Second
)

// noteDPITransitionAndDetectFlap records a WM_DPICHANGED transition and
// classifies it. It detects the mixed-DPI oscillation observed in the field
// (#5701): a window resting near a monitor boundary ping-pongs A→B→A→B because
// each applied suggested rect moves the window's majority back across the
// boundary. Chromium tolerates a handful of GPU-process crashes and then
// INTENTIONALLY kills the browser process ("GPU process isn't usable.
// Goodbye."), so an unbroken storm always ends in a wedged controller.
//
// Returns (stormActive, resolveStraddle): stormActive is telemetry plus the
// settle-backstop trigger these days (the suggested rect is always applied,
// native scale detection owns the rasterization scale, and the
// RasterizationScaleChanged handler re-lays content out per scale change —
// geometry stays OS-owned, Steps 15/16); resolveStraddle=true → the caller
// should run the straddle resolver for this transition (further gated there
// on !inSizeMove: a window the user actively holds must not be repositioned;
// and suppressed entirely within dpiPlacementGrace of an app-driven
// placement). Detection has these tiers:
//   - in-drag: 3 reversals within dpiFlapReversalWindow (a human CAN wiggle a
//     drag across the boundary, so demand a sustained pattern);
//   - parked: ONE reversal within dpiParkedReversalWindow trips immediately —
//     a parked window ping-ponging its DPI is proof of the feedback loop, and
//     field data shows the parked oscillation can run slower than any
//     rapid-reversal threshold while still assert-crashing the GPU process;
//   - resumed: ONE reversal within dpiFlapResumeWindow of the last settle
//     trips immediately even in-drag — a storm restarting right after a lull
//     is a continuation, not a fresh human wiggle.
//
// While suppressed, parked reversals keep resolving: suppression stops the
// resize churn, but only the resolver ends the straddle feeding the storm. A
// settle resync self-reschedules until the transition stream has been quiet
// for dpiFlapQuietSettle, then resyncs and ends suppression. Main thread only.
func (w *windowsWebviewWindow) noteDPITransitionAndDetectFlap(fromDPI, newDPI uint32) (suppressed, resolveStraddle bool) {
	now := time.Now()
	parked := !w.inSizeMove && !w.isMinimizing
	isReversal := newDPI == w.lastDPITransitionFrom
	reversalAge := now.Sub(w.lastDPITransitionAt)
	// A reversal inside the placement grace is the app placing the window
	// (multi-monitor placement legitimately ping-pongs the DPI), not a parked
	// straddle oscillation — the resolver must never fight it. v200.0.22 field
	// trace: it moved a freshly placed fullscreen alarm window onto the wrong
	// monitor with a mid-placement 3456x1944 rect (#5701).
	recentPlacement := !w.lastProgrammaticPlacementAt.IsZero() &&
		now.Sub(w.lastProgrammaticPlacementAt) < dpiPlacementGrace
	parkedReversal := parked && isReversal && reversalAge < dpiParkedReversalWindow && !recentPlacement

	if now.Before(w.dpiFlapSuppressUntil) {
		// Still oscillating while suppressed — extend so suppression outlives
		// the storm instead of re-arming it for another crash window. Keep the
		// transition bookkeeping current (a stale lastDPITransitionAt made the
		// first post-settle reversal look isolated, delaying re-detection into
		// exactly the unsuppressed gap where the v200.0.2 crashes clustered).
		w.dpiFlapSuppressUntil = now.Add(dpiFlapSuppression)
		w.lastDPITransitionFrom = fromDPI
		w.lastDPITransitionAt = now
		return true, parkedReversal
	}

	if isReversal && reversalAge < dpiFlapReversalWindow {
		w.dpiFlapReversals++
	} else {
		w.dpiFlapReversals = 0
	}
	w.lastDPITransitionFrom = fromDPI
	w.lastDPITransitionAt = now

	// A reversal arriving shortly after a settle is the SAME storm resuming
	// after a lull — trip again on the first reversal instead of handing the
	// storm 2-3 unsuppressed transitions while the counter refills.
	resumedStorm := isReversal && reversalAge < dpiFlapReversalWindow &&
		!w.lastDPIFlapSettledAt.IsZero() && now.Sub(w.lastDPIFlapSettledAt) < dpiFlapResumeWindow

	if w.dpiFlapReversals < 3 && !parkedReversal && !resumedStorm {
		return false, false
	}
	w.dpiFlapReversals = 0
	w.dpiFlapSuppressUntil = now.Add(dpiFlapSuppression)
	w.dpiFlapStormStartAt = now
	switch {
	case parkedReversal:
		w.dpiFlapResumeCount = 0
		globalApplication.warning("DPI flap detected (parked reversal %d↔%d): window %d is straddle-oscillating at rest — resolving (#5701)",
			fromDPI, newDPI, w.parent.id)
	case resumedStorm:
		// The ordinal is telemetry only these days: with the settle reduced to
		// one backstop bounds re-assert (no scale-put), settle/resume churn is
		// cheap, so the old per-resume threshold escalation was dropped.
		w.dpiFlapResumeCount++
		globalApplication.warning("DPI flap resumed (%d↔%d within %s of the last settle, resume %d): window %d — re-tripping on the first reversal (#5701)",
			fromDPI, newDPI, dpiFlapResumeWindow, w.dpiFlapResumeCount, w.parent.id)
	default:
		w.dpiFlapResumeCount = 0
		globalApplication.warning("DPI flap detected: window %d oscillating %d↔%d — settle backstop armed (#5701)",
			w.parent.id, fromDPI, newDPI)
	}

	// The webview stays VISIBLE through the storm (Step 16). With geometry
	// OS-applied, native monitor-scale detection owning the scale, and the
	// RasterizationScaleChanged handler re-laying content out per scale
	// change, mid-storm content tracks the correct size on its own. A storm
	// trip now only arms the resolver (except inside the placement grace —
	// the app knows where its window goes) and schedules the belt-and-braces
	// settle sync.
	w.scheduleDPIFlapSettle(dpiFlapQuietSettle + 100*time.Millisecond)
	return true, !recentPlacement
}

// scheduleDPIFlapSettle arms a settle check after delay. Duplicate pending
// checks are harmless: dpiFlapSettleCheck no-ops once the storm has settled
// (dpiFlapSuppressUntil zeroed), and each live check either settles or re-arms
// exactly one successor, so the timer count stays bounded. Callable from any
// context; the check itself hops to the main thread via InvokeAsync.
func (w *windowsWebviewWindow) scheduleDPIFlapSettle(delay time.Duration) {
	time.AfterFunc(delay, func() {
		InvokeAsync(func() {
			w.dpiFlapSettleCheck()
		})
	})
}

// dpiFlapSettleCheck ends the storm once the transition stream has been quiet
// long enough, reconciles the rasterization scale, re-asserts the controller
// bounds as a backstop, and emits the storm/tail breadcrumb that quantifies the
// exposure in field logs. Native detection owns the scale during the storm
// (every settle in the v200.0.21/22 field sessions logged put=false — the
// browser had always corrected first) and the RasterizationScaleChanged
// handler does the per-event re-layout, so the normal path still puts nothing.
// The single Resize below covers a missed event (registration is
// non-fatal-logged, so it CAN be absent), AND a raster==target guard closes the
// hole v200.0.23 exposed: if native detection went silent (no scale-changed
// event after the last flip) the settle would lock in a stale raster, so when
// raster still differs from target at settle we force one corrective scale-put
// before the Resize. Quiet thresholds: mid-drag a lull may just be the user
// hovering near the boundary (dpiFlapQuietSettle); once the drag is released a
// non-straddling window produces no further transitions, so
// dpiFlapReleaseSettle suffices (v200.0.7 field report). Main thread only.
func (w *windowsWebviewWindow) dpiFlapSettleCheck() {
	if w.hwnd == 0 || w.parent.isDestroyed() || globalApplication.performingShutdown {
		return
	}
	if w.dpiFlapSuppressUntil.IsZero() {
		return // already settled (duplicate check, e.g. from the drag-release path)
	}
	threshold := dpiFlapQuietSettle
	if !w.inSizeMove {
		threshold = dpiFlapReleaseSettle
	}
	quiet := time.Since(w.lastDPITransitionAt)
	if quiet < threshold || w.webviewRebuildInProgress {
		// Storm (or a rebuild) still active — check again once the quiet
		// window could next be satisfied.
		next := threshold - quiet + 100*time.Millisecond
		if next < 250*time.Millisecond {
			next = 250 * time.Millisecond
		}
		w.scheduleDPIFlapSettle(next)
		return
	}
	w.dpiFlapSuppressUntil = time.Time{} // storm over — resume normal DPI handling
	w.lastDPIFlapSettledAt = time.Now()
	if dpi, _ := w.DPI(); dpi != 0 {
		w.lastKnownDPI = dpi
	}
	// Backstop bounds re-assert: a scale change without a bounds re-assert
	// does not re-lay out the content (#5677). The event handler normally did
	// this already; one extra Resize per storm is cheap insurance against a
	// missed event (registration is non-fatal, so it CAN be absent).
	if w.chromium.IsReady() {
		w.chromium.Resize()
	}
	// Client size is logged so field data can separate the two wrong-size
	// failure modes: stale content layout (client size correct for the DPI,
	// content wrong) vs stale window size (client size sized for the other
	// monitor — e.g. the v200.0.22 alarm placement defect showed here as
	// client 3456x1945 on dpi 120).
	var cw, ch int32
	if rect := w32.GetClientRect(w.hwnd); rect != nil {
		cw, ch = rect.Right-rect.Left, rect.Bottom-rect.Top
	}
	storm := w.lastDPITransitionAt.Sub(w.dpiFlapStormStartAt)
	globalApplication.warning("DPI flap settled: window %d on dpi %d (storm %dms + %dms tail, client %dx%d px) (#5701)",
		w.parent.id, w.lastKnownDPI, storm.Milliseconds(), quiet.Milliseconds(), cw, ch)
	// The scale story moved to the DPI verify ladder: verify now (absorbing
	// the v200.0.23 settle guard — if the scale writer left a stale raster,
	// this is the corrective put), then re-check at +2s and +10s. Wrong-scale
	// events provably land up to ~6s AFTER the last flip (v200.0.24 field
	// trace: RasterizationScaleChanged → 1.25 committed 19ms after the final
	// flip to dpi 216, then silence for 62s), i.e. a settle-time check alone
	// leaves a post-settle blind window — the probes close it, and every
	// occurrence becomes a logged MISMATCH line instead of a screenshot-only
	// defect.
	w.verifyWebviewScale("settle", true)
	w.scheduleScaleVerify("settle+2s", scaleVerifyProbeShort, true)
	w.scheduleScaleVerify("settle+10s", scaleVerifyProbeLong, true)
}

// resolveDPIFlapStraddle ends the oscillation's root condition: the window
// straddling a mixed-DPI monitor boundary. It takes the OS suggested rect of
// the transition that tripped the breaker, shifts it so it lies ENTIRELY on
// the suggested rect's target monitor (tao/winit ship the same nudge on
// Windows 10; here it is gated to flap detection only, since winit PR #4119
// showed unconditional nudging causes the opposite bug), and applies position
// + size in a single SetWindowPos. With the window fully on one monitor the
// majority-area rule has nothing left to flip. Skipped while the user holds a
// drag (the modal move loop would fight the reposition). Main thread only.
func (w *windowsWebviewWindow) resolveDPIFlapStraddle(suggested *w32.RECT) {
	// Rate limit: with the parked fast path the resolver can be invoked on
	// every transition of a storm; one reposition per second is plenty to
	// converge and guarantees the resolver itself can never become the churn.
	now := time.Now()
	if now.Sub(w.lastStraddleResolveAt) < time.Second {
		return
	}
	monitor := w32.MonitorFromRect(suggested, w32.MONITOR_DEFAULTTONEAREST)
	if monitor == 0 {
		return
	}
	var info w32.MONITORINFO
	info.CbSize = uint32(unsafe.Sizeof(info))
	if !w32.GetMonitorInfo(monitor, &info) {
		return
	}
	work := info.RcWork
	rect := *suggested
	width := rect.Right - rect.Left
	height := rect.Bottom - rect.Top

	// Shift fully inside the work area (size preserved; an oversized window
	// simply pins to the near edge).
	if rect.Right > work.Right {
		rect.Left -= rect.Right - work.Right
	}
	rect.Right = rect.Left + width
	if rect.Left < work.Left {
		rect.Left = work.Left
		rect.Right = rect.Left + width
	}
	if rect.Bottom > work.Bottom {
		rect.Top -= rect.Bottom - work.Bottom
	}
	rect.Bottom = rect.Top + height
	if rect.Top < work.Top {
		rect.Top = work.Top
		rect.Bottom = rect.Top + height
	}

	// Converged guard: if the window already sits exactly at the target rect
	// yet transitions keep arriving (monitor assignment flipping without
	// movement), re-Setting the same rect would only feed the loop.
	if cur := w32.GetWindowRect(w.hwnd); cur != nil &&
		cur.Left == rect.Left && cur.Top == rect.Top &&
		cur.Right == rect.Right && cur.Bottom == rect.Bottom {
		return
	}

	w.lastStraddleResolveAt = now
	globalApplication.warning("DPI flap resolver: window %d moved fully onto its target monitor at (%d,%d) %dx%d (#5701)",
		w.parent.id, rect.Left, rect.Top, width, height)
	w32.SetWindowPos(w.hwnd, 0,
		int(rect.Left), int(rect.Top), int(width), int(height),
		uint(w32.SWP_NOZORDER|w32.SWP_NOACTIVATE))
}

// syncWebviewRasterizationScale puts the given DPI's scale onto the WebView2
// controller and reports whether a re-put was needed. Used with the window's
// current DPI by the settle/unminimise resyncs, and with WM_DPICHANGED's
// incoming DPI BEFORE the suggested rect applies (Step 19, #5701) so the
// rect's WM_SIZE re-layout rasters at the correct scale. It is a no-op when
// the controller is unavailable or the scale is already in sync; a failed
// put is log-only.
func (w *windowsWebviewWindow) syncWebviewRasterizationScale(dpi uint32) bool {
	// The #5605 restore crash is prevented by the DPI-change gate in
	// resyncWebviewDPIAfterUnminimiseIfDPIChanged, which keeps us off the
	// controller entirely when the DPI is unchanged. The GetController nil
	// check below covers the early-initialisation / unavailable-controller
	// window.
	controller := w.chromium.GetController()
	if controller == nil {
		return false
	}
	controller3 := controller.GetICoreWebView2Controller3()
	if controller3 == nil {
		return false
	}
	if dpi == 0 {
		return false
	}
	// dpi/96 × text scale — the text-scale term is non-optional: WebView2
	// defines RasterizationScale as monitor scale × text scale, so putting a
	// bare dpi/96 would shrink a text-scaling user's content (#5701).
	scale := w.targetRasterizationScale(dpi)
	// Compare with a tolerance: GetRasterizationScale returns a float that may
	// not be bit-identical to the target for non-25% DPI steps, and an exact ==
	// would re-Put the scale on every resync.
	if current, err := controller3.GetRasterizationScale(); err == nil && math.Abs(current-scale) < 0.001 {
		return false
	}
	if err := controller3.PutRasterizationScale(scale); err != nil {
		globalApplication.error("failed to update WebView2 rasterization scale: %s", err)
		return false
	}
	return true
}

// currentWebviewRasterizationScale reads the controller's live rasterization
// scale for diagnostics; 0 means unavailable. A property read only — it never
// triggers a re-raster, so it is safe inside DPI storms. Callers must gate on
// !isMinimizing: COM calls into a possibly-suspended controller are the #5605
// restore-crash class.
func (w *windowsWebviewWindow) currentWebviewRasterizationScale() float64 {
	controller := w.chromium.GetController()
	if controller == nil {
		return 0
	}
	controller3 := controller.GetICoreWebView2Controller3()
	if controller3 == nil {
		return 0
	}
	scale, err := controller3.GetRasterizationScale()
	if err != nil {
		return 0
	}
	return scale
}

// onWebviewRasterizationScaleChanged closes the scale/layout race of the
// mixed-DPI wrong-size investigation (#5701). With native monitor-scale
// detection on, the browser re-rasters on its own schedule during monitor
// moves, and a scale change with no following bounds put re-scales the
// displayed frame WITHOUT re-laying it out (#5677) — v200.0.21 field data:
// the native scale-put landed a median 12ms AFTER the flip's WM_SIZE
// bounds-put on 68/91 flips, leaving crisp-but-wrong-sized content visible
// until the settle (median 3.3s). Re-putting the bounds right here re-lays
// the content out at the scale the browser just committed, collapsing that
// exposure to the event latency itself. The re-layout cost rides a re-raster
// the browser has ALREADY done — this is affordable per event where per-flip
// scale-puts were not (the v200.0.16/17 browser kills were a windowed-hosting
// GPU cost; under Window-to-Visual hosting the same session shows zero
// process failures). No feedback loop: a bounds put never changes the scale.
// Skipped while minimised (#5605: no controller COM off a possibly-suspended
// state) and mid-rebuild. Warning level so New Relic receives the breadcrumb.
// Fires on the UI thread — the same thread as the WM_SIZE resize path.
func (w *windowsWebviewWindow) onWebviewRasterizationScaleChanged(scale float64) {
	dpi, _ := w.DPI()
	sinceFlip := int64(-1)
	if !w.lastDPITransitionAt.IsZero() {
		sinceFlip = time.Since(w.lastDPITransitionAt).Milliseconds()
	}
	relayout := !w.isMinimizing && !w.webviewRebuildInProgress && scale > 0
	if relayout {
		w.chromium.Resize()
	}
	// target is text-scale aware (see rasterizationTargetForDPI) and diff is
	// logged explicitly so a field session shows every event where the
	// browser committed a scale that disagrees with the window's monitor —
	// the v200.0.24 wrong-scale stick started with exactly such an event.
	target := w.targetRasterizationScale(uint32(dpi))
	storm := !w.dpiFlapSuppressUntil.IsZero()
	// A mismatched commit OUTSIDE a storm has no settle chain coming to
	// correct it — this is exactly how the v200.0.24 62s stick would recur
	// even with a settle guard (wrong events land up to ~6s after the last
	// flip, i.e. after the settle). Arm a verify probe rather than putting
	// inline: the event handler is the historically crash-adjacent context,
	// and the probe keeps verifyWebviewScale the single put site. No loop: a
	// corrective put echoes this event with scale==target → diff 0 → no
	// probe. During storms the settle chain owns the correction.
	probe := false
	if !storm && !w.isMinimizing && !w.webviewRebuildInProgress &&
		math.Abs(scale-target) > dpiFlapSettleScaleTolerance {
		probe = true
		w.scheduleScaleVerify("scale-event", scaleVerifyRetryProbe, true)
	}
	globalApplication.warning(
		"WebView2 rasterization scale now %.3f: window %d (dpi %d → target %.3f, diff %+.3f), storm=%v, inSizeMove=%v, minimised=%v, relayout=%v, probe=%v, %dms after last DPI flip (#5701)",
		scale, w.parent.id, dpi, target, scale-target,
		storm, w.inSizeMove, w.isMinimizing, relayout, probe, sinceFlip)
}

// resyncWebviewDPIAfterUnminimiseIfDPIChanged is the un-minimise entry point
// for the DPI resync (#5544: while minimised the window is parked off its
// restore monitor, so the rasterization scale can drift). It only proceeds
// when the window's DPI actually differs from the last value seen while
// non-minimised — in the common case (same monitor, unchanged DPI) it makes
// zero COM calls, which is what prevents the restore crash in #5605: after a
// window has been minimised long enough for WebView2 to suspend or for its
// render/GPU process to be torn down, any COM call into the controller can be
// fatal, and the old resync's GetRasterizationScale probe ran on every
// restore. DPI is read with the Win32 GetDpiForWindow, never via COM. On a
// genuine change the verify ladder observes, logs, and (if needed) corrects +
// re-lays-out in one pass. Callers must clear isMinimizing first.
func (w *windowsWebviewWindow) resyncWebviewDPIAfterUnminimiseIfDPIChanged() {
	dpi, _ := w.DPI()
	if dpi == 0 || dpi == w.lastKnownDPI {
		// DPI unchanged (or unreadable): nothing to resync. Critically, do not
		// call into the WebView2 controller here.
		return
	}
	w.verifyWebviewScale("unminimise", true)
}

// enableNativeMonitorScaleDetection re-enables WebView2's built-in monitor
// scale-change detection. edge.NewChromium disables it (PutShouldDetectMonitor
// ScaleChanges(false)) to run in raw-pixels bounds mode and make the app own DPI
// transitions. On a DPI scale-up across a mixed-DPI / dual-GPU boundary the
// manual rasterization-scale resync could not keep WebView2's render surface
// valid, so the render process terminated and every subsequent controller COM
// call returned RESOURCE_NOT_IN_CORRECT_STATE forever (wailsapp/wails#5701).
// Letting WebView2 detect the change natively keeps the surface valid. Only the
// rasterization scale changes hands; BoundsMode stays USE_RAW_PIXELS, so the
// app still sets exact pixel bounds. Called via configureWebviewScaleOwnership
// right after Embed, when the controller is guaranteed live (Embed
// message-pumps until initialised). Returns the outcome for the ownership
// breadcrumb ("on" / "FAILED: …" / "unavailable: …").
func (w *windowsWebviewWindow) enableNativeMonitorScaleDetection() string {
	if w.chromium == nil {
		return "unavailable: no chromium"
	}
	controller := w.chromium.GetController()
	if controller == nil {
		return "unavailable: no controller"
	}
	controller3 := controller.GetICoreWebView2Controller3()
	if controller3 == nil {
		return "unavailable: no controller3"
	}
	if err := controller3.PutShouldDetectMonitorScaleChanges(true); err != nil {
		globalApplication.error("failed to enable WebView2 monitor scale detection: %s", err)
		return "FAILED: " + err.Error()
	}
	return "on"
}

// crossPermissionToWebView2Kind maps a cross-platform PermissionType to the
// WebView2 permission-kind enum. Unknown types map to UnknownPermission.
func crossPermissionToWebView2Kind(p PermissionType) edge.CoreWebView2PermissionKind {
	switch p {
	case PermissionMicrophone:
		return edge.CoreWebView2PermissionKind(CoreWebView2PermissionKindMicrophone)
	case PermissionCamera:
		return edge.CoreWebView2PermissionKind(CoreWebView2PermissionKindCamera)
	case PermissionGeolocation:
		return edge.CoreWebView2PermissionKind(CoreWebView2PermissionKindGeolocation)
	case PermissionNotifications:
		return edge.CoreWebView2PermissionKind(CoreWebView2PermissionKindNotifications)
	case PermissionClipboardRead:
		return edge.CoreWebView2PermissionKind(CoreWebView2PermissionKindClipboardRead)
	default:
		return edge.CoreWebView2PermissionKind(CoreWebView2PermissionKindUnknownPermission)
	}
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
	topSource, err := sender.GetSource()
	if err != nil {
		globalApplication.error("Unable to get source from sender: %s", err.Error())
		topSource = ""
	}

	senderSource, err := args.GetSource()
	if err != nil {
		globalApplication.error("Unable to get source from args: %s", err.Error())
		senderSource = ""
	}

	// We send all messages to the centralised window message buffer
	windowMessageBuffer <- &windowMessage{
		windowId: w.parent.id,
		message:  message,
		originInfo: &OriginInfo{
			Origin:    senderSource,
			TopOrigin: topSource,
		},
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
	// Recorded for the ProcessFailed diagnostic (#5701): correlating process
	// deaths with the runtime version identifies whether a crash is fixed or
	// introduced by a WebView2 runtime rollout.
	webviewRuntimeVersion = webview2version

	// Browser flags apply globally to the shared WebView2 environment
	// Use application-level options, not per-window options
	appOpts := globalApplication.options.Windows

	// We disable this by default. Can be overridden with the `EnableFraudulentWebsiteWarnings` option
	disabledFeatures := append([]string{"msSmartScreenProtection"}, appOpts.DisabledFeatures...)

	if len(disabledFeatures) > 0 {
		disabledFeatures = sliceutil.Unique(disabledFeatures)
		arg := fmt.Sprintf("--disable-features=%s", strings.Join(disabledFeatures, ","))
		chromium.AdditionalBrowserArgs = append(chromium.AdditionalBrowserArgs, arg)
	}

	if len(appOpts.EnabledFeatures) > 0 {
		enabledFeatures := sliceutil.Unique(appOpts.EnabledFeatures)
		arg := fmt.Sprintf("--enable-features=%s", strings.Join(enabledFeatures, ","))
		chromium.AdditionalBrowserArgs = append(chromium.AdditionalBrowserArgs, arg)
	}

	if len(appOpts.AdditionalBrowserArgs) > 0 {
		chromium.AdditionalBrowserArgs = append(chromium.AdditionalBrowserArgs, appOpts.AdditionalBrowserArgs...)
	}

	chromium.DataPath = globalApplication.options.Windows.WebviewUserDataPath
	chromium.BrowserPath = globalApplication.options.Windows.WebviewBrowserPath

	// Apply the cross-platform Permissions map first; the WebView2-specific
	// Windows.Permissions map below can override individual kinds.
	hasPermissionPolicy := len(w.parent.options.Permissions) > 0 || len(opts.Permissions) > 0
	for permission, state := range w.parent.options.Permissions {
		chromium.SetPermission(crossPermissionToWebView2Kind(permission),
			edge.CoreWebView2PermissionState(state))
	}

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
	// Diagnostic only: surface a render/GPU/browser process death (the
	// unrecoverable "controller zombie" root event, wailsapp/wails#5701) with a
	// stack trace, instead of it being invisible. Must be set before Embed so it
	// is registered when AddProcessFailed runs during controller creation.
	chromium.ProcessFailedCallback = w.processFailed
	// Diagnostic only: timestamp every rasterization-scale change, including
	// the ones the BROWSER makes via native monitor-scale detection. Content is
	// wrongly sized exactly between a browser-initiated scale change and the
	// next bounds put (#5677), so these events bracket the visible wrong-size
	// windows in field logs (#5701). Set before Embed, same as processFailed.
	chromium.RasterizationScaleChangedCallback = w.onWebviewRasterizationScaleChanged

	if !chromium.Embed(w.hwnd) {
		// Environment/controller creation failed — edge reports the error via
		// the error callback and returns instead of exiting (#5701; the old
		// fatal path killed the app from a controller rebuild). Bail out:
		// every chromium call below needs the live controller/webview this
		// embed did not produce. rebuildWebview sees the missing controller
		// (IsReady) and retries; the health watchdog is the backstop.
		globalApplication.error("WebView2 Embed failed for window %d — controller not created (#5701)", w.parent.id)
		return
	}

	// Select the rasterization-scale owner and emit the ownership breadcrumb
	// (#5701). Bounds stay in raw pixels — only the scale changes hands.
	w.configureWebviewScaleOwnership()

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

	// File drop handling on Windows:
	// WebView2's AllowExternalDrop controls ALL drag-and-drop (both external file drops
	// AND internal HTML5 drag-and-drop). We cannot disable it without breaking HTML5 DnD.
	//
	// When EnableFileDrop is true:
	// - JS dragenter/dragover/drop events fire for external file drags
	// - JS calls preventDefault() to stop the browser from navigating to the file
	// - JS uses chrome.webview.postMessageWithAdditionalObjects to send file paths to Go
	// - Go receives paths via processMessageWithAdditionalObjects
	//
	// When EnableFileDrop is false:
	// - We cannot use AllowExternalDrag(false) as it breaks HTML5 internal drag-and-drop
	// - JS runtime checks window._wails.flags.enableFileDrop and shows "no drop" cursor
	// - The enableFileDrop flag is injected in navigationCompleted callback

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

	// WebView2's PermissionRequested handler checks the global permission
	// before the per-kind map, so the blanket "allow all" must only be set
	// when no per-kind policy is configured — otherwise it would override the
	// Permissions map and a configured PermissionDeny would be ignored. When a
	// policy is present, unset kinds fall through to PermissionDefault (the
	// platform's native prompt).
	if !hasPermissionPolicy {
		chromium.SetGlobalPermission(edge.CoreWebView2PermissionStateAllow)
	}
	chromium.AddWebResourceRequestedFilter("*", edge.COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)

	// Reset BEFORE either navigation branch: navigationCompleted early-returns
	// when this flag is still true, skipping the WebView2Feedback#1077
	// Hide/Show repaint nudge and focus restore. On a controller rebuild
	// (#5701) the flag is true from the previous controller's first load, so
	// missing this reset leaves the rebuilt webview permanently unpainted —
	// for HTML-mode windows there was previously no reset at all.
	w.webviewNavigationCompleted = false

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

func (w *windowsWebviewWindow) flash(enabled bool) {
	w32.FlashWindow(w.hwnd, enabled)
}

func (w *windowsWebviewWindow) navigationCompleted(
	sender *edge.ICoreWebView2,
	args *edge.ICoreWebView2NavigationCompletedEventArgs,
) {

	// Install the runtime core
	w.execJS(runtime.Core(globalApplication.impl.GetFlags(globalApplication.options)))

	// Set the EnableFileDrop flag for this window (Windows-specific)
	// The JS runtime checks this before processing file drops
	w.execJS(fmt.Sprintf("window._wails.flags.enableFileDrop = %v;", w.parent.options.EnableFileDrop))

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
	if strings.HasPrefix(message, "file:drop:") {
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

		// Extract X/Y coordinates from message - format is "file:drop:x:y"
		var x, y int
		parts := strings.Split(message, ":")
		if len(parts) >= 4 {
			if parsedX, err := strconv.Atoi(parts[2]); err == nil {
				x = parsedX
			}
			if parsedY, err := strconv.Atoi(parts[3]); err == nil {
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

func (w *windowsWebviewWindow) setFullscreenButtonState(_ ButtonState) {
	// Windows has no dedicated fullscreen button in the standard title bar; no-op.
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
