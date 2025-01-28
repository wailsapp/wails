package application

type menuImpl interface {
	update()
}

type ContextMenu struct {
	*Menu
	name string
}

func NewContextMenu(name string) *ContextMenu {
	result := &ContextMenu{
		Menu: NewMenu(),
		name: name,
	}
	result.Update()
	return result
}

func (m *ContextMenu) Update() {
	m.Menu.Update()
	globalApplication.registerContextMenu(m)
}

func (m *ContextMenu) Destroy() {
	globalApplication.unregisterContextMenu(m.name)
}

type Menu struct {
	items []*MenuItem
	label string

	impl menuImpl
}

func NewMenu() *Menu {
	return &Menu{}
}

func (m *Menu) Add(label string) *MenuItem {
	result := NewMenuItem(label)
	m.items = append(m.items, result)
	return result
}

func (m *Menu) AddSeparator() {
	result := NewMenuItemSeparator()
	m.items = append(m.items, result)
}

func (m *Menu) AddCheckbox(label string, enabled bool) *MenuItem {
	result := NewMenuItemCheckbox(label, enabled)
	m.items = append(m.items, result)
	return result
}

func (m *Menu) AddRadio(label string, enabled bool) *MenuItem {
	result := NewMenuItemRadio(label, enabled)
	m.items = append(m.items, result)
	return result
}

func (m *Menu) Update() {
	m.processRadioGroups()
	if m.impl == nil {
		m.impl = newMenuImpl(m)
	}
	m.impl.update()
}

// Clear all menu items
func (m *Menu) Clear() {
	for _, item := range m.items {
		removeMenuItemByID(item.id)
	}
	m.items = nil
}

func (m *Menu) Destroy() {
	for _, item := range m.items {
		item.Destroy()
	}
	m.items = nil
}

func (m *Menu) AddSubmenu(s string) *Menu {
	result := NewSubMenuItem(s)
	m.items = append(m.items, result)
	return result.submenu
}

func (m *Menu) AddRole(role Role) *Menu {
	result := NewRole(role)
	if result != nil {
		m.items = append(m.items, result)
	}
	return m
}

func (m *Menu) processRadioGroups() {
	var radioGroup []*MenuItem

	closeOutRadioGroups := func() {
		if len(radioGroup) > 0 {
			for _, item := range radioGroup {
				item.radioGroupMembers = radioGroup
			}
			radioGroup = []*MenuItem{}
		}
	}

	for _, item := range m.items {
		if item.itemType != radio {
			closeOutRadioGroups()
		}
		if item.itemType == submenu {
			item.submenu.processRadioGroups()
			continue
		}
		if item.itemType == radio {
			radioGroup = append(radioGroup, item)
		}
	}
	closeOutRadioGroups()
}

func (m *Menu) SetLabel(label string) {
	m.label = label
}

func (m *Menu) setContextData(data *ContextMenuData) {
	for _, item := range m.items {
		item.setContextData(data)
	}
}

// FindByLabel recursively searches for a menu item with the given label
// and returns the first match, or nil if not found.
func (m *Menu) FindByLabel(label string) *MenuItem {
	for _, item := range m.items {
		if item.label == label {
			return item
		}
		if item.submenu != nil {
			found := item.submenu.FindByLabel(label)
			if found != nil {
				return found
			}
		}
	}
	return nil
}

// FindByRole recursively searches for a menu item with the given role
// and returns the first match, or nil if not found.
func (m *Menu) FindByRole(role Role) *MenuItem {
	for _, item := range m.items {
		if item.role == role {
			return item
		}
		if item.submenu != nil {
			found := item.submenu.FindByRole(role)
			if found != nil {
				return found
			}
		}
	}
	return nil
}

func (m *Menu) RemoveMenuItem(target *MenuItem) {
	for i, item := range m.items {
		if item == target {
			// Remove the item from the slice
			m.items = append(m.items[:i], m.items[i+1:]...)
			break
		}
		if item.submenu != nil {
			item.submenu.RemoveMenuItem(target)
		}
	}
}

// ItemAt returns the menu item at the given index, or nil if the index is out of bounds.
func (m *Menu) ItemAt(index int) *MenuItem {
	if index < 0 || index >= len(m.items) {
		return nil
	}
	return m.items[index]
}

// Clone recursively clones the menu and all its submenus.
func (m *Menu) Clone() *Menu {
	result := &Menu{
		label: m.label,
	}
	for _, item := range m.items {
		result.items = append(result.items, item.Clone())
	}
	return result
}

// Append menu to an existing menu
func (m *Menu) Append(in *Menu) {
	if in == nil {
		return
	}
	m.items = append(m.items, in.items...)
}

// Prepend menu before an existing menu
func (m *Menu) Prepend(in *Menu) {
	m.items = append(in.items, m.items...)
}

func (a *App) NewMenu() *Menu {
	return &Menu{}
}

func NewMenuFromItems(item *MenuItem, items ...*MenuItem) *Menu {
	result := &Menu{
		items: []*MenuItem{item},
	}
	result.items = append(result.items, items...)
	return result
}

func NewSubmenu(s string, items *Menu) *MenuItem {
	result := NewSubMenuItem(s)
	result.submenu = items
	return result
}
