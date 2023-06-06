//go:build windows && !production

package application

import (
	"github.com/wailsapp/go-webview2/pkg/edge"
)

func init() {
	showDevTools = func(chromium *edge.Chromium) {
		chromium.OpenDevToolsWindow()
	}
}
