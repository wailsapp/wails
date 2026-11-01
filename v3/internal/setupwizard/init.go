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
	"strings"

	"github.com/wailsapp/wails/v3/internal/features"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

// InitTemplate is a selectable project template, sourced from the templates
// package by the caller (the setupwizard package must not import templates).
type InitTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TypeScript  bool   `json:"typescript"`
}

// InitData is the payload exchanged with the init wizard frontend. The caller
// (commands.Init) seeds it with defaults + the available templates, and gets it
// back with the user's choices after Create.
type InitData struct {
	// Mode lets the frontend distinguish the init wizard from the setup wizard.
	Mode        string             `json:"mode"`
	WakeEnabled bool               `json:"wakeEnabled"`
	WakeSetup   *wakeSetupResponse `json:"wakeSetup,omitempty"`
	Wake        *WakeSetupConfig   `json:"wake,omitempty"`

	// BuildConfigFile is display-only; the backend always checks the environment.
	BuildConfigFile string `json:"buildConfigFile"`

	ProjectName        string `json:"projectName"`
	TemplateName       string `json:"templateName"`
	ProductName        string `json:"productName"`
	ProductCompany     string `json:"productCompany"`
	ProductIdentifier  string `json:"productIdentifier"`
	ProductDescription string `json:"productDescription"`
	ProductVersion     string `json:"productVersion"`
	ProductCopyright   string `json:"productCopyright"`
	ProductComments    string `json:"productComments"`
	// UseInterfaces selects interface vs class bindings for TypeScript projects.
	UseInterfaces bool `json:"useInterfaces"`

	// BaseDir is the absolute directory the project folder will be created in
	// (display only). The project lands at BaseDir/<ProjectName>.
	BaseDir string `json:"baseDir"`

	// Catalogue for the picker (not modified by the user).
	Templates       []InitTemplate `json:"templates"`
	DefaultTemplate string         `json:"defaultTemplate"`
}

// NewInitWizard creates a wizard that runs in project-init mode, pre-seeded with
// the given data (defaults + template list).
func NewInitWizard(data InitData) *Wizard {
	d := data
	d.Mode = "init"
	d.BuildConfigFile = "build/config.yml"
	d.WakeEnabled = features.WakeEnabled()
	d.WakeSetup, d.Wake = nil, nil
	if d.WakeEnabled {
		d.BuildConfigFile = "wails.yaml"
		root := filepath.Join(d.BaseDir, d.ProjectName)
		config := WakeSetupConfig{ApplicationType: "gui", Backend: "cgo", Targets: []WakeSetupTarget{{OS: runtime.GOOS, Arch: runtime.GOARCH, Formats: []string{}}}, Project: initProjectState(d)}
		d.WakeSetup = &wakeSetupResponse{Enabled: true, HostOS: runtime.GOOS, HostArch: runtime.GOARCH, ProjectRoot: root, Config: config, Catalog: wakeSetupCatalog()}
	}
	w := New()
	w.initData = &d
	return w
}

// RunInit launches the init wizard and blocks until the user clicks Create (or
// closes the window). On Create it returns the user-edited data; on close it
// returns (nil, nil).
func (w *Wizard) RunInit() (*InitData, error) {
	if err := w.Run(); err != nil {
		return nil, err
	}
	return w.initResult, nil
}

// handleInit returns the seeded init data (and signals init mode). In setup mode
// it returns null so the frontend falls back to the standard OOBE flow.
func (w *Wizard) handleInit(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	if w.initData == nil {
		_, _ = rw.Write([]byte("null"))
		return
	}
	_ = json.NewEncoder(rw).Encode(w.initData)
}

