package buildinfo

import (
	"fmt"
	"runtime/debug"

	wdebug "github.com/wailsapp/wails/v3/internal/debug"
	"github.com/wailsapp/wails/v3/internal/lo"
)

type Info struct {
	// Development is true when the running wails3 binary was built from a
	// local Wails source tree that is still resolvable on disk at runtime,
	// i.e. when [wdebug.LocalModulePath] resolved to a real directory.
	//
	// Earlier versions of this field derived Development from the
	// `vcs=git` build setting that Go's toolchain embeds automatically.
	// That signal is unsafe to rely on for "is this a local dev build?":
	// Go emits `vcs=git` for any binary built inside a git checkout,
	// including release artefacts produced by CI. A user running a
	// tagged `wails3.exe` on Windows would therefore see Development=true
	// even though the Wails source tree is nowhere on their machine.
	// The downstream effect was a malformed `replace github.com/wailsapp/wails/v3
	// => /v3` directive in scaffolded projects, breaking `wails3 init`
	// outright.
	//
	// [wdebug.LocalModulePath] is set by the debug package's init() only
	// when `.git` is found at a path resolved relative to the live
	// source file location — i.e. only when the dev's checkout is
	// actually there to be replaced into. That's exactly the precondition
	// the only consumer of this field (templates.Install) needs.
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
	result.Development = wdebug.LocalModulePath != ""

	return &result, nil

}
