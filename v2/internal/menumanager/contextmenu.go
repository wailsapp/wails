package menumanager

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

type ContextMenu struct {
	ID          string
	JSON        string
	menuItemMap *MenuItemMap
	menu        *menu.Menu
}

func NewContextMenu(ID string, menu *menu.Menu) *ContextMenu {

	result := &ContextMenu{
		ID:          ID,
		JSON:        "",
		menu:        menu,
		menuItemMap: NewMenuItemMap(),
	}

	result.menuItemMap.AddMenu(menu)

	return result
}

func (m *Manager) AddContextMenu(menuID string, menu *menu.Menu) error {
	contextMenu := NewContextMenu(menuID, menu)
	m.contextMenus[menuID] = contextMenu
	return contextMenu.process()
}

func (c *ContextMenu) process() error {

	// Process the menu
	processedApplicationMenu := NewWailsMenu(c.menuItemMap, c.menu)
	JSON, err := processedApplicationMenu.AsJSON()
	if err != nil {
		return err
	}
	c.JSON = JSON
	fmt.Printf("Processed context menu '%s':", c.ID)
	println(JSON)
	return nil
}
