package appng

import (
	"context"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/dispatcher"
	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
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

	// Startup/Shutdown
	startupCallback  func(ctx context.Context)
	shutdownCallback func()
	ctx              context.Context
}

func (a *App) Run() error {
	return a.frontend.Run(a.ctx)
}

// CreateApp creates the app!
func CreateApp(appoptions *options.App) (*App, error) {

	ctx := context.Background()

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

	// Create binding exemptions - Ugly hack. There must be a better way
	bindingExemptions := []interface{}{appoptions.Startup, appoptions.Shutdown}
	appBindings := binding.NewBindings(myLogger, appoptions.Bind, bindingExemptions)
	eventHandler := runtime.NewEvents(myLogger)
	ctx = context.WithValue(ctx, "events", eventHandler)
	messageDispatcher := dispatcher.NewDispatcher(myLogger, appBindings, eventHandler)

	appFrontend := NewFrontend(appoptions, myLogger, appBindings, messageDispatcher)
	eventHandler.SetFrontend(appFrontend)

	result := &App{
		ctx:              ctx,
		frontend:         appFrontend,
		logger:           myLogger,
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
