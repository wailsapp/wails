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
void* createNSMenu(char* label) {
	NSMenu *menu = [[NSMenu alloc] init];
	if( label != NULL && strlen(label) > 0 ) {
		menu.title = [NSString stringWithUTF8String:label];
		free(label);
	}
	[menu setAutoenablesItems:NO];
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

// Add services menu
static void addServicesMenu(void* menu) {
	NSMenu *nsMenu = (__bridge NSMenu *)menu;
	[NSApp setServicesMenu:nsMenu];
}


*/
import "C"
import (
	"fmt"
	"unsafe"
)

type macosMenu struct {
	menu *Menu

	nsMenu unsafe.Pointer
}

func newMenuImpl(menu *Menu) *macosMenu {
	result := &macosMenu{
		menu: menu,
		nsMenu: unsafe.Pointer(C.createNSMenu(C.CString(menu.label))),
	}
	return result
}

func (m *macosMenu) update() {
	fmt.Println("macosMenu.update()")
	if m.nsMenu == nil {
		m.nsMenu = C.createNSMenu(C.CString(m.menu.label))
	} else {
		C.clearMenu(m.nsMenu)
	}
}

func (m *macosMenu) addMenuItem(parent *Menu, menu *MenuItem) {
	C.addMenuItem(unsafe.Pointer(parent.impl.(*macosMenu).nsMenu),
		unsafe.Pointer((menu.impl).(*macosMenuItem).nsMenuItem))
}

func (l *macosMenu) addMenuItemSubMenu(item *MenuItem, menu *Menu) {
	C.setMenuItemSubmenu(unsafe.Pointer((item.impl).(*macosMenuItem).nsMenuItem),
		unsafe.Pointer((menu.impl).(*macosMenu).nsMenu))
}

func (l *macosMenu) addMenuSeparator(menu *Menu) {
	C.addMenuSeparator(unsafe.Pointer(menu.impl.(*macosMenu).nsMenu))
}

func (l *macosMenu) addServicesMenu(menu *Menu) {
	C.addServicesMenu(unsafe.Pointer(menu.impl.(*macosMenu).nsMenu))
}

func (l *macosMenu) createMenu(name string) *Menu {
	impl := newMenuImpl(&Menu{label: name})
	menu := &Menu{
		label:  name,
		items:  []*MenuItem{},
		impl: impl,
	}
	impl.menu = menu
	return menu
}
