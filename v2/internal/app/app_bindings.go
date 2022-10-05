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

	if projectConfig.WailsJSDir == "" {
		projectConfig.WailsJSDir = filepath.Join(cwd, "frontend")
	}
	wrapperDir := filepath.Join(projectConfig.WailsJSDir, "wailsjs", "runtime")
	_ = os.RemoveAll(wrapperDir)
	extractor := gosod.New(wrapper.RuntimeWrapper)
	err = extractor.Extract(wrapperDir, nil)
	if err != nil {
		return err
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

	wailsJSDir := filepath.Join(projectConfig.WailsJSDir, "wailsjs")
	return fs.SetPermissions(wailsJSDir, 0755)
}
