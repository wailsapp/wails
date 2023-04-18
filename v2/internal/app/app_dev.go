//go:build dev

package app

import (
	"context"
	"embed"
	"flag"
	"fmt"
	iofs "io/fs"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/assetserver"

	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop"
	"github.com/wailsapp/wails/v2/internal/frontend/devserver"
	"github.com/wailsapp/wails/v2/internal/frontend/dispatcher"
	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/menumanager"
	pkglogger "github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func (a *App) Run() error {
	err := a.frontend.Run(a.ctx)
	a.frontend.RunMainLoop()
	a.frontend.WindowClose()
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
	devFlags := flag.NewFlagSet("dev", flag.ContinueOnError)

	var assetdirFlag *string
	var devServerFlag *string
	var frontendDevServerURLFlag *string
	var loglevelFlag *string

	assetdir := os.Getenv("assetdir")
	if assetdir == "" {
		assetdirFlag = devFlags.String("assetdir", "", "Directory to serve assets")
	}

	devServer := os.Getenv("devserver")
	if devServer == "" {
		devServerFlag = devFlags.String("devserver", "", "Address to bind the wails dev server to")
	}

	frontendDevServerURL := os.Getenv("frontenddevserverurl")
	if frontendDevServerURL == "" {
		frontendDevServerURLFlag = devFlags.String("frontenddevserverurl", "", "URL of the external frontend dev server")
	}

	loglevel := os.Getenv("loglevel")
	if loglevel == "" {
		loglevelFlag = devFlags.String("loglevel", "debug", "Loglevel to use - Trace, Debug, Info, Warning, Error")
	}

	// If we weren't given the assetdir in the environment variables
	if assetdir == "" {
		// Parse args but ignore errors in case -appargs was used to pass in args for the app.
		_ = devFlags.Parse(os.Args[1:])
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

	assetConfig, err := assetserver.BuildAssetServerConfig(appoptions)
	if err != nil {
		return nil, err
	}

	if assetConfig.Assets == nil && frontendDevServerURL != "" {
		myLogger.Warning("No AssetServer.Assets has been defined but a frontend DevServer, the frontend DevServer will not be used.")
		frontendDevServerURL = ""
		assetdir = ""
	}

	if frontendDevServerURL != "" {
		if os.Getenv("legacyusedevsererinsteadofcustomscheme") != "" {
			startURL, err := url.Parse("http://" + devServer)
			if err != nil {
				return nil, err
			}

			ctx = context.WithValue(ctx, "starturl", startURL)
		}

		ctx = context.WithValue(ctx, "frontenddevserverurl", frontendDevServerURL)

		externalURL, err := url.Parse(frontendDevServerURL)
		if err != nil {
			return nil, err
		}

		if externalURL.Host == "" {
			return nil, fmt.Errorf("Invalid frontend:dev:serverUrl missing protocol scheme?")
		}

		waitCb := func() { myLogger.Debug("Waiting for frontend DevServer '%s' to be ready", externalURL) }
		if !checkPortIsOpen(externalURL.Host, time.Minute, waitCb) {
			myLogger.Error("Timeout waiting for frontend DevServer")
		}

		handler := assetserver.NewExternalAssetsHandler(myLogger, assetConfig, externalURL)
		assetConfig.Assets = nil
		assetConfig.Handler = handler
		assetConfig.Middleware = nil

		myLogger.Info("Serving assets from frontend DevServer URL: %s", frontendDevServerURL)
	} else {
		if assetdir == "" {
			// If no assetdir has been defined, let's try to infer it from the project root and the asset FS.
			assetdir, err = tryInferAssetDirFromFS(assetConfig.Assets)
			if err != nil {
				return nil, fmt.Errorf("unable to infer the AssetDir from your Assets fs.FS: %w", err)
			}
		}

		if assetdir != "" {
			// Let's override the assets to serve from on disk, if needed
			absdir, err := filepath.Abs(assetdir)
			if err != nil {
				return nil, err
			}

			myLogger.Info("Serving assets from disk: %s", absdir)
			assetConfig.Assets = os.DirFS(absdir)

			ctx = context.WithValue(ctx, "assetdir", assetdir)
		}
	}

	// Migrate deprecated options to the new AssetServer option
	appoptions.Assets = nil
	appoptions.AssetsHandler = nil
	appoptions.AssetServer = &assetConfig

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
	bindingExemptions := []interface{}{
		appoptions.OnStartup,
		appoptions.OnShutdown,
		appoptions.OnDomReady,
		appoptions.OnBeforeClose,
	}
	appBindings := binding.NewBindings(myLogger, appoptions.Bind, bindingExemptions, false)

	eventHandler := runtime.NewEvents(myLogger)
	ctx = context.WithValue(ctx, "events", eventHandler)
	messageDispatcher := dispatcher.NewDispatcher(ctx, myLogger, appBindings, eventHandler)

	// Create the frontends and register to event handler
	desktopFrontend := desktop.NewFrontend(ctx, appoptions, myLogger, appBindings, messageDispatcher)
	appFrontend := devserver.NewFrontend(ctx, appoptions, myLogger, appBindings, messageDispatcher, menuManager, desktopFrontend)
	eventHandler.AddFrontend(appFrontend)
	eventHandler.AddFrontend(desktopFrontend)

	ctx = context.WithValue(ctx, "frontend", appFrontend)
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

func checkPortIsOpen(host string, timeout time.Duration, waitCB func()) (ret bool) {
	if timeout == 0 {
		timeout = time.Minute
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, _ := net.DialTimeout("tcp", host, 2*time.Second)
		if conn != nil {
			conn.Close()
			return true
		}

		waitCB()
		time.Sleep(1 * time.Second)
	}
	return false
}
