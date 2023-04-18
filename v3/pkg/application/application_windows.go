//go:build windows

package application

type windowsApp struct {
	//applicationMenu unsafe.Pointer
	parent *App
}

func (m *windowsApp) dispatchOnMainThread(id uint) {
	//TODO implement me
	panic("implement me")
}

func (m *windowsApp) getPrimaryScreen() (*Screen, error) {
	//TODO implement me
	panic("implement me")
}

func (m *windowsApp) getScreens() ([]*Screen, error) {
	//TODO implement me
	panic("implement me")
}

func (m *windowsApp) hide() {
	//C.hide()
}

func (m *windowsApp) show() {
	//C.show()
}

func (m *windowsApp) on(eventID uint) {
	//C.registerListener(C.uint(eventID))
}

func (m *windowsApp) setIcon(icon []byte) {
	//C.setApplicationIcon(unsafe.Pointer(&icon[0]), C.int(len(icon)))
}

func (m *windowsApp) name() string {
	//appName := C.getAppName()
	//defer C.free(unsafe.Pointer(appName))
	//return C.GoString(appName)
	return ""
}

func (m *windowsApp) getCurrentWindowID() uint {
	//return uint(C.getCurrentWindowID())
	return uint(0)
}

func (m *windowsApp) setApplicationMenu(menu *Menu) {
	if menu == nil {
		// Create a default menu for mac
		menu = defaultApplicationMenu()
	}
	menu.Update()

	// Convert impl to macosMenu object
	//m.applicationMenu = (menu.impl).(*macosMenu).nsMenu
	//C.setApplicationMenu(m.applicationMenu)
}

func (m *windowsApp) run() error {
	// Add a hook to the ApplicationDidFinishLaunching event
	//m.parent.On(events.Mac.ApplicationDidFinishLaunching, func() {
	//	C.setApplicationShouldTerminateAfterLastWindowClosed(C.bool(m.parent.options.Mac.ApplicationShouldTerminateAfterLastWindowClosed))
	//	C.setActivationPolicy(C.int(m.parent.options.Mac.ActivationPolicy))
	//	C.activateIgnoringOtherApps()
	//})
	// setup event listeners
	for eventID := range m.parent.applicationEventListeners {
		m.on(eventID)
	}
	//C.run()
	return nil
}

func (m *windowsApp) destroy() {
	//C.destroyApp()
}

func newPlatformApp(app *App) *windowsApp {
	//C.init()
	return &windowsApp{
		parent: app,
	}
}
