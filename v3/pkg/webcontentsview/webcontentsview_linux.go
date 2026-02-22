//go:build linux && cgo && !gtk4 && !android && !server

package webcontentsview

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.1 gdk-3.0
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

static void* createWebContentsView_linux(int x, int y, int w, int h, int devTools, int js, int images) {
    WebKitSettings *settings = webkit_settings_new();
    
    webkit_settings_set_enable_developer_extras(settings, devTools ? TRUE : FALSE);
    webkit_settings_set_enable_javascript(settings, js ? TRUE : FALSE);
    webkit_settings_set_auto_load_images(settings, images ? TRUE : FALSE);
    
    GtkWidget *webview = webkit_web_view_new_with_settings(settings);
    gtk_widget_set_size_request(webview, w, h);
    return webview;
}

static void webContentsViewSetBounds_linux(void* view, void* parentFixed, int x, int y, int w, int h) {
    GtkWidget *webview = (GtkWidget*)view;
    gtk_widget_set_size_request(webview, w, h);
    if (parentFixed != NULL) {
        gtk_fixed_move(GTK_FIXED(parentFixed), webview, x, y);
    }
}

static void webContentsViewSetURL_linux(void* view, const char* url) {
    webkit_web_view_load_uri(WEBKIT_WEB_VIEW((GtkWidget*)view), url);
}

static void webContentsViewExecJS_linux(void* view, const char* js) {
    webkit_web_view_run_javascript(WEBKIT_WEB_VIEW((GtkWidget*)view), js, NULL, NULL, NULL);
}

static void webContentsViewAttach_linux(void* window, void* view) {
    // Attempt to add to the main container. Wails v3 usually uses a vbox.
    GtkWindow *gtkWindow = GTK_WINDOW(window);
    GtkWidget *child = gtk_bin_get_child(GTK_BIN(gtkWindow));
    if (child != NULL && GTK_IS_BOX(child)) {
        gtk_box_pack_start(GTK_BOX(child), GTK_WIDGET(view), FALSE, FALSE, 0);
        gtk_widget_show(GTK_WIDGET(view));
    }
}

static void webContentsViewDetach_linux(void* view) {
    GtkWidget *webview = (GtkWidget*)view;
    GtkWidget *parent = gtk_widget_get_parent(webview);
    if (parent != NULL) {
        gtk_container_remove(GTK_CONTAINER(parent), webview);
    }
}
*/
import "C"
import (
	"unsafe"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type linuxWebContentsView struct {
	parent *WebContentsView
	widget unsafe.Pointer
}

func newWebContentsViewImpl(parent *WebContentsView) webContentsViewImpl {
	devTools := 1
	if parent.options.WebPreferences.DevTools == application.Disabled {
		devTools = 0
	}
	
	js := 1
	if parent.options.WebPreferences.Javascript == application.Disabled {
		js = 0
	}
	
	images := 1
	if parent.options.WebPreferences.Images == application.Disabled {
		images = 0
	}
	
	view := C.createWebContentsView_linux(
		C.int(parent.options.Bounds.X),
		C.int(parent.options.Bounds.Y),
		C.int(parent.options.Bounds.Width),
		C.int(parent.options.Bounds.Height),
		C.int(devTools),
		C.int(js),
		C.int(images),
	)

	result := &linuxWebContentsView{
		parent: parent,
		widget: view,
	}

	return result
}

func (w *linuxWebContentsView) setBounds(bounds application.Rect) {
	C.webContentsViewSetBounds_linux(w.widget, nil, C.int(bounds.X), C.int(bounds.Y), C.int(bounds.Width), C.int(bounds.Height))
}

func (w *linuxWebContentsView) setURL(url string) {
	cUrl := C.CString(url)
	defer C.free(unsafe.Pointer(cUrl))
	C.webContentsViewSetURL_linux(w.widget, cUrl)
}

func (w *linuxWebContentsView) goBack() {
	// TODO: webkit_web_view_go_back
}

func (w *linuxWebContentsView) getURL() string {
	return ""
}

func (w *linuxWebContentsView) execJS(js string) {
	cJs := C.CString(js)
	defer C.free(unsafe.Pointer(cJs))
	C.webContentsViewExecJS_linux(w.widget, cJs)
}

func (w *linuxWebContentsView) attach(window application.Window) {
	if window.NativeWindow() != nil {
		C.webContentsViewAttach_linux(window.NativeWindow(), w.widget)
		if w.parent.options.URL != "" {
			w.setURL(w.parent.options.URL)
		}
	}
}

func (w *linuxWebContentsView) detach() {
	C.webContentsViewDetach_linux(w.widget)
}

func (w *linuxWebContentsView) nativeView() unsafe.Pointer {
	return w.widget
}
