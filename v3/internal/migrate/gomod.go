package migrate

import (
	"os"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

// minGoVersion is the minimum Go directive required by wails v3 projects
// (matches the version in the v3 project template).
const minGoVersion = "1.24"

// TransformGoMod rewrites the project's go.mod for v3: the wails/v2 require
// is dropped, wails/v3 is added at the given version and the go directive is
// raised to the v3 minimum if needed. All other requires are preserved.
func TransformGoMod(proj *V2Project, wailsVersion string) ([]byte, error) {
	data, err := os.ReadFile(proj.GoModPath)
	if err != nil {
		return nil, err
	}
	mod, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, err
	}

	if err := mod.DropRequire("github.com/wailsapp/wails/v2"); err != nil {
		return nil, err
	}
	if err := mod.AddRequire("github.com/wailsapp/wails/v3", wailsVersion); err != nil {
		return nil, err
	}

	if mod.Go == nil || semver.Compare("v"+mod.Go.Version, "v"+minGoVersion) < 0 {
		if err := mod.AddGoStmt(minGoVersion); err != nil {
			return nil, err
		}
	}

	// Drop any replace directives pointing at wails/v2.
	for _, rep := range mod.Replace {
		if rep.Old.Path == "github.com/wailsapp/wails/v2" {
			if err := mod.DropReplace(rep.Old.Path, rep.Old.Version); err != nil {
				return nil, err
			}
		}
	}

	mod.Cleanup()
	return mod.Format()
}
