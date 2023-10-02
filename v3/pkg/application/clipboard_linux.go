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
	clipboardSet(text)
	return true
}

func (m linuxClipboard) text() (string, bool) {
	clipboardLock.RLock()
	defer clipboardLock.RUnlock()
	return clipboardGet(), true
}

func newClipboardImpl() *linuxClipboard {
	return &linuxClipboard{}
}
