//go:build ios

package application

// iOS menu stubs - iOS doesn't have traditional menus

// iOS doesn't have traditional menus like desktop platforms
// These are placeholder implementations

func (m *Menu) handleStyleChange() {}

type iosMenu struct {
	menu *Menu
}

func newMenuImpl(menu *Menu) *iosMenu {
	return &iosMenu{
		menu: menu,
	}
}

func (m *iosMenu) update() {
	// iOS doesn't have traditional menus
}