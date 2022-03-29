package project

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Project holds the data related to a Wails project
type Project struct {

	/*** Application Data ***/
	Name           string `json:"name"`
	AssetDirectory string `json:"assetdir,omitempty"`

	ReloadDirectories string `json:"reloaddirs,omitempty"`

	BuildCommand   string `json:"frontend:build"`
	InstallCommand string `json:"frontend:install"`

	// Commands used in `wails dev`
	DevCommand        string `json:"frontend:dev"`
	DevWatcherCommand string `json:"frontend:dev:watcher"`

	// Directory to generate the API Module
	WailsJSDir string `json:"wailsjsdir"`

	Version string `json:"version"`

	/*** Internal Data ***/

	// The path to the project directory
	Path string

	// Build directory
	BuildDir string

	// The output filename
	OutputFilename string `json:"outputfilename"`

	// The type of application. EG: Desktop, Server, etc
	OutputType string

	// The platform to target
	Platform string

	// RunNonNativeBuildHooks will run build hooks though they are defined for a GOOS which is not equal to the host os
	RunNonNativeBuildHooks bool `json:"runNonNativeBuildHooks"`

	// Post build hooks for different targets, the hooks are executed in the following order
	// Key: GOOS/GOARCH - Executed at build level after a build of the specific platform and arch
	// Key: GOOS/*      - Executed at build level after a build of the specific platform
	// Key: */*         - Executed at build level after a build
	// The following keys are not yet supported.
	// Key: GOOS        - Executed at platform level after all builds of the specific platform
	// Key: *           - Executed at platform level after all builds of a platform
	// Key: [empty]     - Executed at global level after all builds of all platforms
	PostBuildHooks map[string]string `json:"postBuildHooks"`

	// The application author
	Author Author

	// The application information
	Info Info

	// Fully qualified filename
	filename string

	// The debounce time for hot-reload of the built-in dev server. Default 100
	DebounceMS int `json:"debounceMS"`

	// The url to use to server assets. Default "https://localhost:34115"
	DevServerURL string `json:"devserverurl"`

	// Arguments that are forwared to the application in dev mode
	AppArgs string `json:"appargs"`

	// NSISType to be build
	NSISType string `json:"nsisType"`
}

func (p *Project) Save() error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p.filename, data, 0755)
}

// Author stores details about the application author
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Info struct {
	CompanyName    string  `json:"companyName"`
	ProductName    string  `json:"productName"`
	ProductVersion string  `json:"productVersion"`
	Copyright      *string `json:"copyright"`
	Comments       *string `json:"comments"`
}

// Load the project from the current working directory
func Load(projectPath string) (*Project, error) {

	// Attempt to load project.json
	projectFile := filepath.Join(projectPath, "wails.json")
	rawBytes, err := os.ReadFile(projectFile)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON
	var result Project
	err = json.Unmarshal(rawBytes, &result)
	if err != nil {
		return nil, err
	}

	// Fix up our project paths
	result.filename = projectFile

	if result.Version == "" {
		result.Version = "2"
	}

	// Create default name if not given
	if result.Name == "" {
		result.Name = "wailsapp"
	}

	// Fix up OutputFilename
	switch runtime.GOOS {
	case "windows":
		if !strings.HasSuffix(result.OutputFilename, ".exe") {
			result.OutputFilename += ".exe"
		}
	case "darwin", "linux":
		if strings.HasSuffix(result.OutputFilename, ".exe") {
			result.OutputFilename = strings.TrimSuffix(result.OutputFilename, ".exe")
		}
	}

	if result.Info.CompanyName == "" {
		result.Info.CompanyName = result.Name
	}
	if result.Info.ProductName == "" {
		result.Info.ProductName = result.Name
	}
	if result.Info.ProductVersion == "" {
		result.Info.ProductVersion = "1.0.0"
	}
	if result.Info.Copyright == nil {
		v := "Copyright........."
		result.Info.Copyright = &v
	}
	if result.Info.Comments == nil {
		v := "Built using Wails (https://wails.app)"
		result.Info.Comments = &v
	}

	// Return our project data
	return &result, nil
}
