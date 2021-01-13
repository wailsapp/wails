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
