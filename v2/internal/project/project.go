package project

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/samber/lo"
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
	DevBuildCommand   string `json:"frontend:dev:build"`
	DevInstallCommand string `json:"frontend:dev:install"`
	DevWatcherCommand string `json:"frontend:dev:watcher"`
	// The url of the external wails dev server. If this is set, this server is used for the frontend. Default ""
	FrontendDevServerURL string `json:"frontend:dev:serverUrl"`

	// Directory to generate the API Module
	WailsJSDir string `json:"wailsjsdir"`

	Version string `json:"version"`

	/*** Internal Data ***/

	// The path to the project directory
	Path string `json:"projectdir"`

	// Build directory
	BuildDir string `json:"build:dir"`

	// The output filename
	OutputFilename string `json:"outputfilename"`

	// The type of application. EG: Desktop, Server, etc
	OutputType string

	// The platform to target
	Platform string

	// RunNonNativeBuildHooks will run build hooks though they are defined for a GOOS which is not equal to the host os
	RunNonNativeBuildHooks bool `json:"runNonNativeBuildHooks"`

	// Build hooks for different targets, the hooks are executed in the following order
	// Key: GOOS/GOARCH - Executed at build level before/after a build of the specific platform and arch
	// Key: GOOS/*      - Executed at build level before/after a build of the specific platform
	// Key: */*         - Executed at build level before/after a build
	// The following keys are not yet supported.
	// Key: GOOS        - Executed at platform level before/after all builds of the specific platform
	// Key: *           - Executed at platform level before/after all builds of a platform
	// Key: [empty]     - Executed at global level before/after all builds of all platforms
	PostBuildHooks map[string]string `json:"postBuildHooks"`
	PreBuildHooks  map[string]string `json:"preBuildHooks"`

	// The application author
	Author Author

	// The application information
	Info Info

	// Fully qualified filename
	filename string

	// The debounce time for hot-reload of the built-in dev server. Default 100
	DebounceMS int `json:"debounceMS"`

	// The address to bind the wails dev server to. Default "localhost:34115"
	DevServer string `json:"devServer"`

	// Arguments that are forwared to the application in dev mode
	AppArgs string `json:"appargs"`

	// NSISType to be build
	NSISType string `json:"nsisType"`

	// Garble
	Obfuscated bool   `json:"obfuscated"`
	GarbleArgs string `json:"garbleargs"`

	// Frontend directory
	FrontendDir string `json:"frontend:dir"`

	Bindings Bindings `json:"bindings"`
}

func (p *Project) GetFrontendDir() string {
	if filepath.IsAbs(p.FrontendDir) {
		return p.FrontendDir
	}
	return filepath.Join(p.Path, p.FrontendDir)
}

func (p *Project) GetWailsJSDir() string {
	if filepath.IsAbs(p.WailsJSDir) {
		return p.WailsJSDir
	}
	return filepath.Join(p.Path, p.WailsJSDir)
}

func (p *Project) GetBuildDir() string {
	if filepath.IsAbs(p.BuildDir) {
		return p.BuildDir
	}
	return filepath.Join(p.Path, p.BuildDir)
}

func (p *Project) GetDevBuildCommand() string {
	if p.DevBuildCommand != "" {
		return p.DevBuildCommand
	}
	if p.DevCommand != "" {
		return p.DevCommand
	}
	return p.BuildCommand
}

func (p *Project) GetDevInstallerCommand() string {
	if p.DevInstallCommand != "" {
		return p.DevInstallCommand
	}
	return p.InstallCommand
}

func (p *Project) IsFrontendDevServerURLAutoDiscovery() bool {
	return p.FrontendDevServerURL == "auto"
}

func (p *Project) Save() error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p.filename, data, 0o755)
}

func (p *Project) setDefaults() {
	if p.Path == "" {
		p.Path = lo.Must(os.Getwd())
	}
	if p.Version == "" {
		p.Version = "2"
	}
	// Create default name if not given
	if p.Name == "" {
		p.Name = "wailsapp"
	}
	if p.OutputFilename == "" {
		p.OutputFilename = p.Name
	}
	if p.FrontendDir == "" {
		p.FrontendDir = "frontend"
	}
	if p.WailsJSDir == "" {
		p.WailsJSDir = p.FrontendDir
	}
	if p.BuildDir == "" {
		p.BuildDir = "build"
	}
	if p.DebounceMS == 0 {
		p.DebounceMS = 100
	}
	if p.DevServer == "" {
		p.DevServer = "localhost:34115"
	}
	if p.NSISType == "" {
		p.NSISType = "multiple"
	}
	if p.Info.CompanyName == "" {
		p.Info.CompanyName = p.Name
	}
	if p.Info.ProductName == "" {
		p.Info.ProductName = p.Name
	}
	if p.Info.ProductVersion == "" {
		p.Info.ProductVersion = "1.0.0"
	}
	if p.Info.Copyright == nil {
		v := "Copyright........."
		p.Info.Copyright = &v
	}
	if p.Info.Comments == nil {
		v := "Built using Wails (https://wails.io)"
		p.Info.Comments = &v
	}

	// Fix up OutputFilename
	switch runtime.GOOS {
	case "windows":
		if !strings.HasSuffix(p.OutputFilename, ".exe") {
			p.OutputFilename += ".exe"
		}
	case "darwin", "linux":
		p.OutputFilename = strings.TrimSuffix(p.OutputFilename, ".exe")
	}
}

// Author stores details about the application author
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Info struct {
	CompanyName      string            `json:"companyName"`
	ProductName      string            `json:"productName"`
	ProductVersion   string            `json:"productVersion"`
	Copyright        *string           `json:"copyright"`
	Comments         *string           `json:"comments"`
	FileAssociations []FileAssociation `json:"fileAssociations"`
	Protocols        []Protocol        `json:"protocols"`
}

type FileAssociation struct {
	Ext         string `json:"ext"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IconName    string `json:"iconName"`
	Role        string `json:"role"`
}

type Protocol struct {
	Scheme      string `json:"scheme"`
	Description string `json:"description"`
	Role        string `json:"role"`
}

type Bindings struct {
	TsGeneration TsGeneration `json:"ts_generation"`
}

type TsGeneration struct {
	Prefix     string `json:"prefix"`
	Suffix     string `json:"suffix"`
	OutputType string `json:"outputType"`
}

// Parse the given JSON data into a Project struct
func Parse(projectData []byte) (*Project, error) {
	project := &Project{}
	err := json.Unmarshal(projectData, project)
	if err != nil {
		return nil, err
	}
	project.setDefaults()
	return project, nil
}

// Load the project from the current working directory
func Load(projectPath string) (*Project, error) {
	projectFile := filepath.Join(projectPath, "wails.json")
	rawBytes, err := os.ReadFile(projectFile)
	if err != nil {
		return nil, err
	}
	result, err := Parse(rawBytes)
	if err != nil {
		return nil, err
	}
	result.filename = projectFile
	return result, nil
}
