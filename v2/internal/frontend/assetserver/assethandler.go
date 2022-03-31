package assetserver

import (
	"context"
	_ "embed"
	iofs "io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed defaultindex.html
var defaultHTML []byte

type assetHandler struct {
	fs iofs.FS

	logger *logger.Logger

	servingFromDisk bool
}

func NewAsssetHandler(ctx context.Context, options *options.App) (http.Handler, error) {
	vfs := options.Assets
	if vfs != nil {
		if _, err := vfs.Open("."); err != nil {
			return nil, err
		}

		subDir, err := fs.FindPathToFile(vfs, "index.html")
		if err != nil {
			return nil, err
		}

		vfs, err = iofs.Sub(vfs, path.Clean(subDir))
		if err != nil {
			return nil, err
		}
	}

	result := &assetHandler{
		fs: vfs,

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

func (d *assetHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if d.fs == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	filename := strings.TrimPrefix(req.URL.Path, "/")
	if d.logger != nil {
		d.logger.Debug("[AssetHandler] Loading file '%s'", filename)
	}

	var content []byte
	var err error
	switch filename {
	case "", "index.html":
		content, err = d.loadFile("index.html")
		if err != nil {
			err = nil
			content = defaultHTML
		}

	default:
		content, err = d.loadFile(filename)
	}

	if os.IsNotExist(err) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if err == nil {
		mimeType := GetMimetype(filename, content)
		rw.Header().Set(HeaderContentType, mimeType)
		rw.WriteHeader(http.StatusOK)
		_, err = rw.Write(content)
	}

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		if d.logger != nil {
			d.logger.Error("[AssetHandler] Unable to load file '%s': %s", filename, err)
		}
	}
}

// loadFile will try to load the file from disk. If there is an error
// it will retry until eventually it will give up and error.
func (d *assetHandler) loadFile(filename string) ([]byte, error) {
	if !d.servingFromDisk {
		return iofs.ReadFile(d.fs, filename)
	}
	var result []byte
	var err error
	for tries := 0; tries < 50; tries++ {
		result, err = iofs.ReadFile(d.fs, filename)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
		}
	}
	return result, err
}
