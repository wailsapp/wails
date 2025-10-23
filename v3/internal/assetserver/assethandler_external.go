package assetserver

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// NewProxyServer creates a simple reverse proxy for the given URL
func NewProxyServer(proxyURL string) http.Handler {
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(parsedURL)
}

// NewExternalAssetsHandler creates a handler that proxies requests to an external URL
// with fallback support to a base handler
func NewExternalAssetsHandler(logger *slog.Logger, options *Options, externalURL *url.URL) http.Handler {
	baseHandler := options.Handler
	if baseHandler == nil {
		baseHandler = http.NotFoundHandler()
	}

	// Create error for skipping proxy
	errSkipProxy := fmt.Errorf("skip proxying")

	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(externalURL)

	// Configure transport to handle self-signed certificates in development
	// In production, you should use proper certificates
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Skip certificate verification for development
		},
	}

	// Store the original director
	baseDirector := proxy.Director

	// Custom Director to log requests
	proxy.Director = func(r *http.Request) {
		baseDirector(r)
		if logger != nil {
			logger.Debug("Proxying request", "url", r.URL.String())
		}
	}

	// ModifyResponse to handle error cases and fallback
	proxy.ModifyResponse = func(res *http.Response) error {
		// Allow WebSocket upgrades
		if res.StatusCode == http.StatusSwitchingProtocols {
			return nil
		}
		// Fall back to base handler for 404 or 405
		if res.StatusCode == http.StatusNotFound || res.StatusCode == http.StatusMethodNotAllowed {
			return errSkipProxy
		}
		return nil
	}

	// ErrorHandler to manage fallback
	proxy.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
		if errors.Is(err, errSkipProxy) {
			// Use the fallback handler
			baseHandler.ServeHTTP(rw, r)
		} else {
			// Log the error and return bad gateway
			if logger != nil {
				logger.Error("Proxy error", "error", err.Error())
			}
			rw.WriteHeader(http.StatusBadGateway)
		}
	}

	// Only proxy GET requests that are not Wails-specific paths
	var result http.Handler = http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			// Don't proxy Wails-specific paths
			if strings.HasPrefix(req.URL.Path, "/wails/") {
				// Let the base handler (Wails) handle these
				baseHandler.ServeHTTP(rw, req)
				return
			}

			// Proxy GET requests to external server
			if req.Method == http.MethodGet {
				proxy.ServeHTTP(rw, req)
				return
			}

			// Forward non-GET requests to base handler
			baseHandler.ServeHTTP(rw, req)
		})

	// Apply middleware if configured
	if middleware := options.Middleware; middleware != nil {
		result = middleware(result)
	}

	return result
}

// IsExternalURL checks if a URL is external (has http:// or https:// scheme)
func IsExternalURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}