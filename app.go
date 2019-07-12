package wails

import (
	"github.com/wailsapp/wails/cmd"
	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/runtime/go/runtime"
	"github.com/wailsapp/wails/lib/renderer"
	"github.com/wailsapp/wails/lib/binding"
	"github.com/wailsapp/wails/lib/ipc"
	"github.com/wailsapp/wails/lib/event"
	"github.com/wailsapp/wails/lib/interfaces"
)

// -------------------------------- Compile time Flags ------------------------------

// BuildMode indicates what mode we are in
var BuildMode = cmd.BuildModeProd

// ----------------------------------------------------------------------------------

// App defines the main application struct
type App struct {
	config         *Config              // The Application configuration object
	cli            *cmd.Cli             // In debug mode, we have a cli
	renderer       interfaces.Renderer    // The renderer is what we will render the app to
	logLevel       string               // The log level of the app
	ipc            interfaces.IPCManager          // Handles the IPC calls
	log            *logger.CustomLogger // Logger
	bindingManager interfaces.BindingManager     // Handles binding of Go code to renderer
	eventManager   interfaces.EventManager        // Handles all the events
	runtime        interfaces.Runtime     // The runtime object for registered structs
}

// CreateApp creates the application window with the given configuration
// If none given, the defaults are used
func CreateApp(optionalConfig ...*Config) *App {
	var userConfig *Config
	if len(optionalConfig) > 0 {
		userConfig = optionalConfig[0]
	}

	result := &App{
		logLevel:       "info",
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

	// Check if we are to run in headless mode
	if BuildMode == cmd.BuildModeBridge {
		a.renderer = &renderer.Headless{}
	}

	// Initialise the renderer
	err := a.renderer.Initialise(a.config, a.ipc, a.eventManager)
	if err != nil {
		return err
	}

	// Start event manager and give it our renderer
	a.eventManager.Start(a.renderer)

	// Start the IPC Manager and give it the event manager and binding manager
	a.ipc.Start(a.eventManager, a.bindingManager)

	// Create the runtime
	a.runtime = runtime.NewRuntime(a.eventManager, a.renderer)

	// Start binding manager and give it our renderer
	err = a.bindingManager.Start(a.renderer, a.runtime)
	if err != nil {
		return err
	}

	// Run the renderer
	return a.renderer.Run()
}

// Bind allows the user to bind the given object
// with the application
func (a *App) Bind(object interface{}) {
	a.bindingManager.Bind(object)
}
