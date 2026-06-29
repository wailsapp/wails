//go:build ios

package application

/*
#include <stdlib.h>
#include "application_ios.h"
*/
import "C"

import "unsafe"

type iosClipboardImpl struct{}

func newClipboardImpl() clipboardImpl {
	return &iosClipboardImpl{}
}

func (c *iosClipboardImpl) setText(text string) bool {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.ios_clipboard_set_text(ctext)
	return true
}

func (c *iosClipboardImpl) text() (string, bool) {
	ctext := C.ios_clipboard_get_text()
	if ctext == nil {
		return "", false
	}
	defer C.free(unsafe.Pointer(ctext))
	return C.GoString(ctext), true
}
