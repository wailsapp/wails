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

type AssetServer struct {
	indexFile  []byte
	runtimeJS  []byte
	assetdir   string
	appOptions *options.App
}

func NewAssetServer(assetdir string, bindingsJSON string, appOptions *options.App) (*AssetServer, error) {
	result := &AssetServer{
		assetdir:   assetdir,
		appOptions: appOptions,
	}

	err := result.init()
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	buffer.WriteString(`window.wailsbindings='` + bindingsJSON + `';` + "\n")
	buffer.Write(runtime.RuntimeDesktopJS)
	result.runtimeJS = buffer.Bytes()
	err = result.init()
	return result, err
}

func (a *AssetServer) loadFileFromDisk(filename string) ([]byte, error) {
	return os.ReadFile(filepath.Join(a.assetdir, filename))
}

func (a *AssetServer) init() error {
	var err error
	a.indexFile, err = a.loadFileFromDisk("index.html")
	if err != nil {
		return err
	}
	a.indexFile, err = injectHTML(string(a.indexFile), `<div id="wails-spinner"></div>`)
	if err != nil {
		return err
	}
	a.indexFile, err = injectHTML(string(a.indexFile), `<script src="/wails/ipc.js"></script>`)
	if err != nil {
		return err
	}
	a.indexFile, err = injectHTML(string(a.indexFile), `<script src="/wails/runtime.js"></script>`)
	if err != nil {
		return err
	}
	return nil
}

func (a *AssetServer) Load(filename string) ([]byte, string, error) {
	var content []byte
	var err error
	switch filename {
	case "/":
		content = a.indexFile
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
