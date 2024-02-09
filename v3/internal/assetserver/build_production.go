//go:build production

package assetserver

import "net/http"

func defaultIndexHTML() []byte {
	return []byte{}
}

func (a *AssetServer) setupHandler() (http.Handler, error) {
	return NewDefaultAssetHandler(a.options)
}

func GetDevServerURL() string {
	return ""
}
