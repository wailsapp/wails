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
	"unsafe"
)

var (
	// DLL stuff
	user32               = windows.NewLazySystemDLL("User32.dll")
	win32CreateMenu      = user32.NewProc("CreateMenu")
	win32CreatePopupMenu = user32.NewProc("CreatePopupMenu")
	win32AppendMenuW     = user32.NewProc("AppendMenuW")
	win32SetMenu         = user32.NewProc("SetMenu")
	win32CheckMenuItem   = user32.NewProc("CheckMenuItem")
	win32GetMenuState    = user32.NewProc("GetMenuState")
	applicationMenu      *menumanager.WailsMenu
	menuManager          *menumanager.Manager
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

// Credit: https://github.com/getlantern/systray/blob/2c0986dda9aea361e925f90e848d9036be7b5367/systray_windows.go
type menuItemInfo struct {
	Size, Mask, Type, State     uint32
	ID                          uint32
	SubMenu, Checked, Unchecked windows.Handle
	ItemData                    uintptr
	TypeData                    *uint16
	Cch                         uint32
	BMPItem                     windows.Handle
}

const (
	appMenuType menuType = iota
	contextMenuType
	trayMenuType
)

type menuCacheEntry struct {
	parent   uintptr
	menuType menuType
	index    int
	item     *menumanager.ProcessedMenuItem
}

var menuCache = []menuCacheEntry{}

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

//export createApplicationMenu
func createApplicationMenu(hwnd uintptr) {
	if applicationMenu == nil {
		return
	}

	// Create top level menu bar
	menubar, err := createMenu()
	if err != nil {
		log.Fatal("createMenu:", err.Error())
	}

	// Process top level menus
	for index, toplevelmenu := range applicationMenu.Menu.Items {
		err = processMenuItem(menubar, toplevelmenu, appMenuType, index)
		if err != nil {
			log.Fatal(err)
		}
	}

	res, _, err := win32SetMenu.Call(hwnd, menubar)
	if res == 0 {
		log.Fatal("setmenu", err)
	}
}

//export menuClicked
func menuClicked(id uint32) {

	// Get the menu from the cache
	menuitem := menuCache[id]

	if menuitem.item.Type == menu.CheckboxType {

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

func processMenuItem(parent uintptr, menuItem *menumanager.ProcessedMenuItem, menuType menuType, index int) error {

	// Ignore hidden items
	if menuItem.Hidden {
		return nil
	}

	// Add menuitem to cache
	ID := len(menuCache)

	// Calculate the flags for this menu item
	flags := uintptr(calculateFlags(menuItem))

	switch menuItem.Type {
	case menu.SubmenuType:
		submenu, err := createPopupMenu()
		if err != nil {
			return err
		}
		for index, submenuItem := range menuItem.SubMenu.Items {
			err = processMenuItem(submenu, submenuItem, menuType, index)
			if err != nil {
				return err
			}
		}
		err = appendMenuItem(parent, flags, submenu, menuItem.Label)
		if err != nil {
			return err
		}
	case menu.TextType, menu.CheckboxType:
		err := appendMenuItem(parent, flags, uintptr(ID), menuItem.Label)
		if err != nil {
			return err
		}
		menuCacheItem := menuCacheEntry{
			parent:   parent,
			menuType: menuType,
			index:    index,
			item:     menuItem,
		}
		menuCache = append(menuCache, menuCacheItem)
	case menu.SeparatorType:
		err := appendMenuItem(parent, flags, 0, "")
		if err != nil {
			return err
		}
	}
	return nil
}
