package application

// ClipboardManager manages clipboard operations
type ClipboardManager struct {
	app       *App
	clipboard *Clipboard
}

// newClipboardManager creates a new ClipboardManager instance
func newClipboardManager(app *App) *ClipboardManager {
	return &ClipboardManager{
		app: app,
	}
}

// SetText sets text in the clipboard
func (cm *ClipboardManager) SetText(text string) bool {
	return cm.getClipboard().SetText(text)
}

// Text gets text from the clipboard
func (cm *ClipboardManager) Text() (string, bool) {
	return cm.getClipboard().Text()
}

// getClipboard returns the clipboard instance, creating it if needed (lazy initialization)
func (cm *ClipboardManager) getClipboard() *Clipboard {
	if cm.clipboard == nil {
		cm.clipboard = newClipboard()
	}
	return cm.clipboard
}
