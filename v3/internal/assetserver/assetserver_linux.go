//go:build linux && !android

package assetserver

import (
	_ "embed"
	"net/url"
)

var baseURL = url.URL{
	Scheme: "wails",
	Host:   "localhost",
}

//go:embed assetserver_linux.js
var platformJS []byte
