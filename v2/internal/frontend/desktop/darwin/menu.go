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
	"fmt"
	"sync"
	"unsafe"

	"github.com/google/uuid"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
)

var createNSObjectMap = make(map[uint32]chan unsafe.Pointer)
var createNSObjectMapLock sync.RWMutex

func waitNSObjectCreate(id uint32, fn func()) unsafe.Pointer {
	waitchan := make(chan unsafe.Pointer)
	createNSObjectMapLock.Lock()
	createNSObjectMap[id] = waitchan
	createNSObjectMapLock.Unlock()
	fn()
	result := <-waitchan
	createNSObjectMapLock.Lock()
	createNSObjectMap[id] = nil
	createNSObjectMapLock.Unlock()
	return result
}

//export objectCreated
func objectCreated(id uint32, pointer unsafe.Pointer) {
	createNSObjectMapLock.Lock()
	createNSObjectMap[id] <- pointer
	createNSObjectMapLock.Unlock()
}

func NewNSTrayMenu(context unsafe.Pointer, trayMenu *menu.TrayMenu, scalingFactor int) *NSTrayMenu {
	c := NewCalloc()
	defer c.Free()

	id := uuid.New().ID()
	nsStatusItem := waitNSObjectCreate(id, func() {
		C.NewNSStatusItem(C.int(id), C.int(trayMenu.Sizing))
	})
	result := &NSTrayMenu{
		context:       context,
		nsStatusItem:  nsStatusItem,
		scalingFactor: scalingFactor,
	}

	result.SetLabel(trayMenu.Label)
	result.SetMenu(trayMenu.Menu)
	result.SetImage(trayMenu.Image)

	return result
}

func (n *NSTrayMenu) SetImage(image *menu.TrayImage) {
	if image == nil {
		return
	}
	bitmap := image.GetBestBitmap(n.scalingFactor, false)
	if bitmap == nil {
		fmt.Printf("[Warning] No TrayMenu Image available for scaling factor %dx\n", n.scalingFactor)
		return
	}
	C.SetTrayImage(n.nsStatusItem,
		unsafe.Pointer(&bitmap[0]),
		C.int(len(bitmap)),
		bool2Cint(image.IsTemplate),
		C.int(image.Position),
	)
}

func (n *NSTrayMenu) SetMenu(menu *menu.Menu) {
	if menu == nil {
		return
	}
	theMenu := NewNSMenu(n.context, "")
	processMenu(theMenu, menu)
	C.SetTrayMenu(n.nsStatusItem, theMenu.nsmenu)
}

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

func (m *MenuItem) SetChecked(value bool) {
	C.SetMenuItemChecked(m.nsmenuitem, bool2Cint(value))
}

func (m *MenuItem) SetLabel(label string) {
	cLabel := C.CString(label)
	C.SetMenuItemLabel(m.nsmenuitem, cLabel)
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
	menuItem.Impl = result
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
	f.MenuSetApplicationMenu(f.frontendOptions.Menu)
	f.mainWindow.UpdateApplicationMenu()
}
