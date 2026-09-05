//go:build !linux

package runtime

// ResetSignalHandlers resets signal handlers to allow panic recovery.
//
// On Linux, WebKit (used for the webview) may install signal handlers without
// the SA_ONSTACK flag, which prevents Go from properly recovering from panics
// caused by nil pointer dereferences or other memory access violations.
//
// Call this function immediately before code that might panic to ensure
// the signal handlers are properly configured for Go's panic recovery mechanism.
//
// Note: This function only has an effect on Linux. On other platforms,
// it is a no-op.
func ResetSignalHandlers() {
	// No-op on non-Linux platforms
}
