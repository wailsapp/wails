package wv2runtime

import (
	"github.com/leaanthony/webview2runtime"
)

const MinimumRuntimeVersion string = "91.0.864.48"

type installationStatus int

const (
	needsInstalling installationStatus = iota
	needsUpdating
	installed
)

func Process() (*webview2runtime.Info, error) {
	installStatus := needsInstalling
	installedVersion := webview2runtime.GetInstalledVersion()
	if installedVersion != nil {
		installStatus = installed
		updateRequired, err := installedVersion.IsOlderThan(MinimumRuntimeVersion)
		if err != nil {
			_ = webview2runtime.Error(err.Error(), "Error")
			return installedVersion, err
		}
		// Installed and does not require updating
		if !updateRequired {
			return installedVersion, nil
		}

	}
	return installedVersion, doInstallationStrategy(installStatus)
}
