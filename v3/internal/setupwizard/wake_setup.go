package setupwizard

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v3/internal/features"
	"github.com/wailsapp/wails/v3/internal/wake/buildinfo"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

type WakeSetupConfig struct {
	ApplicationType    string              `json:"applicationType"`
	Backend            string              `json:"backend"`
	ExperimentalPureGo bool                `json:"experimentalPureGo"`
	Prebuilt           bool                `json:"prebuilt"`
	Sign               bool                `json:"sign"`
	Targets            []WakeSetupTarget   `json:"targets"`
	Project            *ProjectConfigState `json:"project,omitempty"`
}
type WakeSetupTarget struct {
	OS          string   `json:"os"`
	Arch        string   `json:"arch"`
	Formats     []string `json:"formats"`
	Destination string   `json:"destination,omitempty"`
}
type wakeSetupCatalogEntry struct {
	OS           string   `json:"os"`
	Arch         string   `json:"arch"`
	Formats      []string `json:"formats"`
	Experimental bool     `json:"experimental"`
}
type wakeSetupResponse struct {
	Enabled      bool                    `json:"enabled"`
	HostOS       string                  `json:"hostOS"`
	HostArch     string                  `json:"hostArch"`
	ProjectRoot  string                  `json:"projectRoot"`
	ManifestPath string                  `json:"manifestPath,omitempty"`
	Config       WakeSetupConfig         `json:"config"`
	Catalog      []wakeSetupCatalogEntry `json:"catalog"`
	Error        string                  `json:"error,omitempty"`
}
type wakeRequirementsResponse struct {
	Requirements []DependencyStatus `json:"requirements"`
	Errors       []string           `json:"errors"`
	Warnings     []string           `json:"warnings"`
	Command      string             `json:"command"`
}

func wakeSetupCommand(c WakeSetupConfig) string {
	command := "wails3 build setup"
	if c.Prebuilt {
		command += " --skip-build"
	}
	return command
}

func wakeSetupProject() (string, *manifest.Loaded, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", nil, err
	}
	root, _, err := manifest.Discover(cwd)
	if err == nil {
		loaded, err := manifest.Load(root, "")
		return root, loaded, err
	}
	if !errors.Is(err, manifest.ErrNotFound) {
		return cwd, nil, err
	}
	root = cwd
	for dir := cwd; ; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			root = dir
			break
		}
		if filepath.Dir(dir) == dir {
			break
		}
	}
	for _, name := range []string{manifest.Filename, "Taskfile.yml", "Taskfile.yaml", "taskfile.yml", "taskfile.yaml"} {
		if _, err := os.Lstat(filepath.Join(root, name)); err == nil {
			return root, nil, fmt.Errorf("%s already exists; review wails3 migrate before activating Wake setup", name)
		} else if !errors.Is(err, os.ErrNotExist) {
			return root, nil, err
		}
	}
	return root, nil, nil
}

func wakeSetupDefaults(root string, loaded *manifest.Loaded) WakeSetupConfig {
	c := WakeSetupConfig{ApplicationType: "gui", Backend: "cgo", Targets: []WakeSetupTarget{{OS: runtime.GOOS, Arch: runtime.GOARCH, Formats: []string{}}}}
	if loaded == nil {
		defaults, _ := LoadGlobalDefaults()
		name := filepath.Base(root)
		metadata := ProjectConfigState{Info: ProjectMetadataState{ProjectName: name, ProductName: name, ProductIdentifier: defaults.GenerateProductIdentifier(name), Version: defaults.Project.DefaultVersion, CompanyName: defaults.Author.Company, Description: defaults.GenerateDescription(name), Copyright: defaults.GenerateCopyright()}, Bindings: &BindingState{TypeScript: true, Interfaces: defaults.Project.UseInterfaces}}
		if metadata.Info.Version == "" {
			metadata.Info.Version = "0.1.0"
		}
		c.Project = &metadata
		return c
	}
	config := loaded.Config
	c.ApplicationType = config.Project.Kind
	if c.ApplicationType == "" {
		c.ApplicationType = "gui"
	}
	if slices.Contains(config.Build.Go.Tags, "purego") || config.Build.Environment["CGO_ENABLED"] == "0" {
		c.Backend = "purego"
		c.ExperimentalPureGo = true
	}
	metadata := projectConfigStateFromConfig(config)
	c.Project = &metadata
	if profile, ok := config.Profiles["setup"]; ok {
		c.Targets = nil
		for _, target := range profile.Targets {
			platform, arch, _ := strings.Cut(target.Target, "/")
			c.Targets = append(c.Targets, WakeSetupTarget{platform, arch, append([]string{}, target.Formats...), target.Destination})
			c.Sign = c.Sign || target.Sign
		}
	}
	return c
}

