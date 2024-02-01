package runtime

import _ "embed"

//go:embed runtime.js
var RuntimeJS []byte

//go:embed runtime.debug.js
var RuntimeDebugJS []byte
