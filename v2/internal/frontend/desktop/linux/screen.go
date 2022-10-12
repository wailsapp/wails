//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0
#cgo CFLAGS: -w
#include <stdio.h>
#include "webkit2/webkit2.h"
#include "gtk/gtk.h"
#include "gdk/gdk.h"

typedef struct Screen {
	int isCurrent;
	int isPrimary;
	int height;
	int width;
} Screen;

int GetNMonitors(GtkWindow *window){
	GdkWindow *gdk_window = gtk_widget_get_window(GTK_WIDGET(window));
	GdkDisplay *display = gdk_window_get_display(gdk_window);
	return gdk_display_get_n_monitors(display);
}

Screen GetNThMonitor(int monitor_num, GtkWindow *window){
	GdkWindow *gdk_window = gtk_widget_get_window(GTK_WIDGET(window));
	GdkDisplay *display = gdk_window_get_display(gdk_window);
	GdkMonitor *monitor = gdk_display_get_monitor(display,monitor_num);
	GdkMonitor *currentMonitor = gdk_display_get_monitor_at_window(display,gdk_window);
	Screen screen;
	GdkRectangle geometry;
	gdk_monitor_get_geometry(monitor,&geometry);
	screen.isCurrent = currentMonitor==monitor;
	screen.isPrimary = gdk_monitor_is_primary(monitor);
	screen.height = geometry.height;
	screen.width = geometry.width;
	return screen;
}
*/
import "C"
import (
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"sync"
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
		numMonitors := C.GetNMonitors(window)
		for i := 0; i < int(numMonitors); i++ {
			cMonitor := C.GetNThMonitor(C.int(i), window)
			screen := Screen{
				IsCurrent: cMonitor.isCurrent == 1,
				IsPrimary: cMonitor.isPrimary == 1,
				Width:     int(cMonitor.width),
				Height:    int(cMonitor.height),
			}
			screens = append(screens, screen)
		}

		wg.Done()
	})
	wg.Wait()
	return screens, nil
}
