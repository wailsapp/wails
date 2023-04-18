//go:build windows

package application

type windowsClipboard struct{}

func (m windowsClipboard) setText(text string) bool {
	//clipboardLock.Lock()
	//defer clipboardLock.Unlock()
	//cText := C.CString(text)
	//success := C.setClipboardText(cText)
	//C.free(unsafe.Pointer(cText))
	//return bool(success)
	panic("implement me")
}

func (m windowsClipboard) text() string {
	//clipboardLock.RLock()
	//defer clipboardLock.RUnlock()
	//clipboardText := C.getClipboardText()
	//result := C.GoString(clipboardText)
	//return result
	panic("implement me")
}

func newClipboardImpl() *windowsClipboard {
	return &windowsClipboard{}
}
