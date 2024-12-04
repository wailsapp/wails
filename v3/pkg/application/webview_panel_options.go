package application

type WebviewPanelOptions struct {
	WebviewWindowOptions

	// Floating will make the panel float above other application in every workspace.
	Floating bool

	// ShouldClose is called when the panel is about to close.
	// Return true to allow the panel to close, or false to prevent it from closing.
	ShouldClose func(panel *WebviewPanel) bool

	// KeyBindings is a map of key bindings to functions. Other key bindings provided from
	// the embedded field [WebviewWindowOptions] are still valid and the two maps are merged
	// together, with the [WebviewPanelOptions] values overriding the others if they have the
	// same key.
	KeyBindings map[string]func(panel *WebviewPanel)
}

func processKeyBindingOptionsForPanel(keyBindings map[string]func(panel *WebviewPanel), windowKeyBindings map[string]func(panel *WebviewWindow)) map[string]func(panel *WebviewPanel) {
	result := make(map[string]func(panel *WebviewPanel))

	for key, callback := range windowKeyBindings {
		acc, err := parseAccelerator(key)
		if err != nil {
			globalApplication.error("Invalid keybinding: %s", err.Error())
			continue
		}
		result[acc.String()] = func(panel *WebviewPanel) {
			callback(panel.WebviewWindow)
		}
		globalApplication.debug("Added Keybinding", "accelerator", acc.String())
	}

	for key, callback := range keyBindings {
		// Parse the key to an accelerator
		acc, err := parseAccelerator(key)
		if err != nil {
			globalApplication.error("Invalid keybinding: %s", err.Error())
			continue
		}
		result[acc.String()] = callback
		globalApplication.debug("Added Keybinding", "accelerator", acc.String())
	}
	return result
}
