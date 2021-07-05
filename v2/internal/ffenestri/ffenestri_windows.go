package ffenestri

import "C"

/*

#cgo windows CXXFLAGS: -std=c++11
#cgo windows,amd64 LDFLAGS: -lgdi32 -lole32 -lShlwapi -luser32 -loleaut32 -ldwmapi

#include "ffenestri.h"

extern void DisableWindowIcon(struct Application* app);

*/
import "C"
import (
	"github.com/wailsapp/wails/v2/internal/menumanager"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"golang.org/x/sys/windows"
	"log"
	"strconv"
	"sync"
	"unsafe"
)

var (
	// DLL stuff
	user32                  = windows.NewLazySystemDLL("User32.dll")
	win32CreateMenu         = user32.NewProc("CreateMenu")
	win32CreatePopupMenu    = user32.NewProc("CreatePopupMenu")
	win32AppendMenuW        = user32.NewProc("AppendMenuW")
	win32SetMenu            = user32.NewProc("SetMenu")
	win32CheckMenuItem      = user32.NewProc("CheckMenuItem")
	win32GetMenuState       = user32.NewProc("GetMenuState")
	win32CheckMenuRadioItem = user32.NewProc("CheckMenuRadioItem")
	applicationMenu         *menumanager.WailsMenu
	menuManager             *menumanager.Manager
)

const MF_BITMAP uint32 = 0x00000004
const MF_CHECKED uint32 = 0x00000008
const MF_DISABLED uint32 = 0x00000002
const MF_ENABLED uint32 = 0x00000000
const MF_GRAYED uint32 = 0x00000001
const MF_MENUBARBREAK uint32 = 0x00000020
const MF_MENUBREAK uint32 = 0x00000040
const MF_OWNERDRAW uint32 = 0x00000100
const MF_POPUP uint32 = 0x00000010
const MF_SEPARATOR uint32 = 0x00000800
const MF_STRING uint32 = 0x00000000
const MF_UNCHECKED uint32 = 0x00000000
const MF_BYCOMMAND uint32 = 0x00000000
const MF_BYPOSITION uint32 = 0x00000400

type menuType int

const (
	appMenuType menuType = iota
	contextMenuType
	trayMenuType
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

type menuCacheEntry struct {
	parent   uintptr
	menuType menuType
	item     *menumanager.ProcessedMenuItem
}

var menuCache = map[uint32]menuCacheEntry{}
var menuCacheLock sync.RWMutex

func addMenuCacheEntry(id uint32, entry menuCacheEntry) {
	menuCacheLock.Lock()
	defer menuCacheLock.Unlock()
	menuCache[id] = entry
}

func getMenuCacheEntry(id uint32) menuCacheEntry {
	menuCacheLock.Lock()
	defer menuCacheLock.Unlock()
	return menuCache[id]
}

func (a *Application) processPlatformSettings() error {

	menuManager = a.menuManager
	config := a.config.Windows
	if config == nil {
		return nil
	}

	// Check if the webview should be transparent
	if config.WebviewIsTransparent {
		C.WebviewIsTransparent(a.app)
	}

	if config.WindowBackgroundIsTranslucent {
		C.WindowBackgroundIsTranslucent(a.app)
	}

	if config.DisableWindowIcon {
		C.DisableWindowIcon(a.app)
	}

	// Unfortunately, we need to store this in the package variable so the C callback can see it
	applicationMenu = a.menuManager.GetProcessedApplicationMenu()

	//
	//// Process tray
	//trays, err := a.menuManager.GetTrayMenus()
	//if err != nil {
	//	return err
	//}
	//if trays != nil {
	//	for _, tray := range trays {
	//		C.AddTrayMenu(a.app, a.string2CString(tray))
	//	}
	//}
	//
	//// Process context menus
	//contextMenus, err := a.menuManager.GetContextMenus()
	//if err != nil {
	//	return err
	//}
	//if contextMenus != nil {
	//	for _, contextMenu := range contextMenus {
	//		C.AddContextMenu(a.app, a.string2CString(contextMenu))
	//	}
	//}
	//
	//// Process URL Handlers
	//if a.config.Mac.URLHandlers != nil {
	//	C.HasURLHandlers(a.app)
	//}

	return nil
}

func createMenu() (uintptr, error) {
	res, _, err := win32CreateMenu.Call()
	if res == 0 {
		return 0, err
	}
	return res, nil
}

func createPopupMenu() (uintptr, error) {
	res, _, err := win32CreatePopupMenu.Call()
	if res == 0 {
		return 0, err
	}
	return res, nil
}

func appendMenuItem(menu uintptr, flags uintptr, id uintptr, label string) error {
	menuText, err := windows.UTF16PtrFromString(label)
	if err != nil {
		return err
	}
	res, _, err := win32AppendMenuW.Call(
		menu,
		flags,
		id,
		uintptr(unsafe.Pointer(menuText)),
	)
	if res == 0 {
		return err
	}
	return nil
}

/*
Radio Groups
------------
Radio groups are stored by the ProcessedMenu as a list of menu ids.
Windows only cares about the start and end ids of the group so we
preprocess the radio groups and store this data in a radioGroupCache.
When a radio button is clicked, we use the menu id to read in the
radio group data and call CheckMenuRadioItem to update the group.
*/
type radioGroupCacheEntry struct {
	startID uint32
	endID   uint32
}

var radioGroupCache = map[uint32]*radioGroupCacheEntry{}
var radioGroupCacheLock sync.RWMutex

func addRadioGroupCacheEntry(id uint32, entry *radioGroupCacheEntry) {
	radioGroupCacheLock.Lock()
	defer radioGroupCacheLock.Unlock()
	radioGroupCache[id] = entry
}

func getRadioGroupCacheEntry(id uint32) *radioGroupCacheEntry {
	radioGroupCacheLock.Lock()
	defer radioGroupCacheLock.Unlock()
	return radioGroupCache[id]
}

func mustAtoi(input string) int {
	result, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal("invalid string value for mustAtoi: %s", input)
	}
	return result
}

