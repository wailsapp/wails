//go:build !mcp || ios || android

package application

// startMCPServer is a no-op when the mcp build tag is absent.
func startMCPServer(_ *App) error { return nil }
