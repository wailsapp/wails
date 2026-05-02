package application

// SystemTrayManager manages system tray-related operations
type SystemTrayManager struct {
	app *App
}

// newSystemTrayManager creates a new SystemTrayManager instance
func newSystemTrayManager(app *App) *SystemTrayManager {
	return &SystemTrayManager{
		app: app,
	}
}

// New creates a new system tray
func (stm *SystemTrayManager) New() *SystemTray {
	id := stm.getNextID()
	newSystemTray := newSystemTray(id)

	stm.app.systemTraysLock.Lock()
	stm.app.systemTrays[id] = newSystemTray
	stm.app.systemTraysLock.Unlock()

	stm.app.runOrDeferToAppRun(newSystemTray)

	return newSystemTray
}

// getNextID generates the next system tray ID (internal use)
func (stm *SystemTrayManager) getNextID() uint {
	stm.app.systemTrayIDLock.Lock()
	defer stm.app.systemTrayIDLock.Unlock()
	stm.app.systemTrayID++
	return stm.app.systemTrayID
}

// destroy destroys a system tray (internal use)
func (stm *SystemTrayManager) destroy(tray *SystemTray) {
	// Remove the system tray from the app.systemTrays map
	stm.app.systemTraysLock.Lock()
	delete(stm.app.systemTrays, tray.id)
	stm.app.systemTraysLock.Unlock()
	tray.destroy()
}
