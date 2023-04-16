package assetserver

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"

	"golang.org/x/net/html"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

const (
	runtimeJSPath = "/wails/runtime.js"
	ipcJSPath     = "/wails/ipc.js"
	runtimePath   = "/wails/runtime"
)

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
	wsHandler http.Handler
	runtimeJS []byte
	ipcJS     func(*http.Request) []byte

	logger  Logger
	runtime RuntimeAssets

	servingFromDisk     bool
	appendSpinnerToBody bool

	// Use http based runtime
	runtimeHandler RuntimeHandler

	// plugin scripts
	pluginScripts map[string]string

	assetServerWebView
}

func NewAssetServerMainPage(bindingsJSON string, options *options.App, servingFromDisk bool, logger Logger, runtime RuntimeAssets) (*AssetServer, error) {
	assetOptions, err := BuildAssetServerConfig(options)
	if err != nil {
		return nil, err
	}
	return NewAssetServer(bindingsJSON, assetOptions, servingFromDisk, logger, runtime)
}

func NewAssetServer(bindingsJSON string, options assetserver.Options, servingFromDisk bool, logger Logger, runtime RuntimeAssets) (*AssetServer, error) {
	handler, err := NewAssetHandler(options, logger)
	if err != nil {
		return nil, err
	}

	return NewAssetServerWithHandler(handler, bindingsJSON, servingFromDisk, logger, runtime)
}

func NewAssetServerWithHandler(handler http.Handler, bindingsJSON string, servingFromDisk bool, logger Logger, runtime RuntimeAssets) (*AssetServer, error) {
	var buffer bytes.Buffer
	if bindingsJSON != "" {
		buffer.WriteString(`window.wailsbindings='` + bindingsJSON + `';` + "\n")
	}
	buffer.Write(runtime.RuntimeDesktopJS())

	result := &AssetServer{
		handler:   handler,
		runtimeJS: buffer.Bytes(),

		// Check if we have been given a directory to serve assets from.
		// If so, this means we are in dev mode and are serving assets off disk.
		// We indicate this through the `servingFromDisk` flag to ensure requests
		// aren't cached in dev mode.
		servingFromDisk: servingFromDisk,
		logger:          logger,
		runtime:         runtime,
	}

	return result, nil
}

func (d *AssetServer) UseRuntimeHandler(handler RuntimeHandler) {
	d.runtimeHandler = handler
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
	if isWebSocket(req) {
		// Forward WebSockets to the distinct websocket handler if it exists
		if wsHandler := d.wsHandler; wsHandler != nil {
			wsHandler.ServeHTTP(rw, req)
		} else {
			rw.WriteHeader(http.StatusNotImplemented)
		}
		return
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
		for k, v := range recorder.HeaderMap {
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

	case runtimeJSPath:
		d.writeBlob(rw, path, d.runtimeJS)

	case runtimePath:
		if d.runtimeHandler != nil {
			d.runtimeHandler.HandleRuntimeCall(rw, req)
		} else {
			d.handler.ServeHTTP(rw, req)
		}

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
			return
		}
		d.handler.ServeHTTP(rw, req)
	}
}

func (d *AssetServer) processIndexHTML(indexHTML []byte) ([]byte, error) {
	htmlNode, err := getHTMLNode(indexHTML)
	if err != nil {
		return nil, err
	}

	if d.appendSpinnerToBody {
		err = appendSpinnerToBody(htmlNode)
		if err != nil {
			return nil, err
		}
	}

	if err := insertScriptInHead(htmlNode, runtimeJSPath); err != nil {
		return nil, err
	}

	if err := insertScriptInHead(htmlNode, ipcJSPath); err != nil {
		return nil, err
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

func (d *AssetServer) logDebug(message string, args ...interface{}) {
	if d.logger != nil {
		d.logger.Debug("[AssetServer] "+message, args...)
	}
}

func (d *AssetServer) logError(message string, args ...interface{}) {
	if d.logger != nil {
		d.logger.Error("[AssetServer] "+message, args...)
	}
}
