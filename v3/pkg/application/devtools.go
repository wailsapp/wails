package application

// OpenDevTools opens the developer tools window for the current window.
// This function only works when the application is built with the -runtimedevtools flag.
// On platforms where devtools cannot be opened programmatically (e.g., macOS in production),
// this function will do nothing.
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
