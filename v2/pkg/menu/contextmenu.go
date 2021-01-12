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

func (c *ContextMenus) GetByID(menuID string) *MenuItem {

	// Loop over menu items
	for _, item := range c.Items {
		result := item.GetByID(menuID)
		if result != nil {
			return result
		}
	}
	return nil
}

func (c *ContextMenus) RemoveByID(id string) bool {
	// Loop over menu items
	for _, item := range c.Items {
		result := item.RemoveByID(id)
		if result == true {
			return result
		}
	}
	return false
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
