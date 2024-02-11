//go:build linux

package application

import (
	"sync"
)

func (a *linuxApp) getPrimaryScreen() (*Screen, error) {
	var wg sync.WaitGroup
	var screen *Screen
	var err error
	wg.Add(1)
	InvokeSync(func() {
		screen, err = getPrimaryScreen()
		wg.Done()
	})
	wg.Wait()
	return screen, err
}

func (a *linuxApp) getScreens() ([]*Screen, error) {
	var wg sync.WaitGroup
	var screens []*Screen
	var err error
	wg.Add(1)
	InvokeSync(func() {
		screens, err = getScreens(a.application)
		wg.Done()
	})
	wg.Wait()
	return screens, err
}

func getScreenForWindow(window *linuxWebviewWindow) (*Screen, error) {
	return window.getScreen()
}
