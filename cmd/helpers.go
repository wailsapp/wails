package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/leaanthony/slicer"
	"github.com/leaanthony/spinner"
	wailsruntime "github.com/wailsapp/wails/runtime"
)

const xgoVersion = "1.16.2"

var fs = NewFSHelper()

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

// InstallGoDependencies will run go get in the current directory
func InstallGoDependencies(verbose bool) error {
	var depSpinner *spinner.Spinner
	if !verbose {
		depSpinner = spinner.New("Ensuring Dependencies are up to date...")
		depSpinner.SetSpinSpeed(50)
		depSpinner.Start()
	}
	err := NewProgramHelper(verbose).RunCommand("go get")
	if err != nil {
		if !verbose {
			depSpinner.Error()
		}
		return err
	}
	if !verbose {
		depSpinner.Success()
	}
	return nil
}

func InitializeCrossCompilation(verbose bool) error {
	// Check Docker
	if err := CheckIfInstalled("docker"); err != nil {
		return err
	}

	var packSpinner *spinner.Spinner
	msg := fmt.Sprintf("Pulling wailsapp/xgo:%s docker image... (may take a while)", xgoVersion)
	if !verbose {
		packSpinner = spinner.New(msg)
		packSpinner.SetSpinSpeed(50)
		packSpinner.Start()
	} else {
		println(msg)
	}

	err := NewProgramHelper(verbose).RunCommandArray([]string{"docker",
		"pull", fmt.Sprintf("wailsapp/xgo:%s", xgoVersion)})

	if err != nil {
		if packSpinner != nil {
			packSpinner.Error()
		}
		return err
	}
	if packSpinner != nil {
		packSpinner.Success()
	}

	return nil
}

// BuildDocker builds the project using the cross compiling wailsapp/xgo:<xgoVersion> container
func BuildDocker(binaryName string, buildMode string, projectOptions *ProjectOptions) error {
	var packSpinner *spinner.Spinner
	if buildMode == BuildModeBridge {
		return fmt.Errorf("you cant serve the application in cross-compilation")
	}

	// Check build directory
	buildDirectory := filepath.Join(fs.Cwd(), "build")
	if !fs.DirExists(buildDirectory) {
		err := fs.MkDir(buildDirectory)
		if err != nil {
			return err
		}
	}

	buildCommand := slicer.String()
	userid := 1000
	user, _ := user.Current()
	if i, err := strconv.Atoi(user.Uid); err == nil {
		userid = i
	}
	for _, arg := range []string{
		"docker",
		"run",
		"--rm",
		"-v", fmt.Sprintf("%s:/build", filepath.Join(fs.Cwd(), "build")),
		"-v", fmt.Sprintf("%s:/source", fs.Cwd()),
		"-e", fmt.Sprintf("LOCAL_USER_ID=%v", userid),
		"-e", fmt.Sprintf("FLAG_TAGS=%s", projectOptions.Tags),
		"-e", fmt.Sprintf("FLAG_LDFLAGS=%s", ldFlags(projectOptions, buildMode)),
		"-e", "FLAG_V=false",
		"-e", "FLAG_X=false",
		"-e", "FLAG_RACE=false",
		"-e", "FLAG_BUILDMODE=default",
		"-e", "FLAG_TRIMPATH=false",
		"-e", fmt.Sprintf("TARGETS=%s/%s", projectOptions.Platform, projectOptions.Architecture),
		"-e", "GOPROXY=",
		"-e", "GO111MODULE=on",
	} {
		buildCommand.Add(arg)
	}

	if projectOptions.GoPath != "" {
		buildCommand.Add("-v")
		buildCommand.Add(fmt.Sprintf("%s:/go", projectOptions.GoPath))
	}

	buildCommand.Add(fmt.Sprintf("wailsapp/xgo:%s", xgoVersion))
	buildCommand.Add(".")

	compileMessage := fmt.Sprintf(
		"Packing + Compiling project for %s/%s using docker image wailsapp/xgo:%s",
		projectOptions.Platform, projectOptions.Architecture, xgoVersion)

	if buildMode == BuildModeDebug {
		compileMessage += " (Debug Mode)"
	}

	if !projectOptions.Verbose {
		packSpinner = spinner.New(compileMessage + "...")
		packSpinner.SetSpinSpeed(50)
		packSpinner.Start()
	} else {
		println(compileMessage)
	}

	err := NewProgramHelper(projectOptions.Verbose).RunCommandArray(buildCommand.AsSlice())
	if err != nil {
		if packSpinner != nil {
			packSpinner.Error()
		}
		return err
	}
	if packSpinner != nil {
		packSpinner.Success()
	}

	return nil
}

