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

func (m *Menu) GetByID(menuID string) *MenuItem {

	// Loop over menu items
	for _, item := range m.Items {
		result := item.getByID(menuID)
		if result != nil {
			return result
		}
	}
	return nil
}

func (m *Menu) RemoveByID(id string) bool {
	// Loop over menu items
	for index, item := range m.Items {
		if item.ID == id {
			m.Items = append(m.Items[:index], m.Items[index+1:]...)
			return true
		}
		result := item.removeByID(id)
		if result == true {
			return result
		}
	}
	return false
}
