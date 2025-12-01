//go:build !android && !ios

package browser

import (
	"github.com/pkg/browser"
)

// openURL opens a URL in the default browser on desktop platforms.
// This uses the github.com/pkg/browser package which handles Windows, macOS, and Linux.
func openURL(url string) error {
	return browser.OpenURL(url)
}

// openFile opens a file in the default application on desktop platforms.
// This uses the github.com/pkg/browser package which handles Windows, macOS, and Linux.
func openFile(path string) error {
	return browser.OpenFile(path)
}
