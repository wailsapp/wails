//go:build linux && !production

package runtime

import _ "embed"

//go:embed runtime_debug_desktop_linux.js
var DesktopRuntime []byte
