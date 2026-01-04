//go:build linux && !android && !gtk3

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

	for _, item := range menu.items {
		switch item.itemType {
		case submenu:
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			m.processMenu(item.submenu)
			m.addSubMenuToItem(item.submenu, item)
			m.addMenuItem(menu, item)
		case text:
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			m.addMenuItem(menu, item)
		case checkbox:
			menuItem := newCheckMenuItemImpl(item)
			item.impl = menuItem
			m.addMenuItem(menu, item)
		case radio:
			menuItem := newRadioMenuItemImpl(item)
			item.impl = menuItem
			m.addMenuItem(menu, item)
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
