// +build wv2runtime.browser

package wv2runtime

import (
	"fmt"
	"github.com/leaanthony/webview2runtime"
)

func doInstallationStrategy(installStatus installationStatus) error {
	confirmed, err := webview2runtime.Confirm("This application requires the WebView2 runtime. Press OK to open the download page. Minimum version required: "+minimumRuntimeVersion, "Missing Requirements")
	if err != nil {
		return err
	}
	if confirmed {
		err = webview2runtime.OpenInstallerDownloadWebpage()
		if err != nil {
			return err
		}
	}

	return fmt.Errorf("webview2 runtime not installed")
}
