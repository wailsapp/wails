package debug

import (
	"os"
	"path/filepath"
	"runtime"
)

var LocalModulePath = ""

func init() {
	// Check if .git exists in the relative directory from here: ../../..
	// If it does, we are in a local build
	gitDir := RelativePath("..", "..", "..", ".git")
	if _, err := os.Stat(gitDir); err == nil {
		modulePath := RelativePath("..", "..", "..")
		LocalModulePath, _ = filepath.Abs(modulePath)
	}
}

// RelativePath returns a qualified path created by joining the
// directory of the calling file and the given relative path.
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
