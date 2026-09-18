//go:build !server

package application

// GetBrowserWindow is a stub for non-server builds.
// Returns nil as browser windows are only available in server mode.
func GetBrowserWindow(clientId string) Window {
	return nil
}
