//go:build windows

package assetserver

import "net/url"

var baseURL = url.URL{
	Scheme: "http",
	Host:   "wails.localhost",
}

// platformJS is empty on windows - no platform-specific JS needed.
var platformJS []byte
