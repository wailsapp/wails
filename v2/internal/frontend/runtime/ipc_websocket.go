//go:build dev || server || hybrid
// +build dev server hybrid

package runtime

import _ "embed"

//go:embed ipc_websocket.js
var WebsocketIPC []byte
