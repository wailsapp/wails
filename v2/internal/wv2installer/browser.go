//go:build windows && wv2runtime.browser
// +build windows,wv2runtime.browser

package wv2installer

import (
	"fmt"
	"github.com/wailsapp/wails/v2/internal/webview2runtime"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

func doInstallationStrategy(installStatus installationStatus, messages *windows.Messages) error {
	confirmed, err := webview2runtime.Confirm(messages.DownloadPage+MinimumRuntimeVersion, messages.MissingRequirements)
	if err != nil {
		return err
	}
	if confirmed {
		err = webview2runtime.OpenInstallerDownloadWebpage()
		if err != nil {
			return err
		}
	}

	return fmt.Errorf(messages.FailedToInstall)
}
