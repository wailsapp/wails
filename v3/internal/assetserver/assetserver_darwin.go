//go:build darwin && !ios

package assetserver

import "net/url"

var baseURL = url.URL{
	Scheme: "wails",
	Host:   "localhost",
}

// platformJS is empty on darwin - no platform-specific JS needed.
var platformJS []byte
