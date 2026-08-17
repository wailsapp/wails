//go:build !wails_native

package application

import "net/http"

// AssetServerTransport lets a transport serve the Wails asset handler.
type AssetServerTransport interface {
	Transport
	ServeAssets(assetHandler http.Handler) error
}

// TransportHTTPHandler contributes middleware to the WebView asset server.
type TransportHTTPHandler interface {
	Handler() func(next http.Handler) http.Handler
}
