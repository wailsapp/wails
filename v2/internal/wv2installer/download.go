//go:build windows && !wv2runtime.error && !wv2runtime.browser && !wv2runtime.embed
// +build windows,!wv2runtime.error,!wv2runtime.browser,!wv2runtime.embed

package wv2installer

import (
	"fmt"

	"github.com/wailsapp/wails/v2/internal/webview2runtime"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

func doInstallationStrategy(installStatus installationStatus, messages *windows.Messages) error {
	message := messages.InstallationRequired
	if installStatus == needsUpdating {
		message = messages.UpdateRequired
	}
	confirmed, err := webview2runtime.Confirm(message, messages.MissingRequirements)
	if err != nil {
		return err
	}
	if !confirmed {
		return fmt.Errorf(messages.Webview2NotInstalled)
	}
	installedCorrectly, err := webview2runtime.InstallUsingBootstrapper()
	if err != nil {
		_ = webview2runtime.Error(err.Error(), messages.Error)
		return err
	}
	if !installedCorrectly {
		err = webview2runtime.Error(messages.FailedToInstall, messages.Error)
		return err
	}
	return nil
}
