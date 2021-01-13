package menumanager

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

type Manager struct {

	// The application menu.
	applicationMenu     *menu.Menu
	applicationMenuJSON string

	// Our application menu mappings
	applicationMenuItemMap *MenuItemMap

	// Context menus
	contextMenus        map[string]*ContextMenu
	contextMenuPointers map[*menu.ContextMenu]string

	// Tray menu stores
	trayMenus        map[string]*TrayMenu
	trayMenuPointers map[*menu.TrayMenu]string
}

func NewManager() *Manager {
	return &Manager{
		applicationMenuItemMap: NewMenuItemMap(),
		contextMenus:           make(map[string]*ContextMenu),
		contextMenuPointers:    make(map[*menu.ContextMenu]string),
		trayMenus:              make(map[string]*TrayMenu),
		trayMenuPointers:       make(map[*menu.TrayMenu]string),
	}
}

func (m *Manager) getMenuItemByID(menuMap *MenuItemMap, menuId string) *menu.MenuItem {
	return menuMap.idToMenuItemMap[menuId]
}

func (m *Manager) ProcessClick(menuID string, data string, menuType string, parentID string) error {

	var menuItemMap *MenuItemMap

	switch menuType {
	case "ApplicationMenu":
		menuItemMap = m.applicationMenuItemMap
	case "ContextMenu":
		contextMenu := m.contextMenus[parentID]
		if contextMenu == nil {
			return fmt.Errorf("unknown context menu: %s", parentID)
		}
		menuItemMap = contextMenu.menuItemMap
	case "TrayMenu":
		trayMenu := m.trayMenus[parentID]
		if trayMenu == nil {
			return fmt.Errorf("unknown tray menu: %s", parentID)
		}
		menuItemMap = trayMenu.menuItemMap
	default:
		return fmt.Errorf("unknown menutype: %s", menuType)
	}

	// Get the menu item
	menuItem := menuItemMap.getMenuItemByID(menuID)
	if menuItem == nil {
		return fmt.Errorf("Cannot process menuid %s - unknown", menuID)
	}

	// Is the menu item a checkbox?
	if menuItem.Type == menu.CheckboxType {
		// Toggle state
		menuItem.Checked = !menuItem.Checked
	}

	if menuItem.Click == nil {
		// No callback
		return fmt.Errorf("No callback for menu '%s'", menuItem.Label)
	}

	// Create new Callback struct
	callbackData := &menu.CallbackData{
		MenuItem:    menuItem,
		ContextData: data,
	}

	// Call back!
	go menuItem.Click(callbackData)

	return nil
}
