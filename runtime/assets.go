package runtime

import _ "embed"

//go:embed assets/bridge.js
var BridgeJS []byte

//go:embed assets/wails.js
var WailsJS string

//go:embed assets/wails.css
var WailsCSS string

//go:embed js/runtime/init.js
var InitJS []byte
