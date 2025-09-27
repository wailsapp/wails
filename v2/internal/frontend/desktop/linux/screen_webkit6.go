//go:build linux && webkit_6
// +build linux,webkit_6

package linux

/*
#cgo pkg-config: gtk4
#cgo webkit_6 pkg-config: webkitgtk-6.0

#cgo CFLAGS: -w
#include <stdio.h>

#include "webkit/webkit.h"
#include "gtk/gtk.h"
#include "gdk/gdk.h"

typedef struct Screen {
	int isCurrent;
	int isPrimary;
	int height;
	int width;
	int scale;
} Screen;

GListModel* GetMonitors(GtkWindow *window){
	GdkDisplay *display = gtk_widget_get_display(GTK_WIDGET(window));
	return gdk_display_get_monitors(display);
}

Screen GetNThMonitor(int monitor_num, GListModel *monitors, GtkWindow *window){
	GtkNative *native = gtk_widget_get_native(GTK_WIDGET(window));
	GdkSurface *surface = gtk_native_get_surface(native);

	GdkDisplay *display = gtk_widget_get_display(GTK_WIDGET(window));

	GdkMonitor *monitor = g_list_model_get_item(monitors, monitor_num);
	GdkMonitor *currentMonitor = gdk_display_get_monitor_at_surface(display, surface);

	Screen screen;
	GdkRectangle geometry;

	gdk_monitor_get_geometry(monitor, &geometry);

	screen.isCurrent = currentMonitor == monitor;
	// screen.isPrimary = gdk_monitor_is_primary(monitor); //// TODO: is_primary no longer exists on monitor
	screen.height = geometry.height;
	screen.width = geometry.width;
	screen.scale = gdk_monitor_get_scale_factor(monitor);

	return screen;
}
*/
import "C"
import (
	"sync"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/frontend"
)

type Screen = frontend.Screen

func GetAllScreens(window *C.GtkWindow) ([]Screen, error) {
	if window == nil {
		return nil, errors.New("window is nil, cannot perform screen operations")
	}
	var wg sync.WaitGroup
	var screens []Screen
	wg.Add(1)
	invokeOnMainThread(func() {
		monitors := C.GetMonitors(window)
		numMonitors := C.g_list_model_get_n_items(monitors)

		for i := 0; i < int(numMonitors); i++ {
			cMonitor := C.GetNThMonitor(C.int(i), monitors, window)

			screen := Screen{
				IsCurrent: cMonitor.isCurrent == 1,
				IsPrimary: cMonitor.isPrimary == 1,
				Width:     int(cMonitor.width),
				Height:    int(cMonitor.height),

				Size: frontend.ScreenSize{
					Width:  int(cMonitor.width),
					Height: int(cMonitor.height),
				},
				PhysicalSize: frontend.ScreenSize{
					Width:  int(cMonitor.width * cMonitor.scale),
					Height: int(cMonitor.height * cMonitor.scale),
				},
			}
			screens = append(screens, screen)
		}

		wg.Done()
	})
	wg.Wait()
	return screens, nil
}
