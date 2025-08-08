//go:build linux
// +build linux

package linux

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v2/internal/frontend/utils"
)

// BrowserOpenURL Use the default browser to open the url
func (f *Frontend) BrowserOpenURL(url string) {
	if err := utils.ValidateURL(url); err != nil {
		f.logger.Error(fmt.Sprintf("Invalid URL %s", err.Error()))
	}
	// Specific method implementation
	if err := browser.OpenURL(url); err != nil {
		f.logger.Error("Unable to open default system browser")
	}
}
