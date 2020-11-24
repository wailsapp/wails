package menu

// Type of the menu item
type Type string

const (
	// TextType is the text menuitem type
	TextType Type = "Text"
	// SeparatorType is the Separator menuitem type
	SeparatorType Type = "Separator"
	// SubmenuType is the Submenu menuitem type
	SubmenuType Type = "Submenu"
	// CheckboxType is the Checkbox menuitem type
	CheckboxType Type = "Checkbox"
	// RadioType is the Radio menuitem type
	RadioType Type = "Radio"
)
