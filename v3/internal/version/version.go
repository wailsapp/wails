package version

import (
	_ "embed"
	"github.com/wailsapp/wails/v3/internal/debug"
)

//go:embed version.txt
var versionString string

func String() string {
	if !IsDev() {
		return versionString
	}
	return "v3 dev"
}

func LatestStable() string {
	return versionString
}

func IsDev() bool {
	return debug.LocalModulePath != ""
}
