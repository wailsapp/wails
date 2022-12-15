package application

import "C"
import (
	"log"
	"runtime"
	"sync"

	"github.com/wailsapp/wails/exp/pkg/events"

	"github.com/wailsapp/wails/exp/pkg/options"
)

func init() {
	runtime.LockOSThread()
}

func New() *App {
	return NewWithOptions(nil)
}

func NewWithOptions(appOptions *options.Application) *App {
	if appOptions == nil {
		appOptions = options.ApplicationDefaults
	}
	return &App{
		options:                   appOptions,
		applicationEventListeners: make(map[uint][]func()),
		systemTrays:               make(map[uint]*SystemTray),
	}
}

type platformApp interface {
	run() error
	destroy()
	setApplicationMenu(menu *Menu)
}

// Messages sent from javascript get routed here
type windowMessage struct {
	windowId uint
	message  string
}

var windowMessageBuffer = make(chan *windowMessage)

type App struct {
	options                   *options.Application
	applicationEventListeners map[uint][]func()

	// Windows
	windows           map[uint]*Window
	windowsLock       sync.Mutex
	windowAliases     map[string]uint
	windowAliasesLock sync.Mutex

	// System Trays
	systemTrays      map[uint]*SystemTray
	systemTraysLock  sync.Mutex
	systemTrayID     uint
	systemTrayIDLock sync.RWMutex

	// MenuItems
	menuItems     map[uint]*MenuItem
	menuItemsLock sync.Mutex

	// Running
	running bool

	// platform app
	impl platformApp

	// The main application menu
	ApplicationMenu *Menu
}

func (a *App) getSystemTrayID() uint {
	a.systemTrayIDLock.Lock()
	defer a.systemTrayIDLock.Unlock()
	a.systemTrayID++
	return a.systemTrayID
}
func (a *App) On(eventType events.ApplicationEventType, callback func()) {
	eventID := uint(eventType)
	a.applicationEventListeners[eventID] = append(a.applicationEventListeners[eventID], callback)
}
func (a *App) NewWindow() *Window {
	return a.NewWindowWithOptions(nil)
}

func (a *App) NewWindowWithOptions(windowOptions *options.Window) *Window {
	// Ensure we have sane defaults
	if windowOptions == nil {
		windowOptions = options.WindowDefaults
	}

	newWindow := NewWindow(windowOptions)
	id := newWindow.id
	if a.windows == nil {
		a.windows = make(map[uint]*Window)
	}
	a.windowsLock.Lock()
	a.windows[id] = newWindow
	a.windowsLock.Unlock()

	if windowOptions.Alias != "" {
		if a.windowAliases == nil {
			a.windowAliases = make(map[string]uint)
		}
		a.windowAliasesLock.Lock()
		a.windowAliases[windowOptions.Alias] = id
		a.windowAliasesLock.Unlock()
	}
	if a.running {
		newWindow.Run()
	}

	return newWindow
}

func (a *App) NewSystemTray() *SystemTray {

	id := a.getSystemTrayID()
	newSystemTray := NewSystemTray(id)
	a.systemTraysLock.Lock()
	a.systemTrays[id] = newSystemTray
	a.systemTraysLock.Unlock()

	if a.running {
		newSystemTray.Run()
	}
	return newSystemTray
}

func (a *App) Run() error {
	a.impl = newPlatformApp(a.options)

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
			event := <-windowMessageBuffer
			a.handleWindowMessage(event)
		}
	}()

	go func() {
		for {
			menuItemID := <-menuItemClicked
			a.handleMenuItemClicked(menuItemID)
		}
	}()

	// run windows
	for _, window := range a.windows {
		go window.Run()
	}

	// run system trays
	for _, systray := range a.systemTrays {
		go systray.Run()
	}

	if a.ApplicationMenu != nil {
		a.impl.setApplicationMenu(a.ApplicationMenu)
	}

	return a.impl.run()
}

func (a *App) handleApplicationEvent(event uint) {
	listeners, ok := a.applicationEventListeners[event]
	if !ok {
		return
	}
	for _, listener := range listeners {
		go listener()
	}
}

func (a *App) handleWindowMessage(event *windowMessage) {
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

func (a *App) handleMenuItemClicked(menuItemID uint) {
	menuItem := getMenuItemByID(menuItemID)
	if menuItem == nil {
		log.Printf("MenuItem #%d not found", menuItemID)
		return
	}
	menuItem.handleClick()
}

func (a *App) Quit() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		a.windowsLock.Lock()
		for _, window := range a.windows {
			window.Destroy()
		}
		a.windowsLock.Unlock()
		wg.Done()
	}()
	go func() {
		a.systemTraysLock.Lock()
		for _, systray := range a.systemTrays {
			systray.Destroy()
		}
		a.systemTraysLock.Unlock()
		wg.Done()
	}()
	wg.Wait()
	a.impl.destroy()
}

func (a *App) SetActivationPolicy(accessory options.ActivationPolicy) {
	a.options.Mac.ActivationPolicy = accessory
}

func (a *App) SetMenu(menu *Menu) {
	a.ApplicationMenu = menu
	if a.impl != nil {
		a.impl.setApplicationMenu(menu)
	}
}
