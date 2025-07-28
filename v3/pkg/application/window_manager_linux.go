//go:build linux

package application

// showSnapAssist is a no-op on Linux as SnapAssist is a Windows-only feature
func showSnapAssist(window *WebviewWindow) {
	// SnapAssist is not available on Linux
}