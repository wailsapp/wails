package bundledassets

import _ "embed"

//go:embed runtime.js
var runtimeProductionJS []byte

//go:embed runtime.debug.js
var runtimeDebugJS []byte

// RuntimeFiles returns the production and debug bundles for CLI export.
func RuntimeFiles() (production, debug []byte) {
	return runtimeProductionJS, runtimeDebugJS
}
