package runtime

import _ "embed"

//go:embed runtime.js
var runtimeJS []byte

//go:embed runtime.debug.js
var runtimeDebugJS []byte
