//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0


#include "gtk/gtk.h"
#include <stdio.h>
#include <limits.h>

static GtkWidget* GTKWIDGET(void *pointer) {
	return GTK_WIDGET(pointer);
}

static GtkWindow* GTKWINDOW(void *pointer) {
	return GTK_WINDOW(pointer);
}

static void SetMinSize(GtkWindow* window, int width, int height) {
	GdkGeometry size;
	size.min_height = height;
	size.min_width = width;
	gtk_window_set_geometry_hints(window, NULL, &size, GDK_HINT_MIN_SIZE);
}

static void SetMaxSize(GtkWindow* window, int width, int height) {
	GdkGeometry size;
	if( width == 0 ) {
		width = INT_MAX;
	}
	if( height == 0 ) {
		height = INT_MAX;
	}

	size.max_height = height;
	size.max_width = width;
	gtk_window_set_geometry_hints(window, NULL, &size, GDK_HINT_MAX_SIZE);
}

GdkRectangle getCurrentMonitorGeometry(GtkWindow *window) {
    // Get the monitor that the window is currently on
    GdkDisplay *display = gtk_widget_get_display(GTK_WIDGET(window));
    GdkWindow *gdk_window = gtk_widget_get_window(GTK_WIDGET(window));
    GdkMonitor *monitor = gdk_display_get_monitor_at_window (display, gdk_window);

    // Get the geometry of the monitor
    GdkRectangle result;
    gdk_monitor_get_geometry (monitor,&result);
    return result;
}

void SetPosition(GtkWindow *window, int x, int y) {
	GdkRectangle monitorDimensions = getCurrentMonitorGeometry(window);
	gtk_window_move(window, monitorDimensions.x + x, monitorDimensions.y + y);
}

void Center(GtkWindow *window)
{
    // Get the geometry of the monitor
    GdkRectangle m = getCurrentMonitorGeometry(window);

    // Get the window width/height
    int windowWidth, windowHeight;
    gtk_window_get_size(window, &windowWidth, &windowHeight);

	int newX = ((m.width - windowWidth) / 2) + m.x;
	int newY = ((m.height - windowHeight) / 2) + m.y;

    // Place the window at the center of the monitor
    gtk_window_move(window, newX, newY);
}

int IsFullscreen(GtkWidget *widget) {
	GdkWindow *gdkwindow = gtk_widget_get_window(widget);
	GdkWindowState state = gdk_window_get_state(GDK_WINDOW(gdkwindow));
	return state & GDK_WINDOW_STATE_FULLSCREEN == GDK_WINDOW_STATE_FULLSCREEN;
}

*/
import "C"
import (
	"github.com/wailsapp/wails/v2/pkg/options"
	"unsafe"
)

