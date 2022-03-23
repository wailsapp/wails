package assetserver

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"

	"golang.org/x/net/html"
)

const (
	runtimeJSPath = "/wails/runtime.js"
	ipcJSPath     = "/wails/ipc.js"
)

type AssetServer struct {
	handler   http.Handler
	runtimeJS []byte
	ipcJS     func(*http.Request) []byte

	logger *logger.Logger

	servingFromDisk     bool
	appendSpinnerToBody bool
}

func NewAssetServer(ctx context.Context, options *options.App, bindingsJSON string) (*AssetServer, error) {
	handler, err := NewAsssetHandler(ctx, options)
	if err != nil {
		return nil, err
	}

	return NewAssetServerWithHandler(ctx, handler, bindingsJSON)
}

func NewAssetServerWithHandler(ctx context.Context, handler http.Handler, bindingsJSON string) (*AssetServer, error) {
	var buffer bytes.Buffer
	buffer.WriteString(`window.wailsbindings='` + bindingsJSON + `';` + "\n")
	buffer.Write(runtime.RuntimeDesktopJS)

	result := &AssetServer{
		handler:   handler,
		runtimeJS: buffer.Bytes(),

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

func (d *AssetServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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

		if recorder.Code != http.StatusOK {
			rw.WriteHeader(recorder.Code)
			return
		}

		content, err := d.processIndexHTML(recorder.Body.Bytes())
		if err != nil {
			d.serveError(rw, err, "Unable to processIndexHTML")
			return
		}

		d.writeBlob(rw, "/index.html", content)

	case runtimeJSPath:
		d.writeBlob(rw, path, d.runtimeJS)

	case ipcJSPath:
		content := runtime.DesktopIPC
		if d.ipcJS != nil {
			content = d.ipcJS(req)
		}
		d.writeBlob(rw, path, content)

	default:
		d.handler.ServeHTTP(rw, req)
	}
}

func (d *AssetServer) Load(filename string) ([]byte, string, error) {
	// This will be removed as soon as AssetsHandler have been fully introduced.
	if !strings.HasPrefix(filename, "/") {
		filename = "/" + filename
	}

	req, err := http.NewRequest(http.MethodGet, "wails://wails"+filename, nil)
	if err != nil {
		return nil, "", err
	}

	rw := httptest.NewRecorder()
	d.ServeHTTP(rw, req)

	content := rw.Body.Bytes()
	mimeType := rw.HeaderMap.Get(HeaderContentType)
	if mimeType == "" {
		mimeType = GetMimetype(filename, content)
	}
	return content, mimeType, nil
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

	var buffer bytes.Buffer
	err = html.Render(&buffer, htmlNode)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (d *AssetServer) writeBlob(rw http.ResponseWriter, filename string, blob []byte) {
	header := rw.Header()
	header.Set(HeaderContentLength, fmt.Sprintf("%d", len(blob)))
	if mimeType := header.Get(HeaderContentType); mimeType == "" {
		mimeType = GetMimetype(filename, blob)
		header.Set(HeaderContentType, mimeType)
	}

	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(blob); err != nil {
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
