package build

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/internal/assetdb"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/html"
	"github.com/wailsapp/wails/v2/internal/project"
	"github.com/wailsapp/wails/v2/internal/shell"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
)

// BaseBuilder is the common builder struct
type BaseBuilder struct {
	filesToDelete slicer.StringSlicer
	projectData   *project.Project
	options       *Options
}

// NewBaseBuilder creates a new BaseBuilder
func NewBaseBuilder(options *Options) *BaseBuilder {
	result := &BaseBuilder{
		options: options,
	}
	return result
}

// SetProjectData sets the project data for this builder
func (b *BaseBuilder) SetProjectData(projectData *project.Project) {
	b.projectData = projectData
}

func (b *BaseBuilder) addFileToDelete(filename string) {
	if !b.options.KeepAssets {
		b.filesToDelete.Add(filename)
	}
}

func (b *BaseBuilder) fileExists(path string) bool {
	// if file doesn't exist, ignore
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

// buildCustomAssets will iterate through the projects static directory and add all files
// to the application wide asset database.
func (b *BaseBuilder) buildCustomAssets(projectData *project.Project) error {

	// Add trailing slash to Asset directory
	customAssetsDir := filepath.Join(projectData.Path, "assets", "custom") + "/"
	if !b.fileExists(customAssetsDir) {
		err := fs.MkDirs(customAssetsDir)
		if err != nil {
			return err
		}
	}

	assets := assetdb.NewAssetDB()
	err := filepath.Walk(customAssetsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		normalisedPath := filepath.ToSlash(path)
		localPath := strings.TrimPrefix(normalisedPath, customAssetsDir)
		if len(localPath) == 0 {
			return nil
		}
		if data, err := ioutil.ReadFile(filepath.Join(customAssetsDir, localPath)); err == nil {
			assets.AddAsset(localPath, data)
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Write assetdb out to root directory
	assetsDbFilename := fs.RelativePath("../../../assetsdb.go")
	b.addFileToDelete(assetsDbFilename)
	err = ioutil.WriteFile(assetsDbFilename, []byte(assets.Serialize("assets", "wails")), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseBuilder) convertFileToIntegerString(filename string) (string, error) {
	rawData, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return b.convertByteSliceToIntegerString(rawData), nil
}

func (b *BaseBuilder) convertByteSliceToIntegerString(data []byte) string {

	// Create string builder
	var result strings.Builder

	if len(data) > 0 {

		// Loop over all but 1 bytes
		for i := 0; i < len(data)-1; i++ {
			result.WriteString(fmt.Sprintf("%v,", data[i]))
		}

		result.WriteString(fmt.Sprintf("%v", data[len(data)-1]))

	}

	return result.String()
}

// CleanUp does post-build housekeeping
func (b *BaseBuilder) CleanUp() {

	// Delete all the files
	b.filesToDelete.Each(func(filename string) {

		// if file doesn't exist, ignore
		if !b.fileExists(filename) {
			return
		}

		// Delete file. We ignore errors because these files will be overwritten
		// by the next build anyway.
		_ = os.Remove(filename)

	})
}

// CompileProject compiles the project
func (b *BaseBuilder) CompileProject(options *Options) error {

	// Default go build command
	commands := slicer.String([]string{"build"})

	// TODO: Work out if we can make this more efficient
	// We need to do a full build as CGO doesn't detect updates
	// to .h files, and we package assets into .h file. We could
	// potentially try and see if the assets have changed but will
	// this take as much time as a `-a` build?
	commands.Add("-a")

	var tags slicer.StringSlicer
	tags.Add(options.OutputType)

	if options.Mode == Debug {
		tags.Add("debug")
	}

	// Add the output type build tag
	commands.Add("-tags")
	commands.Add(tags.Join(","))

	// Strip binary in Production mode
	if options.Mode == Production {

		// Different linker flags depending on platform
		commands.Add("-ldflags")

		switch runtime.GOOS {
		case "windows":
			commands.Add("-w -s -H windowsgui")
		default:
			commands.Add("-w -s")
		}
	}

	// Add better debugging flags
	if options.Mode == Debug {
		commands.Add(`-gcflags=all="-N -l"`)
	}

	// Get application build directory
	appDir := options.BuildDirectory
	err := cleanBuildDirectory(options)
	if err != nil {
		return err
	}

	if options.LDFlags != "" {
		commands.Add("-ldflags")
		commands.Add(options.LDFlags)
	}

	// Set up output filename
	outputFile := options.OutputFile
	if outputFile == "" {
		outputFile = b.projectData.OutputFilename
	}
	compiledBinary := filepath.Join(appDir, outputFile)
	commands.Add("-o")
	commands.Add(compiledBinary)

	b.projectData.OutputFilename = strings.TrimPrefix(compiledBinary, options.ProjectData.Path)
	options.CompiledBinary = compiledBinary

	// Create the command
	cmd := exec.Command(options.Compiler, commands.AsSlice()...)

	// Set the directory
	cmd.Dir = b.projectData.Path

	// Set GO111MODULE environment variable
	cmd.Env = append(os.Environ(), "GO111MODULE=on")

	// Add CGO flags
	// We use the project/build dir as a temporary place for our generated c headers
	buildBaseDir, err := fs.RelativeToCwd("build")
	if err != nil {
		return err
	}

	cmd.Env = append(os.Environ(), fmt.Sprintf("CGO_CFLAGS=-I%s", buildBaseDir))

	// Setup buffers
	var stdo, stde bytes.Buffer
	cmd.Stdout = &stdo
	cmd.Stderr = &stde

	// Run command
	err = cmd.Run()

	// Format error if we have one
	if err != nil {
		return fmt.Errorf("%s\n%s", err, string(stde.Bytes()))
	}

	return nil
}

// NpmInstall runs "npm install" in the given directory
func (b *BaseBuilder) NpmInstall(sourceDir string) error {
	return b.NpmInstallUsingCommand(sourceDir, "npm install")
}

// NpmInstallUsingCommand runs the given install command in the specified npm project directory
func (b *BaseBuilder) NpmInstallUsingCommand(sourceDir string, installCommand string) error {

	packageJSON := filepath.Join(sourceDir, "package.json")

	// Check package.json exists
	if !fs.FileExists(packageJSON) {
		return fmt.Errorf("unable to load package.json at '%s'", packageJSON)
	}

	install := false

	// Get the MD5 sum of package.json
	packageJSONMD5 := fs.MustMD5File(packageJSON)

	// Check whether we need to npm install
	packageChecksumFile := filepath.Join(sourceDir, "package.json.md5")
	if fs.FileExists(packageChecksumFile) {
		// Compare checksums
		storedChecksum := fs.MustLoadString(packageChecksumFile)
		if storedChecksum != packageJSONMD5 {
			fs.MustWriteString(packageChecksumFile, packageJSONMD5)
			install = true
		}
	} else {
		install = true
		fs.MustWriteString(packageChecksumFile, packageJSONMD5)
	}

	// Install if node_modules doesn't exist
	nodeModulesDir := filepath.Join(sourceDir, "node_modules")
	if !fs.DirExists(nodeModulesDir) {
		install = true
	}

	// Shortcut installation
	if install == false {
		return nil
	}

	// Split up the InstallCommand and execute it
	cmd := strings.Split(installCommand, " ")
	stdout, stderr, err := shell.RunCommand(sourceDir, cmd[0], cmd[1:]...)
	if err != nil {
		for _, l := range strings.Split(stdout, "\n") {
			fmt.Printf("    %s\n", l)
		}
		for _, l := range strings.Split(stderr, "\n") {
			fmt.Printf("    %s\n", l)
		}
	}

	return err
}

// NpmRun executes the npm target in the provided directory
func (b *BaseBuilder) NpmRun(projectDir, buildTarget string, verbose bool) error {
	stdout, stderr, err := shell.RunCommand(projectDir, "npm", "run", buildTarget)
	if verbose || err != nil {
		for _, l := range strings.Split(stdout, "\n") {
			fmt.Printf("    %s\n", l)
		}
		for _, l := range strings.Split(stderr, "\n") {
			fmt.Printf("    %s\n", l)
		}
	}
	return err
}

// NpmRunWithEnvironment executes the npm target in the provided directory, with the given environment variables
func (b *BaseBuilder) NpmRunWithEnvironment(projectDir, buildTarget string, verbose bool, envvars []string) error {
	cmd := shell.CreateCommand(projectDir, "npm", "run", buildTarget)
	cmd.Env = append(os.Environ(), envvars...)
	var stdo, stde bytes.Buffer
	cmd.Stdout = &stdo
	cmd.Stderr = &stde
	err := cmd.Run()
	if verbose || err != nil {
		for _, l := range strings.Split(stdo.String(), "\n") {
			fmt.Printf("    %s\n", l)
		}
		for _, l := range strings.Split(stde.String(), "\n") {
			fmt.Printf("    %s\n", l)
		}
	}
	return err
}

// BuildFrontend executes the `npm build` command for the frontend directory
func (b *BaseBuilder) BuildFrontend(outputLogger *clilogger.CLILogger) error {

	// TODO: Fix this up from the CLI
	verbose := false

	frontendDir := filepath.Join(b.projectData.Path, "frontend")

	// Check there is an 'InstallCommand' provided in wails.json
	if b.projectData.InstallCommand == "" {
		// No - don't install
		outputLogger.Println("    - No Install command. Skipping.")
	} else {
		// Do install if needed
		outputLogger.Println("    - Installing dependencies...")
		if err := b.NpmInstallUsingCommand(frontendDir, b.projectData.InstallCommand); err != nil {
			return err
		}
	}

	// Check if there is a build command
	if b.projectData.BuildCommand == "" {
		outputLogger.Println("    - No Build command. Skipping.")
		// No - ignore
		return nil
	}

	outputLogger.Println("    - Compiling Frontend Project")
	cmd := strings.Split(b.projectData.BuildCommand, " ")
	stdout, stderr, err := shell.RunCommand(frontendDir, cmd[0], cmd[1:]...)
	if verbose || err != nil {
		for _, l := range strings.Split(stdout, "\n") {
			fmt.Printf("    %s\n", l)
		}
		for _, l := range strings.Split(stderr, "\n") {
			fmt.Printf("    %s\n", l)
		}
	}
	return err
}

// ExtractAssets gets the assets from the index.html file
func (b *BaseBuilder) ExtractAssets() (*html.AssetBundle, error) {

	// Read in html
	return html.NewAssetBundle(b.projectData.HTML)
}
