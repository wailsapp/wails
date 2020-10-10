// +build desktop,!server

package app

import (
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/ffenestri"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/internal/signal"
	"github.com/wailsapp/wails/v2/internal/subsystem"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	dispatcher *messagedispatcher.Dispatcher

	// Indicates if the app is in debug mode
	debug bool

	// This is our binding DB
	bindings *binding.Bindings

	// LogLevel Store
	loglevelStore *runtime.Store
}

// Create App
func CreateApp(options *options.App) *App {

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
	result.Init()

	return result

}

// Run the application
func (a *App) Run() error {

	// Setup signal handler
	signal, err := signal.NewManager(a.servicebus, a.logger)
	if err != nil {
		return err
	}
	a.signal = signal
	a.signal.Start()

	// Start the service bus
	a.servicebus.Debug()
	a.servicebus.Start()

	// Start the runtime
	runtime, err := subsystem.NewRuntime(a.servicebus, a.logger)
	if err != nil {
		return err
	}
	a.runtime = runtime
	a.runtime.Start()

	// Application Stores
	a.loglevelStore = a.runtime.GoRuntime().Store.New("loglevel", a.options.LogLevel)

	// Start the binding subsystem
	binding, err := subsystem.NewBinding(a.servicebus, a.logger, a.bindings, a.runtime.GoRuntime())
	if err != nil {
		return err
	}
	a.binding = binding
	a.binding.Start()

	// Start the logging subsystem
	log, err := subsystem.NewLog(a.servicebus, a.logger, a.loglevelStore)
	if err != nil {
		return err
	}
	a.log = log
	a.log.Start()

	// create the dispatcher
	dispatcher, err := messagedispatcher.New(a.servicebus, a.logger)
	if err != nil {
		return err
	}
	a.dispatcher = dispatcher
	dispatcher.Start()

	// Start the eventing subsystem
	event, err := subsystem.NewEvent(a.servicebus, a.logger)
	if err != nil {
		return err
	}
	a.event = event
	a.event.Start()

	// Start the call subsystem
	call, err := subsystem.NewCall(a.servicebus, a.logger, a.bindings.DB(), a.runtime.GoRuntime())
	if err != nil {
		return err
	}
	a.call = call
	a.call.Start()

	// Dump bindings as a debug
	bindingDump, err := a.bindings.ToJSON()
	if err != nil {
		return err
	}

	result := a.window.Run(dispatcher, bindingDump)
	a.logger.Trace("Ffenestri.Run() exited")
	a.servicebus.Stop()

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
