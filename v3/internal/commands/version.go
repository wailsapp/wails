package commands

import (
	_ "embed"
)

//go:embed version.txt
var VersionString string

type VersionOptions struct{}

func Version(_ *VersionOptions) error {
	println(VersionString)
	return nil
}
