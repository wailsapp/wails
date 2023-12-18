package menu

import (
	"sync"

	"github.com/wailsapp/wails/v2/pkg/menu/keys"
)

// MenuItem represents a menuitem contained in a menu
type MenuItem struct {
	// Label is what appears as the menu text
	Label string
	// Role is a predefined menu type
	Role Role
	// Accelerator holds a representation of a key binding
	Accelerator *keys.Accelerator
	// Type of MenuItem, EG: Checkbox, Text, Separator, Radio, Submenu
	Type Type
	// Disabled makes the item unselectable
	Disabled bool
	// Hidden ensures that the item is not shown in the menu
	Hidden bool
	// Checked indicates if the item is selected (used by Checkbox and Radio types only)
	Checked bool
	// Submenu contains a list of menu items that will be shown as a submenu
	// SubMenu []*MenuItem `json:"SubMenu,omitempty"`
	SubMenu *Menu

	// Callback function when menu clicked
	Click Callback
	/*
		// Text Colour
		RGBA string

		// Font
		FontSize int
		FontName string

		// Image - base64 image data
		Image string

		// MacTemplateImage indicates that on a Mac, this image is a template image
		MacTemplateImage bool

		// MacAlternate indicates that this item is an alternative to the previous menu item
		MacAlternate bool

		// Tooltip
		Tooltip string
	*/
	// This holds the menu item's parent.
	parent *MenuItem

	// Used for locking when removing elements
	removeLock sync.Mutex
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
	if !m.isSubMenu() {
		return false
	}
	item.parent = m
	m.SubMenu.Append(item)
	return true
}

// Prepend will attempt to prepend the given menu item to
// this item's submenu items. If this menu item is not a
// submenu, then this method will not add the item and
// simply return false.
func (m *MenuItem) Prepend(item *MenuItem) bool {
	if !m.isSubMenu() {
		return false
	}
	item.parent = m
	m.SubMenu.Prepend(item)
	return true
}

func (m *MenuItem) Remove() {
	// Iterate my parent's children
	m.Parent().removeChild(m)
}

func (m *MenuItem) removeChild(item *MenuItem) {
	m.removeLock.Lock()
	for index, child := range m.SubMenu.Items {
		if item == child {
			m.SubMenu.Items = append(m.SubMenu.Items[:index], m.SubMenu.Items[index+1:]...)
		}
	}
	m.removeLock.Unlock()
}

// InsertAfter attempts to add the given item after this item in the parent
// menu. If there is no parent menu (we are a top level menu) then false is
// returned
func (m *MenuItem) InsertAfter(item *MenuItem) bool {
	// We need to find my parent
	if m.parent == nil {
		return false
	}

	// Get my parent to insert the item
	return m.parent.insertNewItemAfterGivenItem(m, item)
}

// InsertBefore attempts to add the given item before this item in the parent
// menu. If there is no parent menu (we are a top level menu) then false is
// returned
func (m *MenuItem) InsertBefore(item *MenuItem) bool {
	// We need to find my parent
	if m.parent == nil {
		return false
	}

	// Get my parent to insert the item
	return m.parent.insertNewItemBeforeGivenItem(m, item)
}

// insertNewItemAfterGivenItem will insert the given item after the given target
// in this item's submenu. If we are not a submenu,
// then something bad has happened :/
func (m *MenuItem) insertNewItemAfterGivenItem(target *MenuItem,
	newItem *MenuItem,
) bool {
	if !m.isSubMenu() {
		return false
	}

	// Find the index of the target
	targetIndex := m.getItemIndex(target)
	if targetIndex == -1 {
		return false
	}

	// Insert element into slice
	return m.insertItemAtIndex(targetIndex+1, newItem)
}

// insertNewItemBeforeGivenItem will insert the given item before the given
// target in this item's submenu. If we are not a submenu, then something bad
// has happened :/
func (m *MenuItem) insertNewItemBeforeGivenItem(target *MenuItem,
	newItem *MenuItem,
) bool {
	if !m.isSubMenu() {
		return false
	}

	// Find the index of the target
	targetIndex := m.getItemIndex(target)
	if targetIndex == -1 {
		return false
	}

	// Insert element into slice
	return m.insertItemAtIndex(targetIndex, newItem)
}

