//go:build !production

package application

import (
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/commands"
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
func newApplication(options *Options) *App {
	result := &App{
		isDebugMode: true,
		options:     options.getOptions(true),
	}
	result.init()
	return result
}

func (a *App) logStartup() {
	var args []any

	wailsPackage, _ := lo.Find(BuildInfo.Deps, func(dep *debug.Module) bool {
		return dep.Path == "github.com/wailsapp/wails/v3"
	})

	wailsVersion := commands.VersionString
	if wailsPackage != nil && wailsPackage.Replace != nil {
		wailsVersion = "(local) => " + filepath.ToSlash(wailsPackage.Replace.Path)
	}
	args = append(args, "Wails", wailsVersion)
	args = append(args, "Compiler", BuildInfo.GoVersion)

	for key, value := range BuildSettings {
		args = append(args, key, value)
	}

	a.info("Build Info:", args...)
}
