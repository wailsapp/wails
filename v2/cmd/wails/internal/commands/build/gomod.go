package build

import (
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/cmd/wails/internal"
	"github.com/wailsapp/wails/v2/internal/gomod"
	"github.com/wailsapp/wails/v2/internal/goversion"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
)

func SyncGoMod(logger *clilogger.CLILogger, updateWailsVersion bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	gomodFilename := filepath.Join(cwd, "go.mod")
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

	if outOfSync, err := gomod.GoModOutOfSync(gomodData, internal.Version); err != nil {
		return err
	} else if outOfSync {
		if updateWailsVersion {
			LogGreen("Updating go.mod to use Wails '%s'", internal.Version)
			gomodData, err = gomod.UpdateGoModVersion(gomodData, internal.Version)
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
		return os.WriteFile(gomodFilename, gomodData, 0755)
	}

	return nil
}
