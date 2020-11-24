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

func NewMenuFromItems(first *MenuItem, rest ...*MenuItem) *Menu {

	var result = NewMenu()
	result.Append(first)
	for _, item := range rest {
		result.Append(item)
	}

	return result
}
