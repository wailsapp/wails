//go:build !production

package main

import "github.com/wailsapp/wails/v3/pkg/application"

var Options = application.Options{
	Name:        "WebviewWindow Demo (debug)",
	Description: "A demo of the WebviewWindow API",
	Mac: application.MacOptions{
		ApplicationShouldTerminateAfterLastWindowClosed: true,
	},
}
