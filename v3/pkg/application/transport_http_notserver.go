//go:build !server

package application

import "net/http"

// checkRuntimeRequest permits every request outside server mode: desktop and
// mobile webviews deliver runtime calls through the app's own asset handler
// rather than a network listener.
func checkRuntimeRequest(http.ResponseWriter, *http.Request) bool {
	return true
}
