//go:build dev
// +build dev

package appng

import (
	"context"
	"embed"
	"flag"
	"fmt"
	iofs "io/fs"
	"net/url"
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

func (a *App) Shutdown() {
	if a.shutdownCallback != nil {
		a.shutdownCallback(a.ctx)
	}
	a.frontend.Quit()
}

func (a *App) Run() error {
	err := a.frontend.Run(a.ctx)
	a.Shutdown()
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
	var devServerFlag *string
	var frontendDevServerURLFlag *string
	var loglevelFlag *string

	assetdir := os.Getenv("assetdir")
	if assetdir == "" {
		assetdirFlag = flag.String("assetdir", "", "Directory to serve assets")
	}
	devServer := os.Getenv("devserver")
	if devServer == "" {
		devServerFlag = flag.String("devserver", "", "Address to bind the wails dev server to")
	}
	frontendDevServerURL := os.Getenv("frontenddevserverurl")
	if frontendDevServerURL == "" {
		frontendDevServerURLFlag = flag.String("frontenddevserverurl", "", "URL of the external frontend dev server")
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
		if devServerFlag != nil {
			devServer = *devServerFlag
		}
		if frontendDevServerURLFlag != nil {
			frontendDevServerURL = *frontendDevServerURLFlag
		}
		if loglevelFlag != nil {
			loglevel = *loglevelFlag
		}
	}

	if assetdir == "" {
		// If no assetdir has been defined, let's try to infer it from the project root and the asset FS.
		assetdir, err = tryInferAssetDirFromFS(appoptions.Assets)
		if err != nil {
			return nil, err
		}
	}

	if frontendDevServerURL != "" {
		if devServer == "" {
			return nil, fmt.Errorf("Unable to use FrontendDevServerUrl without a DevServer address")
		}

		startURL, err := url.Parse("http://" + devServer)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, "starturl", startURL)
		ctx = context.WithValue(ctx, "frontenddevserverurl", frontendDevServerURL)

		myLogger.Info("Serving assets from frontend DevServer URL: %s", frontendDevServerURL)
	} else if assetdir != "" {
		// Let's override the assets to serve from on disk, if needed
		absdir, err := filepath.Abs(assetdir)
		if err != nil {
			return nil, err
		}

		myLogger.Info("Serving assets from disk: %s", absdir)
		appoptions.Assets = os.DirFS(absdir)

		ctx = context.WithValue(ctx, "assetdir", assetdir)
	}

	if devServer != "" {
		ctx = context.WithValue(ctx, "devserver", devServer)
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
	ctx = context.WithValue(ctx, "buildtype", "dev")

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
	messageDispatcher := dispatcher.NewDispatcher(ctx, myLogger, appBindings, eventHandler)

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

	if projectConfig.WailsJSDir == "" {
		projectConfig.WailsJSDir = filepath.Join(cwd, "frontend")
	}

	targetDir := filepath.Join(projectConfig.WailsJSDir, "wailsjs", "go")
	err = os.RemoveAll(targetDir)
	if err != nil {
		return err
	}
	_ = fs.MkDirs(targetDir)

	err = bindings.GenerateGoBindings(targetDir)
	if err != nil {
		return err
	}

	return nil

}

func tryInferAssetDirFromFS(assets iofs.FS) (string, error) {
	if _, isEmbedFs := assets.(embed.FS); !isEmbedFs {
		// We only infer the assetdir for embed.FS assets
		return "", nil
	}

	path, err := fs.FindPathToFile(assets, "index.html")
	if err != nil {
		return "", err
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(filepath.Join(path, "index.html")); err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf(
				"inferred assetdir '%s' does not exist or does not contain an 'index.html' file, "+
					"please specify it with -assetdir or set it in wails.json",
				path)
		}
		return "", err
	}

	return path, nil
}
