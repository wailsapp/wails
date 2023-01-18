package runtime

import _ "embed"

//go:embed ipc_websocket.js
var WebsocketIPC []byte

//go:embed ipc.js
var DesktopIPC []byte
