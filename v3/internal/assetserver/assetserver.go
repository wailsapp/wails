package assetserver

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/assetserver/bundledassets"
)

const (
	webViewRequestHeaderWindowId   = "x-wails-window-id"
	webViewRequestHeaderWindowName = "x-wails-window-name"
	HeaderAcceptLanguage           = "accept-language"
)

// Platform-specific options set during application initialization.
var (
	disableGStreamerFix    bool
	enableGStreamerCaching bool
)

// SetGStreamerOptions configures GStreamer workaround options on Linux.
// This is called during application initialization.
func SetGStreamerOptions(disable, enableCaching bool) {
	disableGStreamerFix = disable
	enableGStreamerCaching = enableCaching
}

type RuntimeHandler interface {
	HandleRuntimeCall(w http.ResponseWriter, r *http.Request)
}

type service struct {
	Route   string
	Handler http.Handler
}

type AssetServer struct {
	options  *Options
	handler  http.Handler
	services []service

	assetServerWebView
}

func NewAssetServer(options *Options) (*AssetServer, error) {
	result := &AssetServer{
		options: options,
	}

	userHandler := options.Handler
	if userHandler == nil {
		userHandler = http.NotFoundHandler()
	}

	handler := http.Handler(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				result.serveHTTP(w, r, userHandler)
			}))

	if middleware := options.Middleware; middleware != nil {
		handler = middleware(handler)
	}

	result.handler = handler

	return result, nil
}

func (a *AssetServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	start := time.Now()
	wrapped := newContentTypeSniffer(rw)
	defer func() {
		if _, err := wrapped.complete(); err != nil {
			a.options.Logger.Error("Error writing response data.", "uri", req.RequestURI, "error", err)
		}
	}()

	req = req.WithContext(contextWithLogger(req.Context(), a.options.Logger))
	a.handler.ServeHTTP(wrapped, req)

	a.options.Logger.Info(
		"Asset Request:",
		"windowName", req.Header.Get(webViewRequestHeaderWindowName),
		"windowID", req.Header.Get(webViewRequestHeaderWindowId),
		"code", wrapped.status,
		"method", req.Method,
		"path", req.URL.EscapedPath(),
		"duration", time.Since(start),
	)
}

func (a *AssetServer) serveHTTP(rw http.ResponseWriter, req *http.Request, userHandler http.Handler) {
	if isWebSocket(req) {
		// WebSockets are not supported by the AssetServer
		rw.WriteHeader(http.StatusNotImplemented)
		return
	}

	header := rw.Header()
	// TODO: I don't think this is needed now?
	//if a.servingFromDisk {
	//	header.Add(HeaderCacheControl, "no-cache")
	//}

	reqPath := req.URL.Path

	// Handle internal Wails endpoints - these work regardless of user's handler
	if strings.HasPrefix(reqPath, "/wails/") {
		a.serveWailsEndpoint(rw, reqPath[6:]) // Strip "/wails" prefix
		return
	}

	switch reqPath {
	case "", "/", "/index.html":
		// Cache the accept-language header
		// before passing the request down the chain.
		acceptLanguage := req.Header.Get(HeaderAcceptLanguage)
		if acceptLanguage == "" {
			acceptLanguage = "en"
		}

		wrapped := &fallbackResponseWriter{
			rw:  rw,
			req: req,
			fallback: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Set content type for default index.html
				header.Set(HeaderContentType, "text/html; charset=utf-8")
				a.writeBlob(rw, indexHTML, defaultIndexHTML(acceptLanguage))
			}),
		}
		userHandler.ServeHTTP(wrapped, req)

	default:
		// Check if the path matches a service route
		for _, svc := range a.services {
			if strings.HasPrefix(reqPath, svc.Route) {
				req.URL.Path = strings.TrimPrefix(reqPath, svc.Route)
				svc.Handler.ServeHTTP(rw, req)
				return
			}
		}

		// Forward to the user-provided handler
		userHandler.ServeHTTP(rw, req)
	}
}

func (a *AssetServer) AttachServiceHandler(route string, handler http.Handler) {
	a.services = append(a.services, service{route, handler})
}

// serveWailsEndpoint handles internal /wails/* endpoints.
func (a *AssetServer) serveWailsEndpoint(rw http.ResponseWriter, path string) {
	rw.Header().Set(HeaderContentType, "application/javascript")
	switch path {
	case "/runtime.js":
		rw.Write(bundledassets.RuntimeJS)
	case "/platform.js":
		rw.Write(getPlatformJS())
	default:
		rw.WriteHeader(http.StatusNotFound)
	}
}

// getPlatformJS returns the platform-specific JavaScript based on current options.
func getPlatformJS() []byte {
	if platformJS == nil {
		return nil
	}

	// If the fix is disabled, return empty
	if disableGStreamerFix {
		return nil
	}

	// If caching is not enabled, modify the JS to disable caching
	if !enableGStreamerCaching {
		return bytes.ReplaceAll(platformJS,
			[]byte("const ENABLE_CACHING = true;"),
			[]byte("const ENABLE_CACHING = false;"))
	}

	return platformJS
}

func (a *AssetServer) writeBlob(rw http.ResponseWriter, filename string, blob []byte) {
	err := ServeFile(rw, filename, blob)
	if err != nil {
		a.serveError(rw, err, "Error writing file content.", "filename", filename)
	}
}

func (a *AssetServer) serveError(rw http.ResponseWriter, err error, msg string, args ...interface{}) {
	args = append(args, "error", err)
	a.options.Logger.Error(msg, args...)
	rw.WriteHeader(http.StatusInternalServerError)
}

func GetStartURL(userURL string) (string, error) {
	devServerURL := GetDevServerURL()
	startURL := baseURL.String()
	if devServerURL != "" {
		// Parse the port
		parsedURL, err := url.Parse(devServerURL)
		if err != nil {
			return "", fmt.Errorf("error parsing environment variable `FRONTEND_DEVSERVER_URL`: %w. Please check your `Taskfile.yml` file", err)
		}
		port := parsedURL.Port()
		if port != "" {
			baseURL.Host = net.JoinHostPort(baseURL.Hostname(), port)
			startURL = baseURL.String()
		}
	}

	if userURL != "" {
		parsedURL, err := baseURL.Parse(userURL)
		if err != nil {
			return "", fmt.Errorf("error parsing URL: %w", err)
		}

		startURL = parsedURL.String()
	}

	return startURL, nil
}
