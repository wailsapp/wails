//go:build !server

package application

import (
	"net/http"
	"net/url"
)

// runtimeClientHeader is sent by the bundled runtime on every call. A page
// cannot attach it to a cross-origin request without a CORS preflight, which
// the asset handler never answers, so its presence rules out simple requests
// such as <img> loads and no-cors form or text/plain POSTs.
const runtimeClientHeader = "x-wails-client-id"

// checkRuntimeRequest guards /wails/runtime on desktop and mobile, where every
// document loaded in the webview shares the app's custom scheme and could
// otherwise call any bound method blind. Calls must be POST (a GET would be a
// simple request, and it changes state), carry the runtime's custom header, and
// not declare a different origin or a cross-site fetch.
func checkRuntimeRequest(rw http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		rw.Header().Set("Allow", http.MethodPost)
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}
	if r.Header.Get(runtimeClientHeader) == "" || !runtimeOriginAllowed(r) {
		http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return false
	}
	return true
}

// runtimeOriginAllowed rejects a request whose Origin names a host other than
// the one the webview addressed, or whose Sec-Fetch-Site says it left the
// page's origin. Platforms that send neither (or an opaque "null" Origin for
// the custom scheme) are covered by the header requirement alone.
func runtimeOriginAllowed(r *http.Request) bool {
	switch r.Header.Get("Sec-Fetch-Site") {
	case "cross-site", "same-site":
		return false
	}
	origin := r.Header.Get("Origin")
	if origin == "" || origin == "null" {
		return true
	}
	u, err := url.Parse(origin)
	return err == nil && u.Host != "" && u.Host == r.Host
}
