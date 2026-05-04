//go:build !production

package assetserver

import (
	"context"
	_ "embed"
	"io/fs"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

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
