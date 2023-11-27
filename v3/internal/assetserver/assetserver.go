package assetserver

import (
	"bytes"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	runtimeJSPath    = "/wails/runtime.js"
	ipcJSPath        = "/wails/ipc.js"
	runtimePath      = "/wails/runtime"
	capabilitiesPath = "/wails/capabilities"
	flagsPath        = "/wails/flags"
)

const webViewRequestHeaderWindowId = "x-wails-window-id"
const webViewRequestHeaderWindowName = "x-wails-window-name"

type RuntimeAssets interface {
	DesktopIPC() []byte
	WebsocketIPC() []byte
	RuntimeDesktopJS() []byte
}

type RuntimeHandler interface {
	HandleRuntimeCall(w http.ResponseWriter, r *http.Request)
}

type AssetServer struct {
	handler   http.Handler
	runtimeJS []byte
	debug     bool
	ipcJS     func(*http.Request) []byte

	logger  *slog.Logger
	runtime RuntimeAssets
	options *Options

	servingFromDisk bool

	// Use http based runtime
	runtimeHandler RuntimeHandler

	// plugin scripts
	pluginScripts map[string]string

	// GetCapabilities returns the capabilities of the runtime
	GetCapabilities func() []byte

	// GetFlags returns the application flags
	GetFlags func() []byte

	// External dev server proxy
	wsHandler *httputil.ReverseProxy

	// External dev server URL
	devServerURL string

	assetServerWebView
}

func NewAssetServer(options *Options, servingFromDisk bool, logger *slog.Logger, runtime RuntimeAssets, debug bool, runtimeHandler RuntimeHandler) (*AssetServer, error) {
	handler, err := NewAssetHandler(options, logger)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	buffer.Write(runtime.RuntimeDesktopJS())

	result := &AssetServer{
		handler:        handler,
		runtimeJS:      buffer.Bytes(),
		runtimeHandler: runtimeHandler,
		options:        options,

		// Check if we have been given a directory to serve assets from.
		// If so, this means we are in dev mode and are serving assets off disk.
		// We indicate this through the `servingFromDisk` flag to ensure requests
		// aren't cached in dev mode.
		servingFromDisk: servingFromDisk,
		logger:          logger,
		runtime:         runtime,
		debug:           debug,
	}

	// Check if proxy required
	result.devServerURL = GetDevServerURL()
	if result.devServerURL != "" {
		logger.Info("Using External DevServer", "url", result.devServerURL)
		// Parse devServerURL into url.URL
		devServerURL, err := url.Parse(result.devServerURL)
		if err != nil {
			return nil, err
		}
		err = result.checkDevServerURL(devServerURL)
		if err != nil {
			return nil, err
		}
		result.wsHandler = httputil.NewSingleHostReverseProxy(devServerURL)
	}
	return result, nil
}

func (d *AssetServer) checkDevServerURL(devServerURL *url.URL) error {
	// Open a connection to the devserver URL
	hostPort := devServerURL.Hostname() + ":" + devServerURL.Port()
	_, err := net.DialTimeout("tcp", hostPort, 1*time.Second)
	if err != nil {
		return fmt.Errorf("unable to connect to dev server: %s. Please check it's running", d.devServerURL)
	}
	return nil
}

func (d *AssetServer) LogDetails() {
	if d.debug {
		var info = []any{
			"assetsFS", d.options.Assets != nil,
			"middleware", d.options.Middleware != nil,
			"handler", d.options.Handler != nil,
			"devServerURL", d.devServerURL,
		}
		if d.devServerURL != "" {
			info = append(info, "devServerURL", d.devServerURL)
		}
		d.logger.Info("AssetServer Info:", info...)
	}
}

func (d *AssetServer) AddPluginScript(pluginName string, script string) {
	if d.pluginScripts == nil {
		d.pluginScripts = make(map[string]string)
	}
	pluginName = strings.ReplaceAll(pluginName, "/", "_")
	pluginName = html.EscapeString(pluginName)
	pluginScriptName := fmt.Sprintf("/plugin_%s_%d.js", pluginName, rand.Intn(100000))
	d.pluginScripts[pluginScriptName] = script
}

