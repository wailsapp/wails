package application

import (
	"github.com/leaanthony/u"
	"github.com/wailsapp/wails/v3/pkg/events"
)

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

	// Disables the minimise button
	DisableMinimiseButton bool

	// Disables the maximise button
	DisableMaximiseButton bool

	// Disables the close button
	DisableCloseButton bool
}

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
