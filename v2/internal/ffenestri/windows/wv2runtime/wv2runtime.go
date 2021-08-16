package wv2runtime

import (
	"github.com/jchv/go-webview2/webviewloader"
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
		compareResult, err := webviewloader.CompareBrowserVersions(installedVersion.Version, MinimumRuntimeVersion)
		if err != nil {
			return nil, err
		}
		updateRequired := compareResult == -1
		// Installed and does not require updating
		if !updateRequired {
			return installedVersion, nil
		}

	}
	return installedVersion, doInstallationStrategy(installStatus)
}
