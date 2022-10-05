package assetserver

import (
	"bytes"
	"context"
	"fmt"
	iofs "io/fs"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/internal/logger"

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

func NewAssetServer(ctx context.Context, vfs iofs.FS, assetsHandler http.Handler, bindingsJSON string) (*AssetServer, error) {
	handler, err := NewAssetHandler(ctx, vfs, assetsHandler)
	if err != nil {
		return nil, err
	}

	return NewAssetServerWithHandler(ctx, handler, bindingsJSON)
}

func NewAssetServerWithHandler(ctx context.Context, handler http.Handler, bindingsJSON string) (*AssetServer, error) {
	var buffer bytes.Buffer
	if bindingsJSON != "" {
		buffer.WriteString(`window.wailsbindings='` + bindingsJSON + `';` + "\n")
	}
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

	if isWebSocket(req) {
		// WebSockets can always directly be forwarded to the handler
		d.handler.ServeHTTP(rw, req)
		return
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

// ProcessHTTPRequest processes the HTTP Request by faking a golang HTTP Server.
// The request will be finished with a StatusNotImplemented code if no handler has written to the response.
func (d *AssetServer) ProcessHTTPRequest(logInfo string, rw http.ResponseWriter, reqGetter func() (*http.Request, error)) {
	rw = &contentTypeSniffer{rw: rw} // Make sure we have a Content-Type sniffer

	req, err := reqGetter()
	if err != nil {
		d.logError("Error processing request '%s': %s (HttpResponse=500)", logInfo, err)

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Body == nil {
		req.Body = http.NoBody
	}
	defer req.Body.Close()

	if req.RemoteAddr == "" {
		// 192.0.2.0/24 is "TEST-NET" in RFC 5737
		req.RemoteAddr = "192.0.2.1:1234"
	}

	if req.RequestURI == "" && req.URL != nil {
		req.RequestURI = req.URL.String()
	}

	if req.ContentLength == 0 {
		req.ContentLength, _ = strconv.ParseInt(req.Header.Get(HeaderContentLength), 10, 64)
	} else {
		req.Header.Set(HeaderContentLength, fmt.Sprintf("%d", req.ContentLength))
	}

	if host := req.Header.Get(HeaderHost); host != "" {
		req.Host = host
	}

	d.ServeHTTP(rw, req)
	rw.WriteHeader(http.StatusNotImplemented) // This is a NOP when a handler has already written and set the status
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
