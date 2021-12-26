//go:build !wv2runtime.error && !wv2runtime.browser && !wv2runtime.embed
// +build !wv2runtime.error,!wv2runtime.browser,!wv2runtime.embed

package wv2runtime

import (
	"fmt"
	"github.com/wailsapp/wails/v2/internal/webview2runtime"
)

func doInstallationStrategy(installStatus installationStatus) error {
	message := "The WebView2 runtime is required. "
	if installStatus == needsUpdating {
		message = "The Webview2 runtime needs updating. "
	}
	message += "Press Ok to download and install. Note: The installer will download silently so please wait."
	confirmed, err := webview2runtime.Confirm(message, "Missing Requirements")
	if err != nil {
		return err
	}
	if !confirmed {
		return fmt.Errorf("webview2 runtime not installed")
	}
	installedCorrectly, err := webview2runtime.InstallUsingBootstrapper()
	if err != nil {
		_ = webview2runtime.Error(err.Error(), "Error")
		return err
	}
	if !installedCorrectly {
		err = webview2runtime.Error("The runtime failed to install correctly. Please try again.", "Error")
		return err
	}
	return nil
}
