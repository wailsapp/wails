package application

// KeyBindingManager manages all key binding operations
type KeyBindingManager struct {
	app *App
}

// NewKeyBindingManager creates a new KeyBindingManager instance
func NewKeyBindingManager(app *App) *KeyBindingManager {
	return &KeyBindingManager{
		app: app,
	}
}

// Add adds a key binding
func (kbm *KeyBindingManager) Add(accel string, callback func(window *WebviewWindow)) {
	kbm.app.keyBindingsLock.Lock()
	defer kbm.app.keyBindingsLock.Unlock()
	kbm.app.keyBindings[accel] = callback
}

// Remove removes a key binding
func (kbm *KeyBindingManager) Remove(accel string) {
	kbm.app.keyBindingsLock.Lock()
	defer kbm.app.keyBindingsLock.Unlock()
	delete(kbm.app.keyBindings, accel)
}

// Process processes a key binding and returns true if handled
func (kbm *KeyBindingManager) Process(accel string, window *WebviewWindow) bool {
	kbm.app.keyBindingsLock.RLock()
	callback, exists := kbm.app.keyBindings[accel]
	kbm.app.keyBindingsLock.RUnlock()

	if exists && callback != nil {
		callback(window)
		return true
	}
	return false
}

// GetAll returns all registered key bindings
func (kbm *KeyBindingManager) GetAll() map[string]func(window *WebviewWindow) {
	kbm.app.keyBindingsLock.RLock()
	defer kbm.app.keyBindingsLock.RUnlock()

	result := make(map[string]func(window *WebviewWindow))
	for accel, callback := range kbm.app.keyBindings {
		result[accel] = callback
	}
	return result
}

// HandleWindowKeyEvent handles window key events (internal use)
func (kbm *KeyBindingManager) handleWindowKeyEvent(event *windowKeyEvent) {
	kbm.app.handleWindowKeyEvent(event)
}
