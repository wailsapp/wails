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

	return file.Format()
}
