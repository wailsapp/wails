// +build !server,!desktop,hybrid

package app

import (
	"os"
	"path/filepath"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/ffenestri"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/internal/subsystem"
	"github.com/wailsapp/wails/v2/internal/webserver"
)

// Config defines the Application's configuration
type Config struct {
	Title         string // Title is the value to be displayed in the title bar
	Width         int    // Width is the desired window width
	Height        int    // Height is the desired window height
	DevTools      bool   // DevTools enables or disables the browser development tools
	Resizable     bool   // Resizable when False prevents window resizing
	ServerEnabled bool   // ServerEnabled when True allows remote connections
}

// App defines a Wails application structure
type App struct {
	config    Config
	window    *ffenestri.Application
	webserver *webserver.WebServer
	binding   *subsystem.Binding
	call      *subsystem.Call
	event     *subsystem.Event
	log       *subsystem.Log
	runtime   *subsystem.Runtime

	bindings   *binding.Bindings
	logger     *logger.Logger
	dispatcher *messagedispatcher.Dispatcher
	servicebus *servicebus.ServiceBus

	debug bool
}

// Create App
func CreateApp(options *Options) *App {

	// Merge default options
	options.mergeDefaults()

	// Set up logger
	myLogger := logger.New(os.Stdout)
	myLogger.SetLogLevel(logger.INFO)

	window := ffenestri.NewApplicationWithConfig(&ffenestri.Config{
		Title:       options.Title,
		Width:       options.Width,
		Height:      options.Height,
		MinWidth:    options.MinWidth,
		MinHeight:   options.MinHeight,
		MaxWidth:    options.MaxWidth,
		MaxHeight:   options.MaxHeight,
		StartHidden: options.StartHidden,
		DevTools:    options.DevTools,

		Resizable:  !options.DisableResize,
		Fullscreen: options.Fullscreen,
	}, myLogger)

	app := &App{
		window:     window,
		webserver:  webserver.NewWebServer(myLogger),
		servicebus: servicebus.New(myLogger),
		logger:     myLogger,
		bindings:   binding.NewBindings(myLogger, options.Bind),
	}

	// Initialise the app
	app.Init()

	return app
}

// Run the application
func (a *App) Run() error {

	// Default app options
	var port = 8080
	var ip = "localhost"
	var suppressLogging = false

	// Create CLI
	cli := clir.NewCli(filepath.Base(os.Args[0]), "Desktop/Server Build", "")

	// Setup flags
	cli.IntFlag("p", "Port to serve on", &port)
	cli.StringFlag("i", "IP to serve on", &ip)
	cli.BoolFlag("q", "Suppress logging", &suppressLogging)

	// Setup main action
	cli.Action(func() error {

		// Set IP + Port
		a.webserver.SetPort(port)
		a.webserver.SetIP(ip)
		a.webserver.SetBindings(a.bindings)
		// Log information (if we aren't suppressing it)
		if !suppressLogging {
			cli.PrintBanner()
			a.logger.Info("Running server at %s", a.webserver.URL())
		}

		a.servicebus.Start()
		log, err := subsystem.NewLog(a.servicebus, a.logger)
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

		// Start the runtime
		runtime, err := subsystem.NewRuntime(a.servicebus, a.logger)
		if err != nil {
			return err
		}
		a.runtime = runtime
		a.runtime.Start()

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
		call, err := subsystem.NewCall(a.servicebus, a.logger, a.bindings.DB())
		if err != nil {
			return err
		}
		a.call = call
		a.call.Start()

		// Required so that the WailsInit functions are fired!
		runtime.GoRuntime().Events.Emit("wails:loaded")

		// Set IP + Port
		a.webserver.SetPort(port)
		a.webserver.SetIP(ip)

		// Log information (if we aren't suppressing it)
		if !suppressLogging {
			cli.PrintBanner()
			println("Running server at " + a.webserver.URL())
		}

		// Dump bindings as a debug
		bindingDump, err := a.bindings.ToJSON()
		if err != nil {
			return err
		}

		go func() {
			if err := a.webserver.Start(dispatcher, event); err != nil {
				a.logger.Error("Webserver failed to start %s", err)
			}
		}()

		result := a.window.Run(dispatcher, bindingDump)
		a.servicebus.Stop()

		return result
	})

	return cli.Run()
}
