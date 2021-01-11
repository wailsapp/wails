package menumanager

import "github.com/wailsapp/wails/v2/pkg/menu"

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
	processedApplicationMenu := NewWailsMenu(m.applicationMenuItemMap, m.applicationMenu)
	applicationMenuJSON, err := processedApplicationMenu.AsJSON()
	if err != nil {
		return err
	}
	m.applicationMenuJSON = applicationMenuJSON
	return nil
}
