//go:build darwin && !ios && !server

package application

import "unsafe"

type macosMenu struct {
	menu *Menu

	nsMenu unsafe.Pointer
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

func newMenuImpl(menu *Menu) *macosMenu {
	result := &macosMenu{
		menu: menu,
	}
	return result
}

func (m *macosMenu) update() {
	InvokeSync(func() {
		if m.nsMenu == nil {
			m.nsMenu = newNSMenu(m.menu.label)
		} else {
			clearNSMenu(m.nsMenu)
		}
		m.processMenu(m.nsMenu, m.menu)
	})
}
