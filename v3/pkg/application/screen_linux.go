//go:build linux

package application

import (
	"fmt"
	"sync"
)

func (a *linuxApp) getPrimaryScreen() (*Screen, error) {
	return nil, fmt.Errorf("not implemented")
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
