//go:build dev
// +build dev

package assetserver

import (
	"context"
	"net/http"
	"strings"

	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
)

/*
The assetserver for dev serves assets from disk.
It injects a websocket based IPC script into `index.html`.
*/
func NewBrowserAssetServer(ctx context.Context, handler http.Handler, bindingsJSON string) (*AssetServer, error) {
	result, err := NewAssetServerWithHandler(ctx, handler, bindingsJSON)
	if err != nil {
		return nil, err
	}

	result.appendSpinnerToBody = true
	result.ipcJS = func(req *http.Request) []byte {
		if strings.Contains(req.UserAgent(), WailsUserAgentValue) {
			return runtime.DesktopIPC
		}
		return runtime.WebsocketIPC
	}

	return result, nil
}
