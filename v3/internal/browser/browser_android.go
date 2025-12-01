//go:build android

package browser

import (
	"fmt"
)

// OpenURLFunc is the function that will be set by the application package
// to provide the actual Android implementation via JNI.
var OpenURLFunc func(url string) error

// openURL opens a URL using Android's Intent system.
func openURL(url string) error {
	if OpenURLFunc == nil {
		return fmt.Errorf("Android OpenURL not initialized - application not started")
	}
	return OpenURLFunc(url)
}

// openFile opens a file using Android's Intent system.
// On Android, this typically opens the file with the appropriate app based on MIME type.
func openFile(path string) error {
	// On Android, we can use a file:// URI to open files
	fileURL := "file://" + path
	return openURL(fileURL)
}