func (w *Wizard) handleWakeSetup(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		response := wakeSetupResponse{Enabled: features.WakeEnabled(), HostOS: runtime.GOOS, HostArch: runtime.GOARCH, Catalog: []wakeSetupCatalogEntry{}}
		if response.Enabled {
			root, loaded, err := wakeSetupProject()
			response.ProjectRoot = root
			response.Config = wakeSetupDefaults(root, loaded)
			if loaded != nil && err == nil {
				backend, backendErr := wakeExistingBackend(loaded.Config, response.Config.Targets)
				if backendErr != nil {
					err = backendErr
				} else {
					response.Config.Backend = backend
					response.Config.ExperimentalPureGo = backend == "purego"
				}
			}
			if loaded != nil {
				response.ManifestPath = loaded.Path
			}
			if err != nil {
				response.Error = err.Error()
			}
			for _, capability := range buildinfo.TargetCapabilities() {
				entry := wakeSetupCatalogEntry{OS: capability.Target.OS, Arch: capability.Target.Arch, Formats: []string{}, Experimental: capability.Target.OS == "ios" || capability.Target.OS == "android"}
				for _, format := range buildinfo.FormatCapabilities() {
					if capability.SupportsFormat(format.Name, false) {
						entry.Formats = append(entry.Formats, format.Name)
					}
				}
				response.Catalog = append(response.Catalog, entry)
			}
		}
		_ = json.NewEncoder(rw).Encode(response)
		return
	}
	var config WakeSetupConfig
	if !w.decodeWakeSetup(rw, r, &config) {
		return
	}
	w.wakeMu.Lock()
	defer w.wakeMu.Unlock()
	root, initial, state, _, err := prepareWakeSetup(config)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err := manifest.WriteSetup(root, initial, state); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	_ = json.NewEncoder(rw).Encode(map[string]any{"path": filepath.Join(root, manifest.Filename), "config": config, "command": wakeSetupCommand(config)})
}

