package menu

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
	for _, item := range menu.Items {
		m.Items = append(m.Items, item)
	}
}

func (m *Menu) Prepend(item *MenuItem) {
	m.Items = append([]*MenuItem{item}, m.Items...)
}

func NewMenuFromItems(first *MenuItem, rest ...*MenuItem) *Menu {

	var result = NewMenu()
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
