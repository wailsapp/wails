//go:build darwin && purego && !ios && !server

// Package application - CGO-free macOS menu backend.
//
// This is the purego counterpart of menu_darwin.go. It builds NSMenu/NSMenuItem
// hierarchies by messaging the Objective-C runtime directly (via the helpers in
// darwin_purego_cocoa.go) instead of compiling Objective-C through cgo.
package application

import "unsafe"

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

func (m *macosMenu) processMenu(parent unsafe.Pointer, menu *Menu) {
	parentMenu := idFromPtr(parent)
	for _, item := range menu.items {
		// A rebuild (Menu.Update) replaces every impl: release the previous
		// NSMenuItem's owning reference first, or each update leaks one item
		// (removeAllItems only drops the menu's retain, not ours).
		if old, ok := item.impl.(*macosMenuItem); ok && old != nil && old.nsMenuItem != nil {
			old.destroy()
		}
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
			// The NSMenuItem retains its submenu; drop our creation reference
			// so the whole subtree unwinds when the item goes away.
			nsSubmenu.send("release")
		case text, checkbox, radio, sectionHeader:
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

// menuSectionHeadersSupported reports whether NSMenuItem has the macOS 14
// +sectionHeaderWithTitle: class method. Feature detection replaces the cgo
// build's @available check.
func menuSectionHeadersSupported() bool {
	return respondsTo(class("NSMenuItem"), "sectionHeaderWithTitle:")
}

// newNSMenu creates an NSMenu with the given title.
func newNSMenu(label string) unsafe.Pointer {
	return ptrFromID(createNSMenu(label))
}

// clearNSMenu removes every item from menu.
func clearNSMenu(menu unsafe.Pointer) {
	idFromPtr(menu).send("removeAllItems")
}
