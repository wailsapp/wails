//go:build linux

package application

type linuxMenu struct {
	menu   *Menu
	native pointer
}

func newMenuImpl(menu *Menu) *linuxMenu {
	result := &linuxMenu{
		menu:   menu,
		native: menuBarNew(),
	}
	return result
}

func (m *linuxMenu) run() {
	m.update()
}

func (m *linuxMenu) update() {
	m.processMenu(m.menu)
}

func (m *linuxMenu) processMenu(menu *Menu) {
	if menu.impl == nil {
		menu.impl = &linuxMenu{
			menu:   menu,
			native: menuNew(),
		}
	}
	var currentRadioGroup GSListPointer

	for _, item := range menu.items {
		// Handle Menu items (submenus)
		if submenu, ok := item.(*Menu); ok {
			// Create a temporary MenuItem to represent this submenu
			tempMenuItem := NewMenuItem(submenu.Label())
			tempMenuItem.itemType = submenu
			menuItem := newMenuItemImpl(tempMenuItem)
			tempMenuItem.impl = menuItem
			m.processMenu(submenu)
			m.addSubMenuToItem(submenu, tempMenuItem)
			m.addMenuItem(menu, tempMenuItem)
			continue
		}

		// Handle MenuItem items
		if menuItem, ok := item.(*MenuItem); ok {
			// drop the group if we have run out of radio items
			if menuItem.itemType != radio {
				currentRadioGroup = nilRadioGroup
			}

			switch menuItem.itemType {
			case submenu:
				impl := newMenuItemImpl(menuItem)
				menuItem.impl = impl
				m.processMenu(menuItem.submenu)
				m.addSubMenuToItem(menuItem.submenu, menuItem)
				m.addMenuItem(menu, menuItem)
			case text, checkbox:
				impl := newMenuItemImpl(menuItem)
				menuItem.impl = impl

				// Synchronize all state to the new native menu item
				impl.setLabel(menuItem.Label())
				impl.setDisabled(menuItem.disabled)
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

				m.addMenuItem(menu, menuItem)
			case radio:
				impl := newRadioItemImpl(menuItem, currentRadioGroup)
				menuItem.impl = impl

				// Synchronize state for radio items
				impl.setDisabled(menuItem.disabled)
				impl.setHidden(menuItem.Hidden())
				if menuItem.accelerator != nil {
					impl.setAccelerator(menuItem.accelerator)
				}

				m.addMenuItem(menu, menuItem)
				currentRadioGroup = menuGetRadioGroup(impl)
			case separator:
				m.addMenuSeparator(menu)
			}
		}
	}

	for _, item := range menu.items {
		if menuItem, ok := item.(*MenuItem); ok && menuItem.callback != nil {
			m.attachHandler(menuItem)
		}
	}
}

func (m *linuxMenu) attachHandler(item *MenuItem) {
	(item.impl).(*linuxMenuItem).handlerId = attachMenuHandler(item)
}

func (m *linuxMenu) addSubMenuToItem(menu *Menu, item *MenuItem) {
	if menu.impl == nil {
		menu.impl = &linuxMenu{
			menu:   menu,
			native: menuNew(),
		}
	}
	menuSetSubmenu(item, menu)
}

func (m *linuxMenu) addMenuItem(parent *Menu, menu *MenuItem) {
	menuAppend(parent, menu)
}

func (m *linuxMenu) addMenuSeparator(menu *Menu) {
	menuAddSeparator(menu)

}

func (m *linuxMenu) addServicesMenu(menu *Menu) {
	// FIXME: Should this be required?
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
