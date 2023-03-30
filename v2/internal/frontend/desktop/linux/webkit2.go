//go:build linux

package linux

/*
#cgo linux pkg-config: webkit2gtk-4.0
#include "webkit2/webkit2.h"
*/
import "C"
import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/linux"

	"github.com/wailsapp/wails/v2/pkg/assetserver/webview"
)

func validateWebKit2Version(options *options.App) {
	if C.webkit_get_major_version() == 2 && C.webkit_get_minor_version() >= webview.Webkit2MinMinorVersion {
		return
	}

	msg := linux.DefaultMessages()
	if options.Linux != nil && options.Linux.Messages != nil {
		msg = options.Linux.Messages
	}

	v := fmt.Sprintf("2.%d.0", webview.Webkit2MinMinorVersion)
	showModalDialogAndExit("WebKit2GTK", fmt.Sprintf(msg.WebKit2GTKMinRequired, v))
}
