package commands

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/internal/templates"
)

// v2ProjectConfig mirrors the JSON schema of a Wails v2 wails.json file
// (v2/internal/project/project.go). The v3 module cannot import the v2
// module, so the fields are duplicated here with the same JSON tags.
type v2ProjectConfig struct {
	Name                 string `json:"name"`
	AssetDirectory       string `json:"assetdir"`
	ReloadDirectories    string `json:"reloaddirs"`
	BuildCommand         string `json:"frontend:build"`
	InstallCommand       string `json:"frontend:install"`
	DevCommand           string `json:"frontend:dev"`
	DevBuildCommand      string `json:"frontend:dev:build"`
	DevInstallCommand    string `json:"frontend:dev:install"`
	DevWatcherCommand    string `json:"frontend:dev:watcher"`
	FrontendDevServerURL string `json:"frontend:dev:serverUrl"`
	FrontendDir          string `json:"frontend:dir"`
	WailsJSDir           string `json:"wailsjsdir"`
	Version              string `json:"version"`
	OutputFilename       string `json:"outputfilename"`
	Author               struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
	Info struct {
		CompanyName    string  `json:"companyName"`
		ProductName    string  `json:"productName"`
		ProductVersion string  `json:"productVersion"`
		Copyright      *string `json:"copyright"`
		Comments       *string `json:"comments"`
	} `json:"info"`
}

// loadV2Config parses wails.json into the typed config plus a raw key map
// used to report configuration keys that have no automatic mapping.
func loadV2Config(projectDir string) (*v2ProjectConfig, map[string]any, error) {
	data, err := os.ReadFile(filepath.Join(projectDir, "wails.json"))
	if err != nil {
		return nil, nil, err
	}
	cfg := &v2ProjectConfig{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, nil, fmt.Errorf("could not parse wails.json: %w", err)
	}
	raw := map[string]any{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil, fmt.Errorf("could not parse wails.json: %w", err)
	}
	if cfg.Name == "" {
		cfg.Name = "wailsapp"
	}
	if cfg.OutputFilename == "" {
		cfg.OutputFilename = cfg.Name
	}
	cfg.OutputFilename = strings.TrimSuffix(cfg.OutputFilename, ".exe")
	if cfg.Info.ProductName == "" {
		cfg.Info.ProductName = cfg.Name
	}
	if cfg.Info.CompanyName == "" {
		cfg.Info.CompanyName = cfg.Name
	}
	if cfg.Info.ProductVersion == "" {
		cfg.Info.ProductVersion = "1.0.0"
	}
	return cfg, raw, nil
}

// detectPackageManager infers the package manager from the v2
// frontend:install command. Defaults to npm.
func detectPackageManager(installCommand string) string {
	fields := strings.Fields(installCommand)
	if len(fields) == 0 {
		return "npm"
	}
	switch filepath.Base(fields[0]) {
	case "pnpm", "yarn", "bun":
		return filepath.Base(fields[0])
	default:
		return "npm"
	}
}

// v2ConfigGuidance maps v2 wails.json keys that have no automatic v3
// equivalent to a short pointer for the migration report.
var v2ConfigGuidance = map[string]string{
	"assetdir":               "v3 serves assets via `application.AssetOptions` in main.go; the embed directive in your Go code determines the asset directory.",
	"reloaddirs":             "configure watched paths in `build/config.yml` under `dev_mode`.",
	"wailsjsdir":             "v3 generates bindings into `frontend/bindings` via `wails3 generate bindings` (run automatically by the build tasks).",
	"frontend:dev":           "v3 runs the dev server via the `dev` task in `Taskfile.yml`; adjust `build/Taskfile.yml` if your dev script is not `npm run dev`.",
	"frontend:dev:build":     "v3 builds the frontend via the `build:frontend` task in `build/Taskfile.yml`.",
	"frontend:dev:install":   "v3 installs frontend dependencies via the `install:frontend:deps` task in `build/Taskfile.yml`.",
	"frontend:dev:serverUrl": "in v3, `wails3 dev` manages the dev server; the Vite port is set via the `VITE_PORT` variable in `Taskfile.yml`.",
	"devServer":              "the v3 dev server is configured in `build/config.yml` under `dev_mode` and via `wails3 dev` flags.",
	"debounceMS":             "set `dev_mode.debounce` in `build/config.yml`.",
	"appargs":                "pass application arguments via the `run` task in `Taskfile.yml`.",
	"nsisType":               "v3 generates NSIS configuration under `build/windows/nsis`; customise it there.",
	"obfuscated":             "run the build with `wails3 task build OBFUSCATED=true` (requires garble).",
	"garbleargs":             "set `GARBLE_ARGS` when running the build tasks.",
	"preBuildHooks":          "add steps to the platform build tasks in `build/<platform>/Taskfile.yml`.",
	"postBuildHooks":         "add steps to the platform build tasks in `build/<platform>/Taskfile.yml`.",
	"runNonNativeBuildHooks": "add steps to the platform build tasks in `build/<platform>/Taskfile.yml`.",
	"viteServerTimeout":      "not needed in v3; `wails3 dev` coordinates the Vite server directly.",
	"bindings":               "TypeScript binding generation is configured via flags on `wails3 generate bindings` in `build/Taskfile.yml` (see the `generate:bindings` task).",
	"debounceMs":             "set `dev_mode.debounce` in `build/config.yml`.",
}

