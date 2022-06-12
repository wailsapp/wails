//go:build wv2runtime.manual
// +build wv2runtime.manual

package wv2installer

import "github.com/wailsapp/wails/v2/pkg/options/windows"

func doInstallationStrategy(installStatus installationStatus, messages *windows.Messages) error {
	// fallback for manually specifying webview2
	return nil
}
