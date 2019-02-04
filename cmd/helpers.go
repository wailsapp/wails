package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/leaanthony/slicer"
	"github.com/leaanthony/spinner"
)

// ValidateFrontendConfig checks if the frontend config is valid
func ValidateFrontendConfig(projectOptions *ProjectOptions) error {
	if projectOptions.FrontEnd.Dir == "" {
		return fmt.Errorf("Frontend directory not set in project.json")
	}
	if projectOptions.FrontEnd.Build == "" {
		return fmt.Errorf("Frontend build command not set in project.json")
	}
	if projectOptions.FrontEnd.Install == "" {
		return fmt.Errorf("Frontend install command not set in project.json")
	}
	if projectOptions.FrontEnd.Bridge == "" {
		return fmt.Errorf("Frontend bridge config not set in project.json")
	}

	return nil
}

func InstallGoDependencies() error {
	depSpinner := spinner.New("Installing Dependencies...")
	depSpinner.SetSpinSpeed(50)
	depSpinner.Start()
	err := NewProgramHelper().RunCommand("go get")
	if err != nil {
		depSpinner.Error()
		return err
	}
	depSpinner.Success()
	return nil
}

func BuildApplication(binaryName string, forceRebuild bool, buildMode string) error {
	compileMessage := "Packing + Compiling project"
	if buildMode == "debug" {
		compileMessage += " (Debug Mode)"
	}

	packSpinner := spinner.New(compileMessage + "...")
	packSpinner.SetSpinSpeed(50)
	packSpinner.Start()

	buildCommand := slicer.String()
	buildCommand.AddSlice([]string{"packr", "build"})

	if binaryName != "" {
		buildCommand.Add("-o")
		buildCommand.Add(binaryName)
	}

	// If we are forcing a rebuild
	if forceRebuild {
		buildCommand.Add("-a")
	}

	buildCommand.AddSlice([]string{"-ldflags", "-X github.com/wailsapp/wails.BuildMode=" + buildMode})
	err := NewProgramHelper().RunCommandArray(buildCommand.AsSlice())
	if err != nil {
		packSpinner.Error()
		return err
	}
	packSpinner.Success()
	return nil
}

func PackageApplication(projectOptions *ProjectOptions) error {
	// Package app
	packageSpinner := spinner.New("Packaging Application")
	packageSpinner.SetSpinSpeed(50)
	packageSpinner.Start()
	err := NewPackageHelper().Package(projectOptions)
	if err != nil {
		packageSpinner.Error()
		return err
	}
	packageSpinner.Success()
	return nil
}

func BuildFrontend(buildCommand string) error {
	buildFESpinner := spinner.New("Building frontend...")
	buildFESpinner.SetSpinSpeed(50)
	buildFESpinner.Start()
	err := NewProgramHelper().RunCommand(buildCommand)
	if err != nil {
		buildFESpinner.Error()
		return err
	}
	buildFESpinner.Success()
	return nil
}

func CheckPackr() (err error) {
	programHelper := NewProgramHelper()
	if !programHelper.IsInstalled("packr") {
		buildSpinner := spinner.New()
		buildSpinner.SetSpinSpeed(50)
		buildSpinner.Start("Installing packr...")
		err := programHelper.InstallGoPackage("github.com/gobuffalo/packr/...")
		if err != nil {
			buildSpinner.Error()
			return err
		}
		buildSpinner.Success()
	}
	return nil
}

