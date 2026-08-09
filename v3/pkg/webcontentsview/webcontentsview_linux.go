//go:build linux && cgo && !gtk3 && !android && !server

package webcontentsview

/*
#cgo linux pkg-config: gtk4 webkitgtk-6.0
#include <gtk/gtk.h>
#include <webkit/webkit.h>
#include <stdint.h>

static void* createWebContentsView_linux(int x, int y, int w, int h, int devTools, int js, int images) {
    WebKitSettings *settings = webkit_settings_new();

    webkit_settings_set_enable_developer_extras(settings, devTools ? TRUE : FALSE);
    webkit_settings_set_enable_javascript(settings, js ? TRUE : FALSE);
    webkit_settings_set_auto_load_images(settings, images ? TRUE : FALSE);

    GtkWidget *webview = webkit_web_view_new();
    webkit_web_view_set_settings(WEBKIT_WEB_VIEW(webview), settings);
    g_object_unref(settings);
    // Keep a reference independent of the GtkFixed parent. Detach removes the
    // parent-owned reference, but WebContentsView supports later reattachment.
    g_object_ref_sink(webview);
    gtk_widget_set_size_request(webview, w, h);
    return webview;
}

static void webContentsViewSetBounds_linux(void* view, void* parentOverlay, int x, int y, int w, int h) {
    GtkWidget *webview = (GtkWidget*)view;
    gtk_widget_set_size_request(webview, w, h);
    if (parentOverlay != NULL) {
        // This is an overlay child itself, rather than a child of a full-size
        // GtkFixed overlay. The latter receives pointer hits across its empty
        // area and prevents the host Wails webview from receiving controls.
        gtk_widget_set_margin_start(webview, x);
        gtk_widget_set_margin_top(webview, y);
    }
}

static void webContentsViewSetVisible_linux(void* view, int visible) {
    gtk_widget_set_visible(GTK_WIDGET(view), visible ? TRUE : FALSE);
}

static void webContentsViewSetURL_linux(void* view, const char* url) {
    webkit_web_view_load_uri(WEBKIT_WEB_VIEW((GtkWidget*)view), url);
}

static void webContentsViewSetHTML_linux(void* view, const char* html) {
    webkit_web_view_load_html(WEBKIT_WEB_VIEW((GtkWidget*)view), html, NULL);
}

static void webContentsViewExecJS_linux(void* view, const char* js) {
    webkit_web_view_evaluate_javascript(WEBKIT_WEB_VIEW((GtkWidget*)view), js, -1, NULL, NULL, NULL, NULL, NULL);
}

static void webContentsViewGoBack_linux(void* view) {
    WebKitWebView *webview = WEBKIT_WEB_VIEW((GtkWidget*)view);
    if (webkit_web_view_can_go_back(webview)) webkit_web_view_go_back(webview);
}

static const char* webContentsViewGetURL_linux(void* view) {
    return webkit_web_view_get_uri(WEBKIT_WEB_VIEW((GtkWidget*)view));
}

typedef struct { GMainLoop *loop; char *data; } SnapshotResult;
static void webContentsViewSnapshotReady(GObject *source, GAsyncResult *result, gpointer user_data) {
    SnapshotResult *snapshot = (SnapshotResult*)user_data;
    GError *error = NULL;
    GdkTexture *texture = webkit_web_view_get_snapshot_finish(WEBKIT_WEB_VIEW(source), result, &error);
    if (texture != NULL) {
        GBytes *png = gdk_texture_save_to_png_bytes(texture);
        if (png != NULL) {
            gsize size = 0;
            const guchar *bytes = g_bytes_get_data(png, &size);
            snapshot->data = g_base64_encode(bytes, size);
            g_bytes_unref(png);
        }
        g_object_unref(texture);
    }
    if (error != NULL) g_error_free(error);
    g_main_loop_quit(snapshot->loop);
}
static char* webContentsViewTakeSnapshot_linux(void* view) {
    SnapshotResult snapshot = {0};
    snapshot.loop = g_main_loop_new(NULL, FALSE);
    webkit_web_view_get_snapshot(WEBKIT_WEB_VIEW(view), WEBKIT_SNAPSHOT_REGION_VISIBLE, WEBKIT_SNAPSHOT_OPTIONS_NONE, NULL, webContentsViewSnapshotReady, &snapshot);
    g_main_loop_run(snapshot.loop);
    g_main_loop_unref(snapshot.loop);
    return snapshot.data;
}

static void* webContentsViewAttach_linux(void* window, void* view, int x, int y) {
    GtkWindow *gtkWindow = GTK_WINDOW(window);
    GtkWidget *child = gtk_window_get_child(gtkWindow);
    if (child != NULL && GTK_IS_OVERLAY(child)) {
        GtkWidget *webview = GTK_WIDGET(view);
        gtk_widget_set_halign(webview, GTK_ALIGN_START);
        gtk_widget_set_valign(webview, GTK_ALIGN_START);
        gtk_widget_set_margin_start(webview, x);
        gtk_widget_set_margin_top(webview, y);
        gtk_overlay_add_overlay(GTK_OVERLAY(child), webview);
        gtk_widget_set_visible(GTK_WIDGET(view), TRUE);
        return child;
    }
    return NULL;
}

static void webContentsViewDetach_linux(void* view, void* overlay) {
    GtkWidget *webview = (GtkWidget*)view;
    if (overlay != NULL && GTK_IS_OVERLAY(overlay)) {
        gtk_overlay_remove_overlay(GTK_OVERLAY(overlay), webview);
    }
}

static void webContentsViewDestroy_linux(void* view) {
    if (view != NULL) g_object_unref(view);
}
*/
import "C"
import (
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type linuxWebContentsView struct {
	parent               *WebContentsView
	widget               unsafe.Pointer
	fixed                unsafe.Pointer
	pendingJS            []string
	devTools, js, images int
}

func newWebContentsViewImpl(parent *WebContentsView) webContentsViewImpl {
	options := parent.optionsSnapshot()
	devTools := 1
	if !options.WebPreferences.DevTools.Get() && options.WebPreferences.DevTools.IsSet() {
		devTools = 0
	}

	js := 1
	if !options.WebPreferences.Javascript.Get() && options.WebPreferences.Javascript.IsSet() {
		js = 0
	}

	images := 1
	if !options.WebPreferences.Images.Get() && options.WebPreferences.Images.IsSet() {
		images = 0
	}

	return &linuxWebContentsView{
		parent:   parent,
		devTools: devTools,
		js:       js,
		images:   images,
	}
}

func (w *linuxWebContentsView) createWidget() {
	if w.widget != nil {
		return
	}

	options := w.parent.optionsSnapshot()
	w.widget = C.createWebContentsView_linux(
		C.int(options.Bounds.X),
		C.int(options.Bounds.Y),
		C.int(options.Bounds.Width),
		C.int(options.Bounds.Height),
		C.int(w.devTools),
		C.int(w.js),
		C.int(w.images),
	)
}

func (w *linuxWebContentsView) setBounds(bounds application.Rect) {
	if w.widget == nil {
		return
	}
	application.InvokeSync(func() {
		C.webContentsViewSetBounds_linux(w.widget, w.fixed, C.int(bounds.X), C.int(bounds.Y), C.int(bounds.Width), C.int(bounds.Height))
	})
}

func (w *linuxWebContentsView) setURL(url string) {
	if w.widget == nil {
		return
	}
	cUrl := C.CString(url)
	defer C.free(unsafe.Pointer(cUrl))
	application.InvokeSync(func() {
		C.webContentsViewSetURL_linux(w.widget, cUrl)
	})
}

func (w *linuxWebContentsView) setHTML(html string) {
	if w.widget == nil {
		return
	}
	cHTML := C.CString(html)
	defer C.free(unsafe.Pointer(cHTML))
	application.InvokeSync(func() {
		C.webContentsViewSetHTML_linux(w.widget, cHTML)
	})
}

func (w *linuxWebContentsView) goBack() {
	if w.widget == nil {
		return
	}
	application.InvokeSync(func() {
		C.webContentsViewGoBack_linux(w.widget)
	})
}

func (w *linuxWebContentsView) getURL() string {
	if w.widget == nil {
		return ""
	}
	var url string
	application.InvokeSync(func() {
		if nativeURL := C.webContentsViewGetURL_linux(w.widget); nativeURL != nil {
			// WebKit owns this pointer and may replace it on the next load event;
			// copy it before returning from the GTK thread.
			url = C.GoString(nativeURL)
		}
	})
	return url
}
func (w *linuxWebContentsView) takeSnapshot() string {
	if w.widget == nil {
		return ""
	}
	var snapshot *C.char
	application.InvokeSync(func() {
		snapshot = C.webContentsViewTakeSnapshot_linux(w.widget)
	})
	if snapshot == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(snapshot))
	return "data:image/png;base64," + C.GoString(snapshot)
}

