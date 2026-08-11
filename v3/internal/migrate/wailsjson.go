package migrate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// V2Config mirrors the subset of a Wails v2 wails.json that matters for
// migration. Field names/JSON tags match v2/internal/project/project.go.
type V2Config struct {
	Name           string `json:"name"`
	AssetDirectory string `json:"assetdir"`

	FrontendInstall      string `json:"frontend:install"`
	FrontendBuild        string `json:"frontend:build"`
	FrontendDevWatcher   string `json:"frontend:dev:watcher"`
	FrontendDevServerUrl string `json:"frontend:dev:serverUrl"`
	FrontendDir          string `json:"frontend:dir"`
	WailsJSDir           string `json:"wailsjsdir"`

	BuildDir       string `json:"build:dir"`
	BuildTags      string `json:"build:tags"`
	OutputFilename string `json:"outputfilename"`
	Obfuscated     bool   `json:"obfuscated"`
	GarbleArgs     string `json:"garbleargs"`

	PreBuildHooks  map[string]string `json:"preBuildHooks"`
	PostBuildHooks map[string]string `json:"postBuildHooks"`

	Author struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`

	Info struct {
		CompanyName      string              `json:"companyName"`
		ProductName      string              `json:"productName"`
		ProductVersion   string              `json:"productVersion"`
		Copyright        *string             `json:"copyright"`
		Comments         *string             `json:"comments"`
		FileAssociations []V2FileAssociation `json:"fileAssociations"`
		Protocols        []V2Protocol        `json:"protocols"`
	} `json:"info"`
}

// V2FileAssociation mirrors v2's file association config.
type V2FileAssociation struct {
	Ext         string `json:"ext"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IconName    string `json:"iconName"`
	Role        string `json:"role"`
}

// V2Protocol mirrors v2's custom protocol config.
type V2Protocol struct {
	Scheme      string `json:"scheme"`
	Description string `json:"description"`
	Role        string `json:"role"`
}

// LoadV2Config reads and parses <dir>/wails.json, applying the same defaults
// v2 applies for the fields the migrator uses.
func LoadV2Config(dir string) (V2Config, error) {
	var cfg V2Config
	data, err := os.ReadFile(filepath.Join(dir, "wails.json"))
	if err != nil {
		return cfg, fmt.Errorf("could not read wails.json: %w", err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("could not parse wails.json: %w", err)
	}
	if cfg.Name == "" {
		return cfg, fmt.Errorf("wails.json has no 'name' field")
	}
	if cfg.FrontendDir == "" {
		cfg.FrontendDir = "frontend"
	}
	if cfg.BuildDir == "" {
		cfg.BuildDir = "build"
	}
	if cfg.OutputFilename == "" {
		cfg.OutputFilename = cfg.Name
	}
	cfg.OutputFilename = strings.TrimSuffix(cfg.OutputFilename, ".exe")
	if cfg.Info.ProductName == "" {
		cfg.Info.ProductName = cfg.Name
	}
	if cfg.Info.ProductVersion == "" {
		cfg.Info.ProductVersion = "1.0.0"
	}
	return cfg, nil
}

// PackageManager guesses the frontend package manager from the v2
// frontend:install command. Returns "npm" when unsure (the v3 Taskfile
// default).
func (c V2Config) PackageManager() string {
	fields := strings.Fields(c.FrontendInstall)
	if len(fields) == 0 {
		return "npm"
	}
	switch fields[0] {
	case "npm", "yarn", "pnpm", "bun":
		return fields[0]
	}
	return "npm"
}
