//go:build bindings

package app

import (
	"flag"
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

	// Check for CLI Flags
	bindingFlags := flag.NewFlagSet("bindings", flag.ContinueOnError)

	var tsPrefixFlag *string
	var tsPostfixFlag *string
	var tsOutputTypeFlag *string

	tsPrefix := os.Getenv("tsprefix")
	if tsPrefix == "" {
		tsPrefixFlag = bindingFlags.String("tsprefix", "", "Prefix for generated typescript entities")
	}

	tsSuffix := os.Getenv("tssuffix")
	if tsSuffix == "" {
		tsPostfixFlag = bindingFlags.String("tssuffix", "", "Suffix for generated typescript entities")
	}

	tsOutputType := os.Getenv("tsoutputtype")
	if tsOutputType == "" {
		tsOutputTypeFlag = bindingFlags.String("tsoutputtype", "", "Output type for generated typescript entities (classes|interfaces)")
	}

	_ = bindingFlags.Parse(os.Args[1:])
	if tsPrefixFlag != nil {
		tsPrefix = *tsPrefixFlag
	}
	if tsPostfixFlag != nil {
		tsSuffix = *tsPostfixFlag
	}
	if tsOutputTypeFlag != nil {
		tsOutputType = *tsOutputTypeFlag
	}

	appBindings := binding.NewBindings(a.logger, a.options.Bind, bindingExemptions, IsObfuscated(), a.options.EnumBind)

	appBindings.SetTsPrefix(tsPrefix)
	appBindings.SetTsSuffix(tsSuffix)
	appBindings.SetOutputType(tsOutputType)

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
