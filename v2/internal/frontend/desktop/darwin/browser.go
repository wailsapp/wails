//go:build darwin
// +build darwin

package darwin

import (
	"github.com/pkg/browser"
)

// BrowserOpenURL Use the default browser to open the url
func (f *Frontend) BrowserOpenURL(url string) {
	// Specific method implementation
	if err := browser.OpenURL(url); err != nil {
		f.logger.Error("Unable to open default system browser")
	}
}
