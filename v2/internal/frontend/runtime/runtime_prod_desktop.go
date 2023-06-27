//go:build production && !debug

package runtime

import _ "embed"

//go:embed runtime_prod_desktop.js
var RuntimeDesktopJS []byte
