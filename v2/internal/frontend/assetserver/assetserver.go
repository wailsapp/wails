package assetserver

import (
	"bytes"
	"context"
	_ "embed"
	iofs "io/fs"
	"log"
	"path"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"

	"golang.org/x/net/html"
)

const (
	runtimeJSPath = "/wails/runtime.js"
	ipcJSPath     = "/wails/ipc.js"
)

//go:embed defaultindex.html
var defaultHTML []byte

type AssetServer struct {
	assets    iofs.FS
	runtimeJS []byte
	ipcJS     []byte

	logger *logger.Logger

	servingFromDisk     bool
	appendSpinnerToBody bool
}

func NewAssetServer(ctx context.Context, options *options.App, bindingsJSON string) (*AssetServer, error) {
	assets := options.Assets

	if _, err := assets.Open("."); err != nil {
		return nil, err
	}

	subDir, err := fs.FindPathToFile(assets, "index.html")
	if err != nil {
		return nil, err
	}

	assets, err = iofs.Sub(assets, path.Clean(subDir))
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	buffer.WriteString(`window.wailsbindings='` + bindingsJSON + `';` + "\n")
	buffer.Write(runtime.RuntimeDesktopJS)

	result := &AssetServer{
		assets:    assets,
		runtimeJS: buffer.Bytes(),
		ipcJS:     runtime.DesktopIPC,

		// Check if we have been given a directory to serve assets from.
		// If so, this means we are in dev mode and are serving assets off disk.
		// We indicate this through the `servingFromDisk` flag to ensure requests
		// aren't cached in dev mode.
		servingFromDisk: ctx.Value("assetdir") != nil,
	}

	if _logger := ctx.Value("logger"); _logger != nil {
		result.logger = _logger.(*logger.Logger)
	}

	return result, nil
}

func (d *AssetServer) LogDebug(message string, args ...interface{}) {
	if d.logger != nil {
		d.logger.Debug("[AssetServer] "+message, args...)
	}
}

func (d *AssetServer) Load(filename string) ([]byte, string, error) {
	var content []byte
	var err error
	switch filename {
	case "/":
		content, err = d.loadFile("index.html")
		if err != nil {
			content = defaultHTML
		}

		content, err = d.ProcessIndexHTML(content)
	case runtimeJSPath:
		content = d.runtimeJS
	case ipcJSPath:
		content = d.ipcJS
	default:
		filename = strings.TrimPrefix(filename, "/")
		d.LogDebug("Loading file: %s", filename)
		content, err = d.loadFile(filename)
	}
	if err != nil {
		return nil, "", err
	}
	mimeType := GetMimetype(filename, content)
	return content, mimeType, nil
}

// loadFile will try to load the file from disk. If there is an error
// it will retry until eventually it will give up and error.
func (d *AssetServer) loadFile(filename string) ([]byte, error) {
	if !d.servingFromDisk {
		return iofs.ReadFile(d.assets, filename)
	}
	var result []byte
	var err error
	for tries := 0; tries < 50; tries++ {
		result, err = iofs.ReadFile(d.assets, filename)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
		}
	}
	return result, err
}

func (d *AssetServer) ProcessIndexHTML(indexHTML []byte) ([]byte, error) {
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

	wailsOptions, err := extractOptions(htmlNode)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if !wailsOptions.disableRuntimeInjection {
		err := insertScriptInHead(htmlNode, runtimeJSPath)
		if err != nil {
			return nil, err
		}
	}

	if !wailsOptions.disableIPCInjection {
		err := insertScriptInHead(htmlNode, ipcJSPath)
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
