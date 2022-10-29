//go:build dev
// +build dev

package devserver

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func newExternalDevServerAssetHandler(logger *logger.Logger, url *url.URL, options assetserver.Options) http.Handler {
	handler := newExternalAssetsHandler(logger, url, options.Handler)

	if middleware := options.Middleware; middleware != nil {
		handler = middleware(handler)
	}

	return handler
}

func newExternalAssetsHandler(logger *logger.Logger, url *url.URL, handler http.Handler) http.Handler {
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
		if handler == nil {
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
		if handler != nil && errors.Is(err, errSkipProxy) {
			if logger != nil {
				logger.Debug("[ExternalAssetHandler] Loading '%s' failed, using AssetHandler", r.URL)
			}
			handler.ServeHTTP(rw, r)
		} else {
			if logger != nil {
				logger.Error("[ExternalAssetHandler] Proxy error: %v", err)
			}
			rw.WriteHeader(http.StatusBadGateway)
		}
	}

	return http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			if req.Method == http.MethodGet {
				proxy.ServeHTTP(rw, req)
				return
			}

			if handler != nil {
				handler.ServeHTTP(rw, req)
				return
			}

			rw.WriteHeader(http.StatusMethodNotAllowed)
		})
}
