package runtime

import "github.com/pkg/browser"

// Browser exposes browser methods to the runtime
type Browser struct{}

// NewBrowser creates a new runtime Browser struct
func NewBrowser() *Browser {
	return &Browser{}
}

// OpenURL opens the given url in the system's default browser
func (r *Browser) OpenURL(url string) error {
	return browser.OpenURL(url)
}

// OpenFile opens the given file in the system's default browser
func (r *Browser) OpenFile(filePath string) error {
	return browser.OpenFile(filePath)
}
