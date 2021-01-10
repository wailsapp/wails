package menumanager

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// MenuItemMap holds a mapping between menuIDs and menu items
type MenuItemMap struct {
	idToMenuItemMap map[string]*menu.MenuItem
	menuItemToIDMap map[*menu.MenuItem]string
}

func NewMenuItemMap() *MenuItemMap {
	result := &MenuItemMap{
		idToMenuItemMap: make(map[string]*menu.MenuItem),
		menuItemToIDMap: make(map[*menu.MenuItem]string),
	}

	return result
}

func (m *MenuItemMap) AddMenu(menu *menu.Menu) {
	for _, item := range menu.Items {
		m.processMenuItem(item)
	}
}

func (m *MenuItemMap) Dump() {
	println("idToMenuItemMap:")
	for key, value := range m.idToMenuItemMap {
		fmt.Printf("  %s\t%p\n", key, value)
	}
	println("\nmenuItemToIDMap")
	for key, value := range m.menuItemToIDMap {
		fmt.Printf("  %p\t%s\n", key, value)
	}
}

func (m *MenuItemMap) processMenuItem(item *menu.MenuItem) {

	if item.SubMenu != nil {
		for _, submenuitem := range item.SubMenu.Items {
			m.processMenuItem(submenuitem)
		}
	}

	// Create a unique ID for this menu item
	menuID := fmt.Sprintf("%d", len(m.idToMenuItemMap))

	// Store references
	m.idToMenuItemMap[menuID] = item
	m.menuItemToIDMap[item] = menuID
}