func (w *linuxWebContentsView) setVisible(visible bool) {
	if w.widget == nil {
		return
	}
	application.InvokeSync(func() {
		C.webContentsViewSetVisible_linux(w.widget, C.int(boolToInt(visible)))
	})
}

func (w *linuxWebContentsView) execJS(js string) {
	if w.widget == nil {
		w.pendingJS = append(w.pendingJS, js)
		return
	}
	cJs := C.CString(js)
	defer C.free(unsafe.Pointer(cJs))
	application.InvokeSync(func() {
		C.webContentsViewExecJS_linux(w.widget, cJs)
	})
}

func (w *linuxWebContentsView) attach(window application.Window) {
	if window.NativeWindow() != nil && w.fixed == nil {
		application.InvokeSync(func() {
			w.createWidget()
			options := w.parent.optionsSnapshot()
			visible := w.parent.visibleSnapshot()
			w.fixed = C.webContentsViewAttach_linux(window.NativeWindow(), w.widget, C.int(options.Bounds.X), C.int(options.Bounds.Y))
			C.webContentsViewSetVisible_linux(w.widget, C.int(boolToInt(visible)))
			if options.URL != "" {
				cURL := C.CString(options.URL)
				C.webContentsViewSetURL_linux(w.widget, cURL)
				C.free(unsafe.Pointer(cURL))
			} else if options.HTML != "" {
				cHTML := C.CString(options.HTML)
				C.webContentsViewSetHTML_linux(w.widget, cHTML)
				C.free(unsafe.Pointer(cHTML))
			}
			for _, js := range w.pendingJS {
				cJS := C.CString(js)
				C.webContentsViewExecJS_linux(w.widget, cJS)
				C.free(unsafe.Pointer(cJS))
			}
			w.pendingJS = nil
		})
	}
}

func (w *linuxWebContentsView) detach() {
	if w.widget == nil {
		return
	}
	application.InvokeSync(func() {
		C.webContentsViewDetach_linux(w.widget, w.fixed)
	})
	w.fixed = nil
}

func (w *linuxWebContentsView) destroy() {
	if w.widget == nil {
		return
	}
	application.InvokeSync(func() {
		C.webContentsViewDestroy_linux(w.widget)
	})
	w.widget = nil
	w.fixed = nil
}

func (w *linuxWebContentsView) nativeView() unsafe.Pointer {
	return w.widget
}
