//go:build !production

package assetserver

import (
	"context"
	"io/fs"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

// devProxyMaxIdleConnsPerHost is the idle connection pool the dev proxy keeps to the frontend
// dev server. It only ever targets that one host, so a generous pool is cheap, and
// IdleConnTimeout reaps whatever is not in use.
const devProxyMaxIdleConnsPerHost = 512

// retryTransport implements http.RoundTripper with retry logic for transient connection failures.
// This is particularly useful when the Vite dev server temporarily rejects connections due to
// high concurrency with many dynamic imports.
type retryTransport struct {
	base       http.RoundTripper
	maxRetries int
	delay      time.Duration
}

// RoundTrip executes a single HTTP transaction with retry logic.
func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 0; i < t.maxRetries; i++ {
		resp, err = t.base.RoundTrip(req)
		if err == nil {
			return resp, nil
		}
		// Only retry on connection errors (e.g., connection refused)
		if isConnectionError(err) && i < t.maxRetries-1 {
			time.Sleep(t.delay)
			continue
		}
		break
	}
	return resp, err
}

// isConnectionError checks if the error is a connection-related error that may be transient.
//
// Errors reporting that no local address could be allocated - EADDRNOTAVAIL ("cannot assign
// requested address") on Unix, a failed bind on Windows - must NOT be added here. Those mean the
// host is out of ephemeral ports, and retrying only prolongs an outage that is already affecting
// every process on the machine. Failing fast is the correct response.
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "broken pipe") ||
		strings.Contains(errStr, "connectex")
}

func NewAssetFileServer(vfs fs.FS) http.Handler {
	devServerURL := GetDevServerURL()
	if devServerURL == "" {
		return newAssetFileServerFS(vfs)
	}

	parsedURL, err := url.Parse(devServerURL)
	if err != nil {
		return http.HandlerFunc(
			func(rw http.ResponseWriter, req *http.Request) {
				logError(req.Context(), "[ExternalAssetHandler] Invalid FRONTEND_DEVSERVER_URL. Should be valid URL", "error", err.Error())
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			})

	}

	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	proxy.Transport = &retryTransport{
		base: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				// Force IPv4 for localhost connections to avoid IPv6 issues on Windows
				if parsedURL.Hostname() == "localhost" || parsedURL.Hostname() == "127.0.0.1" {
					return dialer.DialContext(ctx, "tcp4", addr)
				}
				return dialer.DialContext(ctx, network, addr)
			},
			// A dev page pulls hundreds of unbundled ES modules concurrently. Without this,
			// MaxIdleConnsPerHost falls back to DefaultMaxIdleConnsPerHost (2), so most of
			// those requests get a connection that is closed instead of pooled the moment it
			// is returned. Closing it here leaves a TIME_WAIT holding one of the proxy's own
			// ephemeral ports for 2*MSL (30s on macOS). On a large frontend that drains the
			// host's ephemeral range, at which point every process on the machine - not just
			// the app - starts failing to make outbound connections.
			MaxIdleConnsPerHost: devProxyMaxIdleConnsPerHost,
			IdleConnTimeout:     90 * time.Second,
		},
		maxRetries: 50,
		delay:      50 * time.Millisecond,
	}
	proxy.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
		logError(r.Context(), "[ExternalAssetHandler] Proxy error", "error", err.Error())
		rw.WriteHeader(http.StatusBadGateway)
	}

	return proxy
}

func GetDevServerURL() string {
	return os.Getenv("FRONTEND_DEVSERVER_URL")
}
