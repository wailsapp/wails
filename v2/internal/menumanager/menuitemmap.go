package menumanager

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

// MenuItemMap holds a mapping between menuIDs and menu items
type MenuItemMap struct {
	idToMenuItemMap map[string]*menu.MenuItem
	menuItemToIDMap map[*menu.MenuItem]string

	// We use a simple counter to keep track of unique menu IDs
	menuIDCounter      int64
	menuIDCounterMutex sync.Mutex
}

func NewMenuItemMap() *MenuItemMap {
	result := &MenuItemMap{
		idToMenuItemMap: make(map[string]*menu.MenuItem),
		menuItemToIDMap: make(map[*menu.MenuItem]string),
	}

	return result
}

func (m *MenuItemMap) AddMenu(menu *menu.Menu) {
	if menu == nil {
		return
	}
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

// GenerateMenuID returns a unique string ID for a menu item
func (m *MenuItemMap) generateMenuID() string {
	m.menuIDCounterMutex.Lock()
	result := strconv.FormatInt(m.menuIDCounter, 10)
	m.menuIDCounter++
	m.menuIDCounterMutex.Unlock()
	return result
}

func (m *MenuItemMap) processMenuItem(item *menu.MenuItem) {
	if item.SubMenu != nil {
		for _, submenuitem := range item.SubMenu.Items {
			m.processMenuItem(submenuitem)
		}
	}

	// Create a unique ID for this menu item
	menuID := m.generateMenuID()

	// Store references
	m.idToMenuItemMap[menuID] = item
	m.menuItemToIDMap[item] = menuID
}

func (m *MenuItemMap) getMenuItemByID(menuId string) *menu.MenuItem {
	return m.idToMenuItemMap[menuId]
}
