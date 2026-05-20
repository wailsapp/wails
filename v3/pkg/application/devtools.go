package application

// OpenDevTools opens the developer tools window for the current window.
// This function only works when the application is built with the -runtimedevtools flag.
// This function is a no-op (does nothing) when:
// - The global app is nil (called before application initialization)
// - There is no currentWindow (currentWindow == nil)
// - The runtime lacks devtools support
// On macOS, requires macOS 12+ and may be restricted by Apple's private API policies.
func OpenDevTools() {
	app := Get()
	if app == nil {
		return
	}

	currentWindow := app.Window.Current()
	if currentWindow != nil {
		currentWindow.OpenDevTools()
	}
}
