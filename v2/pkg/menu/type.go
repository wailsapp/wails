package menu

// Type of the menu item
type Type string

const (
	// NormalType is the Normal menuitem type
	NormalType Type = "Normal"
	// SeparatorType is the Separator menuitem type
	SeparatorType Type = "Separator"
	// SubmenuType is the Submenu menuitem type
	SubmenuType Type = "Submenu"
	// CheckboxType is the Checkbox menuitem type
	CheckboxType Type = "Checkbox"
	// RadioType is the Radio menuitem type
	RadioType Type = "Radio"
)
