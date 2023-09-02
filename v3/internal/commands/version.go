package commands

import (
	_ "embed"
	"github.com/wailsapp/wails/v3/internal/version"
)

type VersionOptions struct{}

func Version(_ *VersionOptions) error {
	println(version.VersionString)
	return nil
}
