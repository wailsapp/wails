package application

import (
	"log"

	"github.com/wailsapp/wails/exp/pkg/options"
)

// Messages sent from javascript get routed here
type windowMessage struct {
	windowId uint
	message  string
}

var messageBuffer = make(chan *windowMessage)

type Application interface {
	Run() error
}

type App struct {
	options                   *options.Application
	applicationEventListeners map[uint][]func()

	windows map[uint]*Window

	// Running
	running bool
}

func (a *App) On(eventID uint, callback func()) {
	a.applicationEventListeners[eventID] = append(a.applicationEventListeners[eventID], callback)
}

func (a *App) NewWindow(options *options.Window) *Window {
	// Ensure we have sane defaults
	if options.Width == 0 {
		options.Width = 1024
	}
	if options.Height == 0 {
		options.Height = 768
	}

	newWindow := NewWindow(options)
	id := newWindow.id
	if a.windows == nil {
		a.windows = make(map[uint]*Window)
	}
	a.windows[id] = newWindow
	if a.running {
		newWindow.Run()
	}

	return newWindow
}

func (a *App) Run() error {

	a.running = true
	go func() {
		for {
			event := <-applicationEvents
			a.handleApplicationEvent(event)
		}
	}()
	go func() {
		for {
			event := <-windowEvents
			a.handleWindowEvent(event)
		}
	}()
	go func() {
		for {
			event := <-messageBuffer
			a.handleMessage(event)
		}
	}()

	// run windows
	for _, window := range a.windows {
		go window.Run()
	}

	return a.run()
}

func (a *App) handleMessage(event *windowMessage) {
	// Get window from window map
	window, ok := a.windows[event.windowId]
	if !ok {
		log.Printf("Window #%d not found", event.windowId)
		return
	}
	// Get callback from window
	window.handleMessage(event.message)
}

func (a *App) handleWindowEvent(event *WindowEvent) {
	// Get window from window map
	window, ok := a.windows[event.WindowID]
	if !ok {
		log.Printf("Window #%d not found", event.WindowID)
		return
	}
	window.handleWindowEvent(event.EventID)
}
