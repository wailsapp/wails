package wails

import "github.com/pkg/browser"

// GlobalRuntimeBrowser is the global instance of the RuntimeBrowser object
// Why? Because we need to use it in both the runtime and from the frontend
var GlobalRuntimeBrowser = newRuntimeBrowser()

// RuntimeBrowser exposes browser methods to the runtime
type RuntimeBrowser struct {
}

func newRuntimeBrowser() *RuntimeBrowser {
	return &RuntimeBrowser{}
}

// OpenURL opens the given url in the system's default browser
func (r *RuntimeBrowser) OpenURL(url string) error {
	return browser.OpenURL(url)
}

// OpenFile opens the given file in the system's default browser
func (r *RuntimeBrowser) OpenFile(filePath string) error {
	return browser.OpenFile(filePath)
}
