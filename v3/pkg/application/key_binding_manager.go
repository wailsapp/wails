package application

// KeyBindingManager manages all key binding operations
type KeyBindingManager struct {
	app          *App
	menuBindings map[string]*MenuItem
}

// newKeyBindingManager creates a new KeyBindingManager instance
func newKeyBindingManager(app *App) *KeyBindingManager {
	return &KeyBindingManager{
		app:          app,
		menuBindings: make(map[string]*MenuItem),
	}
}

// Add adds a key binding
func (kbm *KeyBindingManager) Add(accelerator string, callback func(window Window)) {
	kbm.app.keyBindingsLock.Lock()
	defer kbm.app.keyBindingsLock.Unlock()
	kbm.app.keyBindings[accelerator] = callback
}

// addMenuBinding keeps native menu shortcuts separate from user key bindings.
func (kbm *KeyBindingManager) addMenuBinding(item *MenuItem) {
	kbm.app.keyBindingsLock.Lock()
	defer kbm.app.keyBindingsLock.Unlock()
	kbm.menuBindings[item.accelerator.String()] = item
}

func (kbm *KeyBindingManager) removeMenuBinding(item *MenuItem) {
	kbm.app.keyBindingsLock.Lock()
	defer kbm.app.keyBindingsLock.Unlock()
	if key := item.accelerator.String(); kbm.menuBindings[key] == item {
		delete(kbm.menuBindings, key)
	}
}

// Remove removes a key binding
func (kbm *KeyBindingManager) Remove(accelerator string) {
	kbm.app.keyBindingsLock.Lock()
	defer kbm.app.keyBindingsLock.Unlock()
	delete(kbm.app.keyBindings, accelerator)
	delete(kbm.menuBindings, accelerator)
}

// Process processes a key binding and returns true if handled
func (kbm *KeyBindingManager) Process(accelerator string, window Window) bool {
	kbm.app.keyBindingsLock.RLock()
	callback, exists := kbm.app.keyBindings[accelerator]
	item := kbm.menuBindings[accelerator]
	kbm.app.keyBindingsLock.RUnlock()

	if exists && callback != nil {
		callback(window)
		return true
	}
	if item != nil && !suppressesApplicationMenuItem(window, kbm.app.applicationMenu, item) {
		item.handleClick()
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

	result := make([]*KeyBinding, 0, len(kbm.app.keyBindings)+len(kbm.menuBindings))
	for accelerator, callback := range kbm.app.keyBindings {
		result = append(result, &KeyBinding{
			Accelerator: accelerator,
			Callback:    callback,
		})
	}
	for accelerator, item := range kbm.menuBindings {
		if _, exists := kbm.app.keyBindings[accelerator]; exists {
			continue
		}
		result = append(result, &KeyBinding{
			Accelerator: accelerator,
			Callback: func(window Window) {
				if !suppressesApplicationMenuItem(window, kbm.app.applicationMenu, item) {
					item.handleClick()
				}
			},
		})
	}
	return result
}

// HandleWindowKeyEvent handles window key events (internal use)
func (kbm *KeyBindingManager) handleWindowKeyEvent(event *windowKeyEvent) {
	kbm.app.handleWindowKeyEvent(event)
}
