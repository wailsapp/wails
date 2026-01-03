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

type WindowStartPosition int

const (
	WindowCentered WindowStartPosition = 0
	WindowXY       WindowStartPosition = 1
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

	// Initial Position
	InitialPosition WindowStartPosition

	// X is the starting X position of the window.
	X int

	// Y is the starting Y position of the window.
	Y int

	// Hidden will hide the window when it is first created.
	Hidden bool

	// Zoom is the zoom level of the window.
	Zoom float64

	// ZoomControlEnabled will enable the zoom control.
	ZoomControlEnabled bool

	// EnableFileDrop enables drag and drop of files onto the window.
	// When enabled, files dragged from the OS onto elements with the
	// `data-file-drop-target` attribute will trigger a FilesDropped event.
	EnableFileDrop bool

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

	// If true, the window's devtools will be available (default true in builds without the `production` build tag)
	DevToolsEnabled bool

	// If true, the window's default context menu will be disabled (default false)
	DefaultContextMenuDisabled bool

	// KeyBindings is a map of key bindings to functions
	KeyBindings map[string]func(window Window)

	// IgnoreMouseEvents will ignore mouse events in the window (Windows + Mac only)
	IgnoreMouseEvents bool

	// ContentProtectionEnabled specifies whether content protection is enabled, preventing screen capture and recording.
	// Effective on Windows and macOS only; no-op on Linux.
	// Best-effort protection with platform-specific caveats (see docs).
	ContentProtectionEnabled bool
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

func NewRGBPtr(red, green, blue uint8) *uint32 {
	result := uint32(red)
	result |= uint32(green) << 8
	result |= uint32(blue) << 16
	return &result
}

type BackgroundType int

const (
	BackgroundTypeSolid BackgroundType = iota
	BackgroundTypeTransparent
	BackgroundTypeTranslucent
)

/******* Windows Options *******/

type BackdropType int32

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
	CustomTheme ThemeSettings

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

	// ResizeDebounceMS is the amount of time to debounce redraws of webview2
	// when resizing the window
	// Default: 0
	ResizeDebounceMS uint16

	// WindowDidMoveDebounceMS is the amount of time to debounce the WindowDidMove event
	// when moving the window
	// Default: 0
	WindowDidMoveDebounceMS uint16

	// Event mapping for the window. This allows you to define a translation from one event to another.
	// Default: nil
	EventMapping map[events.WindowEventType]events.WindowEventType

	// HiddenOnTaskbar hides the window from the taskbar
	// Default: false
	HiddenOnTaskbar bool

	// EnableSwipeGestures enables swipe gestures for the window
	// Default: false
	EnableSwipeGestures bool

	// Menu is the menu to use for the window.
	Menu *Menu

	// Permissions map for WebView2. If empty, default permissions will be granted.
	Permissions map[CoreWebView2PermissionKind]CoreWebView2PermissionState

	// ExStyle is the extended window style
	ExStyle int

	// GeneralAutofillEnabled enables general autofill
	GeneralAutofillEnabled bool

	// PasswordAutosaveEnabled enables autosaving passwords
	PasswordAutosaveEnabled bool

	// EnabledFeatures, DisabledFeatures and AdditionalLaunchArgs are used to enable or disable specific features in the WebView2 browser.
	// Available flags: https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/webview-features-flags?tabs=dotnetcsharp#available-webview2-browser-flags
	// WARNING: Apps in production shouldn't use WebView2 browser flags,
	// because these flags might be removed or altered at any time,
	// and aren't necessarily supported long-term.
	// AdditionalLaunchArgs should always be preceded by "--"
	EnabledFeatures      []string
	DisabledFeatures     []string
	AdditionalLaunchArgs []string
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

type WindowTheme struct {
	// BorderColour is the colour of the window border
	BorderColour *uint32

	// TitleBarColour is the colour of the window title bar
	TitleBarColour *uint32

	// TitleTextColour is the colour of the window title text
	TitleTextColour *uint32
}

type TextTheme struct {
	// Text is the colour of the text
	Text *uint32

	// Background is the background colour of the text
	Background *uint32
}

type MenuBarTheme struct {
	// Default is the default theme
	Default *TextTheme

	// Hover defines the theme to use when the menu item is hovered
	Hover *TextTheme

	// Selected defines the theme to use when the menu item is selected
	Selected *TextTheme
}

// ThemeSettings defines custom colours to use in dark or light mode.
// They may be set using the hex values: 0x00BBGGRR
type ThemeSettings struct {
	// Dark mode active window
	DarkModeActive *WindowTheme

	// Dark mode inactive window
	DarkModeInactive *WindowTheme

	// Light mode active window
	LightModeActive *WindowTheme

	// Light mode inactive window
	LightModeInactive *WindowTheme

	// Dark mode MenuBar
	DarkModeMenuBar *MenuBarTheme

	// Light mode MenuBar
	LightModeMenuBar *MenuBarTheme
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
	// MacBackdropLiquidGlass - The window will use Apple's Liquid Glass effect (macOS 15.0+ with fallback to translucent)
	MacBackdropLiquidGlass
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

// MacLiquidGlassStyle defines the style of the Liquid Glass effect
type MacLiquidGlassStyle int

const (
	// LiquidGlassStyleAutomatic - System determines the best style
	LiquidGlassStyleAutomatic MacLiquidGlassStyle = iota
	// LiquidGlassStyleLight - Light glass appearance
	LiquidGlassStyleLight
	// LiquidGlassStyleDark - Dark glass appearance
	LiquidGlassStyleDark
	// LiquidGlassStyleVibrant - Vibrant glass with enhanced effects
	LiquidGlassStyleVibrant
)

// NSVisualEffectMaterial represents the NSVisualEffectMaterial enum for macOS
type NSVisualEffectMaterial int

const (
	// NSVisualEffectMaterial values from macOS SDK
	NSVisualEffectMaterialAppearanceBased       NSVisualEffectMaterial = 0
	NSVisualEffectMaterialLight                 NSVisualEffectMaterial = 1
	NSVisualEffectMaterialDark                  NSVisualEffectMaterial = 2
	NSVisualEffectMaterialTitlebar              NSVisualEffectMaterial = 3
	NSVisualEffectMaterialSelection             NSVisualEffectMaterial = 4
	NSVisualEffectMaterialMenu                  NSVisualEffectMaterial = 5
	NSVisualEffectMaterialPopover               NSVisualEffectMaterial = 6
	NSVisualEffectMaterialSidebar               NSVisualEffectMaterial = 7
	NSVisualEffectMaterialHeaderView            NSVisualEffectMaterial = 10
	NSVisualEffectMaterialSheet                 NSVisualEffectMaterial = 11
	NSVisualEffectMaterialWindowBackground      NSVisualEffectMaterial = 12
	NSVisualEffectMaterialHUDWindow             NSVisualEffectMaterial = 13
	NSVisualEffectMaterialFullScreenUI          NSVisualEffectMaterial = 15
	NSVisualEffectMaterialToolTip               NSVisualEffectMaterial = 17
	NSVisualEffectMaterialContentBackground     NSVisualEffectMaterial = 18
	NSVisualEffectMaterialUnderWindowBackground NSVisualEffectMaterial = 21
	NSVisualEffectMaterialUnderPageBackground   NSVisualEffectMaterial = 22
	NSVisualEffectMaterialAuto                  NSVisualEffectMaterial = -1 // Use auto-selection based on Style
)

// MacLiquidGlass contains configuration for the Liquid Glass effect
type MacLiquidGlass struct {
	// Style of the glass effect
	Style MacLiquidGlassStyle

	// Material to use for NSVisualEffectView (when NSGlassEffectView is not available)
	// Set to NSVisualEffectMaterialAuto to use automatic selection based on Style
	Material NSVisualEffectMaterial

	// Corner radius for the glass effect (0 for square corners)
	CornerRadius float64

	// Tint color for the glass (optional, nil for no tint)
	TintColor *RGBA

	// Group identifier for merging multiple glass windows
	GroupID string

	// Spacing between grouped glass elements (in points)
	GroupSpacing float64
}

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

	// CollectionBehavior controls how the window behaves across macOS Spaces and fullscreen
	CollectionBehavior MacWindowCollectionBehavior

	// LiquidGlass contains configuration for the Liquid Glass effect
	LiquidGlass MacLiquidGlass
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

// MacWindowCollectionBehavior controls window behavior across macOS Spaces and fullscreen.
// These correspond to NSWindowCollectionBehavior bitmask values and can be combined using bitwise OR.
// For example: MacWindowCollectionBehaviorCanJoinAllSpaces | MacWindowCollectionBehaviorFullScreenAuxiliary
type MacWindowCollectionBehavior int

const (
	// MacWindowCollectionBehaviorDefault is zero value - when set, FullScreenPrimary is used for backwards compatibility
	MacWindowCollectionBehaviorDefault MacWindowCollectionBehavior = 0
	// MacWindowCollectionBehaviorCanJoinAllSpaces allows window to appear on all Spaces
	MacWindowCollectionBehaviorCanJoinAllSpaces MacWindowCollectionBehavior = 1 << 0 // 1
	// MacWindowCollectionBehaviorMoveToActiveSpace moves window to active Space when shown
	MacWindowCollectionBehaviorMoveToActiveSpace MacWindowCollectionBehavior = 1 << 1 // 2
	// MacWindowCollectionBehaviorManaged is the default managed window behavior
	MacWindowCollectionBehaviorManaged MacWindowCollectionBehavior = 1 << 2 // 4
	// MacWindowCollectionBehaviorTransient marks window as temporary/transient
	MacWindowCollectionBehaviorTransient MacWindowCollectionBehavior = 1 << 3 // 8
	// MacWindowCollectionBehaviorStationary keeps window stationary during Space switches
	MacWindowCollectionBehaviorStationary MacWindowCollectionBehavior = 1 << 4 // 16
	// MacWindowCollectionBehaviorParticipatesInCycle includes window in Cmd+` cycling (default for normal windows)
	MacWindowCollectionBehaviorParticipatesInCycle MacWindowCollectionBehavior = 1 << 5 // 32
	// MacWindowCollectionBehaviorIgnoresCycle excludes window from Cmd+` cycling
	MacWindowCollectionBehaviorIgnoresCycle MacWindowCollectionBehavior = 1 << 6 // 64
	// MacWindowCollectionBehaviorFullScreenPrimary allows the window to enter fullscreen
	MacWindowCollectionBehaviorFullScreenPrimary MacWindowCollectionBehavior = 1 << 7 // 128
	// MacWindowCollectionBehaviorFullScreenAuxiliary allows window to overlay fullscreen apps
	MacWindowCollectionBehaviorFullScreenAuxiliary MacWindowCollectionBehavior = 1 << 8 // 256
	// MacWindowCollectionBehaviorFullScreenNone prevents window from entering fullscreen (macOS 10.7+)
	MacWindowCollectionBehaviorFullScreenNone MacWindowCollectionBehavior = 1 << 9 // 512
	// MacWindowCollectionBehaviorFullScreenAllowsTiling allows side-by-side tiling in fullscreen (macOS 10.11+)
	MacWindowCollectionBehaviorFullScreenAllowsTiling MacWindowCollectionBehavior = 1 << 11 // 2048
	// MacWindowCollectionBehaviorFullScreenDisallowsTiling prevents tiling in fullscreen (macOS 10.11+)
	MacWindowCollectionBehaviorFullScreenDisallowsTiling MacWindowCollectionBehavior = 1 << 12 // 4096
)

// MacWebviewPreferences contains preferences for the Mac webview
type MacWebviewPreferences struct {
	// TabFocusesLinks will enable tabbing to links
	TabFocusesLinks u.Bool
	// TextInteractionEnabled will enable text interaction
	TextInteractionEnabled u.Bool
	// FullscreenEnabled will enable fullscreen
	FullscreenEnabled u.Bool
	// AllowsBackForwardNavigationGestures enables horizontal swipe gestures for back/forward navigation
	AllowsBackForwardNavigationGestures u.Bool
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

	// WindowDidMoveDebounceMS is the debounce time in milliseconds for the WindowDidMove event
	WindowDidMoveDebounceMS uint16

	// Menu is the window's menu
	Menu *Menu
}
