//go:build linux

package application

import (
	"fmt"
	"sync"
)

func (m *linuxApp) getPrimaryScreen() (*Screen, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *linuxApp) getScreens() ([]*Screen, error) {
	var wg sync.WaitGroup
	var screens []*Screen
	var err error
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		screens, err = getScreens(m.application)
		wg.Done()
	})
	wg.Wait()
	return screens, err
}

func getScreenForWindow(window *linuxWebviewWindow) (*Screen, error) {
	return window.getScreen()
}