func (w *Wizard) decodeWakeSetup(rw http.ResponseWriter, r *http.Request, config *WakeSetupConfig) bool {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	if !features.WakeEnabled() {
		http.Error(rw, "Wake setup requires WAILS_EXP_USE_WAKE", http.StatusForbidden)
		return false
	}
	if !w.isWizardOrigin(r) {
		http.Error(rw, "Origin not allowed", http.StatusForbidden)
		return false
	}
	r.Body = http.MaxBytesReader(rw, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(config); err != nil {
		http.Error(rw, "Invalid setup request: "+err.Error(), http.StatusBadRequest)
		return false
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		http.Error(rw, "Setup request must contain one JSON value", http.StatusBadRequest)
		return false
	}
	return true
}

func prepareWakeSetup(c WakeSetupConfig) (string, *manifest.InitialState, manifest.SetupState, manifest.Config, error) {
	root, loaded, err := wakeSetupProject()
	if err != nil {
		return root, nil, manifest.SetupState{}, manifest.Config{}, err
	}
	if c.Backend != "cgo" && c.Backend != "purego" {
		return root, nil, manifest.SetupState{}, manifest.Config{}, errors.New("backend must be cgo or purego")
	}
	if c.Backend == "purego" && !c.ExperimentalPureGo {
		return root, nil, manifest.SetupState{}, manifest.Config{}, errors.New("pure Go requires explicit experimental opt-in")
	}
	if len(c.Targets) == 0 || len(c.Targets) > len(buildinfo.TargetCapabilities()) {
		return root, nil, manifest.SetupState{}, manifest.Config{}, errors.New("select at least one supported target")
	}
	var initial *manifest.InitialState
	var existing manifest.Config
	if loaded != nil {
		existing = loaded.Config
	} else {
		if c.Project == nil {
			return root, nil, manifest.SetupState{}, manifest.Config{}, errors.New("project metadata is required before creating wails.yaml")
		}
		value, err := c.Project.ManifestState(manifest.Project{})
		if err != nil {
			return root, nil, manifest.SetupState{}, manifest.Config{}, err
		}
		initial = &value
	}
	state := manifest.SetupState{Kind: c.ApplicationType, PureGo: c.Backend == "purego", FrontendDisabled: c.ApplicationType != "gui", Profile: manifest.Profile{Name: "setup"}, Signing: map[string]manifest.SigningPlatform{}}
	for _, binary := range existing.Binaries {
		if binary.Kind == "gui" {
			state.FrontendDisabled = false
		}
	}
	previous := existing.Profiles["setup"]
	for _, target := range c.Targets {
		capability, ok := buildinfo.LookupTarget(target.OS, target.Arch)
		if !ok {
			return root, nil, state, manifest.Config{}, fmt.Errorf("unsupported target %s/%s", target.OS, target.Arch)
		}
		mobile := target.OS == "ios" || target.OS == "android"
		if mobile && (c.Prebuilt || state.PureGo || c.ApplicationType != "gui") {
			return root, nil, state, manifest.Config{}, errors.New("mobile targets require a GUI source build with CGo")
		}
		for _, format := range target.Formats {
			if !capability.SupportsFormat(format, false) {
				return root, nil, state, manifest.Config{}, fmt.Errorf("format %s is not supported for %s/%s production builds", format, target.OS, target.Arch)
			}
		}
		if c.Sign && target.OS == "linux" {
			if len(target.Formats) == 0 {
				return root, nil, state, manifest.Config{}, errors.New("Linux signing requires a deb, rpm or archlinux package")
			}
			for _, format := range target.Formats {
				if format != "deb" && format != "rpm" && format != "archlinux" {
					return root, nil, state, manifest.Config{}, fmt.Errorf("Linux signing does not support %s; select deb, rpm or archlinux", format)
				}
			}
		}
		selected := manifest.ProfileTarget{Target: target.OS + "/" + target.Arch, Formats: target.Formats, Sign: c.Sign, Destination: target.Destination}
		for _, old := range previous.Targets {
			if old.Target == selected.Target && c.Sign {
				selected.Notarize = old.Notarize
			}
		}
		state.Profile.Targets = append(state.Profile.Targets, selected)
		if c.Sign {
			signing := wakeSigningPlatform(existing.Signing, target.OS)
			signing.Enabled = true
			state.Signing[target.OS] = signing
		}
	}
	config, err := manifest.PreviewSetup(root, initial, state)
	return root, initial, state, config, err
}

func wakeSigningPlatform(s manifest.Signing, platform string) manifest.SigningPlatform {
	switch platform {
	case "windows":
		return s.Windows
	case "darwin":
		return s.Darwin
	case "linux":
		return s.Linux
	case "ios":
		return s.IOS
	default:
		return s.Android
	}
}

func wakeExistingBackend(config manifest.Config, targets []WakeSetupTarget) (string, error) {
	modes := map[string]bool{}
	platforms := map[string]manifest.Platform{"windows": config.Targets.Windows, "darwin": config.Targets.Darwin, "linux": config.Targets.Linux, "ios": config.Targets.IOS, "android": config.Targets.Android}
	for _, target := range targets {
		platform := platforms[target.OS]
		var settings manifest.Target
		switch target.Arch {
		case "amd64":
			settings = platform.AMD64
		case "arm64":
			settings = platform.ARM64
		case "universal":
			settings = platform.Universal
		}
		environment := config.Build.Environment
		if settings.Environment != nil {
			environment = settings.Environment
		}
		mode := environment["CGO_ENABLED"]
		if mode != "" && mode != "0" && mode != "1" {
			return "", fmt.Errorf("target %s/%s has invalid CGO_ENABLED", target.OS, target.Arch)
		}
		pureGo := slices.Contains(config.Build.Go.Tags, "purego") || slices.Contains(settings.Tags, "purego")
		if pureGo && mode == "1" {
			return "", fmt.Errorf("target %s/%s combines purego with CGO_ENABLED=1; resolve this backend conflict in YAML", target.OS, target.Arch)
		}
		gui := config.Project.Kind == "" || config.Project.Kind == "gui"
		if mode == "0" && gui && (target.OS == "linux" || target.OS == "darwin") {
			pureGo = true
		}
		backend := "cgo"
		if pureGo {
			backend = "purego"
		}
		modes[backend] = true
	}
	if len(modes) > 1 {
		return "", errors.New("selected targets use mixed standard and experimental pure Go backends; configure a single backend in YAML before using setup")
	}
	if modes["purego"] {
		return "purego", nil
	}
	return "cgo", nil
}
