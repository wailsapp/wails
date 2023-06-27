//go:build debug || !production

package runtime

import _ "embed"

//go:embed runtime_debug_desktop.js
var RuntimeDesktopJS []byte
