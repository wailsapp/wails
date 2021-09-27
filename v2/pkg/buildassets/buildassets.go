package buildassets

import (
	"embed"
	"github.com/leaanthony/debme"
	"github.com/leaanthony/gosod"
	"os"
	"path/filepath"
)

//go:embed build
var assets embed.FS

type assetData struct {
	Name string
}

// Install will install all default project assets
func Install(targetDir string, projectName string) error {
	templateDir := gosod.New(assets)
	err := templateDir.Extract(targetDir, &assetData{Name: projectName})
	if err != nil {
		return err
	}

	// Rename the manifest file
	windowsDir := filepath.Join(targetDir, "build", "windows")
	manifest := filepath.Join(windowsDir, "wails.exe.manifest")
	targetFile := filepath.Join(windowsDir, projectName+".exe.manifest")
	err = os.Rename(manifest, targetFile)
	if err != nil {
		return err
	}

	return nil
}

func RegenerateManifest(target string) error {
	a, err := debme.FS(assets, "build")
	if err != nil {
		return err
	}
	return a.CopyFile("windows/wails.exe.manifest", target, 0644)
}

func RegenerateAppIcon(target string) error {
	a, err := debme.FS(assets, "build")
	if err != nil {
		return err
	}
	return a.CopyFile("appicon.png", target, 0644)
}
