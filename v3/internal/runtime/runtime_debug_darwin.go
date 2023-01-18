//go:build darwin && !production

package runtime

import _ "embed"

//go:embed runtime_debug_desktop_darwin.js
var DesktopRuntime []byte
