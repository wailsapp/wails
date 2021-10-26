package assets

import _ "embed"

//go:embed desktop_darwin.js
var desktopDarwinJS string

//go:embed desktop_windows.js
var desktopWindowsJS string

//go:embed wails.js
var wailsJS string
