package assetserver

import (
	"github.com/wailsapp/wails/v3/internal/assetserver/bundledassets"
	"io/fs"
	"net/http"
	"strings"
)

type BundledAssetServer struct {
	handler http.Handler
}

func NewBundledAssetFileServer(fs fs.FS) *BundledAssetServer {
	return &BundledAssetServer{
		handler: NewAssetFileServer(fs),
	}
}

func (b *BundledAssetServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/wails/") {
		// Strip the /wails prefix
		req.URL.Path = req.URL.Path[6:]
		switch req.URL.Path {
		case "/runtime.js":
			rw.Header().Set("Content-Type", "application/javascript")
			rw.Write([]byte(bundledassets.RuntimeJS))
			return
		}
		return
	}
	b.handler.ServeHTTP(rw, req)
}
