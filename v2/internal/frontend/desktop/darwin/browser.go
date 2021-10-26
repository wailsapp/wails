//go:build darwin
// +build darwin

package darwin

import (
	"github.com/pkg/browser"
)

// BrowserOpenURL Use the default browser to open the url
func (f *Frontend) BrowserOpenURL(url string) {
	// Specific method implementation
	_ = browser.OpenURL(url)
}
