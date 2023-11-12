package menu

import "github.com/wailsapp/wails/v2/pkg/menu/keys"

type Menu struct {
	Items []*MenuItem
}

func NewMenu() *Menu {
	return &Menu{}
}

func (m *Menu) Append(item *MenuItem) {
	m.Items = append(m.Items, item)
}

// Merge will append the items in the given menu
// into this menu
func (m *Menu) Merge(menu *Menu) {
	m.Items = append(m.Items, menu.Items...)
}

// AddText adds a TextMenu item to the menu
func (m *Menu) AddText(label string, accelerator *keys.Accelerator, click Callback) *MenuItem {
	item := Text(label, accelerator, click)
	m.Append(item)
	return item
}

// AddCheckbox adds a CheckboxMenu item to the menu
func (m *Menu) AddCheckbox(label string, checked bool, accelerator *keys.Accelerator, click Callback) *MenuItem {
	item := Checkbox(label, checked, accelerator, click)
	m.Append(item)
	return item
}

// AddRadio adds a radio item to the menu
func (m *Menu) AddRadio(label string, checked bool, accelerator *keys.Accelerator, click Callback) *MenuItem {
	item := Radio(label, checked, accelerator, click)
	m.Append(item)
	return item
}

// AddSeparator adds a separator to the menu
func (m *Menu) AddSeparator() {
	item := Separator()
	m.Append(item)
}

func (m *Menu) AddSubmenu(label string) *Menu {
	submenu := NewMenu()
	item := SubMenu(label, submenu)
	m.Append(item)
	return submenu
}

func (m *Menu) Prepend(item *MenuItem) {
	m.Items = append([]*MenuItem{item}, m.Items...)
}

func NewMenuFromItems(first *MenuItem, rest ...*MenuItem) *Menu {
	result := NewMenu()
	result.Append(first)
	for _, item := range rest {
		result.Append(item)
	}

	return result
}

func (m *Menu) setParent(menuItem *MenuItem) {
	for _, item := range m.Items {
		item.parent = menuItem
	}
}
