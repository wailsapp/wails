package mac

// AppearanceType is a type of Appearance for Cocoa windows
type AppearanceType string

const (
	// DefaultAppearance uses the default system value
	DefaultAppearance AppearanceType = ""
	// NSAppearanceNameAqua - The standard light system appearance.
	NSAppearanceNameAqua AppearanceType = "NSAppearanceNameAqua"
	// NSAppearanceNameDarkAqua - The standard dark system appearance.
	NSAppearanceNameDarkAqua AppearanceType = "NSAppearanceNameDarkAqua"
	// NSAppearanceNameVibrantLight - The light vibrant appearance
	NSAppearanceNameVibrantLight AppearanceType = "NSAppearanceNameVibrantLight"
	// NSAppearanceNameAccessibilityHighContrastAqua - A high-contrast version of the standard light system appearance.
	NSAppearanceNameAccessibilityHighContrastAqua AppearanceType = "NSAppearanceNameAccessibilityHighContrastAqua"
	// NSAppearanceNameAccessibilityHighContrastDarkAqua - A high-contrast version of the standard dark system appearance.
	NSAppearanceNameAccessibilityHighContrastDarkAqua AppearanceType = "NSAppearanceNameAccessibilityHighContrastDarkAqua"
	// NSAppearanceNameAccessibilityHighContrastVibrantLight - A high-contrast version of the light vibrant appearance.
	NSAppearanceNameAccessibilityHighContrastVibrantLight AppearanceType = "NSAppearanceNameAccessibilityHighContrastVibrantLight"
	// NSAppearanceNameAccessibilityHighContrastVibrantDark - A high-contrast version of the dark vibrant appearance.
	NSAppearanceNameAccessibilityHighContrastVibrantDark AppearanceType = "NSAppearanceNameAccessibilityHighContrastVibrantDark"
)
