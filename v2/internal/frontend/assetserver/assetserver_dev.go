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
The assetserver for the dev mode.
Depending on the UserAgent it injects a websocket based IPC script into `index.html` or the default desktop IPC. The
default desktop IPC is injected when the webview accesses the devserver.
*/
func NewDevAssetServer(ctx context.Context, handler http.Handler, wsHandler http.Handler, bindingsJSON string) (*AssetServer, error) {
	result, err := NewAssetServerWithHandler(ctx, handler, bindingsJSON)
	if err != nil {
		return nil, err
	}

	result.wsHandler = wsHandler
	result.appendSpinnerToBody = true
	result.ipcJS = func(req *http.Request) []byte {
		if strings.Contains(req.UserAgent(), WailsUserAgentValue) {
			return runtime.DesktopIPC
		}
		return runtime.WebsocketIPC
	}

	return result, nil
}
