//go:build dev
// +build dev

package appng

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop"
	"github.com/wailsapp/wails/v2/internal/frontend/devserver"
	"github.com/wailsapp/wails/v2/internal/frontend/dispatcher"
	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/menumanager"
	"github.com/wailsapp/wails/v2/internal/project"
	pkglogger "github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// App defines a Wails application structure
type App struct {
	frontend frontend.Frontend
	logger   *logger.Logger
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
	var assetdirFlag *string
	var devServerURLFlag *string
	var loglevelFlag *string

	assetdir := os.Getenv("assetdir")
	if assetdir == "" {
		assetdirFlag = flag.String("assetdir", "", "Directory to serve assets")
	}
	devServerURL := os.Getenv("devserverurl")
	if devServerURL == "" {
		devServerURLFlag = flag.String("devserverurl", "", "URL of development server")
	}

	loglevel := os.Getenv("loglevel")
	if loglevel == "" {
		loglevelFlag = flag.String("loglevel", "debug", "Loglevel to use - Trace, Debug, Info, Warning, Error")
	}

	// If we weren't given the assetdir in the environment variables
	if assetdir == "" {
		flag.Parse()
		if assetdirFlag != nil {
			assetdir = *assetdirFlag
		}
		if devServerURLFlag != nil {
			devServerURL = *devServerURLFlag
		}
		if loglevelFlag != nil {
			loglevel = *loglevelFlag
		}
	}

	if devServerURL != "" {
		ctx = context.WithValue(ctx, "devserverurl", devServerURL)
	}
	if assetdir != "" {
		ctx = context.WithValue(ctx, "assetdir", assetdir)
	}

	if loglevel != "" {
		level, err := pkglogger.StringToLogLevel(loglevel)
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

	err = generateBindings(appBindings)
	if err != nil {
		return nil, err
	}
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

func generateBindings(bindings *binding.Bindings) error {

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	projectConfig, err := project.Load(cwd)
	if err != nil {
		return err
	}

	targetDir := filepath.Join(projectConfig.WailsJSDir, "wailsjs", "go")
	err = os.RemoveAll(targetDir)
	if err != nil {
		return err
	}
	_ = fs.MkDirs(targetDir)
	modelsFile := filepath.Join(targetDir, "models.ts")
	err = bindings.WriteTS(modelsFile)
	if err != nil {
		return err
	}

	// Write backend method wrappers
	bindingsFilename := filepath.Join(targetDir, "bindings.js")
	err = bindings.GenerateBackendJS(bindingsFilename, false)
	if err != nil {
		return err
	}
	return nil

}
