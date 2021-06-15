package wv2runtime

import (
	"github.com/leaanthony/webview2runtime"
)

const minimumRuntimeVersion string = "91.0.864.48"

type installationStatus int

const (
	needsInstalling installationStatus = iota
	needsUpdating
	installed
)

func Process() error {
	installStatus := needsInstalling
	installedVersion := webview2runtime.GetInstalledVersion()
	if installedVersion != nil {
		installStatus = installed
		updateRequired, err := installedVersion.IsOlderThan(minimumRuntimeVersion)
		if err != nil {
			_ = webview2runtime.Error(err.Error(), "Error")
			return err
		}
		// Installed and does not require updating
		if !updateRequired {
			return nil
		}
		installStatus = needsUpdating
	}
	return doInstallationStrategy(installStatus)
}
