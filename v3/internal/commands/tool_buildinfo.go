package commands

import (
	"fmt"
	"github.com/wailsapp/wails/v3/internal/buildinfo"
)

type BuildInfoOptions struct{}

func BuildInfo(_ *BuildInfoOptions) error {

	info, err := buildinfo.Get()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", info)
	return nil
}
