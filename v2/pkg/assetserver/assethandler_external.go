//go:build dev
// +build dev

package assetserver

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func NewExternalAssetsHandler(logger Logger, options assetserver.Options, url *url.URL) http.Handler {
	baseHandler := options.Handler

	errSkipProxy := fmt.Errorf("skip proxying")

	proxy := httputil.NewSingleHostReverseProxy(url)
	baseDirector := proxy.Director
	proxy.Director = func(r *http.Request) {
		baseDirector(r)
		if logger != nil {
			logger.Debug("[ExternalAssetHandler] Loading '%s'", r.URL)
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
			if logger != nil {
				logger.Debug("[ExternalAssetHandler] Loading '%s' failed, using original AssetHandler", r.URL)
			}
			baseHandler.ServeHTTP(rw, r)
		} else {
			if logger != nil {
				logger.Error("[ExternalAssetHandler] Proxy error: %v", err)
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

	if middleware := options.Middleware; middleware != nil {
		result = middleware(result)
	}

	return result
}
