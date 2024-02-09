package assetserver

import "net/url"

var baseURL = url.URL{
	Scheme: "http",
	Host:   "wails.localhost",
}
