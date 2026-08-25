// Package browser provides functions to open URLs and files in the default browser.
package browser

// OpenURL opens the named URL in the default browser.
func OpenURL(url string) error {
	return open(url)
}

// OpenFile opens the named file in the default browser or file handler.
func OpenFile(path string) error {
	return open(path)
}
