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

// Set the submenu of a menu item
void setMenuItemSubmenu(void* nsMenuItem, void* nsMenu) {
	NSMenuItem *menuItem = (NSMenuItem *)nsMenuItem;
	NSMenu *menu = (NSMenu *)nsMenu;
	[menuItem setSubmenu:menu];
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
	m.processMenu(m.nsMenu, m.menu)
}

func (m *macosMenu) processMenu(parent unsafe.Pointer, menu *Menu) {
	for _, item := range menu.items {
		switch item.itemType {
		case submenu:
			submenu := item.submenu
			nsSubmenu := C.createNSMenu()
			m.processMenu(nsSubmenu, submenu)
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			C.addMenuItem(parent, menuItem.nsMenuItem)
			C.setMenuItemSubmenu(menuItem.nsMenuItem, nsSubmenu)
		case text, checkbox, radio:
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			C.addMenuItem(parent, menuItem.nsMenuItem)
		case separator:
			C.addMenuSeparator(parent)
		}

	}
	//C.setMenu(parent, m.nsMenu)
}
