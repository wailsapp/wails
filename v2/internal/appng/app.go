package appng

import (
	"context"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/menumanager"
	"github.com/wailsapp/wails/v2/internal/signal"
	"github.com/wailsapp/wails/v2/pkg/options"
	"time"
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
		time.Sleep(1 * time.Second)
		println("fullscreen")
		a.frontend.WindowFullscreen()
		time.Sleep(1 * time.Second)
		println("unfullscreen")
		a.frontend.WindowUnFullscreen()
		time.Sleep(1 * time.Second)
		println("hide")
		a.frontend.WindowHide()
		time.Sleep(1 * time.Second)
		println("show")
		a.frontend.WindowShow()
		time.Sleep(1 * time.Second)
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
