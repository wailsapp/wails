//go:build bindings
// +build bindings

package appng

import (
	"github.com/leaanthony/gosod"
	"github.com/wailsapp/wails/v2/internal/binding"
	wailsRuntime "github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/internal/frontend/runtime/wrapper"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/project"
	"github.com/wailsapp/wails/v2/pkg/options"
	"os"
	"path/filepath"
)

// App defines a Wails application structure
type App struct {
	logger     *logger.Logger
	appoptions *options.App
}

func (a *App) Run() error {

	// Create binding exemptions - Ugly hack. There must be a better way
	bindingExemptions := []interface{}{a.appoptions.OnStartup, a.appoptions.OnShutdown, a.appoptions.OnDomReady}
	appBindings := binding.NewBindings(a.logger, a.appoptions.Bind, bindingExemptions)

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
		logger:     myLogger,
		appoptions: appoptions,
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

	//ipcdev.js
	err = os.WriteFile(filepath.Join(wrapperDir, "ipcdev.js"), wailsRuntime.DesktopIPC, 0755)
	if err != nil {
		return err
	}
	//runtimedev.js
	err = os.WriteFile(filepath.Join(wrapperDir, "runtimedev.js"), wailsRuntime.RuntimeDesktopJS, 0755)
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
	err = bindings.GenerateBackendJS(bindingsFilename, true)
	if err != nil {
		return err
	}

	bindingsTypes := filepath.Join(targetDir, "bindings.d.ts")
	err = bindings.GenerateBackendTS(bindingsTypes)
	if err != nil {
		return err
	}

	return nil

}