func (m *MenuItem) isSubMenu() bool {
	return m.Type == SubmenuType
}

// getItemIndex returns the index of the given target relative to this menu
func (m *MenuItem) getItemIndex(target *MenuItem) int {
	// This should only be called on submenus
	if !m.isSubMenu() {
		return -1
	}

	// hunt down that bad boy
	for index, item := range m.SubMenu.Items {
		if item == target {
			return index
		}
	}

	return -1
}

// insertItemAtIndex attempts to insert the given item into the submenu at
// the given index
// Credit: https://stackoverflow.com/a/61822301
func (m *MenuItem) insertItemAtIndex(index int, target *MenuItem) bool {
	// If index is OOB, return false
	if index > len(m.SubMenu.Items) {
		return false
	}

	// Save parent reference
	target.parent = m

	// If index is last item, then just regular append
	if index == len(m.SubMenu.Items) {
		m.SubMenu.Items = append(m.SubMenu.Items, target)
		return true
	}

	m.SubMenu.Items = append(m.SubMenu.Items[:index+1], m.SubMenu.Items[index:]...)
	m.SubMenu.Items[index] = target
	return true
}

func (m *MenuItem) SetLabel(name string) {
	if m.Label == name {
		return
	}
	m.Label = name
}

func (m *MenuItem) IsSeparator() bool {
	return m.Type == SeparatorType
}

func (m *MenuItem) IsCheckbox() bool {
	return m.Type == CheckboxType
}

func (m *MenuItem) Disable() *MenuItem {
	m.Disabled = true
	return m
}

func (m *MenuItem) Enable() *MenuItem {
	m.Disabled = false
	return m
}

func (m *MenuItem) OnClick(click Callback) *MenuItem {
	m.Click = click
	return m
}

func (m *MenuItem) SetAccelerator(acc *keys.Accelerator) *MenuItem {
	m.Accelerator = acc
	return m
}

func (m *MenuItem) SetChecked(value bool) *MenuItem {
	m.Checked = value
	if m.Type != RadioType {
		m.Type = CheckboxType
	}
	return m
}

func (m *MenuItem) Hide() *MenuItem {
	m.Hidden = true
	return m
}

func (m *MenuItem) Show() *MenuItem {
	m.Hidden = false
	return m
}

func (m *MenuItem) IsRadio() bool {
	return m.Type == RadioType
}

func Label(label string) *MenuItem {
	return &MenuItem{
		Type:  TextType,
		Label: label,
	}
}

// Text is a helper to create basic Text menu items
func Text(label string, accelerator *keys.Accelerator, click Callback) *MenuItem {
	return &MenuItem{
		Label:       label,
		Type:        TextType,
		Accelerator: accelerator,
		Click:       click,
	}
}

// Separator provides a menu separator
func Separator() *MenuItem {
	return &MenuItem{
		Type: SeparatorType,
	}
}

// Radio is a helper to create basic Radio menu items with an accelerator
func Radio(label string, selected bool, accelerator *keys.Accelerator, click Callback) *MenuItem {
	return &MenuItem{
		Label:       label,
		Type:        RadioType,
		Checked:     selected,
		Accelerator: accelerator,
		Click:       click,
	}
}

// Checkbox is a helper to create basic Checkbox menu items
func Checkbox(label string, checked bool, accelerator *keys.Accelerator, click Callback) *MenuItem {
	return &MenuItem{
		Label:       label,
		Type:        CheckboxType,
		Checked:     checked,
		Accelerator: accelerator,
		Click:       click,
	}
}

// SubMenu is a helper to create Submenus
func SubMenu(label string, menu *Menu) *MenuItem {
	result := &MenuItem{
		Label:   label,
		SubMenu: menu,
		Type:    SubmenuType,
	}

	menu.setParent(result)

	return result
}
