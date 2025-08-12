package application

// KeyBindingManager manages all key binding operations
type KeyBindingManager struct {
	app *App
}

// newKeyBindingManager creates a new KeyBindingManager instance
func newKeyBindingManager(app *App) *KeyBindingManager {
	return &KeyBindingManager{
		app: app,
	}
}

// Add adds a key binding
func (kbm *KeyBindingManager) Add(accelerator string, callback func(window Window)) {
	kbm.app.keyBindingsLock.Lock()
	defer kbm.app.keyBindingsLock.Unlock()
	kbm.app.keyBindings[accelerator] = callback
}

// Remove removes a key binding
func (kbm *KeyBindingManager) Remove(accelerator string) {
	kbm.app.keyBindingsLock.Lock()
	defer kbm.app.keyBindingsLock.Unlock()
	delete(kbm.app.keyBindings, accelerator)
}

// Process processes a key binding and returns true if handled
func (kbm *KeyBindingManager) Process(accelerator string, window Window) bool {
	kbm.app.keyBindingsLock.RLock()
	callback, exists := kbm.app.keyBindings[accelerator]
	kbm.app.keyBindingsLock.RUnlock()

	if exists && callback != nil {
		callback(window)
		return true
	}
	return false
}

// KeyBinding represents a key binding with its accelerator and callback
type KeyBinding struct {
	Accelerator string
	Callback    func(window Window)
}

// GetAll returns all registered key bindings as a slice
func (kbm *KeyBindingManager) GetAll() []*KeyBinding {
	kbm.app.keyBindingsLock.RLock()
	defer kbm.app.keyBindingsLock.RUnlock()

	result := make([]*KeyBinding, 0, len(kbm.app.keyBindings))
	for accelerator, callback := range kbm.app.keyBindings {
		result = append(result, &KeyBinding{
			Accelerator: accelerator,
			Callback:    callback,
		})
	}
	return result
}

// HandleWindowKeyEvent handles window key events (internal use)
func (kbm *KeyBindingManager) handleWindowKeyEvent(event *windowKeyEvent) {
	kbm.app.handleWindowKeyEvent(event)
}