func gtkBool(input bool) C.gboolean {
	if input {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

type Window struct {
	appoptions *options.App
	debug      bool
	gtkWindow  unsafe.Pointer
}

func NewWindow(appoptions *options.App, debug bool) *Window {

	result := &Window{
		appoptions: appoptions,
		debug:      debug,
	}

	gtkWindow := C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
	C.g_object_ref_sink(C.gpointer(gtkWindow))

	result.gtkWindow = unsafe.Pointer(gtkWindow)
	result.SetKeepAbove(appoptions.AlwaysOnTop)
	result.SetResizable(!appoptions.DisableResize)
	result.SetSize(appoptions.Width, appoptions.Height)
	result.Center()
	result.SetDecorated(!appoptions.Frameless)
	result.SetTitle(appoptions.Title)
	result.SetMinSize(appoptions.MinWidth, appoptions.MinHeight)
	result.SetMaxSize(appoptions.MaxWidth, appoptions.MaxHeight)

	//if appoptions.Linux != nil && appoptions.Linux.Icon != nil {
	//	xpmData := png2XPM(appoptions.Linux.Icon)
	//	xpm := C.CString(xpmData)
	//	defer C.free(unsafe.Pointer(xpm))
	//	appIcon := C.gdk_pixbuf_new_from_xpm_data(
	//	C.gtk_window_set_icon(result.asGTKWindow(), appIcon)
	//}

	//windowStartState := C.int(int(a.appoptions.WindowStartState))

	//if a.appoptions.RGBA != nil {
	//	result.SetRGBA(a.appoptions.RGBA.R, a.appoptions.RGBA.G, a.appoptions.RGBA.B, a.appoptions.RGBA.A)
	//}

	//if a.appoptions.Menu != nil {
	//	result.SetApplicationMenu(a.appoptions.Menu)
	//}

	return result
}

func (w *Window) asGTKWidget() *C.GtkWidget {
	return C.GTKWIDGET(w.gtkWindow)
}

func (w *Window) asGTKWindow() *C.GtkWindow {
	return C.GTKWINDOW(w.gtkWindow)
}

//func (w *Window) Dispatch(f func()) {
//	glib.IdleAdd(f)
//}
//

func (w *Window) Fullscreen() {
	C.gtk_window_fullscreen(w.asGTKWindow())
}

func (w *Window) UnFullscreen() {
	C.gtk_window_unfullscreen(w.asGTKWindow())
}

func (w *Window) Destroy() {
	/*
	   for (gulong connection: {
	       impl_->deleteEventConnection,
	       impl_->focusInEventConnection,
	       impl_->focusOutEventConnection,
	       impl_->configureEventConnection
	   }) {
	       g_signal_handler_disconnect(impl_->gtkWindow, connection);
	   }
	   gtk_widget_destroy(GTK_WIDGET(impl_->gtkWindow));
	*/

	//TODO: Proper shutdown
	C.g_object_unref(C.gpointer(w.gtkWindow))
	C.gtk_widget_destroy(w.asGTKWidget())
}

func (w *Window) Close() {
	C.gtk_window_close(w.asGTKWindow())
}

func (w *Window) Center() {
	C.Center(w.asGTKWindow())
}

func (w *Window) SetPos(x int, y int) {
	cX := C.int(x)
	cY := C.int(y)
	C.gtk_window_move(w.asGTKWindow(), cX, cY)
}

func (w *Window) Size() (int, int) {
	var width, height C.int
	C.gtk_window_get_size(w.asGTKWindow(), &width, &height)
	return int(width), int(height)
}

func (w *Window) Pos() (int, int) {
	var width, height C.int
	C.gtk_window_get_position(w.asGTKWindow(), &width, &height)
	return int(width), int(height)
}

func (w *Window) SetMaxSize(maxWidth int, maxHeight int) {
	C.SetMaxSize(w.asGTKWindow(), C.int(maxWidth), C.int(maxHeight))
}

func (w *Window) SetMinSize(minWidth int, minHeight int) {
	C.SetMinSize(w.asGTKWindow(), C.int(minWidth), C.int(minHeight))
}

func (w *Window) Show() {
	C.gtk_widget_show(w.asGTKWidget())
}

func (w *Window) Hide() {
	C.gtk_widget_hide(w.asGTKWidget())
}

func (w *Window) Maximise() {
	C.gtk_window_maximize(w.asGTKWindow())
}

func (w *Window) UnMaximise() {
	C.gtk_window_unmaximize(w.asGTKWindow())
}

func (w *Window) Minimise() {
	C.gtk_window_iconify(w.asGTKWindow())
}

func (w *Window) UnMinimise() {
	C.gtk_window_present(w.asGTKWindow())
}

func (w *Window) IsFullScreen() bool {
	result := C.IsFullscreen(w.asGTKWidget())
	if result == 1 {
		return true
	}
	return false
}

func (w *Window) SetRGBA(r uint8, g uint8, b uint8, a uint8) {
	//C.SetRGBA(w.context, C.int(r), C.int(g), C.int(b), C.int(a))
}

//func (w *Window) SetApplicationMenu(inMenu *menu.Menu) {
//	//mainMenu := NewNSMenu(w.context, "")
//	//processMenu(mainMenu, inMenu)
//	//C.SetAsApplicationMenu(w.context, mainMenu.nsmenu)
//}

func (w *Window) UpdateApplicationMenu() {
	//C.UpdateApplicationMenu(w.context)
}

func (w *Window) Run() {
	C.gtk_widget_show_all(w.asGTKWidget())
	//switch w.appoptions.WindowStartState {
	//case options.Fullscreen:
	//	w.Fullscreen()
	//case options.Minimised:
	//	w.Minimise()
	//case options.Maximised:
	//	w.Maximise()
	//}

	//println("Fullscreen: ", w.IsFullScreen())
	C.gtk_main()
}

func (w *Window) SetKeepAbove(top bool) {
	C.gtk_window_set_keep_above(w.asGTKWindow(), gtkBool(top))
}

func (w *Window) SetResizable(resizable bool) {
	C.gtk_window_set_resizable(w.asGTKWindow(), gtkBool(resizable))
}

func (w *Window) SetSize(width int, height int) {
	C.gtk_window_resize(w.asGTKWindow(), C.gint(width), C.gint(height))
}

func (w *Window) SetDecorated(frameless bool) {
	C.gtk_window_set_decorated(w.asGTKWindow(), gtkBool(frameless))
}

func (w *Window) SetTitle(title string) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.gtk_window_set_title(w.asGTKWindow(), cTitle)
}
