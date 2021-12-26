//go:build wv2runtime.error
// +build wv2runtime.error

package wv2runtime

import (
	"fmt"
	"github.com/wailsapp/wails/v2/internal/webview2runtime"
)

func doInstallationStrategy(installStatus installationStatus) error {
	_ = webview2runtime.Error("The WebView2 runtime is required to run this application. Please contact your system administrator.", "Error")
	return fmt.Errorf("webview2 runtime not installed")
}
