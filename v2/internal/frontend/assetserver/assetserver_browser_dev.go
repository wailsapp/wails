//go:build dev
// +build dev

package assetserver

import (
	"context"

	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/pkg/options"
)

/*
The assetserver for dev serves assets from disk.
It injects a websocket based IPC script into `index.html`.
*/

func NewBrowserAssetServer(ctx context.Context, options *options.App, bindingsJSON string) (*AssetServer, error) {
	result, err := NewAssetServer(ctx, options, bindingsJSON)
	if err != nil {
		return nil, err
	}

	result.appendSpinnerToBody = true
	result.ipcJS = runtime.WebsocketIPC
	return result, nil
}
