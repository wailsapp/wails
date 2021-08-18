//go:build dev

package appng

import (
	"context"
	"flag"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/devserver"
	"github.com/wailsapp/wails/v2/internal/logger"
	pkglogger "github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"os"
)

func NewFrontend(appoptions *options.App, myLogger *logger.Logger, bindings *binding.Bindings, dispatcher frontend.Dispatcher) frontend.Frontend {
	return devserver.NewFrontend(appoptions, myLogger, bindings, dispatcher)
}

func PreflightChecks(options *options.App, logger *logger.Logger) error {
	return nil
}

func (a *App) Init() {
	// Check for CLI Flags
	assetdir := flag.String("assetdir", "", "Directory to serve assets")
	loglevel := flag.String("loglevel", "debug", "Loglevel to use - Trace, Debug, Info, Warning, Error")
	flag.Parse()
	if assetdir != nil && *assetdir != "" {
		a.ctx = context.WithValue(a.ctx, "assetdir", *assetdir)
	}
	if loglevel != nil && *loglevel != "" {
		level, err := pkglogger.StringToLogLevel(*loglevel)
		if err != nil {
			println("ERROR:", err.Error())
			os.Exit(1)
		}
		a.logger.SetLogLevel(level)
	}

}