// handledV2ConfigKeys are wails.json keys the migrator maps automatically.
var handledV2ConfigKeys = map[string]bool{
	"$schema":              true,
	"name":                 true,
	"outputfilename":       true,
	"frontend:install":     true,
	"frontend:build":       true,
	"frontend:dev:watcher": true,
	"version":              true, // schema version marker, dropped in v3
	"author":               true,
	"info":                 true,
	"frontend:dir":         true, // validated separately
}

// migrateConfig generates the v3 project configuration (root Taskfile.yml and
// the build/ assets directory) from the v2 wails.json, without overwriting
// any existing files.
func (m *migrator) migrateConfig() error {
	binaryName := templates.NormalizeBinaryName(m.cfg.OutputFilename)

	// Root Taskfile.yml, as laid down by `wails3 init`.
	taskfilePath := filepath.Join(m.projectDir, "Taskfile.yml")
	packageManager := detectPackageManager(m.cfg.InstallCommand)
	if _, err := os.Stat(taskfilePath); err == nil {
		m.skipped = append(m.skipped, "Taskfile.yml")
	} else {
		content, err := templates.RootTaskfile(binaryName)
		if err != nil {
			return fmt.Errorf("could not render Taskfile.yml: %w", err)
		}
		if packageManager != "npm" {
			content = []byte(strings.Replace(string(content), `| default "npm"`, `| default "`+packageManager+`"`, 1))
		}
		if err := os.WriteFile(taskfilePath, content, 0o644); err != nil {
			return err
		}
		m.created = append(m.created, "Taskfile.yml")
	}

	// Build assets: generate a pristine v3 build directory into a temporary
	// location, then copy over only the files the project does not already
	// have. v2 projects ship their own build/ directory (icons, Info.plist,
	// installer files) and those must be preserved.
	tempDir, err := os.MkdirTemp("", "wails3-migrate-build-assets")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	buildAssetsOptions := &BuildAssetsOptions{
		Dir:                tempDir,
		Name:               m.cfg.Name,
		BinaryName:         binaryName,
		ProductName:        m.cfg.Info.ProductName,
		ProductCompany:     m.cfg.Info.CompanyName,
		ProductVersion:     m.cfg.Info.ProductVersion,
		ProductDescription: "", // v2 wails.json has no description field
		Silent:             true,
	}
	if m.cfg.Info.Copyright != nil {
		buildAssetsOptions.ProductCopyright = *m.cfg.Info.Copyright
	}
	if m.cfg.Info.Comments != nil {
		buildAssetsOptions.ProductComments = *m.cfg.Info.Comments
	}
	if err := GenerateBuildAssets(buildAssetsOptions); err != nil {
		return fmt.Errorf("could not generate build assets: %w", err)
	}

	buildDir := filepath.Join(m.projectDir, "build")
	err = filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(tempDir, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(buildDir, rel)
		relDisplay := filepath.ToSlash(filepath.Join("build", rel))
		if _, err := os.Stat(dest); err == nil {
			m.skipped = append(m.skipped, relDisplay)
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		srcInfo, err := d.Info()
		if err != nil {
			return err
		}
		if err := os.WriteFile(dest, data, srcInfo.Mode().Perm()); err != nil {
			return err
		}
		m.created = append(m.created, relDisplay)
		return nil
	})
	if err != nil {
		return err
	}

	// Write the v2 product information into the freshly created config.yml.
	// If the project somehow already had one, it is left untouched.
	configYML := filepath.Join(buildDir, "config.yml")
	if m.wasCreated("build/config.yml") {
		values := map[string]string{
			"companyName": m.cfg.Info.CompanyName,
			"productName": m.cfg.Info.ProductName,
			"version":     m.cfg.Info.ProductVersion,
			// GenerateBuildAssets defaulted this and baked it into the
			// generated assets; keep config.yml consistent with them.
			"productIdentifier": buildAssetsOptions.ProductIdentifier,
		}
		if m.cfg.Info.Copyright != nil {
			values["copyright"] = *m.cfg.Info.Copyright
		}
		if m.cfg.Info.Comments != nil {
			values["comments"] = *m.cfg.Info.Comments
		}
		if err := updateConfigYMLInfo(configYML, values); err != nil {
			return err
		}
	}

	// Record the automatic mappings for the report.
	m.configMapped = append(m.configMapped,
		[2]string{"name", "`build/config.yml` (`info.productName`) and application Name in `main_v3.go.example`"},
		[2]string{"outputfilename", "`Taskfile.yml` (`APP_NAME` variable)"},
	)
	if m.cfg.InstallCommand != "" {
		m.configMapped = append(m.configMapped,
			[2]string{"frontend:install", fmt.Sprintf("`build/Taskfile.yml` task `install:frontend:deps` (package manager: %s)", packageManager)})
	}
	if m.cfg.BuildCommand != "" {
		m.configMapped = append(m.configMapped,
			[2]string{"frontend:build", fmt.Sprintf("`build/Taskfile.yml` task `build:frontend` (runs `%s run build`)", packageManager)})
	}
	if m.cfg.DevWatcherCommand != "" {
		m.configMapped = append(m.configMapped,
			[2]string{"frontend:dev:watcher", "`Taskfile.yml` task `dev` (managed by `wails3 dev`)"})
	}
	m.configMapped = append(m.configMapped,
		[2]string{"info.*", "`build/config.yml` (`info` section)"})

	if m.cfg.Author.Name != "" || m.cfg.Author.Email != "" {
		m.configManual = append(m.configManual,
			"`author` has no dedicated v3 field; company and copyright information lives in `build/config.yml`.")
	}
	if m.cfg.FrontendDir != "" && m.cfg.FrontendDir != "frontend" && m.cfg.FrontendDir != "./frontend" {
		m.configManual = append(m.configManual,
			fmt.Sprintf("`frontend:dir` is `%s`, but the generated v3 tasks assume `./frontend`. Update the `dir:` entries in `build/Taskfile.yml` accordingly.", m.cfg.FrontendDir))
	}
	if m.cfg.BuildCommand != "" && !isConventionalScript(m.cfg.BuildCommand, "build") {
		m.configManual = append(m.configManual,
			fmt.Sprintf("`frontend:build` was `%s`. The v3 tasks run the package.json `build` script; adjust the `build:frontend` task in `build/Taskfile.yml` if that does not match.", m.cfg.BuildCommand))
	}
	if m.cfg.DevWatcherCommand != "" && !isConventionalScript(m.cfg.DevWatcherCommand, "dev") {
		m.configManual = append(m.configManual,
			fmt.Sprintf("`frontend:dev:watcher` was `%s`. The v3 dev task runs the package.json `dev` script; adjust `build/Taskfile.yml` if that does not match.", m.cfg.DevWatcherCommand))
	}
	m.configManual = append(m.configManual,
		"Set `productIdentifier` (and a `description`) in `build/config.yml`, then run `wails3 task common:update:build-assets`.")

	// Report any remaining top-level keys that were not migrated.
	if fileAssociations := m.rawInfoKey("fileAssociations"); fileAssociations {
		m.configManual = append(m.configManual,
			"`info.fileAssociations`: copy your file associations into the `fileAssociations` section of `build/config.yml`, then run `wails3 task common:update:build-assets`.")
	}
	if protocols := m.rawInfoKey("protocols"); protocols {
		m.configManual = append(m.configManual,
			"`info.protocols`: copy your custom protocols into a `protocols` section in `build/config.yml`, then run `wails3 task common:update:build-assets`.")
	}
	for _, key := range sortedKeys(m.rawConfig) {
		if handledV2ConfigKeys[key] {
			continue
		}
		guidance, ok := v2ConfigGuidance[key]
		if !ok {
			guidance = "no direct v3 equivalent; see the migration guide."
		}
		m.configManual = append(m.configManual, fmt.Sprintf("`%s`: %s", key, guidance))
	}

	return nil
}

// isConventionalScript reports whether a v2 frontend command matches the
// package.json script convention the v3 Taskfile relies on, e.g.
// "npm run build" / "pnpm build" for script "build".
func isConventionalScript(command, script string) bool {
	fields := strings.Fields(command)
	if len(fields) < 2 {
		return false
	}
	last := fields[len(fields)-1]
	return last == script
}

func (m *migrator) rawInfoKey(key string) bool {
	info, ok := m.rawConfig["info"].(map[string]any)
	if !ok {
		return false
	}
	value, ok := info[key]
	if !ok {
		return false
	}
	if list, ok := value.([]any); ok {
		return len(list) > 0
	}
	return true
}