// handleInitCreate records the user's final choices and unblocks RunInit.
func (w *Wizard) handleInitCreate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if w.initData == nil {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "not in init mode"})
		return
	}

	if !w.isWizardOrigin(r) {
		http.Error(rw, "Origin not allowed", http.StatusForbidden)
		return
	}
	r.Body = http.MaxBytesReader(rw, r.Body, 1<<20)
	var result InitData
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&result); err != nil {
		http.Error(rw, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		http.Error(rw, "Init request must contain one JSON value", http.StatusBadRequest)
		return
	}
	if len(w.initData.Templates) != 0 {
		known := false
		for _, option := range w.initData.Templates {
			known = known || option.Name == result.TemplateName
		}
		if !known {
			http.Error(rw, "Select an available project template", http.StatusBadRequest)
			return
		}
	}
	result.Mode = "init"
	result.BaseDir = w.initData.BaseDir
	result.Templates = w.initData.Templates
	result.WakeEnabled = features.WakeEnabled()
	result.WakeSetup = nil
	if result.Wake != nil {
		if !result.WakeEnabled {
			http.Error(rw, "Wake choices require WAILS_EXP_USE_WAKE", http.StatusForbidden)
			return
		}
		result.Wake.Project = initProjectState(result)
		previewRoot, _, _, _, err := w.prepareInitWake(*result.Wake)
		if previewRoot != "" {
			defer os.RemoveAll(previewRoot)
		}
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	} else if result.WakeEnabled {
		http.Error(rw, "Choose the project build settings before creating it", http.StatusBadRequest)
		return
	}

	// Single-shot: only the first create request records the result and unblocks
	// RunInit. Concurrent requests are no-ops (no result race, no double close).
	w.doneOnce.Do(func() {
		w.initResult = &result
		close(w.done)
	})

	json.NewEncoder(rw).Encode(map[string]any{"success": true})
}

func initProjectState(data InitData) *ProjectConfigState {
	typescript := !strings.HasSuffix(data.TemplateName, "-js")
	for _, template := range data.Templates {
		if template.Name == data.TemplateName {
			typescript = template.TypeScript
		}
	}
	return &ProjectConfigState{Info: ProjectMetadataState{ProjectName: data.ProjectName, ProductName: data.ProductName, ProductIdentifier: data.ProductIdentifier, Version: data.ProductVersion, CompanyName: data.ProductCompany, Description: data.ProductDescription, Copyright: data.ProductCopyright, Comments: data.ProductComments}, Bindings: &BindingState{TypeScript: typescript, Interfaces: data.UseInterfaces}}
}

func (w *Wizard) prepareInitWake(config WakeSetupConfig) (string, *manifest.InitialState, manifest.SetupState, manifest.Config, error) {
	if w.initData == nil || !features.WakeEnabled() {
		return "", nil, manifest.SetupState{}, manifest.Config{}, errors.New("Wake choices are available only during project creation")
	}
	if config.Project == nil {
		return "", nil, manifest.SetupState{}, manifest.Config{}, errors.New("project metadata is required")
	}
	name := config.Project.Info.ProjectName
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\:`) {
		return "", nil, manifest.SetupState{}, manifest.Config{}, errors.New("project name must be a single directory name")
	}
	if config.ApplicationType != "gui" || config.Prebuilt {
		return "", nil, manifest.SetupState{}, manifest.Config{}, errors.New("new project templates create GUI source applications; headless and prebuilt configuration is not supported")
	}
	root, err := os.MkdirTemp("", "wails-init-preview-*")
	if err != nil {
		return "", nil, manifest.SetupState{}, manifest.Config{}, err
	}
	return prepareWakeSetupAt(config, root, nil)
}

// ApplyWakeInit applies validated init choices only to the newly scaffolded manifest.
func ApplyWakeInit(root string, config WakeSetupConfig) error {
	if !features.WakeEnabled() {
		return errors.New("Wake choices require WAILS_EXP_USE_WAKE")
	}
	if config.ApplicationType != "gui" || config.Prebuilt {
		return errors.New("new project templates require a GUI source build")
	}
	loaded, err := manifest.Load(root, "")
	if err != nil {
		return fmt.Errorf("load newly created project: %w", err)
	}
	if loaded.Config.Root != filepath.Clean(root) {
		return errors.New("new project must own its build manifest")
	}
	_, _, state, _, err := prepareWakeSetupAt(config, root, loaded)
	if err != nil {
		return err
	}
	return manifest.WriteSetup(root, nil, state)
}
