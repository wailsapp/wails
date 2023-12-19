package build

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/shell"
	"github.com/wailsapp/wails/v2/internal/webview2runtime"
	"github.com/wailsapp/wails/v2/pkg/buildassets"
)

const (
	nsisTypeSingle   = "single"
	nsisTypeMultiple = "multiple"

	nsisFolder            = "windows/installer"
	nsisProjectFile       = "project.nsi"
	nsisToolsFile         = "wails_tools.nsh"
	nsisWebView2SetupFile = "tmp/MicrosoftEdgeWebview2Setup.exe"
)

func GenerateNSISInstaller(options *Options, amd64Binary string, arm64Binary string) error {
	outputLogger := options.Logger
	outputLogger.Println("Creating NSIS installer\n------------------------------")

	// Ensure the file exists, if not the template will be written.
	projectFile := path.Join(nsisFolder, nsisProjectFile)
	if _, err := buildassets.ReadFile(options.ProjectData, projectFile); err != nil {
		return fmt.Errorf("Unable to generate NSIS installer project template: %w", err)
	}

	// Write the resolved nsis tools
	toolsFile := path.Join(nsisFolder, nsisToolsFile)
	if _, err := buildassets.ReadOriginalFileWithProjectDataAndSave(options.ProjectData, toolsFile); err != nil {
		return fmt.Errorf("Unable to generate NSIS tools file: %w", err)
	}

	// Write the WebView2 SetupFile
	webviewSetup := buildassets.GetLocalPath(options.ProjectData, path.Join(nsisFolder, nsisWebView2SetupFile))
	if dir := filepath.Dir(webviewSetup); !fs.DirExists(dir) {
		if err := fs.MkDirs(dir, 0o755); err != nil {
			return err
		}
	}

	if err := webview2runtime.WriteInstallerToFile(webviewSetup); err != nil {
		return fmt.Errorf("Unable to write WebView2 Bootstrapper Setup: %w", err)
	}

	if !shell.CommandExists("makensis") {
		outputLogger.Println("Warning: Cannot create installer: makensis not found")
		return nil
	}

	nsisType := options.ProjectData.NSISType
	if nsisType == nsisTypeSingle && (amd64Binary == "" || arm64Binary == "") {
		nsisType = ""
	}

	switch nsisType {
	case "":
		fallthrough
	case nsisTypeMultiple:
		if amd64Binary != "" {
			if err := makeNSIS(options, "amd64", amd64Binary, ""); err != nil {
				return err
			}
		}

		if arm64Binary != "" {
			if err := makeNSIS(options, "arm64", "", arm64Binary); err != nil {
				return err
			}
		}

	case nsisTypeSingle:
		if err := makeNSIS(options, "single", amd64Binary, arm64Binary); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupported nsisType: %s", nsisType)
	}

	return nil
}

func makeNSIS(options *Options, installerKind string, amd64Binary string, arm64Binary string) error {
	verbose := options.Verbosity == VERBOSE
	outputLogger := options.Logger

	outputLogger.Print("  - Building '%s' installer: ", installerKind)
	args := []string{}
	if amd64Binary != "" {
		args = append(args, "-DARG_WAILS_AMD64_BINARY="+amd64Binary)
	}
	if arm64Binary != "" {
		args = append(args, "-DARG_WAILS_ARM64_BINARY="+arm64Binary)
	}
	args = append(args, nsisProjectFile)

	if verbose {
		outputLogger.Println("makensis %s", strings.Join(args, " "))
	}

	installerDir := buildassets.GetLocalPath(options.ProjectData, nsisFolder)
	stdOut, stdErr, err := shell.RunCommand(installerDir, "makensis", args...)
	if err != nil || verbose {
		outputLogger.Println(stdOut)
		outputLogger.Println(stdErr)
	}
	if err != nil {
		return fmt.Errorf("Error during creation of the installer: %w", err)
	}
	outputLogger.Println("Done.")
	return nil
}
