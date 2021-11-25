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
	AssetDirectory string `json:"assetdir"`

	BuildCommand   string `json:"frontend:build"`
	InstallCommand string `json:"frontend:install"`
	DevCommand     string `json:"frontend:dev"`

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

	// The application author
	Author Author

	// Fully qualified filename
	filename string

	// The debounce time for hot-reload of the built-in dev server. Default 100
	DebounceMS int `json:"debounceMS"`

	// The url to use to server assets. Default "https://localhost:34115"
	DevServerURL string `json:"devserverurl"`

	// Arguments that are forwared to the application in dev mode
	AppArgs string `json:"appargs"`
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

	// Return our project data
	return &result, nil
}
