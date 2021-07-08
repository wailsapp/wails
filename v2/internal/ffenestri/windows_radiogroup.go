//+build windows

package ffenestri

import (
	"sync"

	"github.com/wailsapp/wails/v2/internal/menumanager"
)

/* ---------------------------------------------------------------------------------

Radio Groups
------------
Radio groups are stored by the ProcessedMenu as a list of menu ids.
Windows only cares about the start and end ids of the group so we
preprocess the radio groups and store this data in a radioGroupMap.
When a radio button is clicked, we use the menu id to read in the
radio group data and call CheckMenuRadioItem to update the group.

*/

type radioGroupStartEnd struct {
	startID win32MenuItemID
	endID   win32MenuItemID
}

type RadioGroupCache struct {
	cache map[*menumanager.ProcessedMenu]map[wailsMenuItemID][]*radioGroupStartEnd
	mutex sync.RWMutex
}

func NewRadioGroupCache() *RadioGroupCache {
	return &RadioGroupCache{
		cache: make(map[*menumanager.ProcessedMenu]map[wailsMenuItemID][]*radioGroupStartEnd),
	}
}

func (c *RadioGroupCache) addToRadioGroupCache(menu *menumanager.ProcessedMenu, item wailsMenuItemID, radioGroupMaps []*radioGroupStartEnd) {

	c.mutex.Lock()

	// Get map for menu
	if c.cache[menu] == nil {
		c.cache[menu] = make(map[wailsMenuItemID][]*radioGroupStartEnd)
	}
	menuMap := c.cache[menu]

	// Ensure we have a slice
	if menuMap[item] == nil {
		menuMap[item] = []*radioGroupStartEnd{}
	}

	menuMap[item] = radioGroupMaps

	c.mutex.Unlock()

}

func (c *RadioGroupCache) removeMenuFromRadioBoxCache(menu *menumanager.ProcessedMenu) {
	c.mutex.Lock()
	delete(c.cache, menu)
	c.mutex.Unlock()
}

func (c *RadioGroupCache) getRadioGroupMappings(wailsMenuID wailsMenuItemID) []*radioGroupStartEnd {
	c.mutex.Lock()
	result := []*radioGroupStartEnd{}
	for _, menugroups := range c.cache {
		groups := menugroups[wailsMenuID]
		if groups != nil {
			result = append(result, groups...)
		}
	}
	c.mutex.Unlock()
	return result
}

type RadioGroupMap struct {
	cache map[*menumanager.ProcessedMenu]map[wailsMenuItemID][]win32MenuItemID
	mutex sync.RWMutex
}

func NewRadioGroupMap() *RadioGroupMap {
	return &RadioGroupMap{
		cache: make(map[*menumanager.ProcessedMenu]map[wailsMenuItemID][]win32MenuItemID),
	}
}

func (m *RadioGroupMap) addRadioGroupMapping(menu *menumanager.ProcessedMenu, item wailsMenuItemID, win32ID win32MenuItemID) {
	m.mutex.Lock()

	// Get map for menu
	if m.cache[menu] == nil {
		m.cache[menu] = make(map[wailsMenuItemID][]win32MenuItemID)
	}
	menuMap := m.cache[menu]

	// Ensure we have a slice
	if menuMap[item] == nil {
		menuMap[item] = []win32MenuItemID{}
	}

	menuMap[item] = append(menuMap[item], win32ID)

	m.mutex.Unlock()
}

func (m *RadioGroupMap) getRadioGroupMapping(wailsMenuID wailsMenuItemID) []win32MenuItemID {
	m.mutex.Lock()
	result := []win32MenuItemID{}
	for _, menuids := range m.cache {
		ids := menuids[wailsMenuID]
		if ids != nil {
			result = append(result, ids...)
		}
	}
	m.mutex.Unlock()
	return result
}

func selectRadioItemFromWailsMenuID(wailsMenuID wailsMenuItemID, win32MenuID win32MenuItemID) {
	radioItemGroups := globalRadioGroupCache.getRadioGroupMappings(wailsMenuID)
	// Figure out offset into group
	var offset win32MenuItemID = 0
	for _, radioItemGroup := range radioItemGroups {
		if win32MenuID >= radioItemGroup.startID && win32MenuID <= radioItemGroup.endID {
			offset = win32MenuID - radioItemGroup.startID
			break
		}
	}
	for _, radioItemGroup := range radioItemGroups {
		selectedMenuID := radioItemGroup.startID + offset
		menuItemDetails := getMenuCacheEntry(selectedMenuID)
		selectRadioItem(selectedMenuID, radioItemGroup.startID, radioItemGroup.endID, menuItemDetails.parent)
	}
}
