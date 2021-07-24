package appng

import (
	"context"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/menumanager"
	"github.com/wailsapp/wails/v2/internal/signal"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// App defines a Wails application structure
type App struct {
	frontend frontend.Frontend
	logger   *logger.Logger
	signal   *signal.Manager
	options  *options.App

	menuManager *menumanager.Manager

	// Indicates if the app is in debug mode
	debug bool

	// This is our binding DB
	bindings *binding.Bindings

	// Startup/Shutdown
	startupCallback  func(ctx context.Context)
	shutdownCallback func()
}

func (a *App) Run() error {

	go func() {
		//time.Sleep(3 * time.Second)
		//println("WindowSetSize(3000,2000)")
		//a.frontend.WindowSetSize(3000,2000)
		//x,y := a.frontend.WindowGetSize()
		//println("X", x, "Y", y)
		//time.Sleep(3 * time.Second)
		//println("a.frontend.WindowSetSize(10,10)")
		//a.frontend.WindowSetSize(10,10)
		//x,y = a.frontend.WindowGetSize()
		//println("X", x, "Y", y)
		//time.Sleep(3 * time.Second)
		//time.Sleep(3 * time.Second)
		//println("WindowSetMaxSize(50,50)")
		//a.frontend.WindowSetMaxSize(200,200)
		//x,y := a.frontend.WindowGetSize()
		//println("X", x, "Y", y)
		//time.Sleep(3 * time.Second)
		//println("WindowSetMinSize(100,100)")
		//a.frontend.WindowSetMinSize(600,600)
		//x,y = a.frontend.WindowGetSize()
		//println("X", x, "Y", y)
		//println("fullscreen")
		//a.frontend.WindowFullscreen()
		//time.Sleep(1 * time.Second)
		//println("unfullscreen")
		//a.frontend.WindowUnFullscreen()
		//time.Sleep(1 * time.Second)
		//println("hide")
		//a.frontend.WindowHide()
		//time.Sleep(1 * time.Second)
		//println("show")
		//a.frontend.WindowShow()
		//time.Sleep(1 * time.Second)
		//println("title 1")
		//a.frontend.WindowSetTitle("title 1")
		//time.Sleep(1 * time.Second)
		//println("title 2")
		//a.frontend.WindowSetTitle("title 2")
		//time.Sleep(1 * time.Second)
		//println("setsize 1")
		//a.frontend.WindowSetSize(100,100)
		//time.Sleep(1 * time.Second)
		//println("setsize 2")
		//a.frontend.WindowSetSize(500,500)
		//time.Sleep(1 * time.Second)
		//println("setpos 1")
		//a.frontend.WindowSetPos(0,0)
		//time.Sleep(1 * time.Second)
		//println("setpos 2")
		//a.frontend.WindowSetPos(500,500)
		//time.Sleep(1 * time.Second)
		//println("Center 1")
		//a.frontend.WindowCenter()
		//time.Sleep(5 * time.Second)
		//println("Center 2")
		//a.frontend.WindowCenter()
		//time.Sleep(1 * time.Second)
		//println("maximise")
		//a.frontend.WindowMaximise()
		//time.Sleep(1 * time.Second)
		//println("UnMaximise")
		//a.frontend.WindowUnmaximise()
		//time.Sleep(1 * time.Second)
		//println("minimise")
		//a.frontend.WindowMinimise()
		//time.Sleep(1 * time.Second)
		//println("unminimise")
		//a.frontend.WindowUnminimise()
		//time.Sleep(1 * time.Second)
	}()

	return a.frontend.Run()
}

// CreateApp
func CreateApp(appoptions *options.App) (*App, error) {

	// Merge default options
	options.MergeDefaults(appoptions)

	// Set up logger
	myLogger := logger.New(appoptions.Logger)
	myLogger.SetLogLevel(appoptions.LogLevel)

	// Create the menu manager
	menuManager := menumanager.NewManager()

	// Process the application menu
	appMenu := options.GetApplicationMenu(appoptions)
	err := menuManager.SetApplicationMenu(appMenu)
	if err != nil {
		return nil, err
	}

	// Process context menus
	contextMenus := options.GetContextMenus(appoptions)
	for _, contextMenu := range contextMenus {
		menuManager.AddContextMenu(contextMenu)
	}

	// Process tray menus
	trayMenus := options.GetTrayMenus(appoptions)
	for _, trayMenu := range trayMenus {
		_, err := menuManager.AddTrayMenu(trayMenu)
		if err != nil {
			return nil, err
		}
	}

	//window := ffenestri.NewApplicationWithConfig(appoptions, myLogger, menuManager)
	appFrontend := frontend.NewFrontend(appoptions, myLogger)

	// Create binding exemptions - Ugly hack. There must be a better way
	bindingExemptions := []interface{}{appoptions.Startup, appoptions.Shutdown}

	result := &App{
		frontend:         appFrontend,
		logger:           myLogger,
		bindings:         binding.NewBindings(myLogger, appoptions.Bind, bindingExemptions),
		menuManager:      menuManager,
		startupCallback:  appoptions.Startup,
		shutdownCallback: appoptions.Shutdown,
	}

	result.options = appoptions

	//// Initialise the app
	//err := result.Init()
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Preflight Checks
	//err = result.PreflightChecks(appoptions)
	//if err != nil {
	//	return nil, err
	//}

	return result, nil

}
