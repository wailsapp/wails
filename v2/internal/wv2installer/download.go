//go:build !wv2runtime.error && !wv2runtime.browser && !wv2runtime.embed
// +build !wv2runtime.error,!wv2runtime.browser,!wv2runtime.embed

package wv2installer

import (
	"fmt"
	"github.com/wailsapp/wails/v2/internal/webview2runtime"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

func doInstallationStrategy(installStatus installationStatus, messages *windows.Messages) error {
	message := messages.InstallationRequiredMsg
	if installStatus == needsUpdating {
		message = messages.UpdateRequiredMsg
	}
	confirmed, err := webview2runtime.Confirm(message, messages.MissingRequirementsMsg)
	if err != nil {
		return err
	}
	if !confirmed {
		return fmt.Errorf(messages.Webview2NotInstalledMsg)
	}
	installedCorrectly, err := webview2runtime.InstallUsingBootstrapper()
	if err != nil {
		_ = webview2runtime.Error(err.Error(), messages.ErrorMsg)
		return err
	}
	if !installedCorrectly {
		err = webview2runtime.Error(messages.FailedToInstallMsg, messages.ErrorMsg)
		return err
	}
	return nil
}
