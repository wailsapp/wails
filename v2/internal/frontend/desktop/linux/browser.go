//go:build linux
// +build linux

package linux

import (
	"fmt"
	"net/url"

	"github.com/pkg/browser"
)

// BrowserOpenURL Use the default browser to open the url
func (f *Frontend) BrowserOpenURL(rawURL string) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		f.logger.Error(fmt.Sprintf("BrowserOpenURL cannot parse url: %s", err.Error()))
		return
	}
	// Specific method implementation
	if err := browser.OpenURL(parsed.String()); err != nil {
		f.logger.Error("Unable to open default system browser")
	}
}
