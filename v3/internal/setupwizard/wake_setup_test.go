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

func wakeRequest(t *testing.T, w *Wizard, method, path string, value any) *httptest.ResponseRecorder {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	w.setupRoutes(mux)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(method, path, bytes.NewReader(data)))
	return response
}

func wakeInitTestData(t *testing.T) InitData {
	t.Helper()
	t.Setenv("WAILS_EXP_USE_WAKE", "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	return InitData{ProjectName: "probe", TemplateName: "vanilla-js", ProductName: "Probe", ProductIdentifier: "org.example.probe", ProductVersion: "1.0.0", BaseDir: t.TempDir()}
}

func wakeInitTestConfig(data InitData) WakeSetupConfig {
	return WakeSetupConfig{ApplicationType: "gui", Backend: "purego", ExperimentalPureGo: true, Targets: []WakeSetupTarget{{OS: "windows", Arch: "amd64", Formats: []string{"msix"}}}, Project: initProjectState(data)}
}

func TestWakeChoicesAreAvailableOnlyDuringInit(t *testing.T) {
	data := wakeInitTestData(t)
	for _, value := range []string{"", "false", "1"} {
		t.Setenv("WAILS_EXP_USE_WAKE", value)
		w := NewInitWizard(data)
		response := wakeRequest(t, w, http.MethodGet, "/api/init", nil)
		var got InitData
		if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil || !got.WakeEnabled || got.WakeSetup == nil || len(got.WakeSetup.Catalog) == 0 {
			t.Fatalf("init seed: %s", response.Body.String())
		}
	}
	global := wakeRequest(t, New(), http.MethodGet, "/api/wake/setup", nil)
	var state wakeSetupResponse
	if err := json.Unmarshal(global.Body.Bytes(), &state); err != nil || state.Enabled || len(state.Catalog) != 0 {
		t.Fatalf("global setup exposes project configuration: %s", global.Body.String())
	}
	for _, w := range []*Wizard{New(), NewInitWizard(data)} {
		response := wakeRequest(t, w, http.MethodPost, "/api/wake/setup", wakeInitTestConfig(data))
		if response.Code != http.StatusForbidden {
			t.Fatalf("old save endpoint remains writable: %s", response.Body.String())
		}
	}
	response := wakeRequest(t, New(), http.MethodPost, "/api/wake/setup/requirements", wakeInitTestConfig(data))
	if response.Code != http.StatusForbidden {
		t.Fatal("global setup accepted project preview")
	}
	if err := os.Unsetenv("WAILS_EXP_USE_WAKE"); err != nil {
		t.Fatal(err)
	}
	off := NewInitWizard(data)
	if off.initData.WakeEnabled || off.initData.WakeSetup != nil || off.initData.BuildConfigFile != "build/config.yml" {
		t.Fatal("flag-off init enabled Wake")
	}
}

func TestWakeInitCreateValidatesWithoutWritingAndUsesAuthoritativeMetadata(t *testing.T) {
	data := wakeInitTestData(t)
	w := NewInitWizard(data)
	config := wakeInitTestConfig(data)
	config.Project.Info.ProductName = "Injected metadata"
	data.Wake = &config
	data.BaseDir = "/untrusted"
	response := wakeRequest(t, w, http.MethodPost, "/api/init/create", data)
	if response.Code != 200 || w.initResult == nil {
		t.Fatalf("create: %s", response.Body.String())
	}
	if w.initResult.BaseDir != w.initData.BaseDir || w.initResult.Wake.Project.Info.ProductName != "Probe" {
		t.Fatal("client overrode authoritative root or metadata")
	}
	if _, err := os.Stat(filepath.Join(w.initData.BaseDir, "probe")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("init request wrote files before scaffolding")
	}
}

func TestWakeInitRejectsUnsupportedChoicesBeforeScaffolding(t *testing.T) {
	for _, mutate := range []func(*InitData){
		func(d *InitData) { d.Wake.Prebuilt = true },
		func(d *InitData) { d.Wake.ApplicationType = "cli" },
		func(d *InitData) { d.Wake.ApplicationType = "server" },
		func(d *InitData) { d.Wake.ExperimentalPureGo = false },
		func(d *InitData) { d.Wake.Targets = nil },
		func(d *InitData) { d.Wake.Targets[0].Formats = []string{"deb"} },
		func(d *InitData) { d.ProjectName = "../outside" },
		func(d *InitData) { d.ProductIdentifier = "" },
		func(d *InitData) { d.Wake = nil },
	} {
		data := wakeInitTestData(t)
		config := wakeInitTestConfig(data)
		data.Wake = &config
		w := NewInitWizard(data)
		mutate(&data)
		response := wakeRequest(t, w, http.MethodPost, "/api/init/create", data)
		if response.Code != 400 || w.initResult != nil {
			t.Fatalf("invalid config accepted: %s", response.Body.String())
		}
		entries, err := os.ReadDir(w.initData.BaseDir)
		if err != nil || len(entries) != 0 {
			t.Fatal("invalid init wrote files", err)
		}
	}
}

func TestWakeInitFlagOffRejectsInjectedChoicesButPreservesLegacyInit(t *testing.T) {
	data := wakeInitTestData(t)
	config := wakeInitTestConfig(data)
	data.Wake = &config
	if err := os.Unsetenv("WAILS_EXP_USE_WAKE"); err != nil {
		t.Fatal(err)
	}
	w := NewInitWizard(data)
	response := wakeRequest(t, w, http.MethodPost, "/api/init/create", data)
	if response.Code != 403 || w.initResult != nil {
		t.Fatal("flag-off accepted Wake config")
	}
	data.Wake = nil
	response = wakeRequest(t, w, http.MethodPost, "/api/init/create", data)
	if response.Code != 200 || w.initResult == nil {
		t.Fatalf("legacy init failed: %s", response.Body.String())
	}
}

func TestWakeInitRejectsCrossOriginAndUnknownFields(t *testing.T) {
	data := wakeInitTestData(t)
	config := wakeInitTestConfig(data)
	data.Wake = &config
	for _, path := range []string{"/api/init/create", "/api/wake/setup/requirements"} {
		for _, origin := range []string{"https://foreign.example", ""} {
			w := NewInitWizard(data)
			w.origin = "http://127.0.0.1:12345"
			request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{"unknown":true}`))
			if origin != "" {
				request.Header.Set("Origin", origin)
			}
			response := httptest.NewRecorder()
			mux := http.NewServeMux()
			w.setupRoutes(mux)
			mux.ServeHTTP(response, request)
			if response.Code == 200 {
				t.Fatalf("unsafe request accepted: %s", path)
			}
		}
	}
}

func TestApplyWakeInitPreservesScaffoldAndDrivesBuild(t *testing.T) {
	data := wakeInitTestData(t)
	root := filepath.Join(data.BaseDir, data.ProjectName)
	if err := os.Mkdir(root, 0700); err != nil {
		t.Fatal(err)
	}
	initial, err := initProjectState(data).ManifestState(manifest.Project{})
	if err != nil {
		t.Fatal(err)
	}
	if err := manifest.WriteInitial(root, initial); err != nil {
		t.Fatal(err)
	}
	config := wakeInitTestConfig(data)
	if err := ApplyWakeInit(root, config); err != nil {
		t.Fatal(err)
	}
	loaded, err := manifest.Load(root, "setup")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Config.Project.ProductName != "Probe" || loaded.Config.Frontend.Build[0] != "npm" || loaded.Config.Selected.Targets[0].Formats[0] != "msix" {
		t.Fatal("scaffold metadata/frontend or selections lost")
	}
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, node := range plan.Nodes {
		if spec, ok := node.Spec.(pipeline.CompileSpec); ok {
			found = true
			if spec.CGOEnabled != "0" || !slices.Contains(spec.Tags, "purego") {
				t.Fatalf("backend choices not used: %+v", spec)
			}
		}
	}
	if !found {
		t.Fatal("missing compile node")
	}
}

func TestWakeInitPureGoRequirementsDoNotExecuteCompilersOrContainers(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX tool probes")
	}
	data := wakeInitTestData(t)
	tools := t.TempDir()
	marker := filepath.Join(tools, "executed")
	for _, name := range []string{"go", "npm", "docker", "podman"} {
		if err := os.WriteFile(filepath.Join(tools, name), []byte("#!/bin/sh\ntouch '"+marker+"'\n"), 0700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", tools)
	config := wakeInitTestConfig(data)
	response := wakeRequest(t, NewInitWizard(data), http.MethodPost, "/api/wake/setup/requirements", config)
	var result wakeRequirementsResponse
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil || len(result.Errors) > 0 {
		t.Fatalf("preview failed: %s", response.Body.String())
	}
	for _, dep := range result.Requirements {
		if slices.Contains([]string{"docker", "podman", "cc", "pkg-config"}, dep.Name) {
			t.Fatalf("pure Go requested %s", dep.Name)
		}
	}
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("requirements executed a compiler or container")
	}
	if _, err := os.Stat(filepath.Join(data.BaseDir, data.ProjectName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("requirements created project")
	}
}

func TestWakeInitUsesGlobalSigningWithoutCopyingCredentials(t *testing.T) {
	data := wakeInitTestData(t)
	keyPath := filepath.Join(t.TempDir(), "owned-test.pfx")
	if err := os.WriteFile(keyPath, []byte("owned inert credential fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	defaults := DefaultGlobalDefaults()
	defaults.Signing.Windows.CertificatePath = keyPath
	if err := SaveGlobalDefaults(defaults); err != nil {
		t.Fatal(err)
	}
	config := wakeInitTestConfig(data)
	config.Sign = true
	response := wakeRequest(t, NewInitWizard(data), http.MethodPost, "/api/wake/setup/requirements", config)
	var checks wakeRequirementsResponse
	if err := json.Unmarshal(response.Body.Bytes(), &checks); err != nil || len(checks.Errors) > 0 {
		t.Fatalf("global signing preview: %s", response.Body.String())
	}
	for _, dep := range checks.Requirements {
		if dep.Name == "signtool.exe" {
			t.Fatal("global PFX should use portable signing")
		}
	}
	root := filepath.Join(data.BaseDir, data.ProjectName)
	if err := os.Mkdir(root, 0700); err != nil {
		t.Fatal(err)
	}
	initial, err := initProjectState(data).ManifestState(manifest.Project{})
	if err != nil {
		t.Fatal(err)
	}
	if err := manifest.WriteInitial(root, initial); err != nil {
		t.Fatal(err)
	}
	if err := ApplyWakeInit(root, config); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(root, manifest.Filename))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(raw, []byte(keyPath)) || bytes.Contains(raw, []byte("owned inert credential fixture")) {
		t.Fatal("credential copied into project")
	}
	loaded, err := manifest.Load(root, "setup")
	if err != nil {
		t.Fatal(err)
	}
	if !loaded.Config.Selected.Targets[0].Sign || !loaded.Config.Signing.Windows.Enabled || loaded.Config.Signing.Windows.Certificate != "" {
		t.Fatal("signing intent or global fallback lost")
	}
}
