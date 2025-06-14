package application

// MenuManager manages menu-related operations
type MenuManager struct {
	app *App
}

// NewMenuManager creates a new MenuManager instance
func NewMenuManager(app *App) *MenuManager {
	return &MenuManager{
		app: app,
	}
}

// Set sets the application menu
func (mm *MenuManager) Set(menu *Menu) {
	mm.app.ApplicationMenu = menu
	if mm.app.impl != nil {
		mm.app.impl.setApplicationMenu(menu)
	}
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
