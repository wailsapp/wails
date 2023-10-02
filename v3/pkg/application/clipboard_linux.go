//go:build linux

package application

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include "gtk/gtk.h"
#include "webkit2/webkit2.h"

static gchar* getClipboardText() {
	GtkClipboard *clip = gtk_clipboard_get(GDK_SELECTION_CLIPBOARD);
	return gtk_clipboard_wait_for_text(clip);
}

static void setClipboardText(gchar* text) {
	GtkClipboard *clip = gtk_clipboard_get(GDK_SELECTION_CLIPBOARD);
	gtk_clipboard_set_text(clip, text, -1);

	clip = gtk_clipboard_get(GDK_SELECTION_PRIMARY);
	gtk_clipboard_set_text(clip, text, -1);
}
*/
import "C"
import (
	"sync"
	"unsafe"
)

var clipboardLock sync.RWMutex

type linuxClipboard struct{}

func (m linuxClipboard) setText(text string) bool {
	clipboardLock.Lock()
	defer clipboardLock.Unlock()
	cText := C.CString(text)
	C.setClipboardText(cText)
	C.free(unsafe.Pointer(cText))
	return true
}

func (m linuxClipboard) text() (string, bool) {
	clipboardLock.RLock()
	defer clipboardLock.RUnlock()
	clipboardText := C.getClipboardText()
	result := C.GoString(clipboardText)
	return result, true
}

func newClipboardImpl() *linuxClipboard {
	return &linuxClipboard{}
}
