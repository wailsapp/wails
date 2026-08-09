//go:build linux && cgo && gtk3 && !android && !server

package webcontentsview

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.1 gdk-3.0
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

static void* createWebContentsView_gtk3(int w, int h, int devTools, int js, int images) {
    WebKitSettings *settings = webkit_settings_new();
    webkit_settings_set_enable_developer_extras(settings, devTools ? TRUE : FALSE);
    webkit_settings_set_enable_javascript(settings, js ? TRUE : FALSE);
    webkit_settings_set_auto_load_images(settings, images ? TRUE : FALSE);
    GtkWidget *webview = webkit_web_view_new_with_settings(settings);
    g_object_unref(settings);
    // Keep a reference independent of the GtkFixed parent so RemoveChildView
    // followed by AddChildView can safely reuse the same native browser.
    g_object_ref_sink(webview);
    gtk_widget_set_size_request(webview, w, h);
    return webview;
}

static void webContentsViewSetBounds_gtk3(void* view, void* parentOverlay, int x, int y, int w, int h) {
    GtkWidget *webview = (GtkWidget*)view;
    gtk_widget_set_size_request(webview, w, h);
    if (parentOverlay != NULL) {
        gtk_widget_set_margin_left(webview, x);
        gtk_widget_set_margin_top(webview, y);
    }
}

static void webContentsViewSetVisible_gtk3(void* view, int visible) {
    if (visible) gtk_widget_show(GTK_WIDGET(view));
    else gtk_widget_hide(GTK_WIDGET(view));
}

static void webContentsViewSetURL_gtk3(void* view, const char* url) {
    webkit_web_view_load_uri(WEBKIT_WEB_VIEW(view), url);
}

static void webContentsViewSetHTML_gtk3(void* view, const char* html) {
    webkit_web_view_load_html(WEBKIT_WEB_VIEW(view), html, NULL);
}

static void webContentsViewExecJS_gtk3(void* view, const char* js) {
    webkit_web_view_evaluate_javascript(WEBKIT_WEB_VIEW(view), js, -1, NULL, NULL, NULL, NULL, NULL);
}

static void webContentsViewGoBack_gtk3(void* view) {
    WebKitWebView *webview = WEBKIT_WEB_VIEW(view);
    if (webkit_web_view_can_go_back(webview)) webkit_web_view_go_back(webview);
}

static const char* webContentsViewGetURL_gtk3(void* view) {
    return webkit_web_view_get_uri(WEBKIT_WEB_VIEW(view));
}

static void* webContentsViewAttach_gtk3(void* window, void* view, int x, int y) {
    GtkWidget *child = gtk_bin_get_child(GTK_BIN(window));
    if (child == NULL || !GTK_IS_OVERLAY(child)) return NULL;
    GtkWidget *webview = GTK_WIDGET(view);
    gtk_widget_set_halign(webview, GTK_ALIGN_START);
    gtk_widget_set_valign(webview, GTK_ALIGN_START);
    gtk_widget_set_margin_left(webview, x);
    gtk_widget_set_margin_top(webview, y);
    gtk_overlay_add_overlay(GTK_OVERLAY(child), webview);
    gtk_widget_show(webview);
    return child;
}

static void webContentsViewDetach_gtk3(void* view, void* overlay) {
    GtkWidget *webview = GTK_WIDGET(view);
    if (overlay != NULL && GTK_IS_OVERLAY(overlay)) gtk_container_remove(GTK_CONTAINER(overlay), webview);
}

static void webContentsViewDestroy_gtk3(void* view) {
    if (view != NULL) g_object_unref(view);
}
*/
import "C"

