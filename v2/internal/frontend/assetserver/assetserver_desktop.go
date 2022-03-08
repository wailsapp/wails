package assetserver

import (
	"bytes"
	"context"
	_ "embed"
	"io/fs"
	"log"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/internal/logger"
)

//go:embed defaultindex.html
var defaultHTML []byte

type DesktopAssetServer struct {
	assets    fs.FS
	runtimeJS []byte
	logger    *logger.Logger
}

func NewDesktopAssetServer(ctx context.Context, assets fs.FS, bindingsJSON string) (*DesktopAssetServer, error) {
	result := &DesktopAssetServer{}

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

func (d *DesktopAssetServer) LogDebug(message string, args ...interface{}) {
	if d.logger != nil {
		d.logger.Debug("[DesktopAssetServer] "+message, args...)
	}
}

func (a *DesktopAssetServer) processIndexHTML() ([]byte, error) {
	var indexHTML []byte
	var err error
	for tries := 0; tries < 10; tries++ {
		indexHTML, err = fs.ReadFile(a.assets, "index.html")
		if err != nil {
			time.Sleep(500 * time.Millisecond)
		}
	}
	if err != nil {
		indexHTML = defaultHTML
	}
	wailsOptions, err := extractOptions(indexHTML)
	if err != nil {
		log.Fatal(err)
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

func (a *DesktopAssetServer) Load(filename string) ([]byte, string, error) {
	var content []byte
	var err error
	switch filename {
	case "/":
		content, err = a.processIndexHTML()
	case "/wails/runtime.js":
		content = a.runtimeJS
	case "/wails/ipc.js":
		content = runtime.DesktopIPC
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
