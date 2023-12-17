package menumanager

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/leaanthony/go-ansi-parser"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

var (
	trayMenuID      int
	trayMenuIDMutex sync.Mutex
)

func generateTrayID() string {
	var idStr string
	trayMenuIDMutex.Lock()
	idStr = strconv.Itoa(trayMenuID)
	trayMenuID++
	trayMenuIDMutex.Unlock()
	return idStr
}

type TrayMenu struct {
	ID               string
	Label            string
	FontSize         int
	FontName         string
	Disabled         bool
	Tooltip          string `json:",omitempty"`
	Image            string
	MacTemplateImage bool
	RGBA             string
	menuItemMap      *MenuItemMap
	menu             *menu.Menu
	ProcessedMenu    *WailsMenu
	trayMenu         *menu.TrayMenu
	StyledLabel      []*ansi.StyledText `json:",omitempty"`
}

func (t *TrayMenu) AsJSON() (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func NewTrayMenu(trayMenu *menu.TrayMenu) *TrayMenu {
	// Parse ANSI text
	var styledLabel []*ansi.StyledText
	tempLabel := trayMenu.Label
	if strings.Contains(tempLabel, "\033[") {
		parsedLabel, err := ansi.Parse(tempLabel)
		if err == nil {
			styledLabel = parsedLabel
		}
	}

	result := &TrayMenu{
		Label:            trayMenu.Label,
		FontName:         trayMenu.FontName,
		FontSize:         trayMenu.FontSize,
		Disabled:         trayMenu.Disabled,
		Tooltip:          trayMenu.Tooltip,
		Image:            trayMenu.Image,
		MacTemplateImage: trayMenu.MacTemplateImage,
		menu:             trayMenu.Menu,
		RGBA:             trayMenu.RGBA,
		menuItemMap:      NewMenuItemMap(),
		trayMenu:         trayMenu,
		StyledLabel:      styledLabel,
	}

	result.menuItemMap.AddMenu(trayMenu.Menu)
	result.ProcessedMenu = NewWailsMenu(result.menuItemMap, result.menu)

	return result
}

func (m *Manager) OnTrayMenuOpen(id string) {
	trayMenu, ok := m.trayMenus[id]
	if !ok {
		return
	}
	if trayMenu.trayMenu.OnOpen == nil {
		return
	}
	go trayMenu.trayMenu.OnOpen()
}

func (m *Manager) OnTrayMenuClose(id string) {
	trayMenu, ok := m.trayMenus[id]
	if !ok {
		return
	}
	if trayMenu.trayMenu.OnClose == nil {
		return
	}
	go trayMenu.trayMenu.OnClose()
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

func (m *Manager) GetTrayID(trayMenu *menu.TrayMenu) (string, error) {
	trayID, exists := m.trayMenuPointers[trayMenu]
	if !exists {
		return "", fmt.Errorf("Unable to find menu ID for tray menu!")
	}
	return trayID, nil
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

func (m *Manager) UpdateTrayMenuLabel(trayMenu *menu.TrayMenu) (string, error) {
	trayID, trayMenuKnown := m.trayMenuPointers[trayMenu]
	if !trayMenuKnown {
		return "", fmt.Errorf("[UpdateTrayMenuLabel] unknown tray id for tray %s", trayMenu.Label)
	}

	type LabelUpdate struct {
		ID               string
		Label            string `json:",omitempty"`
		FontName         string `json:",omitempty"`
		FontSize         int
		RGBA             string `json:",omitempty"`
		Disabled         bool
		Tooltip          string `json:",omitempty"`
		Image            string `json:",omitempty"`
		MacTemplateImage bool
		StyledLabel      []*ansi.StyledText `json:",omitempty"`
	}

	// Parse ANSI text
	var styledLabel []*ansi.StyledText
	tempLabel := trayMenu.Label
	if strings.Contains(tempLabel, "\033[") {
		parsedLabel, err := ansi.Parse(tempLabel)
		if err == nil {
			styledLabel = parsedLabel
		}
	}

	update := &LabelUpdate{
		ID:               trayID,
		Label:            trayMenu.Label,
		FontName:         trayMenu.FontName,
		FontSize:         trayMenu.FontSize,
		Disabled:         trayMenu.Disabled,
		Tooltip:          trayMenu.Tooltip,
		Image:            trayMenu.Image,
		MacTemplateImage: trayMenu.MacTemplateImage,
		RGBA:             trayMenu.RGBA,
		StyledLabel:      styledLabel,
	}

	data, err := json.Marshal(update)
	if err != nil {
		return "", errors.Wrap(err, "[UpdateTrayMenuLabel] ")
	}

	return string(data), nil
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
