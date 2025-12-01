// Package browser provides cross-platform URL and file opening functionality.
// It provides a unified API for opening URLs and files across desktop (Windows, macOS, Linux),
// iOS, and Android platforms.
package browser

// OpenURL opens a URL in the default browser.
// This function is platform-specific and implemented separately for each platform.
func OpenURL(url string) error {
	return openURL(url)
}

// OpenFile opens a file in the default application.
// This function is platform-specific and implemented separately for each platform.
func OpenFile(path string) error {
	return openFile(path)
}
