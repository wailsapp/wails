package build

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/shell"
)

func deleteLinuxPackFiles(appDirBase string) {
	// Delete appdir
	appDir := filepath.Join(appDirBase, "AppDir")
	os.RemoveAll(appDir)
}

func packageApplication(options *Options) error {

	// Check we have AppImage tools

	// Create AppImage build directory
	buildDirectory, err := getApplicationBuildDirectory(options, "linux")
	if err != nil {
		return err
	}

	defer deleteLinuxPackFiles(buildDirectory)

	// Get the name of the application and ensure we lower+kebab case it
	name := filepath.Base(options.ProjectData.OutputFilename)

	// Calculate asset directory
	assetDir := getPackageAssetsDirectory()

	// Copy default icon if one doesn't exist
	baseBuildDirectory, err := getBuildBaseDirectory(options)
	if err != nil {
		return err
	}
	iconFile := filepath.Join(baseBuildDirectory, "icon.png")
	if !fs.FileExists(iconFile) {
		err = fs.CopyFile(defaultIconPath(), iconFile)
		if err != nil {
			return err
		}
	}

	// Copy Icon
	targetIcon := filepath.Join(buildDirectory, name+".png")
	err = fs.CopyFile(iconFile, targetIcon)
	if err != nil {
		return err
	}

	// Copy app.desktop
	dotDesktopFile := filepath.Join(baseBuildDirectory, "linux", name+".desktop")
	if !fs.FileExists(dotDesktopFile) {
		bytes, err := ioutil.ReadFile(filepath.Join(assetDir, "app.desktop"))
		if err != nil {
			return err
		}
		appDesktop := string(bytes)
		appDesktop = strings.ReplaceAll(appDesktop, `{{.Name}}`, name)
		err = ioutil.WriteFile(dotDesktopFile, []byte(appDesktop), 0644)
		if err != nil {
			return err
		}
	}

	// Copy AppRun file
	// targetFilename = filepath.Join(buildDirectory, "AppRun")
	// if !fs.FileExists(targetFilename) {
	// 	bytes, err := ioutil.ReadFile(filepath.Join(assetDir, "AppRun"))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	appRun := string(bytes)
	// 	appRun = strings.ReplaceAll(appRun, `{{.OutputFilename}}`, name)

	// 	err = ioutil.WriteFile(targetFilename, []byte(appRun), 0644)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// Copy Binary
	sourceFile := filepath.Join(options.ProjectData.Path, options.ProjectData.OutputFilename)
	targetFile := filepath.Join(buildDirectory, options.ProjectData.OutputFilename)
	err = fs.CopyFile(sourceFile, targetFile)
	if err != nil {
		return err
	}
	err = os.Chmod(targetFile, 0777)
	if err != nil {
		return err
	}

	/** Pack App **/

	// Make file executable
	// Set environment variable: OUTPUT=outputfilename
	command := shell.NewCommand("linuxdeploy-x86_64.AppImage")
	command.Dir(buildDirectory)

	argslice := slicer.String()
	argslice.Add("--appdir", "AppDir")
	argslice.Add("-d", filepath.Join("..", name+".desktop"))
	argslice.Add("-i", name+".png")
	argslice.Add("-e", name)
	argslice.Add("--output", "appimage")
	command.AddArgs(argslice.AsSlice())

	command.Env("OUTPUT", name+".AppImage")

	err = command.Run()
	if err != nil {
		println(command.Stdout())
		println(command.Stderr())
		return err
	}

	// Copy app to project dir

	println(buildDirectory)

	return nil
}

func deleteDirectory(directory string) {
	os.RemoveAll(directory)
}
