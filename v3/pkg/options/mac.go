package options

type ActivationPolicy int

const (
	ActivationPolicyRegular ActivationPolicy = iota
	// ActivationPolicyAccessory is used for applications that do not have a main window,
	// such as system tray applications or background applications.
	ActivationPolicyAccessory
	ActivationPolicyProhibited
)

type Mac struct {
	// ActivationPolicy is the activation policy for the application. Defaults to
	// applicationActivationPolicyRegular.
	ActivationPolicy ActivationPolicy
	// If set to true, the application will terminate when the last window is closed.
	ApplicationShouldTerminateAfterLastWindowClosed bool
}

type MacBackdrop int

const (
	MacBackdropNormal MacBackdrop = iota
	MacBackdropTransparent
	MacBackdropTranslucent
)

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

// MacWindow contains macOS specific options
type MacWindow struct {
	Backdrop                MacBackdrop
	TitleBar                TitleBar
	Appearance              MacAppearanceType
	InvisibleTitleBarHeight int
}

// TitleBar contains options for the Mac titlebar
type TitleBar struct {
	AppearsTransparent   bool
	Hide                 bool
	HideTitle            bool
	FullSizeContent      bool
	UseToolbar           bool
	HideToolbarSeparator bool
	ToolbarStyle         MacToolbarStyle
}

// TitleBarDefault results in the default Mac TitleBar
var TitleBarDefault = TitleBar{
	AppearsTransparent:   false,
	Hide:                 false,
	HideTitle:            false,
	FullSizeContent:      false,
	UseToolbar:           false,
	HideToolbarSeparator: false,
}

// Credit: Comments from Electron site

// TitleBarHidden results in a hidden title bar and a full size content window,
// yet the title bar still has the standard window controls (“traffic lights”)
// in the top left.
var TitleBarHidden = TitleBar{
	AppearsTransparent:   true,
	Hide:                 false,
	HideTitle:            true,
	FullSizeContent:      true,
	UseToolbar:           false,
	HideToolbarSeparator: false,
}

// TitleBarHiddenInset results in a hidden title bar with an alternative look where
// the traffic light buttons are slightly more inset from the window edge.
var TitleBarHiddenInset = TitleBar{
	AppearsTransparent:   true,
	Hide:                 false,
	HideTitle:            true,
	FullSizeContent:      true,
	UseToolbar:           true,
	HideToolbarSeparator: true,
}

// TitleBarHiddenInsetUnified results in a hidden title bar with an alternative look where
// the traffic light buttons are even more inset from the window edge.
var TitleBarHiddenInsetUnified = TitleBar{
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