import (
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type gtk3WebContentsView struct {
	parent               *WebContentsView
	widget               unsafe.Pointer
	fixed                unsafe.Pointer
	pendingJS            []string
	devTools, js, images int
}

func newWebContentsViewImpl(parent *WebContentsView) webContentsViewImpl {
	options := parent.optionsSnapshot()
	devTools, js, images := 1, 1, 1
	if options.WebPreferences.DevTools.IsSet() && !options.WebPreferences.DevTools.Get() {
		devTools = 0
	}
	if options.WebPreferences.Javascript.IsSet() && !options.WebPreferences.Javascript.Get() {
		js = 0
	}
	if options.WebPreferences.Images.IsSet() && !options.WebPreferences.Images.Get() {
		images = 0
	}
	return &gtk3WebContentsView{parent: parent, devTools: devTools, js: js, images: images}
}

func (w *gtk3WebContentsView) createWidget() {
	if w.widget == nil {
		options := w.parent.optionsSnapshot()
		w.widget = C.createWebContentsView_gtk3(C.int(options.Bounds.Width), C.int(options.Bounds.Height), C.int(w.devTools), C.int(w.js), C.int(w.images))
	}
}

func (w *gtk3WebContentsView) setBounds(bounds application.Rect) {
	if w.widget != nil {
		application.InvokeSync(func() {
			C.webContentsViewSetBounds_gtk3(w.widget, w.fixed, C.int(bounds.X), C.int(bounds.Y), C.int(bounds.Width), C.int(bounds.Height))
		})
	}
}

func (w *gtk3WebContentsView) setURL(url string) {
	if w.widget == nil {
		return
	}
	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))
	application.InvokeSync(func() { C.webContentsViewSetURL_gtk3(w.widget, cURL) })
}

func (w *gtk3WebContentsView) setHTML(html string) {
	if w.widget == nil {
		return
	}
	cHTML := C.CString(html)
	defer C.free(unsafe.Pointer(cHTML))
	application.InvokeSync(func() { C.webContentsViewSetHTML_gtk3(w.widget, cHTML) })
}

func (w *gtk3WebContentsView) execJS(js string) {
	if w.widget == nil {
		w.pendingJS = append(w.pendingJS, js)
		return
	}
	cJS := C.CString(js)
	defer C.free(unsafe.Pointer(cJS))
	application.InvokeSync(func() { C.webContentsViewExecJS_gtk3(w.widget, cJS) })
}

func (w *gtk3WebContentsView) goBack() {
	if w.widget != nil {
		application.InvokeSync(func() { C.webContentsViewGoBack_gtk3(w.widget) })
	}
}

func (w *gtk3WebContentsView) getURL() string {
	if w.widget == nil {
		return ""
	}
	var url string
	application.InvokeSync(func() {
		if nativeURL := C.webContentsViewGetURL_gtk3(w.widget); nativeURL != nil {
			// The WebKit-owned string is only stable for the current GTK call.
			url = C.GoString(nativeURL)
		}
	})
	return url
}

func (w *gtk3WebContentsView) takeSnapshot() string { return "" }

func (w *gtk3WebContentsView) setVisible(visible bool) {
	if w.widget == nil {
		return
	}
	application.InvokeSync(func() {
		C.webContentsViewSetVisible_gtk3(w.widget, C.int(boolToInt(visible)))
	})
}

func (w *gtk3WebContentsView) attach(window application.Window) {
	if window.NativeWindow() == nil || w.fixed != nil {
		return
	}
	application.InvokeSync(func() {
		w.createWidget()
		options := w.parent.optionsSnapshot()
		visible := w.parent.visibleSnapshot()
		w.fixed = C.webContentsViewAttach_gtk3(window.NativeWindow(), w.widget, C.int(options.Bounds.X), C.int(options.Bounds.Y))
		C.webContentsViewSetVisible_gtk3(w.widget, C.int(boolToInt(visible)))
		if options.URL != "" {
			cURL := C.CString(options.URL)
			C.webContentsViewSetURL_gtk3(w.widget, cURL)
			C.free(unsafe.Pointer(cURL))
		} else if options.HTML != "" {
			cHTML := C.CString(options.HTML)
			C.webContentsViewSetHTML_gtk3(w.widget, cHTML)
			C.free(unsafe.Pointer(cHTML))
		}
		for _, js := range w.pendingJS {
			cJS := C.CString(js)
			C.webContentsViewExecJS_gtk3(w.widget, cJS)
			C.free(unsafe.Pointer(cJS))
		}
		w.pendingJS = nil
	})
}

func (w *gtk3WebContentsView) detach() {
	if w.widget != nil {
		application.InvokeSync(func() { C.webContentsViewDetach_gtk3(w.widget, w.fixed) })
		w.fixed = nil
	}
}

func (w *gtk3WebContentsView) destroy() {
	if w.widget == nil {
		return
	}
	application.InvokeSync(func() { C.webContentsViewDestroy_gtk3(w.widget) })
	w.widget = nil
	w.fixed = nil
}

func (w *gtk3WebContentsView) nativeView() unsafe.Pointer { return w.widget }
