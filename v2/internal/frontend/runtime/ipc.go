//go:build darwin || windows

package runtime

import _ "embed"

//go:embed ipc.js
var DesktopIPC []byte
