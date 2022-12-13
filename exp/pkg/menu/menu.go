package menu

type Menu struct {
	Items []*Item
}

func (m *Menu) Label(label string) *Item {
	return &Item{
		Label: label,
	}
}

type CallbackContext struct {
	MenuItem *Item
}

type Item struct {
	Label    string
	Disabled bool
	Click    func(*CallbackContext)
	ID       uint
}
