package options

type ActivationPolicy int

const (
	ActivationPolicyRegular ActivationPolicy = iota
	ActivationPolicyAccessory
	ActivationPolicyProhibited
)

type Mac struct {
	// ActivationPolicy is the activation policy for the application. Defaults to
	// applicationActivationPolicyRegular.
	ActivationPolicy ActivationPolicy
}

type MacBackdrop int

const (
	MacBackdropNormal MacBackdrop = iota
	MacBackdropTransparent
	MacBackdropTranslucent
)

// MacWindow contains macOS specific options
type MacWindow struct {
	Backdrop   MacBackdrop
	TitleBar   *TitleBar
	Appearance MacAppearanceType
}

// TitleBar contains options for the Mac titlebar
type TitleBar struct {
	AppearsTransparent   bool
	Hide                 bool
	HideTitle            bool
	FullSizeContent      bool
	UseToolbar           bool
	HideToolbarSeparator bool
}

// TitleBarDefault results in the default Mac Titlebar
func TitleBarDefault() *TitleBar {
	return &TitleBar{
		AppearsTransparent:   false,
		Hide:                 false,
		HideTitle:            false,
		FullSizeContent:      false,
		UseToolbar:           false,
		HideToolbarSeparator: false,
	}
}

// Credit: Comments from Electron site

// TitleBarHidden results in a hidden title bar and a full size content window,
// yet the title bar still has the standard window controls (“traffic lights”)
// in the top left.
func TitleBarHidden() *TitleBar {
	return &TitleBar{
		AppearsTransparent:   true,
		Hide:                 false,
		HideTitle:            true,
		FullSizeContent:      true,
		UseToolbar:           false,
		HideToolbarSeparator: false,
	}
}

// TitleBarHiddenInset results in a hidden title bar with an alternative look where
// the traffic light buttons are slightly more inset from the window edge.
func TitleBarHiddenInset() *TitleBar {

	return &TitleBar{
		AppearsTransparent:   true,
		Hide:                 false,
		HideTitle:            true,
		FullSizeContent:      true,
		UseToolbar:           true,
		HideToolbarSeparator: true,
	}

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
