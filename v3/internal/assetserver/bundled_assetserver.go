package assetserver

import (
	"github.com/wailsapp/wails/v3/internal/assetserver/bundledassets"
	"io/fs"
	"net/http"
	"strconv"
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
		switch req.URL.Path {
		case "/wails/runtime.js":
			rw.Header().Set("Content-Type", "application/javascript")
			rw.Header().Set(HeaderContentLength, strconv.Itoa(len(bundledassets.RuntimeJS)))
			if req.Method == http.MethodHead {
				rw.WriteHeader(http.StatusOK)
				return
			}
			rw.Write([]byte(bundledassets.RuntimeJS))
			return
		}
		http.NotFound(rw, req)
		return
	}
	b.handler.ServeHTTP(rw, req)
}
