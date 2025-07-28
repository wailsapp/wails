//go:build darwin

package application

// showSnapAssist is a no-op on macOS as SnapAssist is a Windows-only feature
func showSnapAssist(window *WebviewWindow) {
	// SnapAssist is not available on macOS
}