package debug

import (
	"github.com/samber/lo"
	"path/filepath"
	"runtime"
)

import "runtime/debug"

// Why go doesn't provide this as a map already is beyond me.
var buildSettings = map[string]string{}
var LocalModulePath = ""

func init() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	buildSettings = lo.Associate(buildInfo.Settings, func(setting debug.BuildSetting) (string, string) {
		return setting.Key, setting.Value
	})
	if isLocalBuild() {
		modulePath := RelativePath("..", "..", "..")
		LocalModulePath, _ = filepath.Abs(modulePath)
	}
}

func isLocalBuild() bool {
	return buildSettings["vcs.modified"] == "true"
}

// RelativePath returns a qualified path created by joining the
// directory of the calling file and the given relative path.
//
// Example: RelativePath("..") in *this* file would give you '/path/to/wails2/v2/internal`
func RelativePath(relativepath string, optionalpaths ...string) string {
	_, thisFile, _, _ := runtime.Caller(1)
	localDir := filepath.Dir(thisFile)

	// If we have optional paths, join them to the relativepath
	if len(optionalpaths) > 0 {
		paths := []string{relativepath}
		paths = append(paths, optionalpaths...)
		relativepath = filepath.Join(paths...)
	}
	result, err := filepath.Abs(filepath.Join(localDir, relativepath))
	if err != nil {
		// I'm allowing this for 1 reason only: It's fatal if the path
		// supplied is wrong as it's only used internally in Wails. If we get
		// that path wrong, we should know about it immediately. The other reason is
		// that it cuts down a ton of unnecassary error handling.
		panic(err)
	}
	return result
}
