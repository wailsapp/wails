package menumanager

import (
	"encoding/json"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

type ContextMenu struct {
	ID            string
	ProcessedMenu *WailsMenu
	menuItemMap   *MenuItemMap
	menu          *menu.Menu
}

func (t *ContextMenu) AsJSON() (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func NewContextMenu(contextMenu *menu.ContextMenu) *ContextMenu {
	result := &ContextMenu{
		ID:          contextMenu.ID,
		menu:        contextMenu.Menu,
		menuItemMap: NewMenuItemMap(),
	}

	result.menuItemMap.AddMenu(contextMenu.Menu)
	result.ProcessedMenu = NewWailsMenu(result.menuItemMap, result.menu)

	return result
}

func (m *Manager) AddContextMenu(contextMenu *menu.ContextMenu) {
	newContextMenu := NewContextMenu(contextMenu)

	// Save the references
	m.contextMenus[contextMenu.ID] = newContextMenu
	m.contextMenuPointers[contextMenu] = contextMenu.ID
}

func (m *Manager) UpdateContextMenu(contextMenu *menu.ContextMenu) (string, error) {
	contextMenuID, contextMenuKnown := m.contextMenuPointers[contextMenu]
	if !contextMenuKnown {
		return "", fmt.Errorf("unknown Context Menu '%s'. Please add the context menu using AddContextMenu()", contextMenu.ID)
	}

	// Create the updated context menu
	updatedContextMenu := NewContextMenu(contextMenu)

	// Save the reference
	m.contextMenus[contextMenuID] = updatedContextMenu

	return updatedContextMenu.AsJSON()
}
