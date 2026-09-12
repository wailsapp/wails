//go:build linux && cgo && !android && !server

package application

/*
#cgo gtk3 pkg-config: gtk+-3.0 webkit2gtk-4.1
#cgo !gtk3 pkg-config: gtk4 webkitgtk-6.0
#include <gtk/gtk.h>
#if GTK_MAJOR_VERSION >= 4
#include <webkit/webkit.h>
#else
#include <webkit2/webkit2.h>
#endif
#include <stdlib.h>

static gboolean embedded_panel_position(GtkOverlay *overlay, GtkWidget *widget, GdkRectangle *allocation, gpointer data) {
 GdkRectangle *bounds = g_object_get_data(G_OBJECT(widget), "wails-panel-bounds");
 if (!bounds) return FALSE;
 *allocation = *bounds;
 return TRUE;
}

// Wrap only the main webview, preserving the menu and content-area coordinates.
static GtkWidget* embedded_panel_overlay(GtkWidget *main, GtkWidget *box) {
 GtkWidget *parent = gtk_widget_get_parent(main);
 if (GTK_IS_OVERLAY(parent)) return parent;
 GtkWidget *overlay = gtk_overlay_new();
 g_object_ref(main);
#if GTK_MAJOR_VERSION >= 4
 gtk_box_remove(GTK_BOX(box), main);
 gtk_overlay_set_child(GTK_OVERLAY(overlay), main);
 gtk_box_append(GTK_BOX(box), overlay);
#else
 gtk_container_remove(GTK_CONTAINER(box), main);
 gtk_container_add(GTK_CONTAINER(overlay), main);
 gtk_box_pack_start(GTK_BOX(box), overlay, TRUE, TRUE, 0);
#endif
 g_object_unref(main);
 gtk_widget_set_hexpand(overlay, TRUE);
 gtk_widget_set_vexpand(overlay, TRUE);
 g_signal_connect(overlay, "get-child-position", G_CALLBACK(embedded_panel_position), NULL);
 gtk_widget_set_visible(overlay, TRUE);
 gtk_widget_set_visible(main, TRUE);
 return overlay;
}

static GtkWidget* embedded_panel_new(GtkWidget *main, unsigned int windowID) {
 WebKitWebContext *context = webkit_web_view_get_context(WEBKIT_WEB_VIEW(main));
 GtkWidget *view = g_object_new(WEBKIT_TYPE_WEB_VIEW, "web-context", context, NULL);
 // The shared context already has Wails' local asset scheme registered.
 g_object_set_data(G_OBJECT(view), "windowid", GUINT_TO_POINTER(windowID));
#if GTK_MAJOR_VERSION < 4
 // Parent show_all must not undo an explicitly hidden panel.
 gtk_widget_set_no_show_all(view, TRUE);
#endif
 return view;
}
static void embedded_panel_bounds(GtkWidget *view, int x, int y, int width, int height) {
 GdkRectangle *bounds = g_new(GdkRectangle, 1);
 *bounds = (GdkRectangle){x,y,width,height};
 g_object_set_data_full(G_OBJECT(view), "wails-panel-bounds", bounds, g_free);
 GtkWidget *parent = gtk_widget_get_parent(view);
 if (parent) gtk_widget_queue_resize(parent);
}
static void embedded_panel_raise(GtkWidget *overlay, GtkWidget *view) {
#if GTK_MAJOR_VERSION >= 4
 gtk_widget_insert_before(view, overlay, NULL);
#else
 gtk_overlay_reorder_overlay(GTK_OVERLAY(overlay), view, -1);
#endif
}
static void embedded_panel_remove(GtkWidget *overlay, GtkWidget *view) {
 webkit_web_view_stop_loading(WEBKIT_WEB_VIEW(view));
#if GTK_MAJOR_VERSION >= 4
 gtk_overlay_remove_overlay(GTK_OVERLAY(overlay), view);
#else
 gtk_widget_destroy(view);
#endif
}
static void embedded_panel_header(WebKitURIRequest *request, const char *key, const char *value) {
 SoupMessageHeaders *headers = webkit_uri_request_get_http_headers(request);
 if (headers) soup_message_headers_replace(headers, key, value);
}
static void embedded_panel_colour(WebKitWebView *view, int r, int g, int b, int a) {
 GdkRGBA colour = {r/255.0,g/255.0,b/255.0,a/255.0};
 webkit_web_view_set_background_color(view,&colour);
}
*/
import "C"
import (
	"unsafe"
)

type linuxPanelImpl struct {
	panel   *WebviewPanel
	webview *C.GtkWidget
	overlay *C.GtkWidget
	parent  *linuxWebviewWindow
}

