//go:build dev

package runtime

import _ "embed"

//go:embed ipc_dev.js
var IPCJS []byte
