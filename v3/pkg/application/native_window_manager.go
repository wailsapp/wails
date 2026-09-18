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

// New creates a native-content window with default options. The window has
// no content yet, so its native creation is deferred: SetSplitView creates it
// as soon as it supplies a layout (immediately in a running application,
// otherwise when App.Run starts). Until then the window is registered with
// the manager but not visible.
func (m *NativeWindowManager) New() *NativeWindow {
	return m.NewWithOptions(NativeWindowOptions{})
}

// NewWithOptions creates a native-content window. Options-time chrome
// (Toolbar, then SplitView) is claimed before the window is scheduled, so a
// window created after App.Run with a SplitView option is created, and shown
// unless Hidden, before this method returns. Option errors are reported
// through the window's error handler and leave the window deferred, exactly
// as if the option had been omitted.
func (m *NativeWindowManager) NewWithOptions(options NativeWindowOptions) *NativeWindow {
	window := newNativeWindow(options)
	m.app.nativeWindowsLock.Lock()
	m.app.nativeWindows[window.ID()] = window
	m.app.nativeWindowsLock.Unlock()
	if options.Toolbar != nil {
		if err := window.SetToolbar(options.Toolbar); err != nil {
			window.Error("Toolbar option: %s", err)
		}
	}
	if options.SplitView != nil {
		if err := window.SetSplitView(options.SplitView); err != nil {
			window.Error("SplitView option: %s", err)
		}
	}
	m.app.runOrDeferToAppRun(nativeWindowRunnable{window: window})
	return window
}

// nativeWindowRunnable adapts NativeWindow.Run, which reports its error, to
// the application's start-up queue. A window without content is deferred by
// design rather than failed, and every other error has already been reported
// through the window's error handler.
type nativeWindowRunnable struct {
	window *NativeWindow
}

func (r nativeWindowRunnable) Run() { _ = r.window.Run() }

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
