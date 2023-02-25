package application

type menuImpl interface {
	update()
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
	result := newMenuItem(label)
	m.items = append(m.items, result)
	return result
}

func (m *Menu) AddSeparator() {
	result := newMenuItemSeperator()
	m.items = append(m.items, result)
}

func (m *Menu) AddCheckbox(label string, enabled bool) *MenuItem {
	result := newMenuItemCheckbox(label, enabled)
	m.items = append(m.items, result)
	return result
}

func (m *Menu) AddRadio(label string, enabled bool) *MenuItem {
	result := newMenuItemRadio(label, enabled)
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

func (m *Menu) AddSubmenu(s string) *Menu {
	result := newSubMenuItem(s)
	m.items = append(m.items, result)
	return result.submenu
}

func (m *Menu) AddRole(role Role) *Menu {
	result := newRole(role)
	m.items = append(m.items, result)
	return m
}

func (m *Menu) processRadioGroups() {
	var radioGroup []*MenuItem
	for _, item := range m.items {
		if item.itemType == submenu {
			item.submenu.processRadioGroups()
			continue
		}
		if item.itemType == radio {
			radioGroup = append(radioGroup, item)
		} else {
			if len(radioGroup) > 0 {
				for _, item := range radioGroup {
					item.radioGroupMembers = radioGroup
				}
				radioGroup = nil
			}
		}
	}
	if len(radioGroup) > 0 {
		for _, item := range radioGroup {
			item.radioGroupMembers = radioGroup
		}
	}
}

func (m *Menu) SetLabel(label string) {
	m.label = label
}

func (m *Menu) setContextData(data *ContextMenuData) {
	for _, item := range m.items {
		item.setContextData(data)
	}
}

func (a *App) NewMenu() *Menu {
	return &Menu{}
}

func defaultApplicationMenu() *Menu {
	menu := NewMenu()
	menu.AddRole(AppMenu)
	menu.AddRole(FileMenu)
	menu.AddRole(EditMenu)
	menu.AddRole(ViewMenu)
	menu.AddRole(WindowMenu)
	menu.AddRole(HelpMenu)
	return menu
}
