package menumanager

import (
	"encoding/json"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
)

type ProcessedMenuItem struct {
	ID string
	// Label is what appears as the menu text
	Label string `json:",omitempty"`
	// Role is a predefined menu type
	// Role menu.Role `json:",omitempty"`
	// Accelerator holds a representation of a key binding
	Accelerator *keys.Accelerator `json:",omitempty"`
	// Type of MenuItem, EG: Checkbox, Text, Separator, Radio, Submenu
	Type menu.Type
	// Disabled makes the item unselectable
	Disabled bool `json:",omitempty"`
	// Hidden ensures that the item is not shown in the menu
	Hidden bool `json:",omitempty"`
	// Checked indicates if the item is selected (used by Checkbox and Radio types only)
	Checked bool `json:",omitempty"`
	// Submenu contains a list of menu items that will be shown as a submenu
	// SubMenu []*MenuItem `json:"SubMenu,omitempty"`
	SubMenu *ProcessedMenu `json:",omitempty"`
	/*
		// Colour
		RGBA string `json:",omitempty"`

		// Font
		FontSize int    `json:",omitempty"`
		FontName string `json:",omitempty"`

		// Image - base64 image data
		Image            string `json:",omitempty"`
		MacTemplateImage bool   `json:", omitempty"`
		MacAlternate     bool   `json:", omitempty"`

		// Tooltip
		Tooltip string `json:",omitempty"`

		// Styled label
		StyledLabel []*ansi.StyledText `json:",omitempty"`
	*/
}

func NewProcessedMenuItem(menuItemMap *MenuItemMap, menuItem *menu.MenuItem) *ProcessedMenuItem {
	ID := menuItemMap.menuItemToIDMap[menuItem]

	// Parse ANSI text
	//var styledLabel []*ansi.StyledText
	//tempLabel := menuItem.Label
	//if strings.Contains(tempLabel, "\033[") {
	//	parsedLabel, err := ansi.Parse(menuItem.Label)
	//	if err == nil {
	//		styledLabel = parsedLabel
	//	}
	//}

	result := &ProcessedMenuItem{
		ID:    ID,
		Label: menuItem.Label,
		// Role:             menuItem.Role,
		Accelerator: menuItem.Accelerator,
		Type:        menuItem.Type,
		Disabled:    menuItem.Disabled,
		Hidden:      menuItem.Hidden,
		Checked:     menuItem.Checked,
		SubMenu:     nil,
		// BackgroundColour:             menuItem.BackgroundColour,
		// FontSize:         menuItem.FontSize,
		// FontName:         menuItem.FontName,
		// Image:            menuItem.Image,
		// MacTemplateImage: menuItem.MacTemplateImage,
		// MacAlternate:     menuItem.MacAlternate,
		// Tooltip:          menuItem.Tooltip,
		// StyledLabel:      styledLabel,
	}

	if menuItem.SubMenu != nil {
		result.SubMenu = NewProcessedMenu(menuItemMap, menuItem.SubMenu)
	}

	return result
}

type ProcessedMenu struct {
	Items []*ProcessedMenuItem
}

func NewProcessedMenu(menuItemMap *MenuItemMap, menu *menu.Menu) *ProcessedMenu {
	result := &ProcessedMenu{}
	if menu != nil {
		for _, item := range menu.Items {
			processedMenuItem := NewProcessedMenuItem(menuItemMap, item)
			result.Items = append(result.Items, processedMenuItem)
		}
	}

	return result
}

// WailsMenu is the original menu with the addition
// of radio groups extracted from the menu data
type WailsMenu struct {
	Menu              *ProcessedMenu
	RadioGroups       []*RadioGroup
	currentRadioGroup []string
}

// RadioGroup holds all the members of the same radio group
type RadioGroup struct {
	Members []string
	Length  int
}

func NewWailsMenu(menuItemMap *MenuItemMap, menu *menu.Menu) *WailsMenu {
	result := &WailsMenu{}

	// Process the menus
	result.Menu = NewProcessedMenu(menuItemMap, menu)

	// Process the radio groups
	result.processRadioGroups()

	return result
}

func (w *WailsMenu) AsJSON() (string, error) {
	menuAsJSON, err := json.Marshal(w)
	if err != nil {
		return "", err
	}
	return string(menuAsJSON), nil
}

func (w *WailsMenu) processRadioGroups() {
	// Loop over top level menus
	for _, item := range w.Menu.Items {
		// Process MenuItem
		w.processMenuItem(item)
	}

	w.finaliseRadioGroup()
}

func (w *WailsMenu) processMenuItem(item *ProcessedMenuItem) {
	switch item.Type {

	// We need to recurse submenus
	case menu.SubmenuType:

		// Finalise any current radio groups as they don't trickle down to submenus
		w.finaliseRadioGroup()

		// Process each submenu item
		for _, subitem := range item.SubMenu.Items {
			w.processMenuItem(subitem)
		}
	case menu.RadioType:
		// Add the item to the radio group
		w.currentRadioGroup = append(w.currentRadioGroup, item.ID)
	default:
		w.finaliseRadioGroup()
	}
}

func (w *WailsMenu) finaliseRadioGroup() {
	// If we were processing a radio group, fix up the references
	if len(w.currentRadioGroup) > 0 {

		// Create new radiogroup
		group := &RadioGroup{
			Members: w.currentRadioGroup,
			Length:  len(w.currentRadioGroup),
		}
		w.RadioGroups = append(w.RadioGroups, group)

		// Empty the radio group
		w.currentRadioGroup = []string{}
	}
}
