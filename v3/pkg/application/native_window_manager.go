package application

// NativeWindowManager creates and looks up experimental native-content
// windows. It is separate from WindowManager in v3 because the existing
// Window interface includes WebView-only operations.
type NativeWindowManager struct {
	app *App
}

func newNativeWindowManager(app *App) *NativeWindowManager {
	return &NativeWindowManager{app: app}
}

// New creates a native-content window with default options.
func (m *NativeWindowManager) New() *NativeWindow {
	return m.NewWithOptions(NativeWindowOptions{})
}

// NewWithOptions creates a native-content window. Native AppKit resources are
// allocated when the application starts, just like WebviewWindow resources.
func (m *NativeWindowManager) NewWithOptions(options NativeWindowOptions) *NativeWindow {
	window := newNativeWindow(options)
	m.app.nativeWindowsLock.Lock()
	m.app.nativeWindows[window.ID()] = window
	m.app.nativeWindowsLock.Unlock()
	m.app.runOrDeferToAppRun(window)
	return window
}

// Get returns a native window by name.
func (m *NativeWindowManager) Get(name string) (*NativeWindow, bool) {
	m.app.nativeWindowsLock.RLock()
	defer m.app.nativeWindowsLock.RUnlock()
	for _, window := range m.app.nativeWindows {
		if window.Name() == name {
			return window, true
		}
	}
	return nil, false
}

// GetByID returns a native window by Wails window ID.
func (m *NativeWindowManager) GetByID(id uint) (*NativeWindow, bool) {
	m.app.nativeWindowsLock.RLock()
	defer m.app.nativeWindowsLock.RUnlock()
	window, ok := m.app.nativeWindows[id]
	return window, ok
}

func (m *NativeWindowManager) remove(id uint) {
	m.app.nativeWindowsLock.Lock()
	delete(m.app.nativeWindows, id)
	m.app.nativeWindowsLock.Unlock()
}
