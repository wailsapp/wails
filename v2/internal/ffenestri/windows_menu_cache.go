//+build windows

package ffenestri

import (
	"github.com/wailsapp/wails/v2/internal/menumanager"
	"sync"
)

/**

MenuCache
---------
When windows calls back to Go (when an item is clicked), we need to
be able to retrieve information about the menu item:
 - The menu that the menuitem is part of (parent)
 - The original processed menu item
 - The type of the menu (application, context or tray)

This cache is built up when a menu is created.

*/

// TODO: Make this like the other caches

type menuCacheEntry struct {
	parent        win32Menu
	menuType      menuType
	item          *menumanager.ProcessedMenuItem
	processedMenu *menumanager.ProcessedMenu
}

// windowsMenuIDCounter keeps track of the unique windows menu IDs
var windowsMenuIDCounter uint32

var menuCache = map[win32MenuItemID]*menuCacheEntry{}
var menuCacheLock sync.RWMutex
var wailsMenuIDtoWin32IDMap = map[wailsMenuItemID]win32MenuItemID{}

func addMenuCacheEntry(parent win32Menu, typ menuType, wailsMenuItem *menumanager.ProcessedMenuItem, processedMenu *menumanager.ProcessedMenu) win32MenuItemID {
	menuCacheLock.Lock()
	defer menuCacheLock.Unlock()
	menuID := win32MenuItemID(windowsMenuIDCounter)
	windowsMenuIDCounter++
	menuCache[menuID] = &menuCacheEntry{
		parent:        parent,
		menuType:      typ,
		item:          wailsMenuItem,
		processedMenu: processedMenu,
	}
	// save the mapping
	wailsMenuIDtoWin32IDMap[wailsMenuItemID(wailsMenuItem.ID)] = menuID
	return menuID

}

func getMenuCacheEntry(id win32MenuItemID) *menuCacheEntry {
	menuCacheLock.Lock()
	defer menuCacheLock.Unlock()
	return menuCache[id]
}
