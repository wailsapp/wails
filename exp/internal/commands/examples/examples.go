package examples

import _ "embed"

//go:embed info.json
var Info []byte

//go:embed wails.exe.manifest
var Manifest []byte

//go:embed appicon.png
var AppIcon []byte

//go:embed icon.ico
var IconIco []byte
