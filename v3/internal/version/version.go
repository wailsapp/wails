package version

import (
	_ "embed"
	"github.com/wailsapp/wails/v3/internal/debug"
)

//go:embed version.txt
var versionString string

const DevVersion = "v3.0.0-dev"

func String() string {
	if !IsDev() {
		return versionString
	}
	return DevVersion
}

func LatestStable() string {
	return versionString
}

func IsDev() bool {
	return debug.LocalModulePath != ""
}
