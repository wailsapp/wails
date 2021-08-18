//go:build debug && desktop

package runtime

import _ "embed"

//go:embed runtime_debug_desktop.js
var RuntimeJS []byte
