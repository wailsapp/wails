//go:build android

package application

import (
	"os"
)

// setupSignalHandler sets up signal handling for Android
// On Android, we don't handle Unix signals directly as the app lifecycle
// is managed by the Android runtime
func setupSignalHandler() {
	// No-op on Android - lifecycle managed by Android framework
}

// handleSignal processes a signal
func handleSignal(_ os.Signal) {
	// No-op on Android
}
