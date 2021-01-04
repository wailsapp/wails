package wails

import (
	"os"
	"syscall"

	"github.com/syossan27/tebata"
	"github.com/wailsapp/wails/cmd"
	"github.com/wailsapp/wails/lib/binding"
	"github.com/wailsapp/wails/lib/event"
	"github.com/wailsapp/wails/lib/interfaces"
	"github.com/wailsapp/wails/lib/ipc"
	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/lib/renderer"
	wailsruntime "github.com/wailsapp/wails/runtime"
)

// -------------------------------- Compile time Flags ------------------------------

// BuildMode indicates what mode we are in
var BuildMode = cmd.BuildModeProd

// Runtime is the Go Runtime struct
type Runtime = wailsruntime.Runtime

// Store is a state store used for syncing with
// the front end
type Store = wailsruntime.Store

// CustomLogger is a specialised logger
type CustomLogger = logger.CustomLogger

// ----------------------------------------------------------------------------------

// App defines the main application struct
type App struct {
	config         *AppConfig                // The Application configuration object
	cli            *cmd.Cli                  // In debug mode, we have a cli
	renderer       interfaces.Renderer       // The renderer is what we will render the app to
	logLevel       string                    // The log level of the app
	ipc            interfaces.IPCManager     // Handles the IPC calls
	log            *logger.CustomLogger      // Logger
	bindingManager interfaces.BindingManager // Handles binding of Go code to renderer
	eventManager   interfaces.EventManager   // Handles all the events
	runtime        interfaces.Runtime        // The runtime object for registered structs
}

// CreateApp creates the application window with the given configuration
// If none given, the defaults are used
func CreateApp(optionalConfig ...*AppConfig) *App {
	var userConfig *AppConfig
	if len(optionalConfig) > 0 {
		userConfig = optionalConfig[0]
	}

	result := &App{
		logLevel:       "debug",
		renderer:       renderer.NewWebView(),
		ipc:            ipc.NewManager(),
		bindingManager: binding.NewManager(),
		eventManager:   event.NewManager(),
		log:            logger.NewCustomLogger("App"),
	}

	appconfig, err := newConfig(userConfig)
	if err != nil {
		result.log.Fatalf("Cannot use custom HTML: %s", err.Error())
	}
	result.config = appconfig

	// Set up the CLI if not in release mode
	if BuildMode != cmd.BuildModeProd {
		result.cli = result.setupCli()
	} else {
		// Disable Inspector in release mode
		result.config.DisableInspector = true
	}

	// Platform specific init
	platformInit()

	return result
}

// Run the app
func (a *App) Run() error {

	if BuildMode != cmd.BuildModeProd {
		return a.cli.Run()
	}

	a.logLevel = "error"
	err := a.start()
	if err != nil {
		a.log.Error(err.Error())
	}
	return err
}

func (a *App) start() error {

	// Set the log level
	logger.SetLogLevel(a.logLevel)

	// Log starup
	a.log.Info("Starting")

	// Check if we are to run in bridge mode
	if BuildMode == cmd.BuildModeBridge {
		a.renderer = renderer.NewBridge()
	}

	// Initialise the renderer
	err := a.renderer.Initialise(a.config, a.ipc, a.eventManager)
	if err != nil {
		return err
	}

	// Start signal handler
	t := tebata.New(os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	t.Reserve(func() {
		a.log.Debug("SIGNAL CAUGHT! Starting Shutdown")
		a.renderer.Close()
	})

	// Start event manager and give it our renderer
	a.eventManager.Start(a.renderer)

	// Start the IPC Manager and give it the event manager and binding manager
	a.ipc.Start(a.eventManager, a.bindingManager)

	// Create the runtime
	a.runtime = wailsruntime.NewRuntime(a.eventManager, a.renderer)

	// Start binding manager and give it our renderer
	err = a.bindingManager.Start(a.renderer, a.runtime)
	if err != nil {
		return err
	}

	// Defer the shutdown
	defer a.shutdown()

	// Run the renderer
	err = a.renderer.Run()
	if err != nil {
		return err
	}

	return nil
}

// shutdown the app
func (a *App) shutdown() {
	// Make sure this is only called once
	a.log.Debug("Shutting down")

	// Shutdown Binding Manager
	a.bindingManager.Shutdown()

	// Shutdown IPC Manager
	a.ipc.Shutdown()

	// Shutdown Event Manager
	a.eventManager.Shutdown()

	a.log.Debug("Cleanly Shutdown")
}

// Bind allows the user to bind the given object
// with the application
func (a *App) Bind(object interface{}) {
	a.bindingManager.Bind(object)
}
