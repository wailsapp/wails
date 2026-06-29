//go:build ios

package application

// platformSignalHandler is empty on iOS as signal handling is not supported
type platformSignalHandler struct {
	// No signal handler on iOS
}