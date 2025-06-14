package application

// WindowManager manages all window-related operations
type WindowManager struct {
	app *App
}

// NewWindowManager creates a new WindowManager instance
func NewWindowManager(app *App) *WindowManager {
	return &WindowManager{
		app: app,
		// Use the app's existing windows map - don't create a new one
	}
}

// GetByName returns a window by name
func (wm *WindowManager) GetByName(name string) Window {
	wm.app.windowsLock.RLock()
	defer wm.app.windowsLock.RUnlock()

	for _, window := range wm.app.windows {
		if window.Name() == name {
			return window
		}
	}
	return nil
}

// OnCreate registers a callback to be called when a window is created
func (wm *WindowManager) OnCreate(callback func(Window)) {
	wm.app.windowCreatedCallbacks = append(wm.app.windowCreatedCallbacks, callback)
}

// New creates a new webview window
func (wm *WindowManager) New() *WebviewWindow {
	return wm.app.NewWebviewWindow()
}

// NewWithOptions creates a new webview window with options
func (wm *WindowManager) NewWithOptions(windowOptions WebviewWindowOptions) *WebviewWindow {
	return wm.app.NewWebviewWindowWithOptions(windowOptions)
}

// Current returns the current active window
func (wm *WindowManager) Current() *WebviewWindow {
	return wm.app.CurrentWindow()
}

// Add adds a window to the manager (internal use)
func (wm *WindowManager) add(window Window) {
	wm.app.windowsLock.Lock()
	defer wm.app.windowsLock.Unlock()
	wm.app.windows[window.ID()] = window

	// Call registered callbacks
	for _, callback := range wm.app.windowCreatedCallbacks {
		callback(window)
	}
}

// Remove removes a window from the manager (internal use)
func (wm *WindowManager) remove(windowID uint) {
	wm.app.windowsLock.Lock()
	defer wm.app.windowsLock.Unlock()
	delete(wm.app.windows, windowID)
}

// GetAll returns all windows
func (wm *WindowManager) GetAll() []Window {
	wm.app.windowsLock.RLock()
	defer wm.app.windowsLock.RUnlock()

	windows := make([]Window, 0, len(wm.app.windows))
	for _, window := range wm.app.windows {
		windows = append(windows, window)
	}
	return windows
}
