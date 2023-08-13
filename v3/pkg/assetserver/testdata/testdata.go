package testdata

import "embed"

//go:embed index.html main.css main.js
var TopLevelFS embed.FS
