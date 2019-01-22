package wails

import (
	"github.com/wailsapp/wails/cmd"
	"github.com/wailsapp/wails/cmd/frameworks"
)

// -------------------------------- Compile time Flags ------------------------------

// DebugMode indicates if we are in debug Mode
var DebugMode = "true"

// ----------------------------------------------------------------------------------

// App defines the main application struct
type App struct {
	config         *AppConfig      // The Application configuration object
	cli            *cmd.Cli        // In debug mode, we have a cli
	renderer       Renderer        // The renderer is what we will render the app to
	logLevel       string          // The log level of the app
	headless       bool            // Indicates if the app should be started in headless mode
	ipc            *ipcManager     // Handles the IPC calls
	log            *CustomLogger   // Logger
	bindingManager *bindingManager // Handles binding of Go code to renderer
	eventManager   *eventManager   // Handles all the events
	runtime        *Runtime        // The runtime object for registered structs

	// This is a list of all the JS/CSS that needs injecting
	// It will get injected in order
	jsCache  []string
	cssCache []string
}

// CreateApp creates the application window with the given configuration
// If none given, the defaults are used
func CreateApp(optionalConfig ...*AppConfig) *App {
	var userConfig *AppConfig
	if len(optionalConfig) > 0 {
		userConfig = optionalConfig[0]
	}

	result := &App{
		logLevel:       "info",
		renderer:       &webViewRenderer{},
		ipc:            newIPCManager(),
		bindingManager: newBindingManager(),
		eventManager:   newEventManager(),
		log:            newCustomLogger("App"),
	}

	appconfig, err := newAppConfig(userConfig)
	if err != nil {
		result.log.Fatalf("Cannot use custom HTML: %s", err.Error())
	}
	result.config = appconfig

	// Set up the CLI if not in release mode
	if DebugMode == "true" {
		result.cli = result.setupCli()
	} else {
		// Disable Inspector in release mode
		result.config.DisableInspector = true
	}

	return result
}

// Run the app
func (a *App) Run() error {
	if DebugMode == "true" {
		return a.cli.Run()
	}

	a.logLevel = "error"
	return a.start()
}

func (a *App) start() error {

	// Set the log level
	setLogLevel(a.logLevel)

	// Log starup
	a.log.Info("Starting")

	// Check if we are to run in headless mode
	if a.headless {
		a.renderer = &Headless{}
	}

	// Initialise the renderer
	err := a.renderer.Initialise(a.config, a.ipc, a.eventManager)
	if err != nil {
		return err
	}

	// Start event manager and give it our renderer
	a.eventManager.start(a.renderer)

	// Start the IPC Manager and give it the event manager and binding manager
	a.ipc.start(a.eventManager, a.bindingManager)

	// Create the runtime
	a.runtime = newRuntime(a.eventManager, a.renderer)

	// Start binding manager and give it our renderer
	err = a.bindingManager.start(a.renderer, a.runtime)
	if err != nil {
		return err
	}

	// Inject framework, if specified
	if frameworks.FrameworkToUse != nil {
		a.renderer.InjectFramework(frameworks.FrameworkToUse.JS, frameworks.FrameworkToUse.CSS)
	}

	// Inject CSS
	a.renderer.AddCSSList(a.cssCache)

	// Inject JS
	a.renderer.AddJSList(a.jsCache)

	// Run the renderer
	a.renderer.Run()

	return nil
}

// Bind allows the user to bind the given object
// with the application
func (a *App) Bind(object interface{}) {
	a.bindingManager.bind(object)
}

// AddJS adds a piece of Javascript to a cache that
// gets injected at runtime
func (a *App) AddJS(js string) {
	a.jsCache = append(a.jsCache, js)
}

// AddCSS adds a CSS string to a cache that
// gets injected at runtime
func (a *App) AddCSS(js string) {
	a.cssCache = append(a.cssCache, js)
}
