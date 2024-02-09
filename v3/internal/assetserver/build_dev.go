//go:build !production

package assetserver

import (
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

//go:embed defaultindex.html
var defaultHTML []byte

func defaultIndexHTML() []byte {
	return defaultHTML
}

func (a *AssetServer) setupHandler() (http.Handler, error) {

	// Do we have an external dev server URL?
	a.devServerURL = GetDevServerURL()
	if a.devServerURL == "" {
		return NewDefaultAssetHandler(a.options)
	}

	// Parse the URL
	parsedURL, err := url.Parse(a.devServerURL)
	if err != nil {
		return nil, fmt.Errorf("invalid FRONTEND_DEVSERVER_URL. Should be valid URL: %s", err.Error())
	}

	baseHandler := a.options.Handler

	errSkipProxy := fmt.Errorf("skip proxying")

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	baseDirector := proxy.Director
	proxy.Director = func(r *http.Request) {
		baseDirector(r)
		if a.options.Logger != nil {
			a.options.Logger.Debug("ExternalAssetHandler: loading", "url", r.URL)
		}
	}

	proxy.ModifyResponse = func(res *http.Response) error {
		if baseHandler == nil {
			return nil
		}

		if res.StatusCode == http.StatusSwitchingProtocols {
			return nil
		}

		if res.StatusCode == http.StatusNotFound || res.StatusCode == http.StatusMethodNotAllowed {
			return errSkipProxy
		}

		return nil
	}

	proxy.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
		if baseHandler != nil && errors.Is(err, errSkipProxy) {
			if a.options.Logger != nil {
				a.options.Logger.Debug("ExternalAssetHandler: Loading file failed, using original AssetHandler", "url", r.URL)
			}
			baseHandler.ServeHTTP(rw, r)
		} else {
			if a.options.Logger != nil {
				a.options.Logger.Error("ExternalAssetHandler: Proxy error", "error", err.Error())
			}
			rw.WriteHeader(http.StatusBadGateway)
		}
	}

	var result http.Handler = http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			if req.Method == http.MethodGet {
				proxy.ServeHTTP(rw, req)
				return
			}

			if baseHandler != nil {
				baseHandler.ServeHTTP(rw, req)
				return
			}

			rw.WriteHeader(http.StatusMethodNotAllowed)
		})

	if middleware := a.options.Middleware; middleware != nil {
		result = middleware(result)
	}

	return result, nil
}

func GetDevServerURL() string {
	return os.Getenv("FRONTEND_DEVSERVER_URL")
}
