package wv2runtime

import (
	"github.com/leaanthony/go-webview2/webviewloader"
)

const MinimumRuntimeVersion string = "91.0.992.28"

type installationStatus int

const (
	needsInstalling installationStatus = iota
	needsUpdating
	installed
)

func Process() (string, error) {
	installStatus := needsInstalling
	installedVersion, err := webviewloader.GetInstalledVersion()
	if err != nil {
		return "", err
	}
	if installedVersion != "" {
		installStatus = installed
		compareResult, err := webviewloader.CompareBrowserVersions(installedVersion, MinimumRuntimeVersion)
		if err != nil {
			return "", err
		}
		updateRequired := compareResult == -1
		// Installed and does not require updating
		if !updateRequired {
			return installedVersion, nil
		}

	}
	return installedVersion, doInstallationStrategy(installStatus)
}
