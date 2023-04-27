//go:build linux && purego

package application

import (
	"fmt"

	"github.com/ebitengine/purego"
)

type linuxMenu struct {
	menu   *Menu
	native uintptr
}

func newMenuImpl(menu *Menu) *linuxMenu {
	var newMenuBar func() uintptr
	purego.RegisterLibFunc(&newMenuBar, gtk, "gtk_menu_bar_new")
	result := &linuxMenu{
		menu:   menu,
		native: newMenuBar(),
	}
	return result
}

func (m *linuxMenu) update() {
	m.processMenu(m.menu)
}

func (m *linuxMenu) processMenu(menu *Menu) {
	var newMenu func() uintptr
	purego.RegisterLibFunc(&newMenu, gtk, "gtk_menu_new")
	if menu.impl == nil {
		menu.impl = &linuxMenu{
			menu:   menu,
			native: newMenu(),
		}
	}
	var currentRadioGroup uintptr

	for _, item := range menu.items {
		// drop the group if we have run out of radio items
		if item.itemType != radio {
			currentRadioGroup = 0
		}

		switch item.itemType {
		case submenu:
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			m.processMenu(item.submenu)
			m.addSubMenuToItem(item.submenu, item)
			m.addMenuItem(menu, item)
		case text, checkbox:
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			m.addMenuItem(menu, item)
		case radio:
			menuItem := newRadioItemImpl(item, currentRadioGroup)
			item.impl = menuItem
			m.addMenuItem(menu, item)

			var radioGetGroup func(uintptr) uintptr
			purego.RegisterLibFunc(&radioGetGroup, gtk, "gtk_radio_menu_item_get_group")

			currentRadioGroup = radioGetGroup(menuItem.native)
		case separator:
			m.addMenuSeparator(menu)
		}

	}

	for _, item := range menu.items {
		if item.callback != nil {
			m.attachHandler(item)
		}
	}

}

func (m *linuxMenu) attachHandler(item *MenuItem) {
	impl := (item.impl).(*linuxMenuItem)
	widget := impl.native
	flags := 0

	var handleClick = func() {
		item := item
		switch item.itemType {
		case text, checkbox:
			processMenuItemClick(item.id)
		case radio:
			menuItem := (item.impl).(*linuxMenuItem)
			if menuItem.isChecked() {
				processMenuItemClick(item.id)
			}
		default:
			fmt.Println("handleClick", item.itemType, item.id)
		}
	}

	var signalConnectObject func(uintptr, string, uintptr, uintptr, int) uint
	purego.RegisterLibFunc(&signalConnectObject, gtk, "g_signal_connect_object")
	handlerId := signalConnectObject(
		widget,
		"activate",
		purego.NewCallback(handleClick),
		widget,
		flags)

	impl.handlerId = handlerId
}

func (m *linuxMenu) addSubMenuToItem(menu *Menu, item *MenuItem) {
	var newMenu func() uintptr
	purego.RegisterLibFunc(&newMenu, gtk, "gtk_menu_new")
	if menu.impl == nil {
		menu.impl = &linuxMenu{
			menu:   menu,
			native: newMenu(),
		}
	}
	var itemSetSubmenu func(uintptr, uintptr)
	purego.RegisterLibFunc(&itemSetSubmenu, gtk, "gtk_menu_item_set_submenu")

	itemSetSubmenu(
		(item.impl).(*linuxMenuItem).native,
		(menu.impl).(*linuxMenu).native)

	if item.role == ServicesMenu {
		// FIXME: what does this mean?
	}
}

func (m *linuxMenu) addMenuItem(parent *Menu, menu *MenuItem) {
	var shellAppend func(uintptr, uintptr)
	purego.RegisterLibFunc(&shellAppend, gtk, "gtk_menu_shell_append")
	shellAppend(
		(parent.impl).(*linuxMenu).native,
		(menu.impl).(*linuxMenuItem).native,
	)
}

func (m *linuxMenu) addMenuSeparator(menu *Menu) {
	var newSeparator func() uintptr
	purego.RegisterLibFunc(&newSeparator, gtk, "gtk_separator_menu_item_new")
	var shellAppend func(uintptr, uintptr)
	purego.RegisterLibFunc(&shellAppend, gtk, "gtk_menu_shell_append")

	sep := newSeparator()
	native := (menu.impl).(*linuxMenu).native
	shellAppend(native, sep)
}

func (m *linuxMenu) addServicesMenu(menu *Menu) {
	fmt.Println("addServicesMenu - not implemented")
	//C.addServicesMenu(unsafe.Pointer(menu.impl.(*linuxMenu).nsMenu))
}

func (l *linuxMenu) createMenu(name string, items []*MenuItem) *Menu {
	impl := newMenuImpl(&Menu{label: name})
	menu := &Menu{
		label: name,
		items: items,
		impl:  impl,
	}
	impl.menu = menu
	return menu
}

func defaultApplicationMenu() *Menu {
	menu := NewMenu()
	menu.AddRole(AppMenu)
	menu.AddRole(FileMenu)
	menu.AddRole(EditMenu)
	menu.AddRole(ViewMenu)
	menu.AddRole(WindowMenu)
	menu.AddRole(HelpMenu)
	return menu
}
