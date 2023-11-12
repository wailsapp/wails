package gomod

import (
	"fmt"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/cmd/wails/internal"
	"github.com/wailsapp/wails/v2/internal/colour"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/gomod"
	"github.com/wailsapp/wails/v2/internal/goversion"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
)

func SyncGoMod(logger *clilogger.CLILogger, updateWailsVersion bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	gomodFilename := fs.FindFileInParents(cwd, "go.mod")
	if gomodFilename == "" {
		return fmt.Errorf("no go.mod file found")
	}
	gomodData, err := os.ReadFile(gomodFilename)
	if err != nil {
		return err
	}

	gomodData, updated, err := gomod.SyncGoVersion(gomodData, goversion.MinRequirement)
	if err != nil {
		return err
	} else if updated {
		LogGreen("Updated go.mod to use Go '%s'", goversion.MinRequirement)
	}

	internalVersion := strings.TrimSpace(internal.Version)
	if outOfSync, err := gomod.GoModOutOfSync(gomodData, internalVersion); err != nil {
		return err
	} else if outOfSync {
		if updateWailsVersion {
			LogGreen("Updating go.mod to use Wails '%s'", internalVersion)
			gomodData, err = gomod.UpdateGoModVersion(gomodData, internalVersion)
			if err != nil {
				return err
			}
			updated = true
		} else {
			gomodversion, err := gomod.GetWailsVersionFromModFile(gomodData)
			if err != nil {
				return err
			}

			logger.Println("Warning: go.mod is using Wails '%s' but the CLI is '%s'. Consider updating your project's `go.mod` file.\n", gomodversion.String(), internal.Version)
		}
	}

	if updated {
		return os.WriteFile(gomodFilename, gomodData, 0o755)
	}

	return nil
}

func LogGreen(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	println(colour.Green(text))
}
