//go:build windows && !production

package runtime

import _ "embed"

//go:embed runtime_debug_desktop_windows.js
var DesktopRuntime []byte
