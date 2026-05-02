//go:build android

package application

// Android menu stubs - Android doesn't have traditional application menus

func (m *Menu) handleStyleChange() {}

type androidMenu struct {
	menu *Menu
}

func newMenuImpl(menu *Menu) *androidMenu {
	return &androidMenu{
		menu: menu,
	}
}

func (m *androidMenu) update() {
	// Android doesn't have traditional menus
}

func defaultApplicationMenu() *Menu {
	// No application menu on Android
	return nil
}
