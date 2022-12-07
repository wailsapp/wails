//go:build windows

package menu

import (
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// MenuManager manages the menus for the application
var MenuManager = NewManager()

type radioGroup []*menu.MenuItem

// Click updates the radio group state based on the item clicked
func (g *radioGroup) Click(item *menu.MenuItem) {
	for _, radioGroupItem := range *g {
		if radioGroupItem != item {
			radioGroupItem.Checked = false
		}
	}
}

type processedMenu struct {

	// the menu we processed
	menu *menu.Menu

	// updateMenuItemCallback is called when the menu item needs to be updated in the UI
	updateMenuItemCallback func(*menu.MenuItem)

	// items is a map of all menu items in this menu
	items map[*menu.MenuItem]struct{}

	// radioGroups tracks which radiogroup a menu item belongs to
	radioGroups map[*menu.MenuItem][]*radioGroup
}

func newProcessedMenu(topLevelMenu *menu.Menu, updateMenuItemCallback func(*menu.MenuItem)) *processedMenu {
	result := &processedMenu{
		updateMenuItemCallback: updateMenuItemCallback,
		menu:                   topLevelMenu,
		items:                  make(map[*menu.MenuItem]struct{}),
		radioGroups:            make(map[*menu.MenuItem][]*radioGroup),
	}
	result.process(topLevelMenu.Items)
	return result
}

func (p *processedMenu) process(items []*menu.MenuItem) {
	var currentRadioGroup radioGroup
	for index, item := range items {
		// Save the reference to the top level menu for this item
		p.items[item] = struct{}{}

		// If this is a radio item, add it to the radio group
		if item.Type == menu.RadioType {
			currentRadioGroup = append(currentRadioGroup, item)
		}

		// If this is not a radio item, or we are processing the last item in the menu,
		// then we need to add the current radio group to the map if it has items
		if item.Type != menu.RadioType || index == len(items)-1 {
			if len(currentRadioGroup) > 0 {
				p.addRadioGroup(currentRadioGroup)
				currentRadioGroup = nil
			}
		}

		// Process the submenu
		if item.SubMenu != nil {
			p.process(item.SubMenu.Items)
		}
	}
}

func (p *processedMenu) processClick(item *menu.MenuItem) {
	// If this item is not in our menu, then we can't process it
	if _, ok := p.items[item]; !ok {
		return
	}

	// If this is a radio item, then we need to update the radio group
	if item.Type == menu.RadioType {
		// Get the radio groups for this item
		radioGroups := p.radioGroups[item]
		// Iterate each radio group this item belongs to and set the checked state
		// of all items apart from the one that was clicked to false
		for _, thisRadioGroup := range radioGroups {
			thisRadioGroup.Click(item)
			for _, thisRadioGroupItem := range *thisRadioGroup {
				p.updateMenuItemCallback(thisRadioGroupItem)
			}
		}
	}

	if item.Type == menu.CheckboxType {
		p.updateMenuItemCallback(item)
	}

}

func (p *processedMenu) addRadioGroup(r radioGroup) {
	for _, item := range r {
		p.radioGroups[item] = append(p.radioGroups[item], &r)
	}
}

type Manager struct {
	menus map[*menu.Menu]*processedMenu
}

func NewManager() *Manager {
	return &Manager{
		menus: make(map[*menu.Menu]*processedMenu),
	}
}

func (m *Manager) AddMenu(menu *menu.Menu, updateMenuItemCallback func(*menu.MenuItem)) {
	m.menus[menu] = newProcessedMenu(menu, updateMenuItemCallback)
}

func (m *Manager) ProcessClick(item *menu.MenuItem) {

	// if menuitem is a checkbox, then we need to toggle the state
	if item.Type == menu.CheckboxType {
		item.Checked = !item.Checked
	}

	// Set the radio item to checked
	if item.Type == menu.RadioType {
		item.Checked = true
	}

	for _, thisMenu := range m.menus {
		thisMenu.processClick(item)
	}

	if item.Click != nil {
		item.Click(&menu.CallbackData{
			MenuItem: item,
		})
	}
}

func (m *Manager) RemoveMenu(data *menu.Menu) {
	delete(m.menus, data)
}
