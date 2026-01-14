//go:build linux && cgo && !android

package application

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.1 gdk-3.0

#include <gtk/gtk.h>
#include <gdk/gdk.h>
#include <webkit2/webkit2.h>
#include <stdio.h>
#include <stdlib.h>

// Create a new WebKitWebView for a panel
static GtkWidget* panel_new_webview() {
    WebKitUserContentManager *manager = webkit_user_content_manager_new();
    GtkWidget *webView = webkit_web_view_new_with_user_content_manager(manager);
    return webView;
}

// Create a fixed container to hold the panel webview at specific position
static GtkWidget* panel_new_fixed() {
    return gtk_fixed_new();
}

// Add webview to fixed container at position
static void panel_fixed_put(GtkWidget *fixed, GtkWidget *webview, int x, int y) {
    gtk_fixed_put(GTK_FIXED(fixed), webview, x, y);
}

// Move webview in fixed container
static void panel_fixed_move(GtkWidget *fixed, GtkWidget *webview, int x, int y) {
    gtk_fixed_move(GTK_FIXED(fixed), webview, x, y);
}

// Set webview size
static void panel_set_size(GtkWidget *webview, int width, int height) {
    gtk_widget_set_size_request(webview, width, height);
}

// Load URL in webview
static void panel_load_url(GtkWidget *webview, const char *url) {
    webkit_web_view_load_uri(WEBKIT_WEB_VIEW(webview), url);
}

// Load HTML in webview
static void panel_load_html(GtkWidget *webview, const char *html) {
    webkit_web_view_load_html(WEBKIT_WEB_VIEW(webview), html, NULL);
}

// Execute JavaScript
static void panel_exec_js(GtkWidget *webview, const char *js) {
    webkit_web_view_run_javascript(WEBKIT_WEB_VIEW(webview), js, NULL, NULL, NULL);
}

// Reload webview
static void panel_reload(GtkWidget *webview) {
    webkit_web_view_reload(WEBKIT_WEB_VIEW(webview));
}

// Force reload webview (bypass cache)
static void panel_force_reload(GtkWidget *webview) {
    webkit_web_view_reload_bypass_cache(WEBKIT_WEB_VIEW(webview));
}

// Show webview
static void panel_show(GtkWidget *webview) {
    gtk_widget_show(webview);
}

// Hide webview
static void panel_hide(GtkWidget *webview) {
    gtk_widget_hide(webview);
}

// Check if visible
static gboolean panel_is_visible(GtkWidget *webview) {
    return gtk_widget_get_visible(webview);
}

// Set zoom level
static void panel_set_zoom(GtkWidget *webview, double zoom) {
    webkit_web_view_set_zoom_level(WEBKIT_WEB_VIEW(webview), zoom);
}

// Get zoom level
static double panel_get_zoom(GtkWidget *webview) {
    return webkit_web_view_get_zoom_level(WEBKIT_WEB_VIEW(webview));
}

// Open inspector
static void panel_open_devtools(GtkWidget *webview) {
    WebKitWebInspector *inspector = webkit_web_view_get_inspector(WEBKIT_WEB_VIEW(webview));
    webkit_web_inspector_show(inspector);
}

// Focus webview
static void panel_focus(GtkWidget *webview) {
    gtk_widget_grab_focus(webview);
}

// Check if focused
static gboolean panel_is_focused(GtkWidget *webview) {
    return gtk_widget_has_focus(webview);
}

// Set background color
static void panel_set_background_color(GtkWidget *webview, int r, int g, int b, int a) {
    GdkRGBA color;
    color.red = r / 255.0;
    color.green = g / 255.0;
    color.blue = b / 255.0;
    color.alpha = a / 255.0;
    webkit_web_view_set_background_color(WEBKIT_WEB_VIEW(webview), &color);
}

// Enable/disable devtools
static void panel_enable_devtools(GtkWidget *webview, gboolean enable) {
    WebKitSettings *settings = webkit_web_view_get_settings(WEBKIT_WEB_VIEW(webview));
    webkit_settings_set_enable_developer_extras(settings, enable);
}

// Destroy the panel webview
static void panel_destroy(GtkWidget *webview) {
    gtk_widget_destroy(webview);
}

// Get position allocation
static void panel_get_allocation(GtkWidget *webview, int *x, int *y, int *width, int *height) {
    GtkAllocation alloc;
    gtk_widget_get_allocation(webview, &alloc);
    *x = alloc.x;
    *y = alloc.y;
    *width = alloc.width;
    *height = alloc.height;
}

*/
import "C"
import (
	"unsafe"
)

type linuxPanelImpl struct {
	panel   *WebviewPanel
	webview *C.GtkWidget
	fixed   *C.GtkWidget // Fixed container to position the webview
	parent  *linuxWebviewWindow
}

func newPanelImpl(panel *WebviewPanel) webviewPanelImpl {
	parentWindow := panel.parent
	if parentWindow == nil || parentWindow.impl == nil {
		return nil
	}

	linuxParent, ok := parentWindow.impl.(*linuxWebviewWindow)
	if !ok {
		return nil
	}

	return &linuxPanelImpl{
		panel:  panel,
		parent: linuxParent,
	}
}

