//go:build windows
// +build windows

package system

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2/webviewloader"
	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
)

func (i *Info) discover() error {

	var err error
	osinfo, err := operatingsystem.Info()
	if err != nil {
		return err
	}
	i.OS = osinfo

	i.Dependencies = append(i.Dependencies, checkWebView2())
	i.Dependencies = append(i.Dependencies, checkNPM())
	i.Dependencies = append(i.Dependencies, checkUPX())
	i.Dependencies = append(i.Dependencies, checkNSIS())
	//i.Dependencies = append(i.Dependencies, checkDocker())

	return nil
}

func checkWebView2() *packagemanager.Dependancy {
	version, _ := webviewloader.GetWebviewVersion("")
	installed := version != ""

	return &packagemanager.Dependancy{
		Name:           "WebView2 ",
		PackageName:    "N/A",
		Installed:      installed,
		InstallCommand: "Available at https://developer.microsoft.com/en-us/microsoft-edge/webview2/",
		Version:        version,
		Optional:       false,
		External:       true,
	}

}
