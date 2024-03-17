//go:build windows

package wv2installer

import (
	"fmt"

	"github.com/wailsapp/go-webview2/webviewloader"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

const MinimumRuntimeVersion string = "94.0.992.31" // WebView2 SDK 1.0.992.28

type installationStatus int

const (
	needsInstalling installationStatus = iota
	needsUpdating
)

func Process(appoptions *options.App) (string, error) {
	messages := windows.DefaultMessages()
	if appoptions.Windows != nil && appoptions.Windows.Messages != nil {
		messages = appoptions.Windows.Messages
	}

	installStatus := needsInstalling

	// Override version check for manually specified webview path if present
	var webviewPath = ""
	if opts := appoptions.Windows; opts != nil && opts.WebviewBrowserPath != "" {
		webviewPath = opts.WebviewBrowserPath
	}

	installedVersion, err := webviewloader.GetAvailableCoreWebView2BrowserVersionString(webviewPath)
	if err != nil {
		return "", err
	}

	if installedVersion != "" {
		installStatus = needsUpdating
		compareResult, err := webviewloader.CompareBrowserVersions(installedVersion, MinimumRuntimeVersion)
		if err != nil {
			return "", err
		}
		updateRequired := compareResult < 0
		// Installed and does not require updating
		if !updateRequired {
			return installedVersion, nil
		}
	}

	// Force error strategy if webview is manually specified
	if webviewPath != "" {
		return installedVersion, fmt.Errorf(messages.InvalidFixedWebview2)
	}

	return installedVersion, doInstallationStrategy(installStatus, messages)
}
