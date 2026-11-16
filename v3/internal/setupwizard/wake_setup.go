package setupwizard

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

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

func (w *Wizard) handleWakeSetup(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(rw, "Project build choices are saved by wails3 init -ui", http.StatusForbidden)
		return
	}
	if w.initData == nil || !features.WakeEnabled() {
		_ = json.NewEncoder(rw).Encode(wakeSetupResponse{Catalog: []wakeSetupCatalogEntry{}})
		return
	}
	_ = json.NewEncoder(rw).Encode(w.initData.WakeSetup)
}

func wakeSetupCatalog() []wakeSetupCatalogEntry {
	catalog := []wakeSetupCatalogEntry{}
	for _, capability := range buildinfo.TargetCapabilities() {
		entry := wakeSetupCatalogEntry{OS: capability.Target.OS, Arch: capability.Target.Arch, Formats: []string{}, Experimental: capability.Target.OS == "ios" || capability.Target.OS == "android"}
		for _, format := range buildinfo.FormatCapabilities() {
			if capability.SupportsFormat(format.Name, false) {
				entry.Formats = append(entry.Formats, format.Name)
			}
		}
		catalog = append(catalog, entry)
	}
	return catalog
}

func (w *Wizard) decodeWakeSetup(rw http.ResponseWriter, r *http.Request, config *WakeSetupConfig) bool {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	if w.initData == nil || !features.WakeEnabled() {
		http.Error(rw, "Wake project choices require wails3 init -ui with WAILS_EXP_USE_WAKE", http.StatusForbidden)
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

func prepareWakeSetupAt(c WakeSetupConfig, root string, loaded *manifest.Loaded) (string, *manifest.InitialState, manifest.SetupState, manifest.Config, error) {
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
