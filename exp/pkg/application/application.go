package application

import "C"
import "github.com/wailsapp/wails/exp/pkg/options"

type Application interface {
	Run() error
}

type App struct {
	options              *options.Application
	systemEventListeners map[string][]func()

	windows []*Window

	// Running
	running bool
}

func (a *App) On(s string, callback func()) {
	a.systemEventListeners[s] = append(a.systemEventListeners[s], callback)
}

func (a *App) NewWindow(options *options.Window) *Window {
	newWindow := NewWindow(options)
	a.windows = append(a.windows, newWindow)

	if a.running {
		err := newWindow.Run()
		if err != nil {
			panic(err)
		}
	}

	return newWindow
}

func (a *App) Run() error {

	a.running = true
	go func() {
		for {
			event := <-systemEvents
			a.handleSystemEvent(event)
		}
	}()

	// run windows
	for _, window := range a.windows {
		err := window.Run()
		if err != nil {
			return err
		}
	}

	return a.run()
}