//export createApplicationMenu
func createApplicationMenu(hwnd uintptr) {
	if applicationMenu == nil {
		return
	}

	// Process Radio groups
	for _, rg := range applicationMenu.RadioGroups {
		startID := uint32(mustAtoi(rg.Members[0]))
		endID := uint32(mustAtoi(rg.Members[len(rg.Members)-1]))
		thisRG := &radioGroupCacheEntry{
			startID: startID,
			endID:   endID,
		}
		// Set this for each member
		for _, member := range rg.Members {
			id := uint32(mustAtoi(member))
			addRadioGroupCacheEntry(id, thisRG)
		}
	}

	// Create top level menu bar
	menubar, err := createMenu()
	if err != nil {
		log.Fatal("createMenu:", err.Error())
	}

	// Process top level menus
	for _, toplevelmenu := range applicationMenu.Menu.Items {
		err = processMenuItem(menubar, toplevelmenu, appMenuType)
		if err != nil {
			log.Fatal(err)
		}
	}

	res, _, err := win32SetMenu.Call(hwnd, menubar)
	if res == 0 {
		log.Fatal("setmenu", err)
	}
}

func mustSelectRadioItem(id uint32, parent uintptr) {
	rg := getRadioGroupCacheEntry(id)
	res, _, err := win32CheckMenuRadioItem.Call(parent, uintptr(rg.startID), uintptr(rg.endID), uintptr(id), uintptr(MF_BYCOMMAND))
	if int(res) == 0 {
		log.Fatal(err)
	}
}

//export menuClicked
func menuClicked(id uint32) {

	// Get the menu from the cache
	menuitem := getMenuCacheEntry(id)

	switch menuitem.item.Type {
	case menu.CheckboxType:
		res, _, err := win32GetMenuState.Call(menuitem.parent, uintptr(id), uintptr(MF_BYCOMMAND))
		if int(res) == -1 {
			log.Fatal(err)
		}

		if uint32(res) == MF_CHECKED {
			res, _, err = win32CheckMenuItem.Call(menuitem.parent, uintptr(id), uintptr(MF_UNCHECKED))
		} else {
			res, _, err = win32CheckMenuItem.Call(menuitem.parent, uintptr(id), uintptr(MF_CHECKED))
		}
		if int(res) == -1 {
			log.Fatal(err)
		}
	case menu.RadioType:
		mustSelectRadioItem(id, menuitem.parent)
	}

	// Print the click error - it's not fatal
	err := menuManager.ProcessClick(menuitem.item.ID, "", "ApplicationMenu", "")
	if err != nil {
		println(err.Error())
	}
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

func processMenuItem(parent uintptr, menuItem *menumanager.ProcessedMenuItem, menuType menuType) error {

	// Ignore hidden items
	if menuItem.Hidden {
		return nil
	}

	// Calculate the flags for this menu item
	flags := uintptr(calculateFlags(menuItem))

	switch menuItem.Type {
	case menu.SubmenuType:
		submenu, err := createPopupMenu()
		if err != nil {
			return err
		}
		for _, submenuItem := range menuItem.SubMenu.Items {
			err = processMenuItem(submenu, submenuItem, menuType)
			if err != nil {
				return err
			}
		}
		err = appendMenuItem(parent, flags, submenu, menuItem.Label)
		if err != nil {
			return err
		}
	case menu.TextType, menu.CheckboxType, menu.RadioType:
		ID := uint32(mustAtoi(menuItem.ID))
		err := appendMenuItem(parent, flags, uintptr(ID), menuItem.Label)
		if err != nil {
			return err
		}
		menuCacheItem := menuCacheEntry{
			parent:   parent,
			menuType: menuType,
			item:     menuItem,
		}
		addMenuCacheEntry(ID, menuCacheItem)
		if menuItem.Type == menu.RadioType && menuItem.Checked {
			mustSelectRadioItem(ID, parent)
		}
	case menu.SeparatorType:
		err := appendMenuItem(parent, flags, 0, "")
		if err != nil {
			return err
		}
	}
	return nil
}
