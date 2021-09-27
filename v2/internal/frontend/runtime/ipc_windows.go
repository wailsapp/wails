//go:build windows

package runtime

import _ "embed"

//go:embed ipc_windows.js
var DesktopIPC []byte
