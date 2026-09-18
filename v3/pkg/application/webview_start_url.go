//go:build !wails_native

package application

import "github.com/wailsapp/wails/v3/internal/assetserver"

func webviewStartURL(value string) string {
	result, _ := assetserver.GetStartURL(value)
	return result
}
