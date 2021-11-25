//go:build windows
// +build windows

package ffenestri

import (
	"fmt"
	"github.com/leaanthony/slicer"
	"os"
	"sync"
	"text/tabwriter"

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

func (c *RadioGroupCache) Dump() {
	// Start a new tabwriter
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	println("---------------- RadioGroupCache", c, "Dump ----------------")
	for menu, processedMenu := range c.cache {
		println("Menu", menu)
		_, _ = fmt.Fprintf(w, "Wails ID \tWindows ID Pairs\n")
		for wailsMenuItemID, radioGroupStartEnd := range processedMenu {
			menus := slicer.String()
			for _, se := range radioGroupStartEnd {
				menus.Add(fmt.Sprintf("[%d -> %d]", se.startID, se.endID))
			}
			_, _ = fmt.Fprintf(w, "%s\t%s\n", wailsMenuItemID, menus.Join(", "))
			_ = w.Flush()
		}
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

func (c *RadioGroupMap) Dump() {
	// Start a new tabwriter
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	println("---------------- RadioGroupMap", c, "Dump ----------------")
	for _, processedMenu := range c.cache {
		_, _ = fmt.Fprintf(w, "Menu\tWails ID \tWindows IDs\n")
		for wailsMenuItemID, win32menus := range processedMenu {
			menus := slicer.String()
			for _, win32menu := range win32menus {
				menus.Add(fmt.Sprintf("%v", win32menu))
			}
			_, _ = fmt.Fprintf(w, "%p\t%s\t%s\n", processedMenu, wailsMenuItemID, menus.Join(", "))
			_ = w.Flush()
		}
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

func (m *RadioGroupMap) removeMenuFromRadioGroupMapping(menu *menumanager.ProcessedMenu) {
	m.mutex.Lock()
	delete(m.cache, menu)
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

func selectRadioItemFromWailsMenuID(wailsMenuID wailsMenuItemID, win32MenuID win32MenuItemID) error {
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
		if menuItemDetails != nil {
			if menuItemDetails.parent != 0 {
				err := selectRadioItem(selectedMenuID, radioItemGroup.startID, radioItemGroup.endID, menuItemDetails.parent)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
