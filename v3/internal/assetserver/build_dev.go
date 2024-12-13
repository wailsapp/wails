//go:build !production

package assetserver

import (
	_ "embed"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func NewAssetFileServer(vfs fs.FS) http.Handler {
	devServerURL := GetDevServerURL()
	if devServerURL == "" {
		return newAssetFileServerFS(vfs)
	}

	parsedURL, err := url.Parse(devServerURL)
	if err != nil {
		return http.HandlerFunc(
			func(rw http.ResponseWriter, req *http.Request) {
				logError(req.Context(), "[ExternalAssetHandler] Invalid FRONTEND_DEVSERVER_URL. Should be valid URL", "error", err.Error())
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			})

	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	proxy.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
		logError(r.Context(), "[ExternalAssetHandler] Proxy error", "error", err.Error())
		rw.WriteHeader(http.StatusBadGateway)
	}

	return proxy
}

func GetDevServerURL() string {
	return os.Getenv("FRONTEND_DEVSERVER_URL")
}
