package application

import (
	"log"
	"runtime"
	"sync"

	"github.com/wailsapp/wails/exp/pkg/options"
)

func init() {
	runtime.LockOSThread()
}

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

	windows           map[uint]*Window
	windowsLock       sync.Mutex
	windowAliases     map[string]uint
	windowAliasesLock sync.Mutex

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
	a.windowsLock.Lock()
	a.windows[id] = newWindow
	a.windowsLock.Unlock()

	if options.Alias != "" {
		if a.windowAliases == nil {
			a.windowAliases = make(map[string]uint)
		}
		a.windowAliasesLock.Lock()
		a.windowAliases[options.Alias] = id
		a.windowAliasesLock.Unlock()
	}
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
	a.windowsLock.Lock()
	window, ok := a.windows[event.windowId]
	a.windowsLock.Unlock()
	if !ok {
		log.Printf("Window #%d not found", event.windowId)
		return
	}
	// Get callback from window
	window.handleMessage(event.message)
}

func (a *App) handleWindowEvent(event *WindowEvent) {
	// Get window from window map
	a.windowsLock.Lock()
	window, ok := a.windows[event.WindowID]
	a.windowsLock.Unlock()
	if !ok {
		log.Printf("Window #%d not found", event.WindowID)
		return
	}
	window.handleWindowEvent(event.EventID)
}
