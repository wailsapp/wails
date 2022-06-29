package assetserver

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	iofs "io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed defaultindex.html
var defaultHTML []byte

const (
	indexHTML = "index.html"
)

type assetHandler struct {
	fs      iofs.FS
	handler http.Handler

	logger *logger.Logger

	retryMissingFiles bool
}

func NewAssetHandler(ctx context.Context, options *options.App) (http.Handler, error) {
	vfs := options.Assets
	if vfs != nil {
		if _, err := vfs.Open("."); err != nil {
			return nil, err
		}

		subDir, err := fs.FindPathToFile(vfs, indexHTML)
		if err != nil {
			return nil, err
		}

		vfs, err = iofs.Sub(vfs, path.Clean(subDir))
		if err != nil {
			return nil, err
		}
	}

	result := &assetHandler{
		fs:      vfs,
		handler: options.AssetsHandler,
	}

	if _logger := ctx.Value("logger"); _logger != nil {
		result.logger = _logger.(*logger.Logger)
	}

	return result, nil
}

func (d *assetHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	handler := d.handler
	if strings.EqualFold(req.Method, http.MethodGet) {
		filename := strings.TrimPrefix(req.URL.Path, "/")
		if filename == "" {
			filename = indexHTML
		}

		d.logDebug("[AssetHandler] Loading file '%s'", filename)
		if err := d.serveFSFile(rw, filename); err != nil {
			if os.IsNotExist(err) {
				if handler != nil {
					d.logDebug("[AssetHandler] File '%s' not found, serving '%s' by AssetHandler", filename, req.URL)
					handler.ServeHTTP(rw, req)
					err = nil
				} else if filename == indexHTML {
					err = serveFile(rw, filename, defaultHTML)
				} else {
					rw.WriteHeader(http.StatusNotFound)
					err = nil
				}
			}

			if err != nil {
				d.logError("[AssetHandler] Unable to load file '%s': %s", filename, err)
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
		}
	} else if handler != nil {
		d.logDebug("[AssetHandler] No GET request, serving '%s' by AssetHandler", req.URL)
		handler.ServeHTTP(rw, req)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// serveFile will try to load the file from the fs.FS and write it to the response
func (d *assetHandler) serveFSFile(rw http.ResponseWriter, filename string) error {
	if d.fs == nil {
		return os.ErrNotExist
	}

	file, err := d.fs.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	statInfo, err := file.Stat()
	if err != nil {
		return err
	}

	rw.Header().Set(HeaderContentLength, fmt.Sprintf("%d", statInfo.Size()))

	var buf [512]byte
	n, err := file.Read(buf[:])
	if err != nil && err != io.EOF {
		return err
	}

	// Detect MimeType by sniffing the first 512 bytes
	if contentType := GetMimetype(filename, buf[:n]); contentType != "" {
		rw.Header().Set(HeaderContentType, contentType)
	}

	// Write the first bytes
	_, err = io.Copy(rw, bytes.NewReader(buf[:n]))
	if err != nil {
		return err
	}

	// Copy the remaining content of the file
	_, err = io.Copy(rw, file)
	return err
}

func (d *assetHandler) logDebug(message string, args ...interface{}) {
	if d.logger != nil {
		d.logger.Debug("[AssetHandler] "+message, args...)
	}
}

func (d *assetHandler) logError(message string, args ...interface{}) {
	if d.logger != nil {
		d.logger.Error("[AssetHandler] "+message, args...)
	}
}
