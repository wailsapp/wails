//+build windows

package ffenestri

import (
	"github.com/wailsapp/wails/v2/internal/menumanager"
	"sync"
)

/* ---------------------------------------------------------------------------------

Checkbox Cache
--------------
The checkbox cache keeps a list of IDs that are associated with the same checkbox menu item.
This can happen when a checkbox is used in an application menu and a tray menu, eg "start at login".
The cache is used to bulk toggle the menu items when one is clicked.

*/

type CheckboxCache struct {
	cache map[*menumanager.ProcessedMenu]map[wailsMenuItemID][]win32MenuItemID
	mutex sync.RWMutex
}

func NewCheckboxCache() *CheckboxCache {
	return &CheckboxCache{
		cache: make(map[*menumanager.ProcessedMenu]map[wailsMenuItemID][]win32MenuItemID),
	}
}

func (c *CheckboxCache) addToCheckboxCache(menu *menumanager.ProcessedMenu, item wailsMenuItemID, menuID win32MenuItemID) {

	// Get map for menu
	if c.cache[menu] == nil {
		c.cache[menu] = make(map[wailsMenuItemID][]win32MenuItemID)
	}
	menuMap := c.cache[menu]

	// Ensure we have a slice
	if menuMap[item] == nil {
		menuMap[item] = []win32MenuItemID{}
	}

	c.mutex.Lock()
	menuMap[item] = append(menuMap[item], menuID)
	c.mutex.Unlock()

}

func (c *CheckboxCache) removeMenuFromChechboxCache(menu *menumanager.ProcessedMenu) {
	c.mutex.Lock()
	delete(c.cache, menu)
	c.mutex.Unlock()
}

// win32MenuIDsForWailsMenuID returns all win32menuids that are used for a wails menu item id across
// all menus
func (c *CheckboxCache) win32MenuIDsForWailsMenuID(item wailsMenuItemID) []win32MenuItemID {
	c.mutex.Lock()
	result := []win32MenuItemID{}
	for _, menu := range c.cache {
		ids := menu[item]
		if ids != nil {
			result = append(result, ids...)
		}
	}
	c.mutex.Unlock()
	return result
}
