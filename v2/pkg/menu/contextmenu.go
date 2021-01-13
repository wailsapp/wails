package menu

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
