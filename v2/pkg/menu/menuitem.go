package menu

// MenuItem represents a menuitem contained in a menu
type MenuItem struct {
	// The unique identifier of this menu item
	ID string `json:"ID,omitempty"`
	// Label is what appears as the menu text
	Label string
	// Role is a predefined menu type
	Role Role `json:"Role,omitempty"`
	// Accelerator holds a representation of a key binding
	Accelerator string `json:"Accelerator,omitempty"`
	// Type of MenuItem, EG: Checkbox, Text, Separator, Radio, Submenu
	Type Type
	// Disabled makes the item unselectable
	Disabled bool
	// Hidden ensures that the item is not shown in the menu
	Hidden bool
	// Checked indicates if the item is selected (used by Checkbox and Radio types only)
	Checked bool
	// Submenu contains a list of menu items that will be shown as a submenu
	SubMenu []*MenuItem `json:"SubMenu,omitempty"`
}

// Text is a helper to create basic Text menu items
func Text(label string, id string) *MenuItem {
	return &MenuItem{
		ID:    id,
		Label: label,
		Type:  TextType,
	}
}

// Separator provides a menu separator
func Separator() *MenuItem {
	return &MenuItem{
		Type: SeparatorType,
	}
}

// Radio is a helper to create basic Radio menu items
func Radio(label string, id string, selected bool) *MenuItem {
	return &MenuItem{
		ID:      id,
		Label:   label,
		Type:    RadioType,
		Checked: selected,
	}
}

// Checkbox is a helper to create basic Checkbox menu items
func Checkbox(label string, id string, checked bool) *MenuItem {
	return &MenuItem{
		ID:      id,
		Label:   label,
		Type:    CheckboxType,
		Checked: checked,
	}
}

// SubMenu is a helper to create Submenus
func SubMenu(label string, items []*MenuItem) *MenuItem {
	return &MenuItem{
		Label:   label,
		SubMenu: items,
		Type:    SubmenuType,
	}
}
