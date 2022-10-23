package myfrontend

import "embed"

//go:embed all:dist
var Assets embed.FS
