package assetserver

import (
	"fmt"
	"html"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	runtimePath      = "/wails/runtime"
	capabilitiesPath = "/wails/capabilities"
	flagsPath        = "/wails/flags"

	webViewRequestHeaderWindowId   = "x-wails-window-id"
	webViewRequestHeaderWindowName = "x-wails-window-name"
)

type RuntimeHandler interface {
	HandleRuntimeCall(w http.ResponseWriter, r *http.Request)
}

type AssetServer struct {
	options *Options

	handler   http.Handler
	wsHandler *httputil.ReverseProxy

	pluginScripts map[string]string

	devServerURL string

	assetServerWebView
}

func NewAssetServer(options *Options) (*AssetServer, error) {
	result := &AssetServer{
		options: options,
	}

	var err error
	result.handler, err = result.setupHandler()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *AssetServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	start := time.Now()
	wrapped := &contentTypeSniffer{rw: rw}
	a.serveHTTP(wrapped, req)
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

func (a *AssetServer) serveHTTP(rw http.ResponseWriter, req *http.Request) {

	if a.wsHandler != nil {
		a.wsHandler.ServeHTTP(rw, req)
		return
	} else {
		if isWebSocket(req) {
			// WebSockets are not supported by the AssetServer
			rw.WriteHeader(http.StatusNotImplemented)
			return
		}
	}

	header := rw.Header()
	// TODO: I don't think this is needed now?
	//if a.servingFromDisk {
	//	header.Add(HeaderCacheControl, "no-cache")
	//}

	path := req.URL.Path
	switch path {
	case "", "/", "/index.html":
		recorder := httptest.NewRecorder()
		a.handler.ServeHTTP(recorder, req)
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
		return

	case capabilitiesPath:
		var data = a.options.GetCapabilities()
		a.writeBlob(rw, path, data)

	case flagsPath:
		var data = a.options.GetFlags()
		a.writeBlob(rw, path, data)

	case runtimePath:
		a.options.RuntimeHandler.ServeHTTP(rw, req)
		return

	default:
		// Check if this is a plugin script
		if script, ok := a.pluginScripts[path]; ok {
			a.writeBlob(rw, path, []byte(script))
		} else {
			a.handler.ServeHTTP(rw, req)
			return
		}
	}
}

func (a *AssetServer) writeBlob(rw http.ResponseWriter, filename string, blob []byte) {
	err := serveFile(rw, filename, blob)
	if err != nil {
		a.serveError(rw, err, "Unable to write content %s", filename)
	}
}

func (a *AssetServer) serveError(rw http.ResponseWriter, err error, msg string, args ...interface{}) {
	args = append(args, err)
	a.options.Logger.Error(msg+":", args...)
	rw.WriteHeader(http.StatusInternalServerError)
}

func (a *AssetServer) AddPluginScript(pluginName string, script string) {
	if a.pluginScripts == nil {
		a.pluginScripts = make(map[string]string)
	}
	pluginName = strings.ReplaceAll(pluginName, "/", "_")
	pluginName = html.EscapeString(pluginName)
	pluginScriptName := fmt.Sprintf("/wails/plugin/%s.js", pluginName)
	a.pluginScripts[pluginScriptName] = script
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
			baseURL.Host = net.JoinHostPort(baseURL.Host, port)
			startURL = baseURL.String()
		}
	} else {
		if userURL != "" {
			// parse the url
			parsedURL, err := url.Parse(userURL)
			if err != nil {
				return "", fmt.Errorf("Error parsing URL: " + err.Error())
			}
			if parsedURL.Scheme == "" {
				baseURL.Path = path.Join(baseURL.Path, userURL)
				startURL = baseURL.String()
				// if the original URL had a trailing slash, add it back
				if strings.HasSuffix(userURL, "/") && !strings.HasSuffix(startURL, "/") {
					startURL = startURL + "/"
				}
			} else {
				startURL = userURL
			}
		}
	}
	return startURL, nil
}
