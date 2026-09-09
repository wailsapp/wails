//go:build android

package application

// The clipboard is backed by Android's ClipboardManager via the WailsBridge.
// Note: on Android 10+ reading the clipboard only succeeds while the app has
// input focus.

type androidClipboardImpl struct{}

func newClipboardImpl() clipboardImpl {
	return &androidClipboardImpl{}
}

func (c *androidClipboardImpl) setText(text string) bool {
	androidBridgeVoidString("setClipboardText", text)
	return true
}

func (c *androidClipboardImpl) text() (string, bool) {
	return androidBridgeString("getClipboardText")
}
