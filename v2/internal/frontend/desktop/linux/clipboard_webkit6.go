//go:build linux && webkit_6
// +build linux,webkit_6

package linux

/*
#cgo pkg-config: gtk4 webkitgtk-6.0

#include "gtk/gtk.h"
#include "webkit/webkit.h"

static gchar* GetClipboardText() {
	GdkClipboard *clip = gdk_display_get_primary_clipboard(gdk_display_get_default());
	GdkContentProvider *provider = gdk_clipboard_get_content(clip);

	GValue value = G_VALUE_INIT;
	g_value_init(&value, G_TYPE_STRING);

	if(!gdk_content_provider_get_value(provider, &value, NULL)) {
		return "";
	}

	return g_value_get_string(&value);
}

static void SetClipboardText(gchar* text) {
	GdkDisplay *display = gdk_display_get_default();

	GdkClipboard *clip = gdk_display_get_primary_clipboard(display);
	gdk_clipboard_set_text(clip, text);

	clip = gdk_display_get_clipboard(display);
	gdk_clipboard_set_text(clip, text);
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
