//go:build linux && production

package runtime

import _ "embed"

//go:embed runtime_production_desktop_linux.js
var DesktopRuntime []byte
