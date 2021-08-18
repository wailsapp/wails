//go:build desktop && windows

package runtime

import _ "embed"

//go:embed ipc_windows.js
var IPCJS []byte
