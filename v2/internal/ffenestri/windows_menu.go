//+build windows

package ffenestri

import (
	"fmt"
	"github.com/wailsapp/wails/v2/internal/menumanager"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"runtime"
	"strings"
)

//-------------------- Types ------------------------

type win32MenuItemID uint32
type win32Menu uintptr
type win32Window uintptr
type wailsMenuItemID string // The internal menu ID

type Menu struct {
	wailsMenu *menumanager.WailsMenu
	menu      win32Menu
	menuType  menuType

	// A list of all checkbox and radio menuitems we
	// create for this menu
	checkboxes                  []win32MenuItemID
	radioboxes                  []win32MenuItemID
	initiallySelectedRadioItems []win32MenuItemID
}

func createMenu(wailsMenu *menumanager.WailsMenu, menuType menuType) (*Menu, error) {

	mainMenu, err := createWin32Menu()
	if err != nil {
		return nil, err
	}

	result := &Menu{
		wailsMenu: wailsMenu,
		menu:      mainMenu,
		menuType:  menuType,
	}

	// Process top level menus
	for _, toplevelmenu := range applicationMenu.Menu.Items {
		err := result.processMenuItem(result.menu, toplevelmenu)
		if err != nil {
			return nil, err
		}
	}

	err = result.processRadioGroups()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Menu) processMenuItem(parent win32Menu, menuItem *menumanager.ProcessedMenuItem) error {

	// Ignore hidden items
	if menuItem.Hidden {
		return nil
	}

	// Calculate the flags for this menu item
	flags := uintptr(calculateFlags(menuItem))

	switch menuItem.Type {
	case menu.SubmenuType:
		submenu, err := createWin32PopupMenu()
		if err != nil {
			return err
		}
		for _, submenuItem := range menuItem.SubMenu.Items {
			err = m.processMenuItem(submenu, submenuItem)
			if err != nil {
				return err
			}
		}
		err = appendWin32MenuItem(parent, flags, uintptr(submenu), menuItem.Label)
		if err != nil {
			return err
		}
	case menu.TextType, menu.CheckboxType, menu.RadioType:
		win32ID := addMenuCacheEntry(parent, m.menuType, menuItem, m.wailsMenu.Menu)
		if menuItem.Accelerator != nil {
			m.processAccelerator(menuItem)
		}
		label := menuItem.Label
		//label := fmt.Sprintf("%s (%d)", menuItem.Label, win32ID)
		err := appendWin32MenuItem(parent, flags, uintptr(win32ID), label)
		if err != nil {
			return err
		}
		if menuItem.Type == menu.CheckboxType {
			// We need to maintain a list of this menu's checkboxes
			m.checkboxes = append(m.checkboxes, win32ID)
			globalCheckboxCache.addToCheckboxCache(m.wailsMenu.Menu, wailsMenuItemID(menuItem.ID), win32ID)
		}
		if menuItem.Type == menu.RadioType {
			// We need to maintain a list of this menu's radioitems
			m.radioboxes = append(m.radioboxes, win32ID)
			globalRadioGroupMap.addRadioGroupMapping(m.wailsMenu.Menu, wailsMenuItemID(menuItem.ID), win32ID)
			if menuItem.Checked {
				m.initiallySelectedRadioItems = append(m.initiallySelectedRadioItems, win32ID)
			}
		}
	case menu.SeparatorType:
		err := appendWin32MenuItem(parent, flags, 0, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Menu) processRadioGroups() error {

	for _, rg := range applicationMenu.RadioGroups {
		startWailsMenuID := wailsMenuItemID(rg.Members[0])
		endWailsMenuID := wailsMenuItemID(rg.Members[len(rg.Members)-1])

		startIDs := globalRadioGroupMap.getRadioGroupMapping(startWailsMenuID)
		endIDs := globalRadioGroupMap.getRadioGroupMapping(endWailsMenuID)

		var radioGroupMaps = []*radioGroupStartEnd{}
		for index := range startIDs {
			startID := startIDs[index]
			endID := endIDs[index]
			thisRadioGroup := &radioGroupStartEnd{
				startID: startID,
				endID:   endID,
			}
			radioGroupMaps = append(radioGroupMaps, thisRadioGroup)
		}

		// Set this for each member
		for _, member := range rg.Members {
			id := wailsMenuItemID(member)
			globalRadioGroupCache.addToRadioGroupCache(m.wailsMenu.Menu, id, radioGroupMaps)
		}
	}

	// Enable all initially checked radio items
	for _, win32MenuID := range m.initiallySelectedRadioItems {
		menuItemDetails := getMenuCacheEntry(win32MenuID)
		wailsMenuID := wailsMenuItemID(menuItemDetails.item.ID)
		err := selectRadioItemFromWailsMenuID(wailsMenuID, win32MenuID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Menu) Destroy() error {

	// Release the MenuIDs
	releaseMenuIDsForProcessedMenu(m.wailsMenu.Menu)

	// Unload this menu's checkboxes from the cache
	globalCheckboxCache.removeMenuFromCheckboxCache(m.wailsMenu.Menu)

	// Unload this menu's radio groups from the cache
	globalRadioGroupCache.removeMenuFromRadioBoxCache(m.wailsMenu.Menu)

	globalRadioGroupMap.removeMenuFromRadioGroupMapping(m.wailsMenu.Menu)

	// Free up callbacks
	resetCallbacks()

	// Delete menu
	return destroyWin32Menu(m.menu)
}

func (m *Menu) processAccelerator(menuitem *menumanager.ProcessedMenuItem) {

	// Add in shortcut to label if there is no "\t" override
	if !strings.Contains(menuitem.Label, "\t") {
		menuitem.Label += "\t" + keys.Stringify(menuitem.Accelerator, runtime.GOOS)
	}

	// Calculate the modifier
	var modifiers uint8
	for _, mod := range menuitem.Accelerator.Modifiers {
		switch mod {
		case keys.ControlKey, keys.CmdOrCtrlKey:
			modifiers |= 1
		case keys.OptionOrAltKey:
			modifiers |= 2
		case keys.ShiftKey:
			modifiers |= 4
		case keys.SuperKey:
			modifiers |= 8
		}
	}

	var keycode = calculateKeycode(strings.ToLower(menuitem.Accelerator.Key))
	if keycode == 0 {
		fmt.Printf("WARNING: Key '%s' is unsupported in windows. Cannot bind callback.", menuitem.Accelerator.Key)
		return
	}
	addMenuCallback(keycode, modifiers, menuitem.ID, m.menuType)

}

var flagMap = map[menu.Type]uint32{
	menu.TextType:      MF_STRING,
	menu.SeparatorType: MF_SEPARATOR,
	menu.SubmenuType:   MF_STRING | MF_POPUP,
	menu.CheckboxType:  MF_STRING,
	menu.RadioType:     MF_STRING,
}

func calculateFlags(menuItem *menumanager.ProcessedMenuItem) uint32 {
	result := flagMap[menuItem.Type]

	if menuItem.Disabled {
		result |= MF_DISABLED
	}

	if menuItem.Type == menu.CheckboxType && menuItem.Checked {
		result |= MF_CHECKED
	}

	return result
}
