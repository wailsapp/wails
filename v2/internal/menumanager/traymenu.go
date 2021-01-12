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

func NewTrayMenu(trayMenu *menu.TrayMenu) *TrayMenu {

	result := &TrayMenu{
		ID:          generateTrayID(),
		Label:       trayMenu.Label,
		Icon:        trayMenu.Icon,
		menu:        trayMenu.Menu,
		menuItemMap: NewMenuItemMap(),
	}

	result.menuItemMap.AddMenu(trayMenu.Menu)
	result.ProcessedMenu = NewWailsMenu(result.menuItemMap, result.menu)

	return result
}

func (m *Manager) AddTrayMenu(trayMenu *menu.TrayMenu) {
	newTrayMenu := NewTrayMenu(trayMenu)
	m.trayMenus[newTrayMenu.ID] = newTrayMenu
}

func (m *Manager) GetTrayMenus() ([]string, error) {
	result := []string{}
	for _, trayMenu := range m.trayMenus {
		data, err := json.Marshal(trayMenu)
		if err != nil {
			return nil, err
		}
		result = append(result, string(data))
	}

	return result, nil
}
