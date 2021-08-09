//+build debug
//+build windows

package runtime

import _ "embed"

//go:embed runtime_debug_windows.js
var RuntimeJS string
