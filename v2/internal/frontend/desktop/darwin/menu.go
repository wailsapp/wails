//go:build darwin
// +build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit
#import <Foundation/Foundation.h>
#import "Application.h"
#import "WailsContext.h"

#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
)

type NSMenu struct {
	context unsafe.Pointer
	nsmenu  unsafe.Pointer
}

func NewNSMenu(context unsafe.Pointer, name string) *NSMenu {
	c := NewCalloc()
	defer c.Free()
	title := c.String(name)
	nsmenu := C.NewMenu(title)
	return &NSMenu{
		context: context,
		nsmenu:  nsmenu,
	}
}

func (m *NSMenu) AddSubMenu(label string) *NSMenu {
	result := NewNSMenu(m.context, label)
	C.AppendSubmenu(m.nsmenu, result.nsmenu)
	return result
}

func (m *NSMenu) AppendRole(role menu.Role) {
	C.AppendRole(m.context, m.nsmenu, C.int(role))
}

type MenuItem struct {
	id                uint
	nsmenuitem        unsafe.Pointer
	wailsMenuItem     *menu.MenuItem
	radioGroupMembers []*MenuItem
}

func (m *NSMenu) AddMenuItem(menuItem *menu.MenuItem) *MenuItem {
	c := NewCalloc()
	defer c.Free()
	var modifier C.int
	var key *C.char
	if menuItem.Accelerator != nil {
		modifier = C.int(keys.ToMacModifier(menuItem.Accelerator))
		key = c.String(menuItem.Accelerator.Key)
	}

	result := &MenuItem{
		wailsMenuItem: menuItem,
	}

	result.id = createMenuItemID(result)
	result.nsmenuitem = C.AppendMenuItem(m.context, m.nsmenu, c.String(menuItem.Label), key, modifier, bool2Cint(menuItem.Disabled), bool2Cint(menuItem.Checked), C.int(result.id))
	return result
}

//func (w *Window) SetApplicationMenu(menu *menu.Menu) {
//w.applicationMenu = menu
//processMenu(w, menu)
//}

func processMenu(parent *NSMenu, wailsMenu *menu.Menu) {
	var radioGroups []*MenuItem

	for _, menuItem := range wailsMenu.Items {
		if menuItem.SubMenu != nil {
			if len(radioGroups) > 0 {
				processRadioGroups(radioGroups)
				radioGroups = []*MenuItem{}
			}
			submenu := parent.AddSubMenu(menuItem.Label)
			processMenu(submenu, menuItem.SubMenu)
		} else {
			lastMenuItem := processMenuItem(parent, menuItem)
			if menuItem.Type == menu.RadioType {
				radioGroups = append(radioGroups, lastMenuItem)
			} else {
				if len(radioGroups) > 0 {
					processRadioGroups(radioGroups)
					radioGroups = []*MenuItem{}
				}
			}
		}
	}
}

func processRadioGroups(groups []*MenuItem) {
	for _, item := range groups {
		item.radioGroupMembers = groups
	}
}

func processMenuItem(parent *NSMenu, menuItem *menu.MenuItem) *MenuItem {
	if menuItem.Hidden {
		return nil
	}
	if menuItem.Role != 0 {
		parent.AppendRole(menuItem.Role)
		return nil
	}
	if menuItem.Type == menu.SeparatorType {
		C.AppendSeparator(parent.nsmenu)
		return nil
	}

	return parent.AddMenuItem(menuItem)

}

func (f *Frontend) MenuSetApplicationMenu(menu *menu.Menu) {
	f.mainWindow.SetApplicationMenu(menu)
}

func (f *Frontend) MenuUpdateApplicationMenu() {
	f.mainWindow.UpdateApplicationMenu()
}
