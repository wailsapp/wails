//go:build darwin && purego && !ios && !server

// Package application - CGO-free macOS menu backend.
//
// This is the purego counterpart of menu_darwin.go. It builds NSMenu/NSMenuItem
// hierarchies by messaging the Objective-C runtime directly (via the helpers in
// darwin_purego_cocoa.go) instead of compiling Objective-C through cgo.
package application

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

// createNSMenu allocates an NSMenu, optionally titles it, and disables
// automatic item enabling so we control enabled state ourselves. The returned
// object is +1 retained and stays alive for the lifetime of the owning menu.
func createNSMenu(label string) id {
	menu := class("NSMenu").send("alloc").send("init")
	if label != "" {
		menu.send("setTitle:", nsString(label))
	}
	menu.send("setAutoenablesItems:", false)
	return menu
}

func (m *macosMenu) update() {
	InvokeSync(func() {
		if m.nsMenu == nil {
			m.nsMenu = ptrFromID(createNSMenu(m.menu.label))
		} else {
			idFromPtr(m.nsMenu).send("removeAllItems")
		}
		m.processMenu(m.nsMenu, m.menu)
	})
}

func (m *macosMenu) processMenu(parent unsafe.Pointer, menu *Menu) {
	parentMenu := idFromPtr(parent)
	for _, item := range menu.items {
		switch item.itemType {
		case submenu:
			submenu := item.submenu
			nsSubmenu := createNSMenu(item.label)
			m.processMenu(ptrFromID(nsSubmenu), submenu)
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			nsItem := idFromPtr(menuItem.nsMenuItem)
			parentMenu.send("addItem:", nsItem)
			nsItem.send("setSubmenu:", nsSubmenu)
			if item.role == ServicesMenu {
				nsApp().send("setServicesMenu:", nsSubmenu)
			}
			if item.role == WindowMenu {
				nsApp().send("setWindowsMenu:", nsSubmenu)
			}
		case text, checkbox, radio:
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			if item.hidden {
				menuItem.setHidden(true)
			}
			parentMenu.send("addItem:", idFromPtr(menuItem.nsMenuItem))
		case separator:
			parentMenu.send("addItem:", class("NSMenuItem").send("separatorItem"))
		}
		if item.bitmap != nil {
			macMenuItem := item.impl.(*macosMenuItem)
			macMenuItem.setBitmap(item.bitmap)
		}
	}
}

// nsApp returns the shared NSApplication instance.
func nsApp() id {
	return class("NSApplication").send("sharedApplication")
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
