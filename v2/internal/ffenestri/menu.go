package ffenestri

import "github.com/wailsapp/wails/v2/pkg/menu"

// ProcessedMenu is the original menu with the addition
// of radio groups extracted from the menu data
type ProcessedMenu struct {
	Menu              *menu.Menu
	RadioGroups       []*RadioGroup
	currentRadioGroup []string
}

// RadioGroup holds all the members of the same radio group
type RadioGroup struct {
	Members []string
	Length  int
}

// NewProcessedMenu processed the given menu and returns
// the original menu with the extracted radio groups
func NewProcessedMenu(menu *menu.Menu) *ProcessedMenu {
	result := &ProcessedMenu{
		Menu:              menu,
		RadioGroups:       []*RadioGroup{},
		currentRadioGroup: []string{},
	}

	result.processMenu()

	return result
}

func (p *ProcessedMenu) processMenu() {
	// Loop over top level menus
	for _, item := range p.Menu.Items {
		// Process MenuItem
		p.processMenuItem(item)
	}

	p.finaliseRadioGroup()
}

func (p *ProcessedMenu) processMenuItem(item *menu.MenuItem) {

	switch item.Type {

	// We need to recurse submenus
	case menu.SubmenuType:

		// Finalise any current radio groups as they don't trickle down to submenus
		p.finaliseRadioGroup()

		// Process each submenu item
		for _, subitem := range item.SubMenu.Items {
			p.processMenuItem(subitem)
		}
	case menu.RadioType:
		// Add the item to the radio group
		p.currentRadioGroup = append(p.currentRadioGroup, item.ID)
	default:
		p.finaliseRadioGroup()
	}
}

func (p *ProcessedMenu) finaliseRadioGroup() {

	// If we were processing a radio group, fix up the references
	if len(p.currentRadioGroup) > 0 {

		// Create new radiogroup
		group := &RadioGroup{
			Members: p.currentRadioGroup,
			Length:  len(p.currentRadioGroup),
		}
		p.RadioGroups = append(p.RadioGroups, group)

		// Empty the radio group
		p.currentRadioGroup = []string{}
	}
}
