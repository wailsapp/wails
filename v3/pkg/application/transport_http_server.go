//go:build server

package application

import (
	"errors"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// checkRuntimeRequest guards /wails/runtime in server mode, where the
// endpoint is reachable from any page open in the user's browser. Bound
// methods are dispatched only for POST requests from the server's own origin
// or an origin listed in ServerOptions.WebSocketOriginPatterns, mirroring the
// policy applied to the server's WebSocket endpoints. Requests carrying
// neither Origin nor Sec-Fetch-Site come from non-browser clients and are
// allowed.
func checkRuntimeRequest(rw http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		rw.Header().Set("Allow", http.MethodPost)
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var server ServerOptions
	if app := globalApplication; app != nil {
		server = app.options.Server
	}
	if server.WebSocketAllowAllOrigins {
		return true
	}

	if err := authorizeRuntimeOrigin(r, server.WebSocketOriginPatterns); err != nil {
		http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return false
	}
	return true
}

// authorizeRuntimeOrigin applies the same origin matching as coder/websocket:
// the request host is always permitted, and patterns follow path.Match syntax,
// matching scheme and host when the pattern includes a scheme.
func authorizeRuntimeOrigin(r *http.Request, patterns []string) error {
	origin := r.Header.Get("Origin")
	if origin == "" {
		// Browsers send Origin on every cross-origin POST; Sec-Fetch-Site
		// covers any that do not.
		switch r.Header.Get("Sec-Fetch-Site") {
		case "cross-site", "same-site":
			return errors.New("cross-origin runtime request")
		}
		return nil
	}

	u, err := url.Parse(origin)
	if err != nil || u.Host == "" {
		return errors.New("invalid Origin header")
	}
	if strings.EqualFold(r.Host, u.Host) {
		return nil
	}
	for _, pattern := range patterns {
		target := u.Host
		if strings.Contains(pattern, "://") {
			target = u.Scheme + "://" + u.Host
		}
		if matched, _ := path.Match(strings.ToLower(pattern), strings.ToLower(target)); matched {
			return nil
		}
	}
	return errors.New("runtime request origin not permitted")
}
