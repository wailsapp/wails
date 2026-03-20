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
	return a.parent.Screen.LayoutScreens(screens)
}

func (a *linuxApp) getPrimaryScreen() (*Screen, error) {
	return a.parent.Screen.primaryScreen, nil
}

func (a *linuxApp) getScreens() ([]*Screen, error) {
	return a.parent.Screen.screens, nil
}

func getScreenForWindow(window *linuxWebviewWindow) (*Screen, error) {
	return window.getScreen()
}
