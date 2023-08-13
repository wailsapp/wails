//go:build !production

package assetserver

import (
	"log/slog"
	"net/http"
	"strings"
)

/*
The assetserver for the dev mode.
Depending on the UserAgent it injects a websocket based IPC script into `index.html` or the default desktop IPC. The
default desktop IPC is injected when the webview accesses the devserver.
*/
func NewDevAssetServer(handler http.Handler, servingFromDisk bool, logger *slog.Logger, runtime RuntimeAssets, runtimeHandler RuntimeHandler) (*AssetServer, error) {
	result, err := NewAssetServerWithHandler(handler, servingFromDisk, logger, runtime, true, runtimeHandler)
	if err != nil {
		return nil, err
	}

	result.ipcJS = func(req *http.Request) []byte {
		if strings.Contains(req.UserAgent(), WailsUserAgentValue) {
			return runtime.DesktopIPC()
		}
		return runtime.WebsocketIPC()
	}

	return result, nil
}
