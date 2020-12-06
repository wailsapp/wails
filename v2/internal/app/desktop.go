// +build desktop,!server

package app

import (
	"fmt"
	goruntime "runtime"

	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/ffenestri"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher"
	"github.com/wailsapp/wails/v2/internal/runtime"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/internal/signal"
	"github.com/wailsapp/wails/v2/internal/subsystem"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// App defines a Wails application structure
type App struct {
	window     *ffenestri.Application
	servicebus *servicebus.ServiceBus
	logger     *logger.Logger
	signal     *signal.Manager
	options    *options.App

	// Subsystems
	log        *subsystem.Log
	runtime    *subsystem.Runtime
	event      *subsystem.Event
	binding    *subsystem.Binding
	call       *subsystem.Call
	menu       *subsystem.Menu
	tray       *subsystem.Tray
	dispatcher *messagedispatcher.Dispatcher

	// Indicates if the app is in debug mode
	debug bool

	// This is our binding DB
	bindings *binding.Bindings

	// Application Stores
	loglevelStore  *runtime.Store
	appconfigStore *runtime.Store
}

// Create App
func CreateApp(options *options.App) (*App, error) {

	// Merge default options
	options.MergeDefaults()

	// Set up logger
	myLogger := logger.New(options.Logger)
	myLogger.SetLogLevel(options.LogLevel)

	window := ffenestri.NewApplicationWithConfig(options, myLogger)

	result := &App{
		window:     window,
		servicebus: servicebus.New(myLogger),
		logger:     myLogger,
		bindings:   binding.NewBindings(myLogger),
	}

	result.options = options

	// Initialise the app
	err := result.Init()

	return result, err

}

// Run the application
func (a *App) Run() error {

	var err error

	// Setup signal handler
	signalsubsystem, err := signal.NewManager(a.servicebus, a.logger)
	if err != nil {
		return err
	}
	a.signal = signalsubsystem
	a.signal.Start()

	// Start the service bus
	a.servicebus.Debug()
	err = a.servicebus.Start()
	if err != nil {
		return err
	}

	// Start the runtime
	runtimesubsystem, err := subsystem.NewRuntime(a.servicebus, a.logger,
		a.options.Mac.Menu)
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

	// Start the binding subsystem
	bindingsubsystem, err := subsystem.NewBinding(a.servicebus, a.logger,
		a.bindings, a.runtime.GoRuntime())
	if err != nil {
		return err
	}
	a.binding = bindingsubsystem
	err = a.binding.Start()
	if err != nil {
		return err
	}

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

	// Start the eventing subsystem
	event, err := subsystem.NewEvent(a.servicebus, a.logger)
	if err != nil {
		return err
	}
	a.event = event
	err = a.event.Start()
	if err != nil {
		return err
	}

	// Start the menu subsystem
	var applicationMenu *menu.Menu
	var trayMenu *menu.Menu
	switch goruntime.GOOS {
	case "darwin":
		applicationMenu = a.options.Mac.Menu
		trayMenu = a.options.Mac.Tray
	// case "linux":
	// 	applicationMenu = a.options.Linux.Menu
	// case "windows":
	// 	applicationMenu = a.options.Windows.Menu
	default:
		return fmt.Errorf("unsupported OS: %s", goruntime.GOOS)
	}

	// Optionally start the menu subsystem
	if applicationMenu != nil {
		menusubsystem, err := subsystem.NewMenu(applicationMenu, a.servicebus,
			a.logger)
		if err != nil {
			return err
		}
		a.menu = menusubsystem
		err = a.menu.Start()
		if err != nil {
			return err
		}
	}

	// Optionally start the tray subsystem
	if trayMenu != nil {
		traysubsystem, err := subsystem.NewTray(trayMenu, a.servicebus,
			a.logger)
		if err != nil {
			return err
		}
		a.tray = traysubsystem
		err = a.tray.Start()
		if err != nil {
			return err
		}
	}

	// Start the call subsystem
	call, err := subsystem.NewCall(a.servicebus, a.logger, a.bindings.DB(), a.runtime.GoRuntime())
	if err != nil {
		return err
	}
	a.call = call
	err = a.call.Start()
	if err != nil {
		return err
	}

	// Dump bindings as a debug
	bindingDump, err := a.bindings.ToJSON()
	if err != nil {
		return err
	}

	result := a.window.Run(dispatcher, bindingDump, a.debug)
	a.logger.Trace("Ffenestri.Run() exited")
	err = a.servicebus.Stop()
	if err != nil {
		return err
	}

	return result
}

// Bind a struct to the application by passing in
// a pointer to it
func (a *App) Bind(structPtr interface{}) {

	// Add the struct to the bindings
	err := a.bindings.Add(structPtr)
	if err != nil {
		a.logger.Fatal("Error during binding: " + err.Error())
	}
}
