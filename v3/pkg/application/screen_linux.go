//go:build linux && !android && !server

package application

import (
	"sync"
)

func (a *linuxApp) processAndCacheScreens() error {
	var wg sync.WaitGroup
	var screens []*Screen
	var err error
	wg.Add(1)
	InvokeSync(func() {
		screens, err = getScreens(a.application)
		wg.Done()
	})
	wg.Wait()
	if err != nil {
		return err
	}
	// gdk_monitor_is_primary is unreliable on Wayland (always returns false).
	// If no screen reports as primary, default to index 0.
	hasPrimary := false
	for _, s := range screens {
		if s.IsPrimary {
			hasPrimary = true
			break
		}
	}
	if !hasPrimary && len(screens) > 0 {
		screens[0].IsPrimary = true
	}
	// GDK reports monitor geometry and work areas in logical coordinates, which
	// getScreenByIndex carries into Bounds and WorkArea as they are (scaling them
	// separately into the Physical* rectangles). They are in DIP space already,
	// so there is no device-pixel layout to rebuild.
	return a.parent.Screen.layoutDIPScreens(screens)
}

func (a *linuxApp) getPrimaryScreen() (*Screen, error) {
	return a.parent.Screen.GetPrimary(), nil
}

func (a *linuxApp) getScreens() ([]*Screen, error) {
	return a.parent.Screen.GetAll(), nil
}

func getScreenForWindow(window *linuxWebviewWindow) (*Screen, error) {
	return window.getScreen()
}
