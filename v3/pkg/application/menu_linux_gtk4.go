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

var radioGroupCounter uint = 0

func (m *linuxMenu) processMenu(menu *Menu) {
	if menu.impl == nil {
		menu.impl = &linuxMenu{
			menu:   menu,
			native: menuNew(),
		}
	}

	var currentRadioGroup uint = 0
	var checkedRadioId uint = 0

	hasSeparators := false
	for _, item := range menu.items {
		if item.itemType == separator {
			hasSeparators = true
			break
		}
	}

	// GMenu uses sections for visual separators
	// Only use sections if the menu has separators
	var currentSection pointer
	var hasSectionItems bool
	if hasSeparators {
		currentSection = menuNewSection()
		hasSectionItems = false
	}

	for _, item := range menu.items {
		if item.itemType != radio {
			currentRadioGroup = 0
			checkedRadioId = 0
		}

		switch item.itemType {
		case submenu:
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			m.processMenu(item.submenu)
			m.addSubMenuToItem(item.submenu, item)
			if hasSeparators {
				m.addMenuItemToSection(currentSection, item)
				hasSectionItems = true
			} else {
				m.addMenuItem(menu, item)
			}
		case text:
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			if hasSeparators {
				m.addMenuItemToSection(currentSection, item)
				hasSectionItems = true
			} else {
				m.addMenuItem(menu, item)
			}
		case checkbox:
			menuItem := newCheckMenuItemImpl(item)
			item.impl = menuItem
			if hasSeparators {
				m.addMenuItemToSection(currentSection, item)
				hasSectionItems = true
			} else {
				m.addMenuItem(menu, item)
			}
		case radio:
			if currentRadioGroup == 0 {
				radioGroupCounter++
				currentRadioGroup = radioGroupCounter
			}
			if item.checked {
				checkedRadioId = item.id
			}
			menuItem := newRadioMenuItemImpl(item, currentRadioGroup, checkedRadioId)
			item.impl = menuItem
			if hasSeparators {
				m.addMenuItemToSection(currentSection, item)
				hasSectionItems = true
			} else {
				m.addMenuItem(menu, item)
			}
		case separator:
			if hasSectionItems {
				menuAppendSection(menu, currentSection)
				currentSection = menuNewSection()
				hasSectionItems = false
			}
		}
	}

	if hasSeparators && hasSectionItems {
		menuAppendSection(menu, currentSection)
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

func (m *linuxMenu) addMenuItem(parent *Menu, item *MenuItem) {
	menuAppend(parent, item, item.hidden)
}

func (m *linuxMenu) addMenuItemToSection(section pointer, item *MenuItem) {
	menuAppendItemToSection(section, item)
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