func (d *AssetServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	start := time.Now()
	wrapped := &contentTypeSniffer{rw: rw}
	d.serveHTTP(wrapped, req)
	d.logger.Info(
		"Asset Request:",
		"windowName", req.Header.Get(webViewRequestHeaderWindowName),
		"windowID", req.Header.Get(webViewRequestHeaderWindowId),
		"code", wrapped.status,
		"method", req.Method,
		"path", req.URL.EscapedPath(),
		"duration", time.Since(start),
	)
}

func (d *AssetServer) serveHTTP(rw http.ResponseWriter, req *http.Request) {

	if d.wsHandler != nil {
		d.wsHandler.ServeHTTP(rw, req)
		return
	} else {
		if isWebSocket(req) {
			// WebSockets are not supported by the AssetServer
			rw.WriteHeader(http.StatusNotImplemented)
			return
		}
	}

	header := rw.Header()
	if d.servingFromDisk {
		header.Add(HeaderCacheControl, "no-cache")
	}

	path := req.URL.Path
	switch path {
	case "", "/", "/index.html":
		recorder := httptest.NewRecorder()
		d.handler.ServeHTTP(recorder, req)
		for k, v := range recorder.Result().Header {
			header[k] = v
		}

		switch recorder.Code {
		case http.StatusOK:
			content, err := d.processIndexHTML(recorder.Body.Bytes())
			if err != nil {
				d.serveError(rw, err, "Unable to processIndexHTML")
				return
			}
			d.writeBlob(rw, indexHTML, content)

		case http.StatusNotFound:
			d.writeBlob(rw, indexHTML, defaultHTML)

		default:
			rw.WriteHeader(recorder.Code)

		}
		return

	case runtimeJSPath:
		d.writeBlob(rw, path, d.runtimeJS)

	case capabilitiesPath:
		var data = []byte("{}")
		if d.GetCapabilities != nil {
			data = d.GetCapabilities()
		}
		d.writeBlob(rw, path, data)

	case flagsPath:
		var data = []byte("{}")
		if d.GetFlags != nil {
			data = d.GetFlags()
		}
		d.writeBlob(rw, path, data)

	case runtimePath:
		d.runtimeHandler.HandleRuntimeCall(rw, req)
		return

	case ipcJSPath:
		content := d.runtime.DesktopIPC()
		if d.ipcJS != nil {
			content = d.ipcJS(req)
		}
		d.writeBlob(rw, path, content)

	default:
		// Check if this is a plugin script
		if script, ok := d.pluginScripts[path]; ok {
			d.writeBlob(rw, path, []byte(script))
		} else {
			d.handler.ServeHTTP(rw, req)
			return
		}
	}
}

func (d *AssetServer) processIndexHTML(indexHTML []byte) ([]byte, error) {
	htmlNode, err := getHTMLNode(indexHTML)
	if err != nil {
		return nil, err
	}

	if d.debug {
		err = appendSpinnerToBody(htmlNode)
		if err != nil {
			return nil, err
		}
	}

	if err := insertScriptInHead(htmlNode, runtimeJSPath); err != nil {
		return nil, err
	}

	if d.debug {
		if err := insertScriptInHead(htmlNode, ipcJSPath); err != nil {
			return nil, err
		}
	}

	// Inject plugins
	for scriptName := range d.pluginScripts {
		if err := insertScriptInHead(htmlNode, scriptName); err != nil {
			return nil, err
		}
	}

	var buffer bytes.Buffer
	err = html.Render(&buffer, htmlNode)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (d *AssetServer) writeBlob(rw http.ResponseWriter, filename string, blob []byte) {
	err := serveFile(rw, filename, blob)
	if err != nil {
		d.serveError(rw, err, "Unable to write content %s", filename)
	}
}

func (d *AssetServer) serveError(rw http.ResponseWriter, err error, msg string, args ...interface{}) {
	args = append(args, err)
	d.logError(msg+": %s", args...)
	rw.WriteHeader(http.StatusInternalServerError)
}

func (d *AssetServer) logInfo(message string, args ...interface{}) {
	d.logger.Info("Asset Request: "+message, args...)
}

func (d *AssetServer) logError(message string, args ...interface{}) {
	d.logger.Error("Asset Request: "+message, args...)
}
