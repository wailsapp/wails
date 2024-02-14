//go:build production

package assetserver

import (
	"io/fs"
	"net/http"
)

func defaultIndexHTML() []byte {
	return []byte{}
}

func NewAssetFileServer(vfs fs.FS) http.Handler {
	return newAssetFileServerFS(vfs)
}

func GetDevServerURL() string {
	return ""
}
