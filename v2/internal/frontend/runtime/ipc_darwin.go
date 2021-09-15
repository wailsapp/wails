//go:build darwin

package runtime

import _ "embed"

//go:embed ipc_darwin.js
var DesktopIPC []byte
