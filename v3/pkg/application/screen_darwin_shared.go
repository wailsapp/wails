//go:build darwin && !ios && !server

package application

func (m *macosApp) getPrimaryScreen() (*Screen, error) {
	if m.parent.Screen.GetPrimary() == nil {
		if err := m.processAndCacheScreens(); err != nil {
			return nil, err
		}
	}
	return m.parent.Screen.GetPrimary(), nil
}

func (m *macosApp) getScreens() ([]*Screen, error) {
	if len(m.parent.Screen.GetAll()) == 0 {
		if err := m.processAndCacheScreens(); err != nil {
			return nil, err
		}
	}
	return m.parent.Screen.GetAll(), nil
}

func (m *macosApp) processAndCacheScreens() error {
	// NSScreen and other AppKit APIs are not thread-safe and must be accessed on
	// the main thread. Application events (including ApplicationDidChangeScreenParameters)
	// are dispatched on background goroutines and can fire several times in quick
	// succession during a display reconfiguration, so without marshalling this
	// enumerates [NSScreen screens] concurrently off the main thread and crashes
	// (SIGSEGV). Running on the main run loop also serialises the burst of events.
	// Guard against InvokeSync deadlocking when we are already on the main thread.
	var screens []*Screen
	if m.isOnMainThread() {
		screens = allScreens()
	} else {
		InvokeSync(func() { screens = allScreens() })
	}
	return m.parent.Screen.layoutDIPScreens(screens)
}
