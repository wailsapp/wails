package application

// ContextMenuManager manages all context menu operations
type ContextMenuManager struct {
	app *App
}

// NewContextMenuManager creates a new ContextMenuManager instance
func NewContextMenuManager(app *App) *ContextMenuManager {
	return &ContextMenuManager{
		app: app,
	}
}

// Register registers a context menu
func (cmm *ContextMenuManager) Register(menu *ContextMenu) {
	cmm.app.contextMenusLock.Lock()
	defer cmm.app.contextMenusLock.Unlock()
	cmm.app.contextMenus[menu.name] = menu
}

// Unregister removes a context menu by name
func (cmm *ContextMenuManager) Unregister(name string) {
	cmm.app.contextMenusLock.Lock()
	defer cmm.app.contextMenusLock.Unlock()
	delete(cmm.app.contextMenus, name)
}

// Get retrieves a context menu by name
func (cmm *ContextMenuManager) Get(name string) (*ContextMenu, bool) {
	cmm.app.contextMenusLock.Lock()
	defer cmm.app.contextMenusLock.Unlock()
	menu, exists := cmm.app.contextMenus[name]
	return menu, exists
}

// GetAll returns all registered context menus
func (cmm *ContextMenuManager) GetAll() map[string]*ContextMenu {
	cmm.app.contextMenusLock.Lock()
	defer cmm.app.contextMenusLock.Unlock()

	result := make(map[string]*ContextMenu)
	for name, menu := range cmm.app.contextMenus {
		result[name] = menu
	}
	return result
}
