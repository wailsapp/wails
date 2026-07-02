//go:build darwin && purego && !ios && !server

package application

// CGO-free clipboard implementation.
//
// This mirrors the behaviour of the cgo clipboard_darwin.go by driving
// NSPasteboard directly through the Objective-C runtime helpers in
// darwin_purego_cocoa.go instead of compiling Objective-C.
//
//	setText -> [[NSPasteboard generalPasteboard] clearContents];
//	           [pb setString:string forType:NSPasteboardTypeString]
//	text    -> [[NSPasteboard generalPasteboard] stringForType:NSPasteboardTypeString]

import "sync"

// nsPasteboardTypeString is the value of the global NSPasteboardTypeString
// symbol (a UTI). Using the literal avoids a dlsym for the exported NSString.
const nsPasteboardTypeString = "public.utf8-plain-text"

var clipboardLock sync.RWMutex

type macosClipboard struct{}

// setText replaces the clipboard contents with text and reports whether the
// write succeeded, matching -[NSPasteboard setString:forType:].
func (m macosClipboard) setText(text string) bool {
	clipboardLock.Lock()
	defer clipboardLock.Unlock()

	pasteboard := class("NSPasteboard").send("generalPasteboard")
	str := nsString(text)
	typ := nsString(nsPasteboardTypeString)
	pasteboard.send("clearContents")
	return get[bool](pasteboard, "setString:forType:", str, typ)
}

// text returns the current plain-text clipboard contents. The bool return
// mirrors the cgo version, which always reports true (an empty/absent value is
// returned as the empty string).
func (m macosClipboard) text() (string, bool) {
	clipboardLock.RLock()
	defer clipboardLock.RUnlock()

	pasteboard := class("NSPasteboard").send("generalPasteboard")
	typ := nsString(nsPasteboardTypeString)
	result := pasteboard.send("stringForType:", typ).string()
	return result, true
}

func newClipboardImpl() *macosClipboard {
	return &macosClipboard{}
}
