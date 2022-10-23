//go:build bindings

package app

import (
	"os"
	"path/filepath"

	"github.com/leaanthony/gosod"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend/runtime/wrapper"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/project"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func (a *App) Run() error {

	// Create binding exemptions - Ugly hack. There must be a better way
	bindingExemptions := []interface{}{
		a.options.OnStartup,
		a.options.OnShutdown,
		a.options.OnDomReady,
		a.options.OnBeforeClose,
	}

	appBindings := binding.NewBindings(a.logger, a.options.Bind, bindingExemptions, IsObfuscated())

	err := generateBindings(appBindings)
	if err != nil {
		return err
	}
	return nil
}

// CreateApp creates the app!
func CreateApp(appoptions *options.App) (*App, error) {
	// Set up logger
	myLogger := logger.New(appoptions.Logger)
	myLogger.SetLogLevel(appoptions.LogLevel)

	result := &App{
		logger:  myLogger,
		options: appoptions,
	}

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

	wailsjsbasedir := filepath.Join(projectConfig.GetWailsJSDir(), "wailsjs")

	runtimeDir := filepath.Join(wailsjsbasedir, "runtime")
	_ = os.RemoveAll(runtimeDir)
	extractor := gosod.New(wrapper.RuntimeWrapper)
	err = extractor.Extract(runtimeDir, nil)
	if err != nil {
		return err
	}

	goBindingsDir := filepath.Join(wailsjsbasedir, "go")
	err = os.RemoveAll(goBindingsDir)
	if err != nil {
		return err
	}
	_ = fs.MkDirs(goBindingsDir)

	err = bindings.GenerateGoBindings(goBindingsDir)
	if err != nil {
		return err
	}

	return fs.SetPermissions(wailsjsbasedir, 0755)
}
