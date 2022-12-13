//go:build darwin

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.10 -x objective-c
#cgo LDFLAGS: -framework Cocoa

#include "menuitem.h"

extern void setMenuItemChecked(void*, unsigned int, bool);

// Clear and release all menu items in the menu
void clearMenu(void* nsMenu) {
	NSMenu *menu = (NSMenu *)nsMenu;
	[menu removeAllItems];
}


// Create a new NSMenu
void* createNSMenu() {
	NSMenu *menu = [[NSMenu alloc] init];
	return (void*)menu;
}

void addMenuItem(void* nsMenu, void* nsMenuItem) {
	NSMenu *menu = (NSMenu *)nsMenu;
	[menu addItem:nsMenuItem];
}

// add seperator to menu
void addMenuSeparator(void* nsMenu) {
	NSMenu *menu = (NSMenu *)nsMenu;
	[menu addItem:[NSMenuItem separatorItem]];
}

*/
import "C"
import "unsafe"

type macosMenu struct {
	menu *Menu

	nsMenu unsafe.Pointer
}

func newMenuImpl(menu *Menu) *macosMenu {
	result := &macosMenu{
		menu: menu,
	}
	return result
}

func (m *macosMenu) update() {
	if m.nsMenu == nil {
		m.nsMenu = C.createNSMenu()
	} else {
		C.clearMenu(m.nsMenu)
	}
	for _, item := range m.menu.items {
		switch item.itemType {
		case text, checkbox:
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			C.addMenuItem(m.nsMenu, menuItem.nsMenuItem)
		case separator:
			C.addMenuSeparator(m.nsMenu)
		}

	}
}
