//go:build !production

package application

import (
	"github.com/wailsapp/wails/v3/internal/git"
	"github.com/wailsapp/wails/v3/internal/lo"
	"github.com/wailsapp/wails/v3/internal/version"
	"path/filepath"
	"runtime/debug"
)

// BuildSettings contains the build settings for the application
var BuildSettings map[string]string

// BuildInfo contains the build info for the application
var BuildInfo *debug.BuildInfo

func init() {
	var ok bool
	BuildInfo, ok = debug.ReadBuildInfo()
	if !ok {
		return
	}
	BuildSettings = lo.Associate(BuildInfo.Settings, func(setting debug.BuildSetting) (string, string) {
		return setting.Key, setting.Value
	})
}

// We use this to patch the application to production mode.
func newApplication(options Options) *App {
	result := &App{
		isDebugMode: true,
		options:     options,
	}
	result.init()
	return result
}

func (a *App) logStartup() {
	var args []any

	// BuildInfo is nil when build with garble
	if BuildInfo == nil {
		return
	}

	wailsPackage, _ := lo.Find(BuildInfo.Deps, func(dep *debug.Module) bool {
		return dep.Path == "github.com/wailsapp/wails/v3"
	})

	wailsVersion := version.String()
	if wailsPackage != nil && wailsPackage.Replace != nil {
		wailsVersion = "(local) => " + filepath.ToSlash(wailsPackage.Replace.Path)
		// Get the latest commit hash
		if hash, err := git.HeadHash(filepath.Join(wailsPackage.Replace.Path, "..")); err == nil {
			wailsVersion += " (" + hash + ")"
		}
	}
	args = append(args, "Wails", wailsVersion)
	args = append(args, "Compiler", BuildInfo.GoVersion)
	for key, value := range BuildSettings {
		args = append(args, key, value)
	}

	a.info("Build Info:", args...)
}
