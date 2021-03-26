package build

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/wailsapp/wails/v2/internal/fs"

	"github.com/wailsapp/wails/v2/internal/shell"

	"github.com/wailsapp/wails/v2/internal/project"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
)

// Mode is the type used to indicate the build modes
type Mode int

const (
	// Debug mode
	Debug Mode = iota
	// Production mode
	Production
)

var modeMap = []string{"Debug", "Production"}

// Options contains all the build options as well as the project data
type Options struct {
	LDFlags             string               // Optional flags to pass to linker
	Logger              *clilogger.CLILogger // All output to the logger
	OutputType          string               // EG: desktop, server....
	Mode                Mode                 // release or debug
	ProjectData         *project.Project     // The project data
	Pack                bool                 // Create a package for the app after building
	Platform            string               // The platform to build for
	Arch                string               // The architecture to build for
	Compiler            string               // The compiler command to use
	IgnoreFrontend      bool                 // Indicates if the frontend does not need building
	OutputFile          string               // Override the output filename
	BuildDirectory      string               // Directory to use for building the application
	CleanBuildDirectory bool                 // Indicates if the build directory should be cleaned before building
	CompiledBinary      string               // Fully qualified path to the compiled binary
	KeepAssets          bool                 // /Keep the generated assets/files
	Verbosity           int                  // Verbosity level (0 - silent, 1 - default, 2 - verbose)
	AppleIdentity       string
}

// GetModeAsString returns the current mode as a string
func GetModeAsString(mode Mode) string {
	return modeMap[mode]
}

// Build the project!
func Build(options *Options) (string, error) {

	// Extract logger
	outputLogger := options.Logger

	// Get working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Load project
	projectData, err := project.Load(cwd)
	if err != nil {
		return "", err
	}
	options.ProjectData = projectData

	// Set build directory
	options.BuildDirectory = filepath.Join(options.ProjectData.Path, "build", options.Platform, options.OutputType)

	// Save the project type
	projectData.OutputType = options.OutputType

	// Create builder
	var builder Builder

	switch projectData.OutputType {
	case "desktop":
		builder = newDesktopBuilder(options)
	case "hybrid":
		builder = newHybridBuilder(options)
	case "server":
		builder = newServerBuilder(options)
	case "dev":
		builder = newDesktopBuilder(options)
	default:
		return "", fmt.Errorf("cannot build assets for output type %s", projectData.OutputType)
	}

	// Set up our clean up method
	defer builder.CleanUp()

	// Initialise Builder
	builder.SetProjectData(projectData)

	// Generate Frontend JS Package
	// outputLogger.Println("  - Generating Backend JS Package")
	// // Ignore the parser report coming back
	// _, err = parser.GenerateWailsFrontendPackage()
	// if err != nil {
	// 	return "", err
	// }
	if !options.IgnoreFrontend {
		err = builder.BuildFrontend(outputLogger)
		if err != nil {
			return "", err
		}
	}

	// Build the base assets
	err = builder.BuildAssets(options)
	if err != nil {
		return "", err
	}

	// Compile the application
	outputLogger.Print("Compiling application: ")

	if options.Platform == "darwin" && options.Arch == "universal" {
		outputFile := builder.OutputFilename(options)
		amd64Filename := outputFile + "-amd64"
		arm64Filename := outputFile + "-arm64"

		// Build amd64 first
		options.Arch = "amd64"
		options.OutputFile = amd64Filename
		err = builder.CompileProject(options)
		if err != nil {
			return "", err
		}
		// Build arm64
		options.Arch = "arm64"
		options.OutputFile = arm64Filename
		err = builder.CompileProject(options)
		if err != nil {
			return "", err
		}
		// Run lipo
		_, stderr, err := shell.RunCommand(options.BuildDirectory, "lipo", "-create", "-output", outputFile, amd64Filename, arm64Filename)
		if err != nil {
			return "", fmt.Errorf("%s - %s", err.Error(), stderr)
		}
		// Remove temp binaries
		fs.DeleteFile(filepath.Join(options.BuildDirectory, amd64Filename))
		fs.DeleteFile(filepath.Join(options.BuildDirectory, arm64Filename))
		projectData.OutputFilename = outputFile
		options.CompiledBinary = filepath.Join(options.BuildDirectory, outputFile)
	} else {
		err = builder.CompileProject(options)
		if err != nil {
			return "", err
		}
	}
	outputLogger.Println("Done.")

	// Do we need to pack the app?
	if options.Pack {

		outputLogger.Print("Packaging application: ")

		// TODO: Allow cross platform build
		err = packageProject(options, runtime.GOOS)
		if err != nil {
			return "", err
		}
		outputLogger.Println("Done.")
	}

	return projectData.OutputFilename, nil

}
