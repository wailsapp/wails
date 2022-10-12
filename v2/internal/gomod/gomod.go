package gomod

import (
	"fmt"

	"github.com/Masterminds/semver"
	"golang.org/x/mod/modfile"
)

func GetWailsVersionFromModFile(goModText []byte) (*semver.Version, error) {
	file, err := modfile.Parse("", goModText, nil)
	if err != nil {
		return nil, err
	}

	for _, req := range file.Require {
		if req.Syntax == nil {
			continue
		}
		tokenPosition := 0
		if !req.Syntax.InBlock {
			tokenPosition = 1
		}
		if req.Syntax.Token[tokenPosition] == "github.com/wailsapp/wails/v2" {
			version := req.Syntax.Token[tokenPosition+1]
			return semver.NewVersion(version)
		}
	}

	return nil, nil
}

func GoModOutOfSync(goModData []byte, currentVersion string) (bool, error) {
	gomodversion, err := GetWailsVersionFromModFile(goModData)
	if err != nil {
		return false, err
	}
	if gomodversion == nil {
		return false, fmt.Errorf("Unable to find Wails in go.mod")
	}

	result, err := semver.NewVersion(currentVersion)
	if err != nil || result == nil {
		return false, fmt.Errorf("Unable to parse Wails version: %s", currentVersion)
	}

	return !gomodversion.Equal(result), nil
}

func UpdateGoModVersion(goModText []byte, currentVersion string) ([]byte, error) {
	file, err := modfile.Parse("", goModText, nil)
	if err != nil {
		return nil, err
	}

	err = file.AddRequire("github.com/wailsapp/wails/v2", currentVersion)
	if err != nil {
		return nil, err
	}

	// Replace
	if len(file.Replace) == 0 {
		return file.Format()
	}

	for _, req := range file.Replace {
		if req.Syntax == nil {
			continue
		}
		tokenPosition := 0
		if !req.Syntax.InBlock {
			tokenPosition = 1
		}
		if req.Syntax.Token[tokenPosition] == "github.com/wailsapp/wails/v2" {
			version := req.Syntax.Token[tokenPosition+1]
			_, err := semver.NewVersion(version)
			if err == nil {
				req.Syntax.Token[tokenPosition+1] = currentVersion
			}
		}
	}

	return file.Format()
}

func SyncGoVersion(goModText []byte, goVersion string) ([]byte, bool, error) {
	file, err := modfile.Parse("", goModText, nil)
	if err != nil {
		return nil, false, err
	}

	modVersion, err := semver.NewVersion(file.Go.Version)
	if err != nil {
		return nil, false, fmt.Errorf("Unable to parse Go version from go mod file: %s", err)
	}

	targetVersion, err := semver.NewVersion(goVersion)
	if err != nil {
		return nil, false, fmt.Errorf("Unable to parse Go version: %s", targetVersion)
	}

	if !targetVersion.GreaterThan(modVersion) {
		return goModText, false, nil
	}

	file.Go.Version = goVersion
	file.Go.Syntax.Token[1] = goVersion
	goModText, err = file.Format()
	if err != nil {
		return nil, false, err
	}

	return goModText, true, nil
}
