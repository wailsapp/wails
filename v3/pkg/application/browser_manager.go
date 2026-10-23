package application

// BrowserManager manages browser-related operations
type BrowserManager struct {
	app *App
}

// newBrowserManager creates a new BrowserManager instance
func newBrowserManager(app *App) *BrowserManager {
	return &BrowserManager{
		app: app,
	}
}

// OpenURL opens a URL in the default browser
func (bm *BrowserManager) OpenURL(url string) error {
	return openBrowserURL(url)
}

// OpenFile opens a file in the default browser
func (bm *BrowserManager) OpenFile(path string) error {
	return openBrowserFile(path)
}
