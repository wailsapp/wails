//go:build windows

package capabilities

import "github.com/wailsapp/go-webview2/webviewloader"

type version string

func (v version) IsAtLeast(input string) bool {
	result, err := webviewloader.CompareBrowserVersions(string(v), input)
	if err != nil {
		return false
	}
	return result >= 0
}

func newCapabilities(webview2version string) Capabilities {
	webview2 := version(webview2version)
	c := Capabilities{}
	c.HasNativeDrag = webview2.IsAtLeast("113.0.0.0")
	return c
}
