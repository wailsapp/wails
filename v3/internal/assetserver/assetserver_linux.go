//go:build linux && !android && !openharmony

package assetserver

import "net/url"

var baseURL = url.URL{
	Scheme: "wails",
	Host:   "localhost",
}
