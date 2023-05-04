//go:build windows

package application

import "unsafe"

type windowsMenu struct {
	menu *Menu

	menuImpl unsafe.Pointer
}

func newMenuImpl(menu *Menu) *windowsMenu {
	result := &windowsMenu{
		menu: menu,
	}
	return result
}

func (m *windowsMenu) update() {
	//if m.menuImpl == nil {
	//	m.menuImpl = C.createNSMenu(C.CString(m.menu.label))
	//} else {
	//	C.clearMenu(m.menuImpl)
	//}
	m.processMenu(m.menuImpl, m.menu)
}

func (m *windowsMenu) processMenu(parent unsafe.Pointer, menu *Menu) {
	//for _, item := range menu.items {
	//	switch item.itemType {
	//	case submenu:
	//		submenu := item.submenu
	//		nsSubmenu := C.createNSMenu(C.CString(item.label))
	//		m.processMenu(nsSubmenu, submenu)
	//		menuItem := newMenuItemImpl(item)
	//		item.impl = menuItem
	//		C.addMenuItem(parent, menuItem.nsMenuItem)
	//		C.setMenuItemSubmenu(menuItem.nsMenuItem, nsSubmenu)
	//		if item.role == ServicesMenu {
	//			C.addServicesMenu(nsSubmenu)
	//		}
	//	case text, checkbox, radio:
	//		menuItem := newMenuItemImpl(item)
	//		item.impl = menuItem
	//		C.addMenuItem(parent, menuItem.nsMenuItem)
	//	case separator:
	//		C.addMenuSeparator(parent)
	//	}
	//
	//}
}

func defaultApplicationMenu() *Menu {
	menu := NewMenu()
	menu.AddRole(FileMenu)
	menu.AddRole(EditMenu)
	menu.AddRole(ViewMenu)
	menu.AddRole(WindowMenu)
	menu.AddRole(HelpMenu)
	return menu
}
