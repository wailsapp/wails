package application

// MenuManager manages menu-related operations
type MenuManager struct {
	app *App
}

// newMenuManager creates a new MenuManager instance
func newMenuManager(app *App) *MenuManager {
	return &MenuManager{
		app: app,
	}
}

// Set sets the application menu
func (mm *MenuManager) Set(menu *Menu) {
	mm.SetApplicationMenu(menu)
}

// SetApplicationMenu sets the application menu
func (mm *MenuManager) SetApplicationMenu(menu *Menu) {
	mm.app.applicationMenu = menu
	if mm.app.impl != nil {
		mm.app.impl.setApplicationMenu(menu)
	}
}

// GetApplicationMenu returns the current application menu
func (mm *MenuManager) GetApplicationMenu() *Menu {
	return mm.app.applicationMenu
}

// New creates a new menu
func (mm *MenuManager) New() *Menu {
	return &Menu{}
}

// ShowAbout shows the about dialog
func (mm *MenuManager) ShowAbout() {
	if mm.app.impl != nil {
		mm.app.impl.showAboutDialog(mm.app.options.Name, mm.app.options.Description, mm.app.options.Icon)
	}
}

// handleMenuItemClicked handles menu item click events (internal use)
func (mm *MenuManager) handleMenuItemClicked(menuItemID uint) {
	defer handlePanic()

	menuItem := getMenuItemByID(menuItemID)
	if menuItem == nil {
		mm.app.warning("MenuItem #%d not found", menuItemID)
		return
	}
	menuItem.handleClick()
}
