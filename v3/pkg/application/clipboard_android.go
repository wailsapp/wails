//go:build android

package application

type androidClipboardImpl struct{}

func newClipboardImpl() clipboardImpl {
	return &androidClipboardImpl{}
}

func (c *androidClipboardImpl) setText(text string) bool {
	AndroidSetClipboardText(text)
	return true
}

func (c *androidClipboardImpl) text() (string, bool) {
	text := AndroidGetClipboardText()
	return text, text != ""
}

// SetClipboardText sets the clipboard text on Android
func (c *ClipboardManager) SetClipboardText(text string) error {
	// Android clipboard implementation would go here
	// For now, return nil as a placeholder
	return nil
}

// GetClipboardText gets the clipboard text on Android
func (c *ClipboardManager) GetClipboardText() (string, error) {
	// Android clipboard implementation would go here
	// For now, return empty string
	return "", nil
}
