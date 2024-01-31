//go:build darwin

package darwin

import (
	"fmt"
)

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.13

#import <Cocoa/Cocoa.h>
#import <stdlib.h>

bool setClipboardText(const char* text) {
	NSPasteboard *pasteBoard = [NSPasteboard generalPasteboard];
	NSError *error = nil;
	NSString *string = [NSString stringWithUTF8String:text];
	[pasteBoard clearContents];
	return [pasteBoard setString:string forType:NSPasteboardTypeString];
}

const char* getClipboardText() {
	NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
	NSString *text = [pasteboard stringForType:NSPasteboardTypeString];
	return [text UTF8String];
}

*/
import "C"
import (
	"sync"
	"unsafe"
)

var clipboardLock sync.RWMutex

func (f *Frontend) ClipboardGetText() (string, error) {
	clipboardLock.RLock()
	defer clipboardLock.RUnlock()
	clipboardText := C.getClipboardText()
	result := C.GoString(clipboardText)
	return result, nil
}

func (f *Frontend) ClipboardSetText(text string) error {
	clipboardLock.Lock()
	defer clipboardLock.Unlock()
	cText := C.CString(text)
	success := C.setClipboardText(cText)
	C.free(unsafe.Pointer(cText))
	if !success {
		return fmt.Errorf("unable to set clipboard text")
	}
	return nil
}
