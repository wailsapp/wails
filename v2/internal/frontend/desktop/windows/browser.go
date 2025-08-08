//go:build windows
// +build windows

package windows

import (
	"fmt"
	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v2/internal/frontend/utils"
	"golang.org/x/sys/windows"
)

var fallbackBrowserPaths = []string{
	`\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
	`\Program Files\Google\Chrome\Application\chrome.exe`,
	`\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
	`\Program Files\Mozilla Firefox\firefox.exe`,
}

// BrowserOpenURL Use the default browser to open the url
func (f *Frontend) BrowserOpenURL(rawURL string) {
	url, err := utils.ValidateAndSanitizeURL(rawURL)
	if err != nil {
		f.logger.Error(fmt.Sprintf("Invalid URL %s", err.Error()))
		return
	}

	// Specific method implementation
	err = browser.OpenURL(url)
	if err == nil {
		return
	}
	for _, fallback := range fallbackBrowserPaths {
		if err := openBrowser(fallback, url); err == nil {
			return
		}
	}
	f.logger.Error("Unable to open default system browser")
}

func openBrowser(path, url string) error {
	return windows.ShellExecute(0, nil, windows.StringToUTF16Ptr(path), windows.StringToUTF16Ptr(url), nil, windows.SW_SHOWNORMAL)
}
