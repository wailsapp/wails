//go:build dev
// +build dev

package assetserver

import (
	"bytes"
	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/pkg/options"
	"path/filepath"
)

/*

The assetserver for dev serves assets from disk.
It injects a websocket based IPC script into `index.html`.

*/

import (
	"os"
)

type BrowserAssetServer struct {
	runtimeJS  []byte
	assetdir   string
	appOptions *options.App
}

func NewBrowserAssetServer(assetdir string, bindingsJSON string, appOptions *options.App) (*BrowserAssetServer, error) {
	result := &BrowserAssetServer{
		assetdir:   assetdir,
		appOptions: appOptions,
	}

	var buffer bytes.Buffer
	buffer.WriteString(`window.wailsbindings='` + bindingsJSON + `';` + "\n")
	buffer.Write(runtime.RuntimeDesktopJS)
	result.runtimeJS = buffer.Bytes()
	return result, nil
}

func (a *BrowserAssetServer) loadFileFromDisk(filename string) ([]byte, error) {
	return os.ReadFile(filepath.Join(a.assetdir, filename))
}

func (a *BrowserAssetServer) processIndexHTML() ([]byte, error) {
	indexHTML, err := a.loadFileFromDisk("index.html")
	if err != nil {
		return nil, err
	}
	indexHTML, err = injectHTML(string(indexHTML), `<div id="wails-spinner"></div>`)
	if err != nil {
		return nil, err
	}
	wailsOptions, err := extractOptions(indexHTML)
	if err != nil {
		return nil, err
	}
	if wailsOptions.disableRuntimeInjection == false {
		indexHTML, err = injectHTML(string(indexHTML), `<script src="/wails/runtime.js"></script>`)
		if err != nil {
			return nil, err
		}
	}
	if wailsOptions.disableIPCInjection == false {
		indexHTML, err = injectHTML(string(indexHTML), `<script src="/wails/ipc.js"></script>`)
		if err != nil {
			return nil, err
		}
	}
	return indexHTML, nil
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
		content, err = a.loadFileFromDisk(filename)
	}
	if err != nil {
		return nil, "", err
	}
	mimeType := GetMimetype(filename, content)
	return content, mimeType, nil
}
