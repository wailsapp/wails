package application

import "time"

const (
	defaultNotchWindowWidth          = 660
	defaultNotchWindowHeight         = 92
	defaultNotchWindowAnimationSpeed = 420 * time.Millisecond
	notchWindowHorizontalInset       = 30
	notchWindowTopInset              = 34
	notchWindowBottomInset           = 12
)

// NotchWindowOptions configures a window attached to the camera housing at the
// top of a macOS display. Width and Height describe the usable webview area;
// Wails adds the native shaped edges around that area.
type NotchWindowOptions struct {
	// Width is the usable webview width inside the native shaped edges.
	// The default is 660 points.
	Width int

	// Height is the usable webview height inside the native shaped edges.
	// The default is 92 points.
	Height int

	// Animated slides the window down when shown and back up when hidden.
	Animated bool

	// AnimationSpeed is the show animation duration. A zero value uses 420ms.
	// Hide uses two thirds of this duration to preserve the faster dismissal.
	AnimationSpeed time.Duration

	// Screen selects the display whose camera housing receives the window. When
	// nil, Wails targets the primary display.
	Screen *Screen

	// WindowOptions supplies the web content and optional standard Wails window
	// settings. NewNotchWindow owns the native geometry and appearance fields.
	WindowOptions WebviewWindowOptions
}

type notchWindowConfig struct {
	contentWidth   int
	contentHeight  int
	animated       bool
	animationSpeed time.Duration
	screenID       string
}

func normaliseNotchWindowOptions(options NotchWindowOptions) NotchWindowOptions {
	if options.Width <= 0 {
		options.Width = defaultNotchWindowWidth
	}
	if options.Height <= 0 {
		options.Height = defaultNotchWindowHeight
	}
	if options.AnimationSpeed <= 0 {
		options.AnimationSpeed = defaultNotchWindowAnimationSpeed
	}
	return options
}

func webviewOptionsForNotchWindow(options NotchWindowOptions) WebviewWindowOptions {
	result := options.WindowOptions
	result.Width = options.Width + 2*notchWindowHorizontalInset
	result.Height = options.Height + notchWindowTopInset + notchWindowBottomInset
	result.Frameless = true
	result.DisableResize = true
	result.MinWidth = 0
	result.MinHeight = 0
	result.MaxWidth = 0
	result.MaxHeight = 0
	result.StartState = WindowStateNormal
	result.Screen = options.Screen
	result.BackgroundType = BackgroundTypeTransparent
	result.BackgroundColour = NewRGBA(0, 0, 0, 0)
	result.Mac.Backdrop = MacBackdropTransparent
	result.Mac.DisableShadow = true
	result.Mac.CornerType = MacWindowCornerTypeSquare
	result.Mac.WindowClass = MacWindowClassPanel
	result.Mac.PanelPreferences = MacPanelPreferences{
		NonActivating: true,
		FloatingPanel: true,
	}
	result.Mac.WindowLevel = MacWindowLevelPopUpMenu
	result.Mac.CollectionBehavior = MacWindowCollectionBehaviorCanJoinAllSpaces |
		MacWindowCollectionBehaviorFullScreenAuxiliary |
		MacWindowCollectionBehaviorStationary |
		MacWindowCollectionBehaviorIgnoresCycle
	return result
}

// NotchWindow is an independently managed camera-housing window. Its narrow
// lifecycle API keeps native placement and geometry under Wails' control.
type NotchWindow struct {
	window *WebviewWindow
}

// Show displays the window, using its configured animation when enabled.
func (w *NotchWindow) Show() *NotchWindow {
	if w == nil || w.window == nil {
		return w
	}
	w.window.Show()
	return w
}

// Hide hides the window, using its configured animation when enabled.
func (w *NotchWindow) Hide() *NotchWindow {
	if w == nil || w.window == nil {
		return w
	}
	w.window.Hide()
	return w
}

// Visibility reports whether the native window is currently visible.
func (w *NotchWindow) Visibility() bool {
	if w == nil || w.window == nil {
		return false
	}
	return w.window.IsVisible()
}

// Close permanently destroys the native window.
func (w *NotchWindow) Close() {
	if w == nil || w.window == nil {
		return
	}
	w.window.Close()
}

// NewNotchWindow creates an independent camera-housing window. Multiple notch
// windows may coexist; the most recently shown window is ordered in front.
func (wm *WindowManager) NewNotchWindow(options NotchWindowOptions) *NotchWindow {
	if !notchWindowsSupported {
		return &NotchWindow{}
	}
	options = normaliseNotchWindowOptions(options)
	window := NewWindow(webviewOptionsForNotchWindow(options))
	screenID := ""
	if options.Screen != nil {
		screenID = options.Screen.ID
	}
	window.notchWindow = &notchWindowConfig{
		contentWidth:   options.Width,
		contentHeight:  options.Height,
		animated:       options.Animated,
		animationSpeed: options.AnimationSpeed,
		screenID:       screenID,
	}
	wm.addAndRun(window)
	return &NotchWindow{window: window}
}
