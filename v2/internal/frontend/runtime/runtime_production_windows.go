//+build !debug
//+build windows

package runtime

import _ "embed"

//go:embed runtime_production_windows.js
var RuntimeJS string
