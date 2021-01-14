package menumanager

import (
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"sync"
)

var trayMenuID int
var trayMenuIDMutex sync.Mutex

func generateTrayID() string {
	trayMenuIDMutex.Lock()
	result := fmt.Sprintf("%d", trayMenuID)
	trayMenuID++
	trayMenuIDMutex.Unlock()
	return result
}

type TrayMenu struct {
	ID            string
	Label         string
	Icon          string
	menuItemMap   *MenuItemMap
	menu          *menu.Menu
	ProcessedMenu *WailsMenu
}

func (t *TrayMenu) AsJSON() (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func NewTrayMenu(trayMenu *menu.TrayMenu) *TrayMenu {

	result := &TrayMenu{
		Label:       trayMenu.Label,
		Icon:        trayMenu.Icon,
		menu:        trayMenu.Menu,
		menuItemMap: NewMenuItemMap(),
	}

	result.menuItemMap.AddMenu(trayMenu.Menu)
	result.ProcessedMenu = NewWailsMenu(result.menuItemMap, result.menu)

	return result
}

func (m *Manager) AddTrayMenu(trayMenu *menu.TrayMenu) (string, error) {
	newTrayMenu := NewTrayMenu(trayMenu)

	// Hook up a new ID
	trayID := generateTrayID()
	newTrayMenu.ID = trayID

	// Save the references
	m.trayMenus[trayID] = newTrayMenu
	m.trayMenuPointers[trayMenu] = trayID

	return newTrayMenu.AsJSON()
}

// SetTrayMenu updates or creates a menu
func (m *Manager) SetTrayMenu(trayMenu *menu.TrayMenu) (string, error) {
	trayID, trayMenuKnown := m.trayMenuPointers[trayMenu]
	if !trayMenuKnown {
		return m.AddTrayMenu(trayMenu)
	}

	// Create the updated tray menu
	updatedTrayMenu := NewTrayMenu(trayMenu)
	updatedTrayMenu.ID = trayID

	// Save the reference
	m.trayMenus[trayID] = updatedTrayMenu

	return updatedTrayMenu.AsJSON()
}

func (m *Manager) GetTrayMenus() ([]string, error) {
	result := []string{}
	for _, trayMenu := range m.trayMenus {
		JSON, err := trayMenu.AsJSON()
		if err != nil {
			return nil, err
		}
		result = append(result, JSON)
	}

	return result, nil
}

func (m *Manager) GetContextMenus() ([]string, error) {
	result := []string{}
	for _, contextMenu := range m.contextMenus {
		JSON, err := contextMenu.AsJSON()
		if err != nil {
			return nil, err
		}
		result = append(result, JSON)
	}

	return result, nil
}