func (p *linuxPanelImpl) create() {
	options := p.panel.options

	// Create the webview
	p.webview = C.panel_new_webview()

	// Set size
	C.panel_set_size(p.webview, C.int(options.Width), C.int(options.Height))

	// Create a fixed container if the parent's vbox doesn't have one for panels
	// For simplicity, we'll use an overlay approach - add the webview directly to the vbox
	// and use CSS/GTK positioning

	// Actually, we need to use GtkFixed or GtkOverlay for absolute positioning
	// For now, let's use the overlay approach with GtkFixed
	p.fixed = C.panel_new_fixed()

	// Add the webview to the fixed container at the specified position
	C.panel_fixed_put(p.fixed, p.webview, C.int(options.X), C.int(options.Y))

	// Add the fixed container to the parent's vbox (above the main webview)
	vbox := (*C.GtkBox)(p.parent.vbox)
	C.gtk_box_pack_start(vbox, p.fixed, 0, 0, 0) // Don't expand

	// Enable devtools if in debug mode
	debugMode := globalApplication.isDebugMode
	devToolsEnabled := debugMode
	if options.DevToolsEnabled != nil {
		devToolsEnabled = *options.DevToolsEnabled
	}
	C.panel_enable_devtools(p.webview, C.gboolean(boolToInt(devToolsEnabled)))

	// Set background color
	if options.Transparent {
		C.panel_set_background_color(p.webview, 0, 0, 0, 0)
	} else {
		C.panel_set_background_color(p.webview,
			C.int(options.BackgroundColour.Red),
			C.int(options.BackgroundColour.Green),
			C.int(options.BackgroundColour.Blue),
			C.int(options.BackgroundColour.Alpha),
		)
	}

	// Set zoom if specified
	if options.Zoom > 0 && options.Zoom != 1.0 {
		C.panel_set_zoom(p.webview, C.double(options.Zoom))
	}

	// Set initial visibility
	if options.Visible == nil || *options.Visible {
		C.gtk_widget_show_all(p.fixed)
	}

	// Load initial content
	if options.HTML != "" {
		html := C.CString(options.HTML)
		defer C.free(unsafe.Pointer(html))
		C.panel_load_html(p.webview, html)
	} else if options.URL != "" {
		url := C.CString(options.URL)
		defer C.free(unsafe.Pointer(url))
		C.panel_load_url(p.webview, url)
	}

	// Open inspector if requested
	if debugMode && options.OpenInspectorOnStartup {
		C.panel_open_devtools(p.webview)
	}

	// Mark runtime as loaded
	p.panel.markRuntimeLoaded()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (p *linuxPanelImpl) destroy() {
	if p.fixed != nil {
		C.panel_destroy(p.fixed)
		p.fixed = nil
		p.webview = nil
	}
}

func (p *linuxPanelImpl) setBounds(bounds Rect) {
	if p.webview == nil || p.fixed == nil {
		return
	}
	C.panel_fixed_move(p.fixed, p.webview, C.int(bounds.X), C.int(bounds.Y))
	C.panel_set_size(p.webview, C.int(bounds.Width), C.int(bounds.Height))
}

func (p *linuxPanelImpl) bounds() Rect {
	if p.webview == nil {
		return Rect{}
	}
	var x, y, width, height C.int
	C.panel_get_allocation(p.webview, &x, &y, &width, &height)
	return Rect{
		X:      int(x),
		Y:      int(y),
		Width:  int(width),
		Height: int(height),
	}
}

func (p *linuxPanelImpl) setZIndex(_ int) {
	// GTK doesn't have a direct z-index concept
	// We could use gtk_box_reorder_child to change ordering
	// For now, this is a no-op
}

func (p *linuxPanelImpl) setURL(url string) {
	if p.webview == nil {
		return
	}
	urlStr := C.CString(url)
	defer C.free(unsafe.Pointer(urlStr))
	C.panel_load_url(p.webview, urlStr)
}

func (p *linuxPanelImpl) setHTML(html string) {
	if p.webview == nil {
		return
	}
	htmlStr := C.CString(html)
	defer C.free(unsafe.Pointer(htmlStr))
	C.panel_load_html(p.webview, htmlStr)
}

func (p *linuxPanelImpl) execJS(js string) {
	if p.webview == nil {
		return
	}
	jsStr := C.CString(js)
	defer C.free(unsafe.Pointer(jsStr))
	C.panel_exec_js(p.webview, jsStr)
}

func (p *linuxPanelImpl) reload() {
	if p.webview == nil {
		return
	}
	C.panel_reload(p.webview)
}

func (p *linuxPanelImpl) forceReload() {
	if p.webview == nil {
		return
	}
	C.panel_force_reload(p.webview)
}

func (p *linuxPanelImpl) show() {
	if p.fixed == nil {
		return
	}
	C.gtk_widget_show_all(p.fixed)
}

func (p *linuxPanelImpl) hide() {
	if p.fixed == nil {
		return
	}
	C.gtk_widget_hide(p.fixed)
}

func (p *linuxPanelImpl) isVisible() bool {
	if p.fixed == nil {
		return false
	}
	return C.gtk_widget_get_visible(p.fixed) != 0
}

func (p *linuxPanelImpl) setZoom(zoom float64) {
	if p.webview == nil {
		return
	}
	C.panel_set_zoom(p.webview, C.double(zoom))
}

func (p *linuxPanelImpl) getZoom() float64 {
	if p.webview == nil {
		return 1.0
	}
	return float64(C.panel_get_zoom(p.webview))
}

func (p *linuxPanelImpl) openDevTools() {
	if p.webview == nil {
		return
	}
	C.panel_open_devtools(p.webview)
}

func (p *linuxPanelImpl) focus() {
	if p.webview == nil {
		return
	}
	C.panel_focus(p.webview)
}

func (p *linuxPanelImpl) isFocused() bool {
	if p.webview == nil {
		return false
	}
	return C.panel_is_focused(p.webview) != 0
}
