//go:build dev
// +build dev

package assetserver

import (
	"bytes"
	"context"
	"io/fs"
	"strings"

	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/internal/logger"
	"golang.org/x/net/html"
)

/*

The assetserver for dev serves assets from disk.
It injects a websocket based IPC script into `index.html`.

*/

type BrowserAssetServer struct {
	assets    fs.FS
	runtimeJS []byte
	logger    *logger.Logger
}

func NewBrowserAssetServer(ctx context.Context, assets fs.FS, bindingsJSON string) (*BrowserAssetServer, error) {
	result := &BrowserAssetServer{}
	_logger := ctx.Value("logger")
	if _logger != nil {
		result.logger = _logger.(*logger.Logger)
	}

	var err error
	result.assets, err = prepareAssetsForServing(assets)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	buffer.WriteString(`window.wailsbindings='` + bindingsJSON + `';` + "\n")
	buffer.Write(runtime.RuntimeDesktopJS)
	result.runtimeJS = buffer.Bytes()

	return result, nil
}

func (d *BrowserAssetServer) LogDebug(message string, args ...interface{}) {
	if d.logger != nil {
		d.logger.Debug("[BrowserAssetServer] "+message, args...)
	}
}

func (a *BrowserAssetServer) processIndexHTML() ([]byte, error) {
	indexHTML, err := fs.ReadFile(a.assets, "index.html")
	if err != nil {
		return nil, err
	}
	htmlNode, err := getHTMLNode(indexHTML)
	if err != nil {
		return nil, err
	}
	err = appendSpinnerToBody(htmlNode)
	if err != nil {
		return nil, err
	}
	wailsOptions, err := extractOptions(indexHTML)
	if err != nil {
		return nil, err
	}

	if wailsOptions.disableIPCInjection == false {
		err := insertScriptInHead(htmlNode, "/wails/ipc.js")
		if err != nil {
			return nil, err
		}
	}

	if wailsOptions.disableRuntimeInjection == false {
		err := insertScriptInHead(htmlNode, "/wails/runtime.js")
		if err != nil {
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

func (a *BrowserAssetServer) Load(filename string) ([]byte, string, error) {
	var content []byte
	var err error
	switch filename {
	case "/":
		content, err = a.processIndexHTML()
	case "/wails/runtime.js":
		content = a.runtimeJS
	case "/wails/ipc.js":
		content = runtime.WebsocketIPC
	default:
		filename = strings.TrimPrefix(filename, "/")
		a.LogDebug("Loading file: %s", filename)
		content, err = fs.ReadFile(a.assets, filename)
	}
	if err != nil {
		return nil, "", err
	}
	mimeType := GetMimetype(filename, content)
	return content, mimeType, nil
}
