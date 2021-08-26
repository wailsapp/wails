package assetserver

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"github.com/leaanthony/debme"
	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/internal/logger"
	"io/fs"
	"path/filepath"
	"strings"
)

type DesktopAssetServer struct {
	assets    debme.Debme
	indexFile []byte
	runtimeJS []byte
	assetdir  string
	logger    *logger.Logger
}

func NewDesktopAssetServer(ctx context.Context, assets embed.FS, bindingsJSON string) (*DesktopAssetServer, error) {
	result := &DesktopAssetServer{}

	_logger := ctx.Value("logger")
	if _logger != nil {
		result.logger = _logger.(*logger.Logger)
	}

	_assetdir := ctx.Value("assetdir")
	if _assetdir != nil {
		result.assetdir = _assetdir.(string)
		absdir, err := filepath.Abs(result.assetdir)
		if err != nil {
			return nil, err
		}
		result.LogInfo("Loading assets from: %s", absdir)
	}

	var buffer bytes.Buffer
	buffer.WriteString(`window.wailsbindings='` + bindingsJSON + `';` + "\n")
	buffer.Write(runtime.RuntimeJS)
	result.runtimeJS = buffer.Bytes()
	err := result.init(assets)
	return result, err
}

func (d *DesktopAssetServer) LogInfo(message string, args ...interface{}) {
	if d.logger != nil {
		d.logger.Info("[DesktopAssetServer] "+message, args...)
	}
}

func (d *DesktopAssetServer) SetAssetDir(assetdir string) {
	d.assetdir = assetdir
}

func processAssets(assets embed.FS) (debme.Debme, error) {

	result, err := debme.FS(assets, ".")
	if err != nil {
		return result, err
	}
	// Find index.html
	stat, err := fs.Stat(assets, "index.html")
	if stat != nil {
		return debme.FS(assets, ".")
	}
	var indexFiles slicer.StringSlicer
	err = fs.WalkDir(result, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, "index.html") {
			indexFiles.Add(path)
		}
		return nil
	})
	if err != nil {
		return debme.Debme{}, err
	}

	if indexFiles.Length() > 1 {
		return debme.Debme{}, fmt.Errorf("multiple 'index.html' files found in assets")
	}

	path, _ := filepath.Split(indexFiles.AsSlice()[0])
	return debme.FS(assets, path)
}

func (a *DesktopAssetServer) init(assets embed.FS) error {

	var err error
	a.assets, err = processAssets(assets)
	if err != nil {
		return err
	}
	indexHTML, err := a.assets.ReadFile("index.html")
	if err != nil {
		return err
	}
	a.indexFile, err = injectHTML(string(indexHTML), `<script src="/wails/runtime.js"></script>`)
	if err != nil {
		return err
	}
	a.indexFile, err = injectHTML(string(a.indexFile), `<script src="/wails/ipc.js"></script>`)
	if err != nil {
		return err
	}
	return nil
}

func (a *DesktopAssetServer) Load(filename string) ([]byte, string, error) {
	var content []byte
	var err error
	switch filename {
	case "/":
		content = a.indexFile
	case "/wails/runtime.js":
		content = a.runtimeJS
	case "/wails/ipc.js":
		content = runtime.DesktopIPC
	default:
		content, err = a.ReadFile(filename)
	}
	if err != nil {
		return nil, "", err
	}
	mimeType := GetMimetype(filename, content)
	return content, mimeType, nil
}
