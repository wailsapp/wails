//go:build !((android || ios) && !server)

package application

import "github.com/wailsapp/wails/v3/internal/browser"

func openBrowserURL(url string) error {
	return browser.OpenURL(url)
}

func openBrowserFile(path string) error {
	return browser.OpenFile(path)
}
