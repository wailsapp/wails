package application

import (
	"github.com/leaanthony/u"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type WindowState int

const (
	WindowStateNormal WindowState = iota
	WindowStateMinimised
	WindowStateMaximised
	WindowStateFullscreen
)

type ButtonState int

const (
	ButtonEnabled  ButtonState = 0
	ButtonDisabled ButtonState = 1
	ButtonHidden   ButtonState = 2
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

	// Linux options
	Linux LinuxWindow

	// Toolbar button states
	MinimiseButtonState ButtonState
	MaximiseButtonState ButtonState
	CloseButtonState    ButtonState

	// ShouldClose is called when the window is about to close.
	// Return true to allow the window to close, or false to prevent it from closing.
	ShouldClose func(window *WebviewWindow) bool

	// If true, the window's devtools will be available (default true in builds without the `production` build tag)
	DevToolsEnabled bool

	// If true, the window's default context menu will be disabled (default false)
	DefaultContextMenuDisabled bool

	// KeyBindings is a map of key bindings to functions
	KeyBindings map[string]func(window *WebviewWindow)

	// IgnoreMouseEvents will ignore mouse events in the window (Windows + Mac only)
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

/******* Windows Options *******/

type BackdropType int32
type DragEffect int32

const (
	// DragEffectNone is used to indicate that the drop target cannot accept the data.
	DragEffectNone DragEffect = 1
	// DragEffectCopy is used to indicate that the data is copied to the drop target.
	DragEffectCopy DragEffect = 2
	// DragEffectMove is used to indicate that the data is removed from the drag source.
	DragEffectMove DragEffect = 3
	// DragEffectLink is used to indicate that a link to the original data is established.
	DragEffectLink DragEffect = 4
	// DragEffectScroll is used to indicate that the target can be scrolled while dragging to locate a drop position that is not currently visible in the target.

)

const (
	Auto    BackdropType = 0
	None    BackdropType = 1
	Mica    BackdropType = 2
	Acrylic BackdropType = 3
	Tabbed  BackdropType = 4
)

type CoreWebView2PermissionKind uint32

const (
	CoreWebView2PermissionKindUnknownPermission CoreWebView2PermissionKind = iota
	CoreWebView2PermissionKindMicrophone
	CoreWebView2PermissionKindCamera
	CoreWebView2PermissionKindGeolocation
	CoreWebView2PermissionKindNotifications
	CoreWebView2PermissionKindOtherSensors
	CoreWebView2PermissionKindClipboardRead
)

type CoreWebView2PermissionState uint32

const (
	CoreWebView2PermissionStateDefault CoreWebView2PermissionState = iota
	CoreWebView2PermissionStateAllow
	CoreWebView2PermissionStateDeny
)

type WindowsWindow struct {
	// Select the type of translucent backdrop. Requires Windows 11 22621 or later.
	// Only used when window's `BackgroundType` is set to `BackgroundTypeTranslucent`.
	// Default: Auto
	BackdropType BackdropType

	// Disable the icon in the titlebar
	// Default: false
	DisableIcon bool

	// Theme (Dark / Light / SystemDefault)
	// Default: SystemDefault - The application will follow system theme changes.
	Theme Theme

	// Specify custom colours to use for dark/light mode
	// Default: nil
	CustomTheme *ThemeSettings

	// Disable all window decorations in Frameless mode, which means no "Aero Shadow" and no "Rounded Corner" will be shown.
	// "Rounded Corners" are only available on Windows 11.
	// Default: false
	DisableFramelessWindowDecorations bool

	// WindowMask is used to set the window shape. Use a PNG with an alpha channel to create a custom shape.
	// Default: nil
	WindowMask []byte

	// WindowMaskDraggable is used to make the window draggable by clicking on the window mask.
	// Default: false
	WindowMaskDraggable bool

	// WebviewGpuIsDisabled is used to enable / disable GPU acceleration for the webview
	// Default: false
	WebviewGpuIsDisabled bool

	// ResizeDebounceMS is the amount of time to debounce redraws of webview2
	// when resizing the window
	// Default: 0
	ResizeDebounceMS uint16

	// Disable the menu bar for this window
	// Default: false
	DisableMenu bool

	// Event mapping for the window. This allows you to define a translation from one event to another.
	// Default: nil
	EventMapping map[events.WindowEventType]events.WindowEventType

	// HiddenOnTaskbar hides the window from the taskbar
	// Default: false
	HiddenOnTaskbar bool

	// EnableSwipeGestures enables swipe gestures for the window
	// Default: false
	EnableSwipeGestures bool

	// EnableFraudulentWebsiteWarnings will enable warnings for fraudulent websites.
	// Default: false
	EnableFraudulentWebsiteWarnings bool

	// Menu is the menu to use for the window.
	Menu *Menu

	// Drag Cursor Effects
	OnEnterEffect DragEffect
	OnOverEffect  DragEffect

	// Permissions map for WebView2. If empty, default permissions will be granted.
	Permissions map[CoreWebView2PermissionKind]CoreWebView2PermissionState

	// ExStyle is the extended window style
	ExStyle int
}

type Theme int

const (
	// SystemDefault will use whatever the system theme is. The application will follow system theme changes.
	SystemDefault Theme = 0
	// Dark Mode
	Dark Theme = 1
	// Light Mode
	Light Theme = 2
)

// ThemeSettings defines custom colours to use in dark or light mode.
// They may be set using the hex values: 0x00BBGGRR
type ThemeSettings struct {
	DarkModeTitleBar           int32
	DarkModeTitleBarInactive   int32
	DarkModeTitleText          int32
	DarkModeTitleTextInactive  int32
	DarkModeBorder             int32
	DarkModeBorderInactive     int32
	LightModeTitleBar          int32
	LightModeTitleBarInactive  int32
	LightModeTitleText         int32
	LightModeTitleTextInactive int32
	LightModeBorder            int32
	LightModeBorderInactive    int32
}

/****** Mac Options *******/

// MacBackdrop is the backdrop type for macOS
type MacBackdrop int

const (
	// MacBackdropNormal - The default value. The window will have a normal opaque background.
	MacBackdropNormal MacBackdrop = iota
	// MacBackdropTransparent - The window will have a transparent background, with the content underneath it being visible
	MacBackdropTransparent
	// MacBackdropTranslucent - The window will have a translucent background, with the content underneath it being "fuzzy" or "frosted"
	MacBackdropTranslucent
)

// MacToolbarStyle is the style of toolbar for macOS
type MacToolbarStyle int

const (
	// MacToolbarStyleAutomatic - The default value. The style will be determined by the window's given configuration
	MacToolbarStyleAutomatic MacToolbarStyle = iota
	// MacToolbarStyleExpanded - The toolbar will appear below the window title
	MacToolbarStyleExpanded
	// MacToolbarStylePreference - The toolbar will appear below the window title and the items in the toolbar will attempt to have equal widths when possible
	MacToolbarStylePreference
	// MacToolbarStyleUnified - The window title will appear inline with the toolbar when visible
	MacToolbarStyleUnified
	// MacToolbarStyleUnifiedCompact - Same as MacToolbarStyleUnified, but with reduced margins in the toolbar allowing more focus to be on the contents of the window
	MacToolbarStyleUnifiedCompact
)

// MacWindow contains macOS specific options for Webview Windows
type MacWindow struct {
	// Backdrop is the backdrop type for the window
	Backdrop MacBackdrop
	// DisableShadow will disable the window shadow
	DisableShadow bool
	// TitleBar contains options for the Mac titlebar
	TitleBar MacTitleBar
	// Appearance is the appearance type for the window
	Appearance MacAppearanceType
	// InvisibleTitleBarHeight defines the height of an invisible titlebar which responds to dragging
	InvisibleTitleBarHeight int
	// Maps events from platform specific to common event types
	EventMapping map[events.WindowEventType]events.WindowEventType

	// EnableFraudulentWebsiteWarnings will enable warnings for fraudulent websites.
	// Default: false
	EnableFraudulentWebsiteWarnings bool

	// WebviewPreferences contains preferences for the webview
	WebviewPreferences MacWebviewPreferences

	// WindowLevel sets the window level to control the order of windows in the screen
	WindowLevel MacWindowLevel
}

type MacWindowLevel string

const (
	MacWindowLevelNormal      MacWindowLevel = "normal"
	MacWindowLevelFloating    MacWindowLevel = "floating"
	MacWindowLevelTornOffMenu MacWindowLevel = "tornOffMenu"
	MacWindowLevelModalPanel  MacWindowLevel = "modalPanel"
	MacWindowLevelMainMenu    MacWindowLevel = "mainMenu"
	MacWindowLevelStatus      MacWindowLevel = "status"
	MacWindowLevelPopUpMenu   MacWindowLevel = "popUpMenu"
	MacWindowLevelScreenSaver MacWindowLevel = "screenSaver"
)

// MacWebviewPreferences contains preferences for the Mac webview
type MacWebviewPreferences struct {
	// TabFocusesLinks will enable tabbing to links
	TabFocusesLinks u.Bool
	// TextInteractionEnabled will enable text interaction
	TextInteractionEnabled u.Bool
	// FullscreenEnabled will enable fullscreen
	FullscreenEnabled u.Bool
}

// MacTitleBar contains options for the Mac titlebar
type MacTitleBar struct {
	// AppearsTransparent will make the titlebar transparent
	AppearsTransparent bool
	// Hide will hide the titlebar
	Hide bool
	// HideTitle will hide the title
	HideTitle bool
	// FullSizeContent will extend the window content to the full size of the window
	FullSizeContent bool
	// UseToolbar will use a toolbar instead of a titlebar
	UseToolbar bool
	// HideToolbarSeparator will hide the toolbar separator
	HideToolbarSeparator bool
	// ShowToolbarWhenFullscreen will keep the toolbar visible when the window is in fullscreen mode
	ShowToolbarWhenFullscreen bool
	// ToolbarStyle is the style of toolbar to use
	ToolbarStyle MacToolbarStyle
}

// MacTitleBarDefault results in the default Mac MacTitleBar
var MacTitleBarDefault = MacTitleBar{
	AppearsTransparent:   false,
	Hide:                 false,
	HideTitle:            false,
	FullSizeContent:      false,
	UseToolbar:           false,
	HideToolbarSeparator: false,
}

// Credit: Comments from Electron site

// MacTitleBarHidden results in a hidden title bar and a full size content window,
// yet the title bar still has the standard window controls (“traffic lights”)
// in the top left.
var MacTitleBarHidden = MacTitleBar{
	AppearsTransparent:   true,
	Hide:                 false,
	HideTitle:            true,
	FullSizeContent:      true,
	UseToolbar:           false,
	HideToolbarSeparator: false,
}

// MacTitleBarHiddenInset results in a hidden title bar with an alternative look where
// the traffic light buttons are slightly more inset from the window edge.
var MacTitleBarHiddenInset = MacTitleBar{
	AppearsTransparent:   true,
	Hide:                 false,
	HideTitle:            true,
	FullSizeContent:      true,
	UseToolbar:           true,
	HideToolbarSeparator: true,
}

// MacTitleBarHiddenInsetUnified results in a hidden title bar with an alternative look where
// the traffic light buttons are even more inset from the window edge.
var MacTitleBarHiddenInsetUnified = MacTitleBar{
	AppearsTransparent:   true,
	Hide:                 false,
	HideTitle:            true,
	FullSizeContent:      true,
	UseToolbar:           true,
	HideToolbarSeparator: true,
	ToolbarStyle:         MacToolbarStyleUnified,
}

// MacAppearanceType is a type of Appearance for Cocoa windows
type MacAppearanceType string

const (
	// DefaultAppearance uses the default system value
	DefaultAppearance MacAppearanceType = ""
	// NSAppearanceNameAqua - The standard light system appearance.
	NSAppearanceNameAqua MacAppearanceType = "NSAppearanceNameAqua"
	// NSAppearanceNameDarkAqua - The standard dark system appearance.
	NSAppearanceNameDarkAqua MacAppearanceType = "NSAppearanceNameDarkAqua"
	// NSAppearanceNameVibrantLight - The light vibrant appearance
	NSAppearanceNameVibrantLight MacAppearanceType = "NSAppearanceNameVibrantLight"
	// NSAppearanceNameAccessibilityHighContrastAqua - A high-contrast version of the standard light system appearance.
	NSAppearanceNameAccessibilityHighContrastAqua MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastAqua"
	// NSAppearanceNameAccessibilityHighContrastDarkAqua - A high-contrast version of the standard dark system appearance.
	NSAppearanceNameAccessibilityHighContrastDarkAqua MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastDarkAqua"
	// NSAppearanceNameAccessibilityHighContrastVibrantLight - A high-contrast version of the light vibrant appearance.
	NSAppearanceNameAccessibilityHighContrastVibrantLight MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastVibrantLight"
	// NSAppearanceNameAccessibilityHighContrastVibrantDark - A high-contrast version of the dark vibrant appearance.
	NSAppearanceNameAccessibilityHighContrastVibrantDark MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastVibrantDark"
)

/******** Linux Options ********/

// WebviewGpuPolicy values used for determining the webview's hardware acceleration policy.
type WebviewGpuPolicy int

const (
	// WebviewGpuPolicyAlways Hardware acceleration is always enabled.
	WebviewGpuPolicyAlways WebviewGpuPolicy = iota
	// WebviewGpuPolicyOnDemand Hardware acceleration is enabled/disabled as request by web contents.
	WebviewGpuPolicyOnDemand
	// WebviewGpuPolicyNever Hardware acceleration is always disabled.
	WebviewGpuPolicyNever
)

// LinuxWindow specific to Linux windows
type LinuxWindow struct {
	// Icon Sets up the icon representing the window. This icon is used when the window is minimized
	// (also known as iconified).
	Icon []byte

	// WindowIsTranslucent sets the window's background to transparent when enabled.
	WindowIsTranslucent bool

	// WebviewGpuPolicy used for determining the hardware acceleration policy for the webview.
	//   - WebviewGpuPolicyAlways
	//   - WebviewGpuPolicyOnDemand
	//   - WebviewGpuPolicyNever
	//
	// Due to https://github.com/wailsapp/wails/issues/2977, if options.Linux is nil
	// in the call to wails.Run(), WebviewGpuPolicy is set by default to WebviewGpuPolicyNever.
	// Client code may override this behavior by passing a non-nil Options and set
	// WebviewGpuPolicy as needed.
	WebviewGpuPolicy WebviewGpuPolicy
}
