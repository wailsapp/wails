package assetserver

import (
	"fmt"
	"io/fs"
	"net/http"
)

// Options defines the configuration of the AssetServer.
type Options struct {
	// Assets defines the static assets to be used. A GET request is first tried to be served from this Assets. If the Assets returns
	// `os.ErrNotExist` for that file, the request handling will fallback to the Handler and tries to serve the GET
	// request from it.
	//
	// If set to nil, all GET requests will be forwarded to Handler.
	Assets fs.FS

	// Handler will be called for every GET request that can't be served from Assets, due to `os.ErrNotExist`. Furthermore all
	// non GET requests will always be served from this Handler.
	//
	// If not defined, the result is the following in cases where the Handler would have been called:
	//   GET request:   `http.StatusNotFound`
	//   Other request: `http.StatusMethodNotAllowed`
	Handler http.Handler

	// Middleware is a HTTP Middleware which allows to hook into the AssetServer request chain. It allows to skip the default
	// request handler dynamically, e.g. implement specialized Routing etc.
	// The Middleware is called to build a new `http.Handler` used by the AssetSever and it also receives the default
	// handler used by the AssetServer as an argument.
	//
	// If not defined, the default AssetServer request chain is executed.
	//
	// Multiple Middlewares can be chained together with:
	//   ChainMiddleware(middleware ...Middleware) Middleware
	Middleware Middleware
}

// Validate the options
func (o Options) Validate() error {
	if o.Assets == nil && o.Handler == nil && o.Middleware == nil {
		return fmt.Errorf("AssetServer options invalid: either Assets, Handler or Middleware must be set")
	}

	return nil
}
