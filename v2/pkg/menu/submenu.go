package menu

// SubMenu creates a new submenu which may be added to other
// menus
func SubMenu(label string, items []*MenuItem) *MenuItem {
	return &MenuItem{
		Label:   label,
		SubMenu: items,
	}
}
