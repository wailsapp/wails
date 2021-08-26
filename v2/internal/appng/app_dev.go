//go:build dev
// +build dev

package appng

import (
	"context"
	"flag"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop"
	"github.com/wailsapp/wails/v2/internal/frontend/devserver"
	"github.com/wailsapp/wails/v2/internal/frontend/dispatcher"
	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/menumanager"
	"github.com/wailsapp/wails/v2/internal/signal"
	pkglogger "github.com/wailsapp/wails/v2/pkg/logger"
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

	// OnStartup/OnShutdown
	startupCallback  func(ctx context.Context)
	shutdownCallback func(ctx context.Context)
	ctx              context.Context
}

func (a *App) Run() error {
	err := a.frontend.Run(a.ctx)
	if a.shutdownCallback != nil {
		a.shutdownCallback(a.ctx)
	}
	return err
}

// CreateApp creates the app!
func CreateApp(appoptions *options.App) (*App, error) {
	var err error

	ctx := context.WithValue(context.Background(), "debug", true)

	// Set up logger
	myLogger := logger.New(appoptions.Logger)
	myLogger.SetLogLevel(appoptions.LogLevel)

	// Check for CLI Flags
	assetdir := flag.String("assetdir", "", "Directory to serve assets")
	loglevel := flag.String("loglevel", "debug", "Loglevel to use - Trace, Debug, Info, Warning, Error")
	flag.Parse()
	if assetdir != nil && *assetdir != "" {
		ctx = context.WithValue(ctx, "assetdir", *assetdir)
	}
	if loglevel != nil && *loglevel != "" {
		level, err := pkglogger.StringToLogLevel(*loglevel)
		if err != nil {
			return nil, err
		}
		myLogger.SetLogLevel(level)
	}

	// Attach logger to context
	ctx = context.WithValue(ctx, "logger", myLogger)

	// Preflight checks
	err = PreflightChecks(appoptions, myLogger)
	if err != nil {
		return nil, err
	}

	// Merge default options
	options.MergeDefaults(appoptions)

	var menuManager *menumanager.Manager

	// Process the application menu
	if appoptions.Menu != nil {
		// Create the menu manager
		menuManager = menumanager.NewManager()
		err = menuManager.SetApplicationMenu(appoptions.Menu)
		if err != nil {
			return nil, err
		}
	}

	// Create binding exemptions - Ugly hack. There must be a better way
	bindingExemptions := []interface{}{appoptions.OnStartup, appoptions.OnShutdown, appoptions.OnDomReady}
	appBindings := binding.NewBindings(myLogger, appoptions.Bind, bindingExemptions)
	eventHandler := runtime.NewEvents(myLogger)
	ctx = context.WithValue(ctx, "events", eventHandler)
	messageDispatcher := dispatcher.NewDispatcher(myLogger, appBindings, eventHandler)

	// Create the frontends and register to event handler
	desktopFrontend := desktop.NewFrontend(ctx, appoptions, myLogger, appBindings, messageDispatcher)
	appFrontend := devserver.NewFrontend(ctx, appoptions, myLogger, appBindings, messageDispatcher, menuManager, desktopFrontend)
	eventHandler.AddFrontend(appFrontend)
	eventHandler.AddFrontend(desktopFrontend)

	result := &App{
		ctx:              ctx,
		frontend:         appFrontend,
		logger:           myLogger,
		menuManager:      menuManager,
		startupCallback:  appoptions.OnStartup,
		shutdownCallback: appoptions.OnShutdown,
		debug:            true,
	}

	result.options = appoptions

	return result, nil

}
