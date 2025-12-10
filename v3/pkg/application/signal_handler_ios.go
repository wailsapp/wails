//go:build ios

package application

// setupSignalHandler is a no-op on iOS as signal handling is not supported
// iOS apps run in a sandboxed environment where signal handling is restricted
// and can cause crashes if attempted
func (app *App) setupSignalHandler(options Options) {
	// No signal handling on iOS - the OS manages app lifecycle
	// Signal handlers would cause crashes due to iOS sandbox restrictions
}