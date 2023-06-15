//go:build linux

package application

import (
	"sync"
)

var clipboardLock sync.RWMutex

type linuxClipboard struct{}

func (m linuxClipboard) setText(text string) bool {
	clipboardLock.Lock()
	defer clipboardLock.Unlock()
	// cText := C.CString(text)
	// success := C.setClipboardText(cText)
	// C.free(unsafe.Pointer(cText))
	success := false
	return bool(success)
}

func (m linuxClipboard) text() (string, bool) {
	clipboardLock.RLock()
	defer clipboardLock.RUnlock()
	//	clipboardText := C.getClipboardText()
	//	result := C.GoString(clipboardText)
	return "", false
}

func newClipboardImpl() *linuxClipboard {
	return &linuxClipboard{}
}
