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
}

func NewManager() *Manager {
	return &Manager{
		applicationMenuItemMap: NewMenuItemMap(),
	}
}

func (m *Manager) SetApplicationMenu(applicationMenu *menu.Menu) error {
	if applicationMenu == nil {
		return nil
	}
	m.applicationMenu = applicationMenu

	// Reset the menu map
	m.applicationMenuItemMap = NewMenuItemMap()

	// Add the menu to the menu map
	m.applicationMenuItemMap.AddMenu(applicationMenu)

	return m.processApplicationMenu()
}

func (m *Manager) GetApplicationMenuJSON() string {
	return m.applicationMenuJSON
}

// UpdateApplicationMenu reprocesses the application menu to pick up structure
// changes etc
// Returns the JSON representation of the updated menu
func (m *Manager) UpdateApplicationMenu() (string, error) {
	m.applicationMenuItemMap = NewMenuItemMap()
	m.applicationMenuItemMap.AddMenu(m.applicationMenu)
	err := m.processApplicationMenu()
	return m.applicationMenuJSON, err
}

func (m *Manager) processApplicationMenu() error {

	// Process the menu
	processedApplicationMenu := m.NewWailsMenu(m.applicationMenuItemMap, m.applicationMenu)
	applicationMenuJSON, err := processedApplicationMenu.AsJSON()
	if err != nil {
		return err
	}
	m.applicationMenuJSON = applicationMenuJSON
	return nil
}

func (m *Manager) getMenuItemByID(menuMap *MenuItemMap, menuId string) *menu.MenuItem {
	return menuMap.idToMenuItemMap[menuId]
}

func (m *Manager) ProcessClick(menuID string, data string, menuType string) error {

	var menuItemMap *MenuItemMap

	switch menuType {
	case "ApplicationMenu":
		menuItemMap = m.applicationMenuItemMap
	//case "ContextMenu":
	//	// TBD
	//case "TrayMenu":
	//	// TBD
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
