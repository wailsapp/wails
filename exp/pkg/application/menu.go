package application

type menuImpl interface {
	update()
}

type Menu struct {
	items []*MenuItem

	impl menuImpl
}

func (m *Menu) Add(label string) *MenuItem {
	result := newMenuItem(label)
	m.items = append(m.items, result)
	return result
}

func (m *Menu) AddSeperator() {
	result := newMenuItemSeperator()
	m.items = append(m.items, result)
}

func (m *Menu) AddCheckbox(label string, enabled bool) *MenuItem {
	result := newMenuItemCheckbox(label, enabled)
	m.items = append(m.items, result)
	return result
}

func (m *Menu) Update() {
	if m.impl == nil {
		m.impl = newMenuImpl(m)
	}
	m.impl.update()
}

func (a *App) NewMenu() *Menu {
	return &Menu{}
}
