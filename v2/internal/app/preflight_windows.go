//go:build windows
// +build windows

package app

import (
	"github.com/wailsapp/wails/v2/internal/ffenestri/windows/wv2runtime"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func (a *App) PreflightChecks(options *options.App) error {

	_ = options

	// Process the webview2 runtime situation. We can pass a strategy in via the `webview2` flag for `wails build`.
	// This will determine how wv2runtime.Process will handle a lack of valid runtime.
	installedVersion, err := wv2runtime.Process()
	if installedVersion != nil {
		a.logger.Debug("WebView2 Runtime installed: Name: '%s' Version:'%s' Location:'%s'. Minimum version required: %s.",
			installedVersion.Name, installedVersion.Version, installedVersion.Location, wv2runtime.MinimumRuntimeVersion)
	}
	if err != nil {
		return err
	}

	return nil
}
