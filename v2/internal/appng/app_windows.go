package appng

import (
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/ffenestri/windows/wv2runtime"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/windows"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func NewFrontend(appoptions *options.App, myLogger *logger.Logger, bindings *binding.Bindings, dispatcher frontend.Dispatcher) frontend.Frontend {
	return windows.NewFrontend(appoptions, myLogger, bindings, dispatcher)
}

func PreflightChecks(options *options.App, logger *logger.Logger) error {

	_ = options

	// Process the webview2 runtime situation. We can pass a strategy in via the `webview2` flag for `wails build`.
	// This will determine how wv2runtime.Process will handle a lack of valid runtime.
	installedVersion, err := wv2runtime.Process()
	if installedVersion != nil {
		logger.Debug("WebView2 Runtime installed: Name: '%s' Version:'%s' Location:'%s'. Minimum version required: %s.",
			installedVersion.Name, installedVersion.Version, installedVersion.Location, wv2runtime.MinimumRuntimeVersion)
	}
	if err != nil {
		return err
	}

	return nil
}
