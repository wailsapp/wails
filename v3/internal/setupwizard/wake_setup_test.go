package setupwizard

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func wakeTestProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	t.Chdir(root)
	t.Setenv("WAILS_EXP_USE_WAKE", "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	for name, data := range map[string]string{"go.mod": "module example.test/probe\n\ngo 1.25\n", "main.go": "package main\nfunc main(){}\n"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}
func wakeTestConfig() WakeSetupConfig {
	return WakeSetupConfig{ApplicationType: "cli", Backend: "purego", ExperimentalPureGo: true, Targets: []WakeSetupTarget{{OS: "linux", Arch: "amd64", Formats: []string{"deb"}}}, Project: &ProjectConfigState{Info: ProjectMetadataState{ProjectName: "probe", ProductName: "Probe", ProductIdentifier: "org.example.probe", Version: "1.0.0"}}}
}
func wakeRequest(t *testing.T, w *Wizard, method, path string, value any) *httptest.ResponseRecorder {
	t.Helper()
	var data []byte
	var err error
	if value != nil {
		data, err = json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
	}
	mux := http.NewServeMux()
	w.setupRoutes(mux)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(method, path, bytes.NewReader(data)))
	return response
}
func TestWakeSetupFlagIsServerAuthoritative(t *testing.T) {
	root := wakeTestProject(t)
	for _, value := range []string{"", "false", "1"} {
		t.Setenv("WAILS_EXP_USE_WAKE", value)
		response := wakeRequest(t, New(), http.MethodGet, "/api/state", nil)
		var state WizardState
		if err := json.Unmarshal(response.Body.Bytes(), &state); err != nil || !state.WakeEnabled {
			t.Fatalf("presence %q: %s", value, response.Body.String())
		}
	}
	if err := os.Unsetenv("WAILS_EXP_USE_WAKE"); err != nil {
		t.Fatal(err)
	}
	response := wakeRequest(t, New(), http.MethodGet, "/api/wake/setup", nil)
	var state wakeSetupResponse
	if err := json.Unmarshal(response.Body.Bytes(), &state); err != nil || state.Enabled {
		t.Fatalf("disabled state: %s", response.Body.String())
	}
	for _, path := range []string{"/api/wake/setup", "/api/wake/setup/requirements"} {
		if response := wakeRequest(t, New(), http.MethodPost, path, wakeTestConfig()); response.Code != 403 {
			t.Fatalf("disabled POST: %s", response.Body.String())
		}
	}
	if _, err := os.Stat(filepath.Join(root, "wails.yaml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("disabled request wrote a manifest")
	}
}
func TestWakeSetupCreatesOnlyOnSaveAndPlansSelectedBackend(t *testing.T) {
	root := wakeTestProject(t)
	config := wakeTestConfig()
	w := New()
	get := wakeRequest(t, w, http.MethodGet, "/api/wake/setup", nil)
	var state wakeSetupResponse
	if err := json.Unmarshal(get.Body.Bytes(), &state); err != nil {
		t.Fatal(err)
	}
	if state.Error != "" || state.ProjectRoot != root || state.Config.Project == nil || len(state.Catalog) == 0 {
		t.Fatalf("bad state: %s", get.Body.String())
	}
	preview := wakeRequest(t, w, http.MethodPost, "/api/wake/setup/requirements", config)
	var requirements wakeRequirementsResponse
	if err := json.Unmarshal(preview.Body.Bytes(), &requirements); err != nil {
		t.Fatal(preview.Body.String())
	}
	if len(requirements.Errors) != 0 {
		t.Fatal(requirements.Errors)
	}
	if _, err := os.Stat(filepath.Join(root, "wails.yaml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("preview wrote manifest")
	}
	response := wakeRequest(t, w, http.MethodPost, "/api/wake/setup", config)
	if response.Code != 200 {
		t.Fatal(response.Body.String())
	}
	loaded, err := manifest.Load(root, "setup")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Config.Project.Kind != "cli" || !loaded.Config.Frontend.Disabled || loaded.Config.Build.Environment["CGO_ENABLED"] != "0" || !slices.Contains(loaded.Config.Build.Go.Tags, "purego") {
		t.Fatal("backend or kind not persisted")
	}
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, node := range plan.Nodes {
		if spec, ok := node.Spec.(pipeline.CompileSpec); ok {
			found = true
			if spec.CGOEnabled != "0" || !spec.Console || !slices.Contains(spec.Tags, "purego") {
				t.Fatalf("wrong compile spec: %+v", spec)
			}
		}
	}
	if !found {
		t.Fatal("no compilation planned")
	}
	if !strings.Contains(response.Body.String(), "wails3 build setup") {
		t.Fatal(response.Body.String())
	}
}
func TestWakeSetupPreservesUnownedSettingsAndInheritedEnvironment(t *testing.T) {
	root := wakeTestProject(t)
	original := `# user configuration
version: 3
project:
  name: probe
  product_name: Existing
  identifier: org.example.existing
  version: 9.1.0
  binary_name: keep-binary
build:
  tags: [enterprise]
  environment: {CUSTOM_FLAG: keep}
frontend:
  install: [pnpm, install]
  build: [pnpm, run, custom]
  output: web/dist
targets:
  windows/amd64:
    tags: [windows_custom]
    environment: {TARGET_ONLY: keep}
profiles:
  release:
    targets:
      windows/amd64: {formats: [msix]}
windows:
  signing: {thumbprint: keep-thumbprint}
`
	path := filepath.Join(root, "wails.yaml")
	if err := os.WriteFile(path, []byte(original), 0640); err != nil {
		t.Fatal(err)
	}
	config := wakeTestConfig()
	response := wakeRequest(t, New(), http.MethodPost, "/api/wake/setup", config)
	if response.Code != 200 {
		t.Fatal(response.Body.String())
	}
	loaded, err := manifest.Load(root, "setup")
	if err != nil {
		t.Fatal(err)
	}
	c := loaded.Config
	if c.Project.ProductName != "Existing" || c.Project.Version != "9.1.0" || c.Project.BinaryName != "keep-binary" || c.Frontend.Build[2] != "custom" || c.Profiles["release"].Targets[0].Formats[0] != "msix" || c.Signing.Windows.Thumbprint != "keep-thumbprint" {
		t.Fatal("unowned settings changed")
	}
	if c.Targets.Linux.AMD64.Environment["CUSTOM_FLAG"] != "keep" || c.Targets.Windows.AMD64.Environment["TARGET_ONLY"] != "keep" || !slices.Contains(c.Build.Go.Tags, "enterprise") {
		t.Fatal("unrelated environment or tags lost")
	}
	raw, _ := os.ReadFile(path)
	if !bytes.Contains(raw, []byte("# user configuration")) {
		t.Fatal("comment lost")
	}
	if runtime.GOOS != "windows" {
		info, _ := os.Stat(path)
		if info.Mode().Perm() != 0640 {
			t.Fatal("permissions lost")
		}
	}
	config.Backend = "cgo"
	config.ExperimentalPureGo = false
	response = wakeRequest(t, New(), http.MethodPost, "/api/wake/setup", config)
	if response.Code != 200 {
		t.Fatal(response.Body.String())
	}
	loaded, err = manifest.Load(root, "setup")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Config.Build.Environment["CGO_ENABLED"] != "" || slices.Contains(loaded.Config.Build.Go.Tags, "purego") || slices.Contains(loaded.Config.Targets.Linux.AMD64.Tags, "purego") {
		t.Fatal("CGo switch retained purego backend")
	}
}
func TestWakeSetupRejectsUnsafeOrInvalidRequestsWithoutWriting(t *testing.T) {
	root := wakeTestProject(t)
	for _, mutate := range []func(*WakeSetupConfig){
		func(c *WakeSetupConfig) { c.ExperimentalPureGo = false }, func(c *WakeSetupConfig) { c.ApplicationType = "unknown" }, func(c *WakeSetupConfig) { c.Targets[0].Arch = "invalid" }, func(c *WakeSetupConfig) { c.Targets[0].Formats = []string{"msix"} }, func(c *WakeSetupConfig) { c.Targets = append(c.Targets, c.Targets[0]) }, func(c *WakeSetupConfig) {
			c.Targets[0] = WakeSetupTarget{OS: "ios", Arch: "arm64", Formats: []string{"ipa"}}
		}, func(c *WakeSetupConfig) { c.Project = nil },
		func(c *WakeSetupConfig) { c.Sign = true; c.Targets[0].Formats = []string{"appimage"} },
		func(c *WakeSetupConfig) { c.Sign = true; c.Targets[0].Formats = []string{} },
	} {
		config := wakeTestConfig()
		mutate(&config)
		response := wakeRequest(t, New(), http.MethodPost, "/api/wake/setup", config)
		if response.Code != 400 {
			t.Fatalf("invalid request accepted: %s", response.Body.String())
		}
	}
	w := New()
	w.origin = "http://127.0.0.1:1234"
	mux := http.NewServeMux()
	w.setupRoutes(mux)
	for _, test := range []struct {
		body, origin string
		status       int
	}{{`{"unknown":true}`, "", 400}, {strings.Repeat(" ", 1<<20) + `{}`, "", 400}, {`{}`, "https://attacker.invalid", 403}, {`{} {}`, "", 400}} {
		request := httptest.NewRequest(http.MethodPost, "/api/wake/setup", strings.NewReader(test.body))
		request.Header.Set("Origin", test.origin)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		if response.Code != test.status {
			t.Fatalf("guard status %d: %s", response.Code, response.Body.String())
		}
	}
	if _, err := os.Stat(filepath.Join(root, "wails.yaml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("rejected request published manifest")
	}
}
func TestWakeSetupRefusesLegacyCutoverAndInvalidManifest(t *testing.T) {
	for _, test := range []struct{ name, data string }{{"wails.yaml", "info: {productName: Legacy}\n"}, {"wails.yaml", "version: [broken\n"}, {"Taskfile.yml", "version: '3'\ntasks: {}\n"}} {
		t.Run(test.name+test.data[:4], func(t *testing.T) {
			root := wakeTestProject(t)
			path := filepath.Join(root, test.name)
			if err := os.WriteFile(path, []byte(test.data), 0600); err != nil {
				t.Fatal(err)
			}
			response := wakeRequest(t, New(), http.MethodPost, "/api/wake/setup", wakeTestConfig())
			if response.Code != 400 {
				t.Fatal(response.Body.String())
			}
			after, err := os.ReadFile(path)
			if err != nil || string(after) != test.data {
				t.Fatal("existing config changed")
			}
		})
	}
}

func TestWakeSetupUsesGlobalSigningForChecksWithoutPersistingKeyPaths(t *testing.T) {
	root := wakeTestProject(t)
	keyPath := filepath.Join(t.TempDir(), "owned-test.pfx")
	if err := os.WriteFile(keyPath, []byte("owned inert credential fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	defaults := DefaultGlobalDefaults()
	defaults.Signing.Windows.CertificatePath = keyPath
	if err := SaveGlobalDefaults(defaults); err != nil {
		t.Fatal(err)
	}
	config := wakeTestConfig()
	config.Sign = true
	config.Prebuilt = true
	config.Targets = []WakeSetupTarget{{OS: "windows", Arch: "amd64", Formats: []string{"msix"}}}
	response := wakeRequest(t, New(), http.MethodPost, "/api/wake/setup/requirements", config)
	var checks wakeRequirementsResponse
	if err := json.Unmarshal(response.Body.Bytes(), &checks); err != nil {
		t.Fatal(response.Body.String())
	}
	if len(checks.Errors) != 0 {
		t.Fatal(checks.Errors)
	}
	for _, requirement := range checks.Requirements {
		if requirement.Name == "signtool.exe" {
			t.Fatal("exported global PFX should use portable signing")
		}
	}
	response = wakeRequest(t, New(), http.MethodPost, "/api/wake/setup", config)
	if response.Code != 200 {
		t.Fatal(response.Body.String())
	}
	raw, err := os.ReadFile(filepath.Join(root, "wails.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(raw, []byte(keyPath)) || bytes.Contains(raw, []byte("owned inert credential fixture")) {
		t.Fatal("global credentials leaked into manifest")
	}
	loaded, err := manifest.Load(root, "setup")
	if err != nil {
		t.Fatal(err)
	}
	if !loaded.Config.Selected.Targets[0].Sign || !loaded.Config.Signing.Windows.Enabled || loaded.Config.Signing.Windows.Certificate != "" {
		t.Fatal("signing intent or global fallback not persisted correctly")
	}
}
