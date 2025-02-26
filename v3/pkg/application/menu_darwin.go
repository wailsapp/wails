//go:build darwin

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.10 -x objective-c
#cgo LDFLAGS: -framework Cocoa

#include "menuitem_darwin.h"

extern void setMenuItemChecked(void*, unsigned int, bool);
extern void setMenuItemBitmap(void*, unsigned char*, int);

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
	InvokeSync(func() {
		if m.nsMenu == nil {
			m.nsMenu = C.createNSMenu(C.CString(m.menu.label))
		} else {
			C.clearMenu(m.nsMenu)
		}
		m.processMenu(m.nsMenu, m.menu)
	})
}

func (m *macosMenu) processMenu(parent unsafe.Pointer, menu *Menu) {
	for _, item := range menu.items {
		if item.hidden {
			continue
		}
		switch item.itemType {
		case submenu:
			submenu := item.submenu
			nsSubmenu := C.createNSMenu(C.CString(item.label))
			m.processMenu(nsSubmenu, submenu)
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			C.addMenuItem(parent, menuItem.nsMenuItem)
			C.setMenuItemSubmenu(menuItem.nsMenuItem, nsSubmenu)
			if item.role == ServicesMenu {
				C.addServicesMenu(nsSubmenu)
			}
		case text, checkbox, radio:
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			C.addMenuItem(parent, menuItem.nsMenuItem)
		case separator:
			C.addMenuSeparator(parent)
		}
		if item.bitmap != nil {
			macMenuItem := item.impl.(*macosMenuItem)
			C.setMenuItemBitmap(macMenuItem.nsMenuItem, (*C.uchar)(&item.bitmap[0]), C.int(len(item.bitmap)))
		}

	}
}

func DefaultApplicationMenu() *Menu {
	menu := NewMenu()
	menu.AddRole(AppMenu)
	menu.AddRole(FileMenu)
	menu.AddRole(EditMenu)
	menu.AddRole(ViewMenu)
	menu.AddRole(WindowMenu)
	menu.AddRole(HelpMenu)
	return menu
}
