//go:build dev
// +build dev

package devserver

import (
	"net"
	"net/http"
	"net/netip"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

// allowedHostsEnvVar names additional hostnames the dev server should trust,
// beyond loopback and its own bind address. The value is a comma-separated
// list, matched against the request Host's hostname only (the port is not
// checked, so a reverse proxy listening on a different public port still
// works). This is the escape hatch for LAN-by-hostname and proxied setups.
const allowedHostsEnvVar = "WAILS_DEV_ALLOWED_HOSTS"

// localhostName is the only hostname that is inherently loopback: unlike an
// arbitrary name it cannot be repointed elsewhere by a DNS answer.
const localhostName = "localhost"

// defaultHTTPPort is assumed when a Host header carries no port. The dev server
// speaks plain HTTP, so a port-less Host can only have meant port 80.
const defaultHTTPPort = "80"

// hostAllowed reports whether a request carrying requestHost may be served.
//
// DNS rebinding requires a resolvable name in the Host header, so the only
// hosts accepted are ones that cannot be repointed by DNS: localhost, any IP
// literal, and the exact address the server was told to bind. Extra hostnames
// can be trusted through allowedHostsEnvVar for reverse-proxy or
// LAN-by-hostname setups. Anything else — including exotic spellings such as a
// trailing dot, shorthand IPs, or unbracketed IPv6 — is rejected, so the guard
// fails closed and the env var is the one way to widen it.
func hostAllowed(requestHost, bindHost, bindPort string, extraHosts []string) bool {
	host, port := splitRequestHost(requestHost)

	lowerHost := strings.ToLower(host)
	for _, extra := range extraHosts {
		if lowerHost == extra {
			return true
		}
	}

	if port != bindPort {
		return false
	}

	if strings.EqualFold(host, localhostName) {
		return true
	}
	if _, err := netip.ParseAddr(host); err == nil {
		return true
	}
	// An empty bind host (e.g. a ":34115" wildcard bind) only matches the empty
	// host the CLI's own poller sends; a browser cannot emit an empty Host.
	return strings.EqualFold(host, bindHost)
}

// splitRequestHost splits a Host header into host and port, treating a missing
// port as defaultHTTPPort. A bracketed IPv6 literal with no port (e.g. "[::1]")
// has its brackets stripped so the caller sees the bare address.
func splitRequestHost(requestHost string) (host, port string) {
	if h, p, err := net.SplitHostPort(requestHost); err == nil {
		return h, p
	}
	return strings.Trim(requestHost, "[]"), defaultHTTPPort
}

// parseAllowedHosts turns the comma-separated allowedHostsEnvVar value into a
// slice of trimmed, lower-cased hostnames, dropping empty entries.
func parseAllowedHosts(value string) []string {
	var hosts []string
	for _, entry := range strings.Split(value, ",") {
		if entry = strings.ToLower(strings.TrimSpace(entry)); entry != "" {
			hosts = append(hosts, entry)
		}
	}
	return hosts
}

// requireAllowedHost builds middleware that rejects any request whose Host
// header is not allowed by hostAllowed. It is registered ahead of the cookie
// middleware so a rebound page is turned away before it is handed the
// capability or reaches any route. The bind address and env var are read once,
// at registration; an unparseable bind address leaves both parts empty and so
// rejects everything, which is the safe default (and unreachable in practice,
// since no listener starts without a valid address).
func (d *DevWebServer) requireAllowedHost() echo.MiddlewareFunc {
	bindHost, bindPort, _ := net.SplitHostPort(d.devServerAddr)
	extraHosts := parseAllowedHosts(os.Getenv(allowedHostsEnvVar))
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !hostAllowed(c.Request().Host, bindHost, bindPort, extraHosts) {
				d.logger.Error("Dev server request rejected: Host %q is not an allowed host; set "+allowedHostsEnvVar+" to trust additional hostnames", c.Request().Host)
				return c.NoContent(http.StatusForbidden)
			}
			return next(c)
		}
	}
}
