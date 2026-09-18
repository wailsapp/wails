//go:build wails_native

package application

type applicationUpdater struct{}

// The built-in updater owns a WebView-based user interface and is therefore
// not installed in a wails_native build. Applications that need a headless
// updater can import pkg/updater directly; doing so intentionally brings its
// HTTP and cryptographic dependencies into their binary.
func prepareUpdaterApplicationProcess() {}

func newApplicationUpdater(*App) *applicationUpdater { return nil }
