package menu

type ContextMenus struct {
	Items map[string]*Menu
}

func NewContextMenus() *ContextMenus {
	return &ContextMenus{
		Items: make(map[string]*Menu),
	}
}

func (c *ContextMenus) AddMenu(ID string, menu *Menu) {
	c.Items[ID] = menu
}

type ContextMenu struct {
	ID   string
	Menu *Menu
}

func NewContextMenu(ID string, menu *Menu) *ContextMenu {
	return &ContextMenu{
		ID:   ID,
		Menu: menu,
	}
}
