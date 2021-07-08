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
	"log"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

// Setup the global caches
var globalCheckboxCache = NewCheckboxCache()
var globalRadioGroupCache = NewRadioGroupCache()
var globalRadioGroupMap = NewRadioGroupMap()
var globalApplicationMenu *Menu

type menuType string

const (
	appMenuType menuType = "ApplicationMenu"
	contextMenuType
	trayMenuType
)

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

func (c *Client) updateApplicationMenu() {
	applicationMenu = c.app.menuManager.GetProcessedApplicationMenu()
	createApplicationMenu(uintptr(C.GetWindowHandle(c.app.app)))
}

/* ---------------------------------------------------------------------------------

Application Menu
----------------
There's only 1 application menu and this is where we create it. This method
is called from C after the window is created and the WM_CREATE message has
been sent.

*/

func checkFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//export createApplicationMenu
func createApplicationMenu(hwnd uintptr) {

	if applicationMenu == nil {
		return
	}

	if globalApplicationMenu != nil {
		checkFatal(globalApplicationMenu.Destroy())
	}

	var err error
	globalApplicationMenu, err = createMenu(applicationMenu, appMenuType)
	checkFatal(err)

	err = setWindowMenu(win32Window(hwnd), globalApplicationMenu.menu)
	checkFatal(err)
}

/*
This method is called by C when a menu item is pressed
*/

//export menuClicked
func menuClicked(id uint32) {
	win32MenuID := win32MenuItemID(id)
	//println("Got click from menu id", win32MenuID)

	// Get the menu from the cache
	menuItemDetails := getMenuCacheEntry(win32MenuID)
	wailsMenuID := wailsMenuItemID(menuItemDetails.item.ID)

	//println("Got click from menu id", win32MenuID, "- wails menu ID", wailsMenuID)
	//spew.Dump(menuItemDetails)

	switch menuItemDetails.item.Type {
	case menu.CheckboxType:

		// Determine if the menu is set or not
		res, _, err := win32GetMenuState.Call(uintptr(menuItemDetails.parent), uintptr(id), uintptr(MF_BYCOMMAND))
		if int(res) == -1 {
			log.Fatal(err)
		}

		flag := MF_CHECKED
		if uint32(res) == MF_CHECKED {
			flag = MF_UNCHECKED
		}

		for _, menuid := range globalCheckboxCache.win32MenuIDsForWailsMenuID(wailsMenuID) {
			//println("setting menuid", menuid, "with flag", flag)
			menuItemDetails := getMenuCacheEntry(menuid)
			res, _, err = win32CheckMenuItem.Call(uintptr(menuItemDetails.parent), uintptr(menuid), uintptr(flag))
			if int(res) == -1 {
				log.Fatal(err)
			}
		}
	case menu.RadioType:
		selectRadioItemFromWailsMenuID(wailsMenuID, win32MenuID)
	}

	// Print the click error - it's not fatal
	err := menuManager.ProcessClick(menuItemDetails.item.ID, "", string(menuItemDetails.menuType), "")
	if err != nil {
		println(err.Error())
	}
}
