//go:build linux && purego

package application

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

func (m *linuxApp) getPrimaryScreen() (*Screen, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *linuxApp) getScreenByIndex(display uintptr, index int) *Screen {
	fmt.Println("getScreenByIndex")
	var getMonitor func(uintptr, int) uintptr
	purego.RegisterLibFunc(&getMonitor, gtk, "gdk_display_get_monitor")

	monitor := getMonitor(display, index)

	// TODO: Do we need to update Screen to contain current info?
	//	currentMonitor := C.gdk_display_get_monitor_at_window(display, window)

	var getGeometry func(uintptr, uintptr)
	purego.RegisterLibFunc(&getGeometry, gtk, "gdk_monitor_get_geometry")

	//var geometry C.GdkRectangle
	/*
		 struct GdkRectangle {
			  int x;
			  int y;
			  int width;
			  int height;
			}
	*/
	geometry := make([]byte, 16)
	getGeometry(monitor, uintptr(unsafe.Pointer(&geometry[0])))
	fmt.Println("geometry: %v\n", geometry)

	var isPrimary func(uintptr) int
	purego.RegisterLibFunc(&isPrimary, gtk, "gdk_monitor_is_primary")

	primary := false
	if isPrimary(monitor) == 1 {
		primary = true
	}

	return &Screen{
		IsPrimary: primary,
		Scale:     1.0,
		X:         0, //int(geometry.x),
		Y:         0, //int(geometry.y),
		Size: Size{
			Height: 1024, //int(geometry.height),
			Width:  1024, //int(geometry.width),
		},
	}
}

func (m *linuxApp) getScreens() ([]*Screen, error) {
	fmt.Println("getScreens")
	var wg sync.WaitGroup
	var screens []*Screen
	wg.Add(1)

	var getWindow func(uintptr) uintptr
	purego.RegisterLibFunc(&getWindow, gtk, "gtk_application_get_active_window")
	var getDisplay func(uintptr) uintptr
	purego.RegisterLibFunc(&getDisplay, gtk, "gdk_window_get_display")
	var getMonitorCount func(uintptr) int
	purego.RegisterLibFunc(&getMonitorCount, gtk, "getNMonitors")
	globalApplication.dispatchOnMainThread(func() {
		window := getWindow(m.application)
		display := getDisplay(window)
		count := getMonitorCount(display)
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
