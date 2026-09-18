//go:build server

package application

// setupCommonEvents sets up common application events for server mode.
// In server mode, there are no platform-specific events to map,
// so this is a no-op.
func (h *serverApp) setupCommonEvents() {
	// No-op: server mode has no platform-specific events to map
}
