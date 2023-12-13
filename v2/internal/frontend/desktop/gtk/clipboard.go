//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include "gtk/gtk.h"
#include "webkit2/webkit2.h"

static gchar* GetClipboardText() {
	GtkClipboard *clip = gtk_clipboard_get(GDK_SELECTION_CLIPBOARD);
	return gtk_clipboard_wait_for_text(clip);
}

static void SetClipboardText(gchar* text) {
	GtkClipboard *clip = gtk_clipboard_get(GDK_SELECTION_CLIPBOARD);
	gtk_clipboard_set_text(clip, text, -1);

	clip = gtk_clipboard_get(GDK_SELECTION_PRIMARY);
	gtk_clipboard_set_text(clip, text, -1);
}
*/
import "C"
import "sync"

func (f *Frontend) ClipboardGetText() (string, error) {
	var text string
	var wg sync.WaitGroup
	wg.Add(1)
	invokeOnMainThread(func() {
		ctxt := C.GetClipboardText()
		defer C.g_free(C.gpointer(ctxt))
		text = C.GoString(ctxt)
		wg.Done()
	})
	wg.Wait()
	return text, nil
}

func (f *Frontend) ClipboardSetText(text string) error {
	invokeOnMainThread(func() {
		ctxt := (*C.gchar)(C.CString(text))
		defer C.g_free(C.gpointer(ctxt))
		C.SetClipboardText(ctxt)
	})
	return nil
}
