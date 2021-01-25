// +build server,!desktop

package app

import (
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/options"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher"
	"github.com/wailsapp/wails/v2/internal/runtime"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/internal/subsystem"
	"github.com/wailsapp/wails/v2/internal/webserver"
)

// App defines a Wails application structure
type App struct {
	appType string

	binding *subsystem.Binding
	call    *subsystem.Call
	event   *subsystem.Event
	log     *subsystem.Log
	runtime *subsystem.Runtime

	options *options.App

	bindings   *binding.Bindings
	logger     *logger.Logger
	dispatcher *messagedispatcher.Dispatcher
	servicebus *servicebus.ServiceBus
	webserver  *webserver.WebServer

	debug bool

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

	result := &App{
		appType:          "server",
		bindings:         binding.NewBindings(myLogger, options.Bind),
		logger:           myLogger,
		servicebus:       servicebus.New(myLogger),
		webserver:        webserver.NewWebServer(myLogger),
		startupCallback:  appoptions.Startup,
		shutdownCallback: appoptions.Shutdown,
	}

	// Initialise app
	result.Init()

	return result, nil
}

// Run the application
func (a *App) Run() error {

	// Default app options
	var port = 8080
	var ip = "localhost"
	var supressLogging = false
	var debugMode = false

	// Create CLI
	cli := clir.NewCli(filepath.Base(os.Args[0]), "Server Build", "")

	// Setup flags
	cli.IntFlag("p", "Port to serve on", &port)
	cli.StringFlag("i", "IP to serve on", &ip)
	cli.BoolFlag("d", "Debug mode", &debugMode)
	cli.BoolFlag("q", "Supress logging", &supressLogging)

	// Setup main action
	cli.Action(func() error {

		// Set IP + Port
		a.webserver.SetPort(port)
		a.webserver.SetIP(ip)
		a.webserver.SetBindings(a.bindings)
		// Log information (if we aren't supressing it)
		if !supressLogging {
			cli.PrintBanner()
			a.logger.Info("Running server at %s", a.webserver.URL())
		}

		if debugMode {
			a.servicebus.Debug()
		}

		// Start the runtime
		runtime, err := subsystem.NewRuntime(a.servicebus, a.logger, a.startupCallback, a.shutdownCallback)
		if err != nil {
			return err
		}
		a.runtime = runtime
		a.runtime.Start()

		// Application Stores
		a.loglevelStore = a.runtime.GoRuntime().Store.New("wails:loglevel", a.options.LogLevel)
		a.appconfigStore = a.runtime.GoRuntime().Store.New("wails:appconfig", a.options)

		a.servicebus.Start()
		log, err := subsystem.NewLog(a.servicebus, a.logger, a.loglevelStore)
		if err != nil {
			return err
		}
		a.log = log
		a.log.Start()
		dispatcher, err := messagedispatcher.New(a.servicebus, a.logger)
		if err != nil {
			return err
		}
		a.dispatcher = dispatcher
		a.dispatcher.Start()

		// Start the binding subsystem
		binding, err := subsystem.NewBinding(a.servicebus, a.logger, a.bindings, runtime.GoRuntime())
		if err != nil {
			return err
		}
		a.binding = binding
		a.binding.Start()

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

		// Required so that the WailsInit functions are fired!
		runtime.GoRuntime().Events.Emit("wails:loaded")

		if err := a.webserver.Start(dispatcher, event); err != nil {
			a.logger.Error("Webserver failed to start %s", err)
			return err
		}

		return nil
	})

	return cli.Run()
}
