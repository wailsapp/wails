package application

import "github.com/wailsapp/wails/exp/pkg/options"

type Application interface {
	Run() error
}

type App struct {
	options              *options.Application
	systemEventListeners map[string][]func()

	windows []*Window
}

func (a *App) On(s string, callback func()) {
	a.systemEventListeners[s] = append(a.systemEventListeners[s], callback)
}

func (a *App) NewWindow(options *options.Window) *Window {
	newWindow := NewWindow(options)
	a.windows = append(a.windows, newWindow)
	return newWindow
}
