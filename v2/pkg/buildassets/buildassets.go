package buildassets

import (
	"embed"
	"github.com/leaanthony/debme"
	"github.com/leaanthony/gosod"
	"os"
	"path/filepath"
	"text/template"
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

	// Rename the rc file
	rc := filepath.Join(windowsDir, "wails.rc")
	targetFile = filepath.Join(windowsDir, projectName+".rc")
	err = os.Rename(rc, targetFile)
	if err != nil {
		return err
	}
	return nil
}

// RegenerateRCFile will recreate the RC file
func RegenerateRCFile(projectDir string, projectName string) error {
	targetFile, err := os.OpenFile(filepath.Join(projectDir, "build", "windows", projectName+".rc"), os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	data := &assetData{Name: projectName}
	templateData, err := assets.ReadFile("build/windows/wails.rc.tmpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("rc").Parse(string(templateData))
	if err != nil {
		return err
	}
	err = tmpl.Execute(targetFile, data)
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
