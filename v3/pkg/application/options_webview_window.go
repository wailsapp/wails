package application

type WindowState int

const (
	WindowStateNormal WindowState = iota
	WindowStateMinimised
	WindowStateMaximised
	WindowStateFullscreen
)

type WebviewWindowOptions struct {
	// Name is a unique identifier that can be given to a window.
	Name string

	// Title is the title of the window.
	Title string

	// Width is the starting width of the window.
	Width int

	// Height is the starting height of the window.
	Height int

	// AlwaysOnTop will make the window float above other windows.
	AlwaysOnTop bool

	// URL is the URL to load in the window.
	URL string

	// DisableResize will disable the ability to resize the window.
	DisableResize bool

	// Frameless will remove the window frame.
	Frameless bool

	// MinWidth is the minimum width of the window.
	MinWidth int

	// MinHeight is the minimum height of the window.
	MinHeight int

	// MaxWidth is the maximum width of the window.
	MaxWidth int

	// MaxHeight is the maximum height of the window.
	MaxHeight int

	// StartState indicates the state of the window when it is first shown.
	// Default: WindowStateNormal
	StartState WindowState

	// Centered will center the window on the screen.
	Centered bool

	// BackgroundType is the type of background to use for the window.
	// Default: BackgroundTypeSolid
	BackgroundType BackgroundType

	// BackgroundColour is the colour to use for the window background.
	BackgroundColour RGBA

	// HTML is the HTML to load in the window.
	HTML string

	// JS is the JavaScript to load in the window.
	JS string

	// CSS is the CSS to load in the window.
	CSS string

	// X is the starting X position of the window.
	X int

	// Y is the starting Y position of the window.
	Y int

	// TransparentTitlebar will make the titlebar transparent.
	// TODO: Move to mac window options
	FullscreenButtonEnabled bool

	// Hidden will hide the window when it is first created.
	Hidden bool

	// Zoom is the zoom level of the window.
	Zoom float64

	// ZoomControlEnabled will enable the zoom control.
	ZoomControlEnabled bool

	// EnableDragAndDrop will enable drag and drop.
	EnableDragAndDrop bool

	// OpenInspectorOnStartup will open the inspector when the window is first shown.
	OpenInspectorOnStartup bool

	// Mac options
	Mac MacWindow

	// Windows options
	Windows WindowsWindow

	// ShouldClose is called when the window is about to close.
	// Return true to allow the window to close, or false to prevent it from closing.
	ShouldClose func(window *WebviewWindow) bool

	// If true, the window's devtools will be available (default true in builds without the `production` build tag)
	DevToolsEnabled bool

	// If true, the window's default context menu will be disabled (default false)
	DefaultContextMenuDisabled bool

	// KeyBindings is a map of key bindings to functions
	KeyBindings map[string]func(window *WebviewWindow)

	// IgnoreMouseEvents will ignore mouse events in the window
	IgnoreMouseEvents bool
}

var WebviewWindowDefaults = &WebviewWindowOptions{
	Title:  "",
	Width:  800,
	Height: 600,
	URL:    "",
	BackgroundColour: RGBA{
		Red:   255,
		Green: 255,
		Blue:  255,
		Alpha: 255,
	},
}

type RGBA struct {
	Red, Green, Blue, Alpha uint8
}

func NewRGBA(red, green, blue, alpha uint8) RGBA {
	return RGBA{
		Red:   red,
		Green: green,
		Blue:  blue,
		Alpha: alpha,
	}
}

func NewRGB(red, green, blue uint8) RGBA {
	return RGBA{
		Red:   red,
		Green: green,
		Blue:  blue,
		Alpha: 255,
	}
}

type BackgroundType int

const (
	BackgroundTypeSolid BackgroundType = iota
	BackgroundTypeTransparent
	BackgroundTypeTranslucent
)
