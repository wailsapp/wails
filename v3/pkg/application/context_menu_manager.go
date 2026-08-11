package application

// ContextMenuManager manages all context menu operations
type ContextMenuManager struct {
	app *App
}

// newContextMenuManager creates a new ContextMenuManager instance
func newContextMenuManager(app *App) *ContextMenuManager {
	return &ContextMenuManager{
		app: app,
	}
}

// New creates a new context menu
func (cmm *ContextMenuManager) New() *ContextMenu {
	return &ContextMenu{
		Menu: NewMenu(),
	}
}

// Add adds a context menu (replaces Register for consistency)
func (cmm *ContextMenuManager) Add(name string, menu *ContextMenu) {
	cmm.app.contextMenusLock.Lock()
	defer cmm.app.contextMenusLock.Unlock()
	cmm.app.contextMenus[name] = menu
}

// Remove removes a context menu by name (replaces Unregister for consistency)
func (cmm *ContextMenuManager) Remove(name string) {
	cmm.app.contextMenusLock.Lock()
	defer cmm.app.contextMenusLock.Unlock()
	delete(cmm.app.contextMenus, name)
}

// Get retrieves a context menu by name
func (cmm *ContextMenuManager) Get(name string) (*ContextMenu, bool) {
    cmm.app.contextMenusLock.RLock()
    defer cmm.app.contextMenusLock.RUnlock()
    menu, exists := cmm.app.contextMenus[name]
    return menu, exists
}

// GetAll returns all registered context menus as a slice
func (cmm *ContextMenuManager) GetAll() []*ContextMenu {
	cmm.app.contextMenusLock.RLock()
	defer cmm.app.contextMenusLock.RUnlock()

	result := make([]*ContextMenu, 0, len(cmm.app.contextMenus))
	for _, menu := range cmm.app.contextMenus {
		result = append(result, menu)
	}
	return result
}
