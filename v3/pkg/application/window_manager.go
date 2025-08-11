package application

// WindowManager manages all window-related operations
type WindowManager struct {
	app *App
}

// newWindowManager creates a new WindowManager instance
func newWindowManager(app *App) *WindowManager {
	return &WindowManager{
		app: app,
		// Use the app's existing windows map - don't create a new one
	}
}

// GetByName returns a window by name and whether it exists
func (wm *WindowManager) GetByName(name string) (Window, bool) {
	wm.app.windowsLock.RLock()
	defer wm.app.windowsLock.RUnlock()

	for _, window := range wm.app.windows {
		if window.Name() == name {
			return window, true
		}
	}
	return nil, false
}

// Get is an alias for GetByName for consistency
func (wm *WindowManager) Get(name string) (Window, bool) {
	return wm.GetByName(name)
}

// GetByID returns a window by ID and whether it exists
func (wm *WindowManager) GetByID(id uint) (Window, bool) {
	wm.app.windowsLock.RLock()
	defer wm.app.windowsLock.RUnlock()

	window, exists := wm.app.windows[id]
	return window, exists
}

// OnCreate registers a callback to be called when a window is created
func (wm *WindowManager) OnCreate(callback func(Window)) {
	wm.app.windowCreatedCallbacks = append(wm.app.windowCreatedCallbacks, callback)
}

// New creates a new webview window
func (wm *WindowManager) New() *WebviewWindow {
	return wm.NewWithOptions(WebviewWindowOptions{})
}

// NewWithOptions creates a new webview window with options
func (wm *WindowManager) NewWithOptions(windowOptions WebviewWindowOptions) *WebviewWindow {
	newWindow := NewWindow(windowOptions)
	id := newWindow.ID()

	wm.app.windowsLock.Lock()
	wm.app.windows[id] = newWindow
	wm.app.windowsLock.Unlock()

	// Call hooks
	for _, hook := range wm.app.windowCreatedCallbacks {
		hook(newWindow)
	}

	wm.app.runOrDeferToAppRun(newWindow)

	return newWindow
}

// Current returns the current active window (may be nil)
func (wm *WindowManager) Current() Window {
	if wm.app.impl == nil {
		return nil
	}
	id := wm.app.impl.getCurrentWindowID()
	wm.app.windowsLock.RLock()
	defer wm.app.windowsLock.RUnlock()
	result := wm.app.windows[id]
	return result
}

// Add adds a window to the manager
func (wm *WindowManager) Add(window Window) {
	wm.app.windowsLock.Lock()
	defer wm.app.windowsLock.Unlock()
	wm.app.windows[window.ID()] = window

	// Call registered callbacks
	for _, callback := range wm.app.windowCreatedCallbacks {
		callback(window)
	}
}

// Remove removes a window from the manager by ID
func (wm *WindowManager) Remove(windowID uint) {
	wm.app.windowsLock.Lock()
	defer wm.app.windowsLock.Unlock()
	delete(wm.app.windows, windowID)
}

// RemoveByName removes a window from the manager by name
func (wm *WindowManager) RemoveByName(name string) bool {
	window, exists := wm.GetByName(name)
	if exists {
		wm.Remove(window.ID())
		return true
	}
	return false
}

// Internal methods for backward compatibility
func (wm *WindowManager) add(window Window) {
	wm.Add(window)
}

func (wm *WindowManager) remove(windowID uint) {
	wm.Remove(windowID)
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
