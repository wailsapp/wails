//go:build dev || bindings || (!dev && !production && !bindings)
// +build dev bindings !dev,!production,!bindings

package runtime

import _ "embed"

//go:embed runtime_dev_desktop.js
var RuntimeDesktopJS []byte
