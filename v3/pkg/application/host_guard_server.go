//go:build server

package application

import (
	"net"
	"net/http"
	"net/netip"
	"path"
	"strings"
)

// localhostName is the only hostname that is inherently loopback: unlike an
// arbitrary name it cannot be repointed elsewhere by a DNS answer.
const localhostName = "localhost"

// isLoopbackBindHost reports whether the server binds only a loopback
// address. Only "localhost" and loopback IP literals count; any other name is
// treated as a deliberate network deployment and left unguarded.
func isLoopbackBindHost(bindHost string) bool {
	host := strings.Trim(bindHost, "[]")
	if strings.EqualFold(host, localhostName) {
		return true
	}
	addr, err := netip.ParseAddr(host)
	return err == nil && addr.IsLoopback()
}

// serverHostAllowed reports whether a request carrying requestHost may be
// served by a loopback-bound server.
//
// DNS rebinding requires a resolvable name in the Host header, so the hosts
// accepted are ones that cannot be repointed by DNS: localhost, any IP
// literal, and the configured bind host. Hosts matching the host part of a
// WebSocketOriginPatterns entry are also accepted so a local reverse proxy
// that forwards its public Host keeps working. Anything else, including
// spellings such as a trailing dot or shorthand IPs, is rejected. The port is
// not checked: localhost and IP literals cannot be rebound whatever the port.
func serverHostAllowed(requestHost, bindHost string, patterns []string) bool {
	if requestHost == "" {
		// Browsers always send Host; only HTTP/1.0 clients omit it.
		return true
	}
	host := requestHost
	if h, _, err := net.SplitHostPort(requestHost); err == nil {
		host = h
	} else {
		host = strings.Trim(requestHost, "[]")
	}

	if strings.EqualFold(host, localhostName) {
		return true
	}
	if _, err := netip.ParseAddr(host); err == nil {
		return true
	}
	if bindHost != "" && strings.EqualFold(host, strings.Trim(bindHost, "[]")) {
		return true
	}

	lowerRequestHost := strings.ToLower(requestHost)
	lowerHost := strings.ToLower(host)
	for _, pattern := range patterns {
		if _, rest, ok := strings.Cut(pattern, "://"); ok {
			pattern = rest
		}
		pattern = strings.ToLower(pattern)
		if matched, _ := path.Match(pattern, lowerRequestHost); matched {
			return true
		}
		if matched, _ := path.Match(pattern, lowerHost); matched {
			return true
		}
	}
	return false
}

// requireAllowedHost wraps next with a Host header check when the server is
// bound to a loopback address, so a page whose name has been rebound to
// 127.0.0.1 cannot pass the same-origin checks on /wails/runtime and the
// WebSocket endpoints. Non-loopback binds are returned unchanged. The check
// applies even when WebSocketAllowAllOrigins is set; WebSocketOriginPatterns
// is the way to trust additional hostnames.
func requireAllowedHost(next http.Handler, app *App, bindHost string, server ServerOptions) http.Handler {
	if !isLoopbackBindHost(bindHost) {
		return next
	}
	patterns := append([]string(nil), server.WebSocketOriginPatterns...)
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !serverHostAllowed(r.Host, bindHost, patterns) {
			if app != nil {
				app.warning("Server request rejected: Host %q is not allowed for a loopback bind; add it to ServerOptions.WebSocketOriginPatterns to trust it", r.Host)
			}
			http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(rw, r)
	})
}
