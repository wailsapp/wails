//go:build linux

package application

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include <gtk/gtk.h>
#include <gdk/gdk.h>
#include <stdlib.h>
#include <stdbool.h>

typedef struct Screen {
	const char* id;
	const char* name;
	int p_width;
	int p_height;
	int width;
	int height;
	int x;
	int y;
	int w_width;
	int w_height;
	int w_x;
	int w_y;
	float scale;
	double rotation;
	bool isPrimary;
} Screen;


int GetNumScreens(){
    return 0;
}

*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

func (m *linuxApp) getPrimaryScreen() (*Screen, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *linuxApp) getScreenByIndex(display *C.struct__GdkDisplay, index int) *Screen {
	monitor := C.gdk_display_get_monitor(display, C.int(index))

	// TODO: Do we need to update Screen to contain current info?
	//	currentMonitor := C.gdk_display_get_monitor_at_window(display, window)

	var geometry C.GdkRectangle
	C.gdk_monitor_get_geometry(monitor, &geometry)
	primary := false
	if C.gdk_monitor_is_primary(monitor) == 1 {
		primary = true
	}

	return &Screen{
		IsPrimary: primary,
		Scale:     1.0,
		X:         int(geometry.x),
		Y:         int(geometry.y),
		Size: Size{
			Height: int(geometry.height),
			Width:  int(geometry.width),
		},
	}
}

func (m *linuxApp) getScreens() ([]*Screen, error) {
	var wg sync.WaitGroup
	var screens []*Screen
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		window := C.gtk_application_get_active_window((*C.GtkApplication)(m.application))
		display := C.gdk_window_get_display((*C.GdkWindow)(unsafe.Pointer(window)))
		count := C.gdk_display_get_n_monitors(display)
		for i := 0; i < int(count); i++ {
			screens = append(screens,
				m.getScreenByIndex(display, i),
			)
		}
		wg.Done()
	})
	wg.Wait()
	return screens, nil
}

func getScreenForWindow(window *linuxWebviewWindow) (*Screen, error) {
	return window.getScreen()
}
