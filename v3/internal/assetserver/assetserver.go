package assetserver

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"
)

const (
	webViewRequestHeaderWindowId   = "x-wails-window-id"
	webViewRequestHeaderWindowName = "x-wails-window-name"
	servicePrefix                  = "wails/services"
)

type RuntimeHandler interface {
	HandleRuntimeCall(w http.ResponseWriter, r *http.Request)
}

type AssetServer struct {
	options *Options

	handler http.Handler

	services map[string]http.Handler

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
	wrapped := &contentTypeSniffer{rw: rw}

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
	switch reqPath {
	case "", "/", "/index.html":
		recorder := httptest.NewRecorder()
		userHandler.ServeHTTP(recorder, req)
		for k, v := range recorder.Result().Header {
			header[k] = v
		}

		switch recorder.Code {
		case http.StatusOK:
			a.writeBlob(rw, indexHTML, recorder.Body.Bytes())

		case http.StatusNotFound:
			a.writeBlob(rw, indexHTML, defaultIndexHTML())

		default:
			rw.WriteHeader(recorder.Code)
		}

	default:

		// Check if the path matches the keys in the services map
		for route, handler := range a.services {
			if strings.HasPrefix(reqPath, route) {
				req.URL.Path = strings.TrimPrefix(reqPath, route)
				handler.ServeHTTP(rw, req)
				return
			}
		}

		// Check if it can be served by the user-provided handler
		if !strings.HasPrefix(reqPath, servicePrefix) {
			userHandler.ServeHTTP(rw, req)
			return
		}

		rw.WriteHeader(http.StatusNotFound)
		return
	}
}

func (a *AssetServer) AttachServiceHandler(prefix string, handler http.Handler) {
	if a.services == nil {
		a.services = make(map[string]http.Handler)
	}
	a.services[prefix] = handler
}

func (a *AssetServer) writeBlob(rw http.ResponseWriter, filename string, blob []byte) {
	err := ServeFile(rw, filename, blob)
	if err != nil {
		a.serveError(rw, err, "Unable to write content %s", filename)
	}
}

func (a *AssetServer) serveError(rw http.ResponseWriter, err error, msg string, args ...interface{}) {
	args = append(args, err)
	a.options.Logger.Error(msg+":", args...)
	rw.WriteHeader(http.StatusInternalServerError)
}

func GetStartURL(userURL string) (string, error) {
	devServerURL := GetDevServerURL()
	startURL := baseURL.String()
	if devServerURL != "" {
		// Parse the port
		parsedURL, err := url.Parse(devServerURL)
		if err != nil {
			return "", fmt.Errorf("Error parsing environment variable 'FRONTEND_DEVSERVER_URL`: " + err.Error() + ". Please check your `Taskfile.yml` file")
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
			return "", fmt.Errorf("Error parsing URL: " + err.Error())
		}

		startURL = parsedURL.String()
	}

	return startURL, nil
}
