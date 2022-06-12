//go:build wv2runtime.manual
// +build wv2runtime.manual

package wv2installer

func doInstallationStrategy(installStatus installationStatus, messages *windows.Messages) error {
	// fallback for manually specifying webview2
	return nil
}
