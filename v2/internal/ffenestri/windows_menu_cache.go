//go:build windows
// +build windows

package ffenestri

import (
	"github.com/leaanthony/idgen"
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

var idGenerator = idgen.New()

var menuCache = map[win32MenuItemID]*menuCacheEntry{}
var menuCacheLock sync.RWMutex
var wailsMenuIDtoWin32IDMap = map[wailsMenuItemID]win32MenuItemID{}

// This releases the menuIDs back to the id generator
var winIDsOwnedByProcessedMenu = map[*menumanager.ProcessedMenu][]win32MenuItemID{}

func releaseMenuIDsForProcessedMenu(processedMenu *menumanager.ProcessedMenu) {
	for _, menuID := range winIDsOwnedByProcessedMenu[processedMenu] {
		idGenerator.ReleaseID(uint(menuID))
	}
	delete(winIDsOwnedByProcessedMenu, processedMenu)
}

func addMenuCacheEntry(parent win32Menu, typ menuType, wailsMenuItem *menumanager.ProcessedMenuItem, processedMenu *menumanager.ProcessedMenu) win32MenuItemID {
	menuCacheLock.Lock()
	defer menuCacheLock.Unlock()
	id, err := idGenerator.NewID()
	checkFatal(err)
	menuID := win32MenuItemID(id)
	menuCache[menuID] = &menuCacheEntry{
		parent:        parent,
		menuType:      typ,
		item:          wailsMenuItem,
		processedMenu: processedMenu,
	}
	// save the mapping
	wailsMenuIDtoWin32IDMap[wailsMenuItemID(wailsMenuItem.ID)] = menuID
	// keep track of menuids owned by this menu (so we can release the ids)
	winIDsOwnedByProcessedMenu[processedMenu] = append(winIDsOwnedByProcessedMenu[processedMenu], menuID)
	return menuID

}

func getMenuCacheEntry(id win32MenuItemID) *menuCacheEntry {
	menuCacheLock.Lock()
	defer menuCacheLock.Unlock()
	return menuCache[id]
}
