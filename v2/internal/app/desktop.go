// +build desktop,!server

package app

import (
	"context"
	"sync"

	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/ffenestri"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/menumanager"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher"
	"github.com/wailsapp/wails/v2/internal/runtime"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/internal/signal"
	"github.com/wailsapp/wails/v2/internal/subsystem"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// App defines a Wails application structure
type App struct {
	appType string

	window     *ffenestri.Application
	servicebus *servicebus.ServiceBus
	logger     *logger.Logger
	signal     *signal.Manager
	options    *options.App

	// Subsystems
	log     *subsystem.Log
	runtime *subsystem.Runtime
	event   *subsystem.Event
	//binding    *subsystem.Binding
	call       *subsystem.Call
	menu       *subsystem.Menu
	url        *subsystem.URL
	dispatcher *messagedispatcher.Dispatcher

	menuManager *menumanager.Manager

	// Indicates if the app is in debug mode
	debug bool

	// This is our binding DB
	bindings *binding.Bindings

	// Application Stores
	loglevelStore  *runtime.Store
	appconfigStore *runtime.Store

	// Startup/Shutdown
	startupCallback  func(*runtime.Runtime)
	shutdownCallback func()
}

// Create App
func CreateApp(appoptions *options.App) (*App, error) {

	// Merge default options
	options.MergeDefaults(appoptions)

	// Set up logger
	myLogger := logger.New(appoptions.Logger)
	myLogger.SetLogLevel(appoptions.LogLevel)

	// Create the menu manager
	menuManager := menumanager.NewManager()

	// Process the application menu
	menuManager.SetApplicationMenu(options.GetApplicationMenu(appoptions))

	// Process context menus
	contextMenus := options.GetContextMenus(appoptions)
	for _, contextMenu := range contextMenus {
		menuManager.AddContextMenu(contextMenu)
	}

	// Process tray menus
	trayMenus := options.GetTrayMenus(appoptions)
	for _, trayMenu := range trayMenus {
		menuManager.AddTrayMenu(trayMenu)
	}

	window := ffenestri.NewApplicationWithConfig(appoptions, myLogger, menuManager)

	// Create binding exemptions - Ugly hack. There must be a better way
	bindingExemptions := []interface{}{appoptions.Startup, appoptions.Shutdown}

	result := &App{
		appType:          "desktop",
		window:           window,
		servicebus:       servicebus.New(myLogger),
		logger:           myLogger,
		bindings:         binding.NewBindings(myLogger, appoptions.Bind, bindingExemptions),
		menuManager:      menuManager,
		startupCallback:  appoptions.Startup,
		shutdownCallback: appoptions.Shutdown,
	}

	result.options = appoptions

	// Initialise the app
	err := result.Init()

	return result, err

}

// Run the application
func (a *App) Run() error {

	var err error

	// Setup a context
	var subsystemWaitGroup sync.WaitGroup
	parentContext := context.WithValue(context.Background(), "waitgroup", &subsystemWaitGroup)
	ctx, cancel := context.WithCancel(parentContext)

	// Start the service bus
	a.servicebus.Debug()
	err = a.servicebus.Start()
	if err != nil {
		return err
	}

	runtimesubsystem, err := subsystem.NewRuntime(ctx, a.servicebus, a.logger, a.startupCallback)
	if err != nil {
		return err
	}
	a.runtime = runtimesubsystem
	err = a.runtime.Start()
	if err != nil {
		return err
	}

	// Application Stores
	a.loglevelStore = a.runtime.GoRuntime().Store.New("wails:loglevel", a.options.LogLevel)
	a.appconfigStore = a.runtime.GoRuntime().Store.New("wails:appconfig", a.options)

	// Start the logging subsystem
	log, err := subsystem.NewLog(a.servicebus, a.logger, a.loglevelStore)
	if err != nil {
		return err
	}
	a.log = log
	err = a.log.Start()
	if err != nil {
		return err
	}

	// create the dispatcher
	dispatcher, err := messagedispatcher.New(a.servicebus, a.logger)
	if err != nil {
		return err
	}
	a.dispatcher = dispatcher
	err = dispatcher.Start()
	if err != nil {
		return err
	}

	if a.options.Mac.URLHandlers != nil {
		// Start the url handler subsystem
		url, err := subsystem.NewURL(a.servicebus, a.logger, a.options.Mac.URLHandlers)
		if err != nil {
			return err
		}
		a.url = url
		err = a.url.Start()
		if err != nil {
			return err
		}
	}

	// Start the eventing subsystem
	eventsubsystem, err := subsystem.NewEvent(ctx, a.servicebus, a.logger)
	if err != nil {
		return err
	}
	a.event = eventsubsystem
	err = a.event.Start()
	if err != nil {
		return err
	}

	// Start the menu subsystem
	menusubsystem, err := subsystem.NewMenu(ctx, a.servicebus, a.logger, a.menuManager)
	if err != nil {
		return err
	}
	a.menu = menusubsystem
	err = a.menu.Start()
	if err != nil {
		return err
	}

	// Start the call subsystem
	callSubsystem, err := subsystem.NewCall(ctx, a.servicebus, a.logger, a.bindings.DB(), a.runtime.GoRuntime())
	if err != nil {
		return err
	}
	a.call = callSubsystem
	err = a.call.Start()
	if err != nil {
		return err
	}

	// Dump bindings as a debug
	bindingDump, err := a.bindings.ToJSON()
	if err != nil {
		return err
	}

	// Setup signal handler
	signalsubsystem, err := signal.NewManager(ctx, cancel, a.servicebus, a.logger)
	if err != nil {
		return err
	}
	a.signal = signalsubsystem
	a.signal.Start()

	err = a.window.Run(dispatcher, bindingDump, a.debug)
	a.logger.Trace("Ffenestri.Run() exited")
	if err != nil {
		return err
	}

	// Close down all the subsystems
	a.logger.Trace("Cancelling subsystems")
	cancel()
	subsystemWaitGroup.Wait()

	a.logger.Trace("Cancelling dispatcher")
	dispatcher.Close()

	// Close log
	a.logger.Trace("Stopping log")
	log.Close()

	a.logger.Trace("Stopping Service bus")
	err = a.servicebus.Stop()
	if err != nil {
		return err
	}

	// Shutdown callback
	if a.shutdownCallback != nil {
		a.shutdownCallback()
	}

	return nil
}
