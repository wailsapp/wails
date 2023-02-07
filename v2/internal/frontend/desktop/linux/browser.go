//go:build linux
// +build linux

package linux

import "github.com/pkg/browser"

// BrowserOpenURL Use the default browser to open the url
func (f *Frontend) BrowserOpenURL(url string) {
	// Specific method implementation
	_ = browser.OpenURL(url)
}

func (f *Frontend) OpenDevToolsWindow() {
	// TODO implement linux dev tools
}