func newPanelImpl(panel *WebviewPanel) webviewPanelImpl {
	parent, ok := panel.parent.impl.(*linuxWebviewWindow)
	if !ok || parent.webview == nil {
		return nil
	}
	return &linuxPanelImpl{panel: panel, parent: parent}
}
func (p *linuxPanelImpl) view() *C.WebKitWebView {
	return (*C.WebKitWebView)(unsafe.Pointer(p.webview))
}
func (p *linuxPanelImpl) create() {
	options := p.panel.snapshotOptions()
	p.overlay = C.embedded_panel_overlay((*C.GtkWidget)(p.parent.webview), (*C.GtkWidget)(p.parent.vbox))
	p.webview = C.embedded_panel_new((*C.GtkWidget)(p.parent.webview), C.uint(p.panel.parent.id))
	C.gtk_overlay_add_overlay((*C.GtkOverlay)(unsafe.Pointer(p.overlay)), p.webview)
	p.setBounds(Rect{X: options.X, Y: options.Y, Width: options.Width, Height: options.Height})
	settings := C.webkit_web_view_get_settings(p.view())
	enabled := globalApplication.isDebugMode
	if options.DevToolsEnabled != nil {
		enabled = *options.DevToolsEnabled
	}
	C.webkit_settings_set_enable_developer_extras(settings, gtkBool(enabled))
	if options.UserAgent != "" {
		agent := C.CString(options.UserAgent)
		C.webkit_settings_set_user_agent(settings, agent)
		C.free(unsafe.Pointer(agent))
	}
	colour := options.BackgroundColour
	if options.Transparent {
		colour = RGBA{}
	}
	C.embedded_panel_colour(p.view(), C.int(colour.Red), C.int(colour.Green), C.int(colour.Blue), C.int(colour.Alpha))
	p.setZoom(options.Zoom)
	if options.URL != "" {
		uri := C.CString(options.URL)
		request := C.webkit_uri_request_new(uri)
		C.free(unsafe.Pointer(uri))
		for key, value := range options.Headers {
			k, v := C.CString(key), C.CString(value)
			C.embedded_panel_header(request, k, v)
			C.free(unsafe.Pointer(k))
			C.free(unsafe.Pointer(v))
		}
		C.webkit_web_view_load_request(p.view(), request)
		C.g_object_unref(C.gpointer(request))
	}
	if *options.Visible {
		p.show()
	} else {
		p.hide()
	}
	if enabled && options.OpenInspectorOnStartup {
		p.openDevTools()
	}
}
func (p *linuxPanelImpl) destroy() {
	if p.webview != nil {
		C.embedded_panel_remove(p.overlay, p.webview)
		p.webview = nil
		p.overlay = nil
	}
}
func (p *linuxPanelImpl) setBounds(bounds Rect) {
	if p.webview != nil {
		C.embedded_panel_bounds(p.webview, C.int(bounds.X), C.int(bounds.Y), C.int(bounds.Width), C.int(bounds.Height))
	}
}
func (p *linuxPanelImpl) bounds() Rect {
	options := p.panel.snapshotOptions()
	return Rect{X: options.X, Y: options.Y, Width: options.Width, Height: options.Height}
}
func (p *linuxPanelImpl) setZIndex(_ int) {
	panels := p.panel.sortedSiblings()
	for _, panel := range panels {
		panel.destroyedLock.RLock()
		native, ok := panel.impl.(*linuxPanelImpl)
		panel.destroyedLock.RUnlock()
		if ok && native.webview != nil {
			C.embedded_panel_raise(native.overlay, native.webview)
		}
	}
}
func (p *linuxPanelImpl) setURL(url string) {
	if p.webview == nil {
		return
	}
	uri := C.CString(url)
	defer C.free(unsafe.Pointer(uri))
	C.webkit_web_view_load_uri(p.view(), uri)
}
func (p *linuxPanelImpl) reload() {
	if p.webview != nil {
		C.webkit_web_view_reload(p.view())
	}
}
func (p *linuxPanelImpl) forceReload() {
	if p.webview != nil {
		C.webkit_web_view_reload_bypass_cache(p.view())
	}
}
func (p *linuxPanelImpl) show() {
	if p.webview != nil {
		C.gtk_widget_set_visible(p.webview, 1)
	}
}
func (p *linuxPanelImpl) hide() {
	if p.webview != nil {
		C.gtk_widget_set_visible(p.webview, 0)
	}
}
func (p *linuxPanelImpl) isVisible() bool {
	return p.webview != nil && C.gtk_widget_get_visible(p.webview) != 0
}
func (p *linuxPanelImpl) setZoom(zoom float64) {
	if p.webview != nil {
		C.webkit_web_view_set_zoom_level(p.view(), C.double(zoom))
	}
}
func (p *linuxPanelImpl) getZoom() float64 {
	if p.webview != nil {
		return float64(C.webkit_web_view_get_zoom_level(p.view()))
	}
	return 1
}
func (p *linuxPanelImpl) openDevTools() {
	if p.webview != nil {
		C.webkit_web_inspector_show(C.webkit_web_view_get_inspector(p.view()))
	}
}
func (p *linuxPanelImpl) focus() {
	if p.webview != nil {
		C.gtk_widget_grab_focus(p.webview)
	}
}
func (p *linuxPanelImpl) isFocused() bool {
	return p.webview != nil && C.gtk_widget_has_focus(p.webview) != 0
}
