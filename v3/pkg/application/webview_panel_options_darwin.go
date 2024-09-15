package application

type WebviewPanelOptions struct {
	WebviewWindowOptions

	// Floating will make the panel float above other application in every workspace.
	Floating bool

	// ShouldClose is called when the panel is about to close.
	// Return true to allow the panel to close, or false to prevent it from closing.
	ShouldClose func(panel *WebviewPanel) bool

	// KeyBindings is a map of key bindings to functions
	KeyBindings map[string]func(panel *WebviewPanel)
}

var WebviewPanelDefaults = &WebviewPanelOptions{
	WebviewWindowOptions: *WebviewWindowDefaults,
}

func processKeyBindingOptionsForPanel(keyBindings map[string]func(panel *WebviewPanel)) map[string]func(panel *WebviewPanel) {
	result := make(map[string]func(panel *WebviewPanel))
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
