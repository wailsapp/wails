//go:build dev

package runtime

import _ "embed"

//go:embed runtime_dev_desktop.js
var RuntimeDesktopJS []byte
