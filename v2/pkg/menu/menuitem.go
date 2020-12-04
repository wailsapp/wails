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
	Accelerator *Accelerator `json:"Accelerator,omitempty"`
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

	// This holds the menu item's parent.
	parent *MenuItem
}

// Parent returns the parent of the menu item.
// If it is a top level menu then it returns nil.
func (m *MenuItem) Parent() *MenuItem {
	return m.parent
}

// Append will attempt to append the given menu item to
// this item's submenu items. If this menu item is not a
// submenu, then this method will not add the item and
// simply return false.
func (m *MenuItem) Append(item *MenuItem) bool {
	if m.Type != SubmenuType {
		return false
	}
	m.SubMenu = append(m.SubMenu, item)
	return true
}

// Prepend will attempt to prepend the given menu item to
// this item's submenu items. If this menu item is not a
// submenu, then this method will not add the item and
// simply return false.
func (m *MenuItem) Prepend(item *MenuItem) bool {
	if m.Type != SubmenuType {
		return false
	}
	m.SubMenu = append([]*MenuItem{item}, m.SubMenu...)
	return true
}

func (m *MenuItem) getByID(id string) *MenuItem {

	// If I have the ID return me!
	if m.ID == id {
		return m
	}

	// Check submenus
	for _, submenu := range m.SubMenu {
		result := submenu.getByID(id)
		if result != nil {
			return result
		}
	}

	return nil
}

func (m *MenuItem) removeByID(id string) bool {

	for index, item := range m.SubMenu {
		if item.ID == id {
			m.SubMenu = append(m.SubMenu[:index], m.SubMenu[index+1:]...)
			return true
		}
		if item.Type == SubmenuType {
			result := item.removeByID(id)
			if result == true {
				return result
			}
		}
	}
	return false
}

// Text is a helper to create basic Text menu items
func Text(label string, id string) *MenuItem {
	return TextWithAccelerator(label, id, nil)
}

// TextWithAccelerator is a helper to create basic Text menu items with an accelerator
func TextWithAccelerator(label string, id string, accelerator *Accelerator) *MenuItem {
	return &MenuItem{
		ID:          id,
		Label:       label,
		Type:        TextType,
		Accelerator: accelerator,
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
	return RadioWithAccelerator(label, id, selected, nil)
}

// RadioWithAccelerator is a helper to create basic Radio menu items with an accelerator
func RadioWithAccelerator(label string, id string, selected bool, accelerator *Accelerator) *MenuItem {
	return &MenuItem{
		ID:          id,
		Label:       label,
		Type:        RadioType,
		Checked:     selected,
		Accelerator: accelerator,
	}
}

// Checkbox is a helper to create basic Checkbox menu items
func Checkbox(label string, id string, checked bool) *MenuItem {
	return CheckboxWithAccelerator(label, id, checked, nil)
}

// CheckboxWithAccelerator is a helper to create basic Checkbox menu items with an accelerator
func CheckboxWithAccelerator(label string, id string, checked bool, accelerator *Accelerator) *MenuItem {
	return &MenuItem{
		ID:          id,
		Label:       label,
		Type:        CheckboxType,
		Checked:     checked,
		Accelerator: accelerator,
	}
}

// SubMenu is a helper to create Submenus
func SubMenu(label string, items []*MenuItem) *MenuItem {
	result := &MenuItem{
		Label:   label,
		SubMenu: items,
		Type:    SubmenuType,
	}

	// Fix up parent pointers
	for _, item := range items {
		item.parent = result
	}

	return result
}