// BuildNative builds on the target platform itself.
func BuildNative(binaryName string, forceRebuild bool, buildMode string, projectOptions *ProjectOptions) error {

	if err := CheckWindres(); err != nil {
		return err
	}

	compileMessage := "Packing + Compiling project"

	if buildMode == BuildModeDebug {
		compileMessage += " (Debug Mode)"
	}

	var packSpinner *spinner.Spinner
	if !projectOptions.Verbose {
		packSpinner = spinner.New(compileMessage + "...")
		packSpinner.SetSpinSpeed(50)
		packSpinner.Start()
	} else {
		println(compileMessage)
	}

	buildCommand := slicer.String()
	buildCommand.Add("go")

	buildCommand.Add("build")

	if binaryName != "" {
		// Alter binary name based on OS
		switch projectOptions.Platform {
		case "windows":
			if !strings.HasSuffix(binaryName, ".exe") {
				binaryName += ".exe"
			}
		default:
			if strings.HasSuffix(binaryName, ".exe") {
				binaryName = strings.TrimSuffix(binaryName, ".exe")
			}
		}
		buildCommand.Add("-o", filepath.Join("build", binaryName))
	}

	// If we are forcing a rebuild
	if forceRebuild {
		buildCommand.Add("-a")
	}

	buildCommand.AddSlice([]string{"-ldflags", ldFlags(projectOptions, buildMode)})

	if projectOptions.Tags != "" {
		buildCommand.AddSlice([]string{"--tags", projectOptions.Tags})
	}

	if projectOptions.Verbose {
		fmt.Printf("Command: %v\n", buildCommand.AsSlice())
	}

	err := NewProgramHelper(projectOptions.Verbose).RunCommandArray(buildCommand.AsSlice())
	if err != nil {
		if packSpinner != nil {
			packSpinner.Error()
		}
		return err
	}
	if packSpinner != nil {
		packSpinner.Success()
	}

	return nil
}

// BuildApplication will attempt to build the project based on the given inputs
func BuildApplication(binaryName string, forceRebuild bool, buildMode string, packageApp bool, projectOptions *ProjectOptions) error {
	var err error

	if projectOptions.CrossCompile {
		if err := InitializeCrossCompilation(projectOptions.Verbose); err != nil {
			return err
		}
	}

	helper := NewPackageHelper(projectOptions.Platform)

	// Generate windows resources
	if projectOptions.Platform == "windows" {
		if err := helper.PackageWindows(projectOptions, false); err != nil {
			return err
		}
	}

	if projectOptions.CrossCompile {
		err = BuildDocker(binaryName, buildMode, projectOptions)
	} else {
		err = BuildNative(binaryName, forceRebuild, buildMode, projectOptions)
	}
	if err != nil {
		return err
	}

	if packageApp {
		err = PackageApplication(projectOptions)
		if err != nil {
			return err
		}
	}

	return nil
}

// PackageApplication will attempt to package the application in a platform dependent way
func PackageApplication(projectOptions *ProjectOptions) error {
	var packageSpinner *spinner.Spinner
	if projectOptions.Verbose {
		packageSpinner = spinner.New("Packaging application...")
		packageSpinner.SetSpinSpeed(50)
		packageSpinner.Start()
	}

	err := NewPackageHelper(projectOptions.Platform).Package(projectOptions)
	if err != nil {
		if packageSpinner != nil {
			packageSpinner.Error()
		}
		return err
	}
	if packageSpinner != nil {
		packageSpinner.Success()
	}
	return nil
}

// BuildFrontend runs the given build command
func BuildFrontend(projectOptions *ProjectOptions) error {
	var buildFESpinner *spinner.Spinner
	if !projectOptions.Verbose {
		buildFESpinner = spinner.New("Building frontend...")
		buildFESpinner.SetSpinSpeed(50)
		buildFESpinner.Start()
	} else {
		println("Building frontend...")
	}
	err := NewProgramHelper(projectOptions.Verbose).RunCommand(projectOptions.FrontEnd.Build)
	if err != nil {
		if buildFESpinner != nil {
			buildFESpinner.Error()
		}
		return err
	}
	if buildFESpinner != nil {
		buildFESpinner.Success()
	}
	return nil
}

// CheckWindres checks if Windres is installed and if not, aborts
func CheckWindres() (err error) {
	if runtime.GOOS != "windows" { // FIXME: Handle windows cross-compile for windows!
		return nil
	}
	programHelper := NewProgramHelper()
	if !programHelper.IsInstalled("windres") {
		return fmt.Errorf("windres not installed. It comes by default with mingw. Ensure you have installed mingw correctly")
	}
	return nil
}

// CheckIfInstalled returns if application is installed
func CheckIfInstalled(application string) (err error) {
	programHelper := NewProgramHelper()
	if !programHelper.IsInstalled(application) {
		return fmt.Errorf("%s not installed. Ensure you have installed %s correctly", application, application)
	}
	return nil
}

