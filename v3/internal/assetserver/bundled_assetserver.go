package assetserver

import (
	"io/fs"
	"net/http"
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
	b.handler.ServeHTTP(rw, req)
}
