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
			m.nsMenu = C.createNSMenu(C.CString(m.menu.Label()))
		} else {
			C.clearMenu(m.nsMenu)
		}
		m.processMenu(m.nsMenu, m.menu)
	})
}

func (m *macosMenu) processMenu(parent unsafe.Pointer, menu *Menu) {
	for _, item := range menu.items {
		// Handle Menu items (submenus)
		if newSubmenu, ok := item.(*Menu); ok {
			nsSubmenu := C.createNSMenu(C.CString(newSubmenu.Label()))
			m.processMenu(nsSubmenu, newSubmenu)

			// Create a temporary MenuItem to represent this submenu in the native menu
			tempMenuItem := NewMenuItem(newSubmenu.Label())
			tempMenuItem.itemType = submenu
			menuItem := newMenuItemImpl(tempMenuItem)
			tempMenuItem.impl = menuItem
			C.addMenuItem(parent, menuItem.nsMenuItem)
			C.setMenuItemSubmenu(menuItem.nsMenuItem, nsSubmenu)
			continue
		}

		// Handle MenuItem items
		if menuItem, ok := item.(*MenuItem); ok {
			switch menuItem.itemType {
			case submenu:
				submenu := menuItem.submenu
				nsSubmenu := C.createNSMenu(C.CString(menuItem.Label()))
				m.processMenu(nsSubmenu, submenu)
				impl := newMenuItemImpl(menuItem)
				menuItem.impl = impl
				C.addMenuItem(parent, impl.nsMenuItem)
				C.setMenuItemSubmenu(impl.nsMenuItem, nsSubmenu)
				if menuItem.role == ServicesMenu {
					C.addServicesMenu(nsSubmenu)
				}
			case text, checkbox, radio:
				// Debug: Log original state
				if globalApplication != nil {
					globalApplication.debug("Processing menu item: %s (Go ID: %d, Disabled: %v, Hidden: %v, Callback: %v)",
						menuItem.Label(), menuItem.id, menuItem.disabled, menuItem.Hidden(), menuItem.callback != nil)
				}

				// Create a temporary enabled version for proper click handler setup
				tempDisabled := menuItem.disabled
				menuItem.disabled = false // Temporarily enable for creation

				impl := newMenuItemImpl(menuItem)
				menuItem.impl = impl

				// Debug: Log after native creation
				if globalApplication != nil {
					globalApplication.debug("Native menu item created for: %s (Go ID: %d)", menuItem.Label(), menuItem.id)
				}

				// Synchronize all state to the new native menu item
				impl.setDisabled(tempDisabled)
				impl.setLabel(menuItem.Label())
				impl.setDisabled(menuItem.disabled) // Apply correct disabled state
				impl.setHidden(menuItem.Hidden())
				if menuItem.checked {
					impl.setChecked(menuItem.checked)
				}
				if menuItem.accelerator != nil {
					impl.setAccelerator(menuItem.accelerator)
				}
				if len(menuItem.bitmap) > 0 {
					impl.setBitmap(menuItem.bitmap)
				}

				C.addMenuItem(parent, impl.nsMenuItem)
				// Debug: Track menu item recreation
				if globalApplication != nil {
					globalApplication.debug(
						"Recreated native menu item: %s (Go ID: %d, Hidden: %v, Disabled: %v, Callback: %v)",
						menuItem.Label(),
						menuItem.id,
						menuItem.Hidden(),
						menuItem.disabled,
						menuItem.callback != nil,
					)
				}
			case separator:
				C.addMenuSeparator(parent)
			}

			if menuItem.bitmap != nil {
				macMenuItem := menuItem.impl.(*macosMenuItem)
				C.setMenuItemBitmap(
					macMenuItem.nsMenuItem,
					(*C.uchar)(&menuItem.bitmap[0]),
					C.int(len(menuItem.bitmap)),
				)
			}
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