func InstallFrontendDeps(projectDir string, projectOptions *ProjectOptions, forceRebuild bool) error {

	// Install frontend deps
	err := os.Chdir(projectOptions.FrontEnd.Dir)
	if err != nil {
		return err
	}

	// Check if frontend deps have been updated
	feSpinner := spinner.New("Installing frontend dependencies (This may take a while)...")
	feSpinner.SetSpinSpeed(50)
	feSpinner.Start()

	requiresNPMInstall := true

	// Read in package.json MD5
	fs := NewFSHelper()
	packageJSONMD5, err := fs.FileMD5("package.json")
	if err != nil {
		return err
	}

	const md5sumFile = "package.json.md5"

	// If we aren't forcing the install and the md5sum file exists
	if !forceRebuild && fs.FileExists(md5sumFile) {
		// Yes - read contents
		savedMD5sum, err := fs.LoadAsString(md5sumFile)
		// File exists
		if err == nil {
			// Compare md5
			if savedMD5sum == packageJSONMD5 {
				// Same - no need for reinstall
				requiresNPMInstall = false
				feSpinner.Success("Skipped frontend dependencies (-f to force rebuild)")
			}
		}
	}

	// Md5 sum package.json
	// Different? Build
	if requiresNPMInstall || forceRebuild {
		// Install dependencies
		err = NewProgramHelper().RunCommand(projectOptions.FrontEnd.Install)
		if err != nil {
			feSpinner.Error()
			return err
		}
		feSpinner.Success()

		// Update md5sum file
		ioutil.WriteFile(md5sumFile, []byte(packageJSONMD5), 0644)
	}

	bridgeFile := "wailsbridge.prod.js"

	// Copy bridge to project
	_, filename, _, _ := runtime.Caller(1)
	bridgeFileSource := filepath.Join(path.Dir(filename), "..", "..", "assets", "default", bridgeFile)
	bridgeFileTarget := filepath.Join(projectDir, projectOptions.FrontEnd.Dir, projectOptions.FrontEnd.Bridge, "wailsbridge.js")
	err = fs.CopyFile(bridgeFileSource, bridgeFileTarget)
	if err != nil {
		return err
	}

	// Build frontend
	err = BuildFrontend(projectOptions.FrontEnd.Build)
	if err != nil {
		return err
	}
	return nil
}

// func CopyBridgeFile(projectDir string, projectOptions ProjectOptions, bridgeMode bool) error {
// 	// Copy bridge to project
// 	fs := NewFSHelper()
// 	var bridgeFile = "wailsbridge.prod.js"
// 	if bridgeMode {
// 		bridgeFile = "wailsbridge.js"
// 	}
// 	_, filename, _, _ := runtime.Caller(1)
// 	bridgeFileSource := filepath.Join(path.Dir(filename), "..", "assets", "default", bridgeFile)
// 	bridgeFileTarget := filepath.Join(projectDir, projectOptions.FrontEnd.Dir, projectOptions.FrontEnd.Bridge, bridgeFile)
// 	err := fs.CopyFile(bridgeFileSource, bridgeFileTarget)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func InstallFrontend(projectOptions *ProjectOptions) error {
// 	// Install frontend deps
// 	err := os.Chdir(projectOptions.FrontEnd.Dir)
// 	if err != nil {
// 		return err
// 	}

// 	// Check if frontend deps have been updated
// 	feSpinner := spinner.New("Installing frontend dependencies (This may take a while)...")
// 	feSpinner.SetSpinSpeed(50)
// 	feSpinner.Start()

// 	requiresNPMInstall := true

// 	// Read in package.json MD5
// 	fs := NewFSHelper()
// 	packageJSONMD5, err := fs.FileMD5("package.json")
// 	if err != nil {
// 		return err
// 	}

// 	const md5sumFile = "package.json.md5"

// 	// If we aren't forcing the install and the md5sum file exists
// 	if !forceRebuild && fs.FileExists(md5sumFile) {
// 		// Yes - read contents
// 		savedMD5sum, err := fs.LoadAsString(md5sumFile)
// 		// File exists
// 		if err == nil {
// 			// Compare md5
// 			if savedMD5sum == packageJSONMD5 {
// 				// Same - no need for reinstall
// 				requiresNPMInstall = false
// 				feSpinner.Success("Skipped frontend dependencies (-f to force rebuild)")
// 			}
// 		}
// 	}

// 	// Md5 sum package.json
// 	// Different? Build
// 	if requiresNPMInstall || forceRebuild {
// 		// Install dependencies
// 		err = program.RunCommand(projectOptions.FrontEnd.Install)
// 		if err != nil {
// 			feSpinner.Error()
// 			return err
// 		}
// 		feSpinner.Success()

// 		// Update md5sum file
// 		ioutil.WriteFile(md5sumFile, []byte(packageJSONMD5), 0644)
// 	}
// 	return nil
// }
