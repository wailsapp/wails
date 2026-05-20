//go:build !darwin || ios || server

package application

// captureLaunchURL is a no-op on non-macOS platforms.  On Windows and Linux
// the URL is already present in os.Args, so no extra capture is needed.
func captureLaunchURL() string {
	return ""
}
