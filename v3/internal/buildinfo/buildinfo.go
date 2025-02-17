package buildinfo

import (
	"fmt"
	"runtime/debug"
	"slices"

	"github.com/samber/lo"
)

type Info struct {
	Development   bool
	Version       string
	BuildSettings map[string]string
	wailsPackage  *debug.Module
}

func Get() (*Info, error) {

	var result Info

	// BuildInfo contains the build info for the application
	var BuildInfo *debug.BuildInfo

	var ok bool
	BuildInfo, ok = debug.ReadBuildInfo()
	if !ok {
		return nil, fmt.Errorf("could not read build info from binary")
	}
	result.BuildSettings = lo.Associate(BuildInfo.Settings, func(setting debug.BuildSetting) (string, string) {
		return setting.Key, setting.Value
	})
	result.Version = BuildInfo.Main.Version
	result.Development = -1 != slices.IndexFunc(BuildInfo.Settings, func(setting debug.BuildSetting) bool {
		return setting.Key == "vcs" && setting.Value == "git"
	})

	return &result, nil

}
