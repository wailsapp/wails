package version

import (
	_ "embed"
)

//go:embed version.txt
var VersionString string
