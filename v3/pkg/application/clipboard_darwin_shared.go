//go:build darwin && !ios && !server

package application

func newClipboardImpl() *macosClipboard {
	return &macosClipboard{}
}
