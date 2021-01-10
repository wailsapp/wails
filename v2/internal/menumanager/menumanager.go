package menumanager

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

type Manager struct {

	// The application menu.
	applicationMenu     *menu.Menu
	applicationMenuJSON string

	// Our menu mappings
	menuItemMap *MenuItemMap
}

func NewManager() *Manager {
	return &Manager{
		menuItemMap: NewMenuItemMap(),
	}
}

func (m *Manager) SetApplicationMenu(applicationMenu *menu.Menu) error {
	if applicationMenu == nil {
		return nil
	}
	m.applicationMenu = applicationMenu
	m.menuItemMap.AddMenu(applicationMenu)
	return m.processApplicationMenu()
}

func (m *Manager) GetApplicationMenuJSON() string {
	return m.applicationMenuJSON
}

func (m *Manager) processApplicationMenu() error {

	// Process the menu
	processedApplicationMenu := m.NewWailsMenu(m.applicationMenu)
	applicationMenuJSON, err := processedApplicationMenu.AsJSON()
	if err != nil {
		return err
	}
	m.applicationMenuJSON = applicationMenuJSON
	return nil
}

func (m *Manager) getMenuItemByID(menuId string) *menu.MenuItem {
	return m.menuItemMap.idToMenuItemMap[menuId]
}

func (m *Manager) ProcessClick(menuID string, data string) error {

	// Get the menu item
	menuItem := m.getMenuItemByID(menuID)
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
