//go:build production

package bundledassets

import _ "embed"

//go:embed runtime.js
var RuntimeJS []byte
