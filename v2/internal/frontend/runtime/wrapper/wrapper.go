package wrapper

import "embed"

//go:embed runtime.js runtime.d.ts package.json
var RuntimeWrapper embed.FS