// InstallFrontendDeps attempts to install the frontend dependencies based on the given options
func InstallFrontendDeps(projectDir string, projectOptions *ProjectOptions, forceRebuild bool, caller string) error {

	// Install frontend deps
	err := os.Chdir(projectOptions.FrontEnd.Dir)
	if err != nil {
		return err
	}

	// Check if frontend deps have been updated
	var feSpinner *spinner.Spinner
	if !projectOptions.Verbose {
		feSpinner = spinner.New("Ensuring frontend dependencies are up to date (This may take a while)")
		feSpinner.SetSpinSpeed(50)
		feSpinner.Start()
	} else {
		println("Ensuring frontend dependencies are up to date (This may take a while)")
	}

	requiresNPMInstall := true

	// Read in package.json MD5
	fs := NewFSHelper()
	packageJSONMD5, err := fs.FileMD5("package.json")
	if err != nil {
		return err
	}

	const md5sumFile = "package.json.md5"

	// If node_modules does not exist, force a rebuild.
	nodeModulesPath, err := filepath.Abs(filepath.Join(".", "node_modules"))
	if err != nil {
		return err
	}
	if !fs.DirExists(nodeModulesPath) {
		forceRebuild = true
	}

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
				if feSpinner != nil {
					feSpinner.Success("Skipped frontend dependencies (-f to force rebuild)")
				} else {
					println("Skipped frontend dependencies (-f to force rebuild)")
				}
			}
		}
	}

	// Md5 sum package.json
	// Different? Build
	if requiresNPMInstall || forceRebuild {
		// Install dependencies
		err = NewProgramHelper(projectOptions.Verbose).RunCommand(projectOptions.FrontEnd.Install)
		if err != nil {
			if feSpinner != nil {
				feSpinner.Error()
			}
			return err
		}
		if feSpinner != nil {
			feSpinner.Success()
		}

		// Update md5sum file
		err := os.WriteFile(md5sumFile, []byte(packageJSONMD5), 0644)
		if err != nil {
			return err
		}
	}

	// Install the runtime
	if caller == "build" {
		err = InstallProdRuntime(projectDir, projectOptions)
	} else {
		err = InstallBridge(projectDir, projectOptions)
	}
	if err != nil {
		return err
	}

	// Build frontend
	err = BuildFrontend(projectOptions)
	if err != nil {
		return err
	}
	return nil
}

// InstallBridge installs the relevant bridge javascript library
func InstallBridge(projectDir string, projectOptions *ProjectOptions) error {
	bridgeFileTarget := filepath.Join(projectDir, projectOptions.FrontEnd.Dir, "node_modules", "@wailsapp", "runtime", "init.js")
	err := fs.CreateFile(bridgeFileTarget, wailsruntime.BridgeJS)
	return err
}

// InstallProdRuntime installs the production runtime
func InstallProdRuntime(projectDir string, projectOptions *ProjectOptions) error {
	bridgeFileTarget := filepath.Join(projectDir, projectOptions.FrontEnd.Dir, "node_modules", "@wailsapp", "runtime", "init.js")
	err := fs.CreateFile(bridgeFileTarget, wailsruntime.InitJS)
	return err
}

// ServeProject attempts to serve up the current project so that it may be connected to
// via the Wails bridge
func ServeProject(projectOptions *ProjectOptions, logger *Logger) error {
	go func() {
		time.Sleep(2 * time.Second)
		if projectOptions.Platform == "windows" {
			logger.Yellow("*** Please note: Windows builds use mshtml which is only compatible with IE11. We strongly recommend only using IE11 when running 'wails serve'! For more information, please read https://wails.app/guides/windows/ ***")
		}
		logger.Green(">>>>> To connect, you will need to run '" + projectOptions.FrontEnd.Serve + "' in the '" + projectOptions.FrontEnd.Dir + "' directory <<<<<")
	}()
	location, err := filepath.Abs(filepath.Join("build", projectOptions.BinaryName))
	if err != nil {
		return err
	}

	logger.Yellow("Serving Application: " + location)
	cmd := exec.Command(location)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func ldFlags(po *ProjectOptions, buildMode string) string {
	// Setup ld flags
	ldflags := "-w -s "
	if buildMode == BuildModeDebug {
		ldflags = ""
	}

	// Add windows flags
	if po.Platform == "windows" && buildMode == BuildModeProd {
		ldflags += "-H windowsgui "
	}

	if po.UseFirebug {
		ldflags += "-X github.com/wailsapp/wails/lib/renderer.UseFirebug=true "
	}

	ldflags += "-X github.com/wailsapp/wails.BuildMode=" + buildMode

	// Add additional ldflags passed in via the `ldflags` cli flag
	if len(po.LdFlags) > 0 {
		ldflags += " " + po.LdFlags
	}

	// If we wish to generate typescript
	if po.typescriptDefsFilename != "" {
		cwd, err := os.Getwd()
		if err == nil {
			filename := filepath.Join(cwd, po.FrontEnd.Dir, po.typescriptDefsFilename)
			ldflags += " -X github.com/wailsapp/wails/lib/binding.typescriptDefinitionFilename=" + filename
		}
	}
	return ldflags
}

func getGitConfigValue(key string) (string, error) {
	output, err := exec.Command("git", "config", "--get", "--null", key).Output()
	// When using --null git appends a null character (\u0000) to the command output
	return strings.TrimRight(string(output), "\u0000"), err
}
