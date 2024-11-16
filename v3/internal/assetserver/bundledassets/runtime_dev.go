//go:build !production

package bundledassets

import _ "embed"

//go:embed runtime.debug.js
var RuntimeJS []byte
