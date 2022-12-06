package application

import (
	"context"
	"github.com/wailsapp/wails/v2/internal/app"
	"github.com/wailsapp/wails/v2/internal/signal"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"sync"
)

// Application is the main Wails application
type Application struct {
	application *app.App
	options     *options.App

	// System Trays
	systemTrays []*SystemTray

	// running flag
	running bool

	shutdown sync.Once
}

// NewWithOptions creates a new Application with the given options
func NewWithOptions(options *options.App) *Application {
	if options == nil {
		return New()
	}
	return &Application{
		options: options,
	}
}

// New creates a new Application with the default options
func New() *Application {
	return &Application{
		options: &options.App{},
	}
}

// SetApplicationMenu sets the application menu
func (a *Application) SetApplicationMenu(appMenu *menu.Menu) {
	if a.running {
		a.application.SetApplicationMenu(appMenu)
		return
	}

	a.options.Menu = appMenu
}

// Run starts the application
func (a *Application) Run() error {

	for _, systemtray := range a.systemTrays {
		go systemtray.run()
	}

	err := applicationInit()
	if err != nil {
		return err
	}

	application, err := app.CreateApp(a.options)
	if err != nil {
		return err
	}

	a.application = application

	// Control-C handlers
	signal.OnShutdown(func() {
		a.application.Shutdown()
	})
	signal.Start()

	a.running = true

	err = a.application.Run()
	a.Quit()
	return err
}

// Quit will shut down the application
func (a *Application) Quit() {
	a.shutdown.Do(func() {
		for _, systray := range a.systemTrays {
			systray.Close()
		}
		a.application.Shutdown()
	})
}

// Bind the given struct to the application
func (a *Application) Bind(boundStruct any) {
	a.options.Bind = append(a.options.Bind, boundStruct)
}

func (a *Application) On(eventType EventType, callback func()) {

	c := func(ctx context.Context) {
		callback()
	}

	switch eventType {
	case StartUp:
		a.options.OnStartup = c
	case ShutDown:
		a.options.OnShutdown = c
	case DomReady:
		a.options.OnDomReady = c
	}
}

func (a *Application) NewSystemTray(options *options.SystemTray) *SystemTray {
	systemTray := newSystemTray(options)
	a.systemTrays = append(a.systemTrays, systemTray)
	return systemTray
}
