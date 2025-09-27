//go:build linux
// +build linux

package linux

/*
#cgo webkit_6 CFLAGS: -DWEBKIT_6
#cgo !webkit_6 pkg-config: gtk+-3.0
#cgo webkit_6 pkg-config: gtk4
#cgo !(webkit2_41 || webkit_6) pkg-config: webkit2gtk-4.0
#cgo webkit2_41 pkg-config: webkit2gtk-4.1
#cgo webkit_6 pkg-config: webkitgtk-6.0

#include "gtk/gtk.h"

#ifdef WEBKIT_6
#include "webkit/webkit.h"
#else
#include "webkit2/webkit2.h"
#endif

static gchar* GetClipboardText() {
	#ifdef WEBKIT_6
		GdkClipboard *clip = gdk_display_get_primary_clipboard(gdk_display_get_default());
		GdkContentProvider *provider = gdk_clipboard_get_content(clip);

		GValue value = G_VALUE_INIT;
		g_value_init(&value, G_TYPE_STRING);

		if(!gdk_content_provider_get_value(provider, &value, NULL)) {
			return "";
		}

		return g_value_get_string(&value);
	#else
		GtkClipboard *clip = gtk_clipboard_get(GDK_SELECTION_CLIPBOARD);
		return gtk_clipboard_wait_for_text(clip);
	#endif
}

static void SetClipboardText(gchar* text) {
	#ifdef WEBKIT_6
		GdkDisplay *display = gdk_display_get_default();

		GdkClipboard *clip = gdk_display_get_primary_clipboard(display);
		gdk_clipboard_set_text(clip, text);

		clip = gdk_display_get_clipboard(display);
		gdk_clipboard_set_text(clip, text);
	#else
		GtkClipboard *clip = gtk_clipboard_get(GDK_SELECTION_CLIPBOARD);
		gtk_clipboard_set_text(clip, text, -1);

		clip = gtk_clipboard_get(GDK_SELECTION_PRIMARY);
		gtk_clipboard_set_text(clip, text, -1);
	#endif
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
