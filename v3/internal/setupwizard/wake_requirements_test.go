package setupwizard

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestWakeToolRequirementsFollowResolvedCompiler(t *testing.T) {
	for _, test := range []struct {
		name, host, arch, target, cgo string
		console                       bool
		environment                   map[string]string
		installed, want, absent       []string
	}{
		{name: "native CGo", host: "linux", arch: "amd64", target: "linux", cgo: "1", installed: []string{"go", "cc"}, want: []string{"go", "cc", "pkg-config"}, absent: []string{"docker", "npm"}},
		{name: "pure Go cross", host: "linux", arch: "amd64", target: "darwin", cgo: "0", want: []string{"go"}, absent: []string{"docker", "cc", "clang", "zig", "npm"}},
		{name: "CGo cross container", host: "linux", arch: "amd64", target: "darwin", cgo: "1", want: []string{"docker"}, absent: []string{"go", "npm", "cc"}},
		{name: "CGo cross Zig", host: "linux", arch: "amd64", target: "windows", cgo: "1", installed: []string{"zig"}, want: []string{"go", "zig"}, absent: []string{"docker", "npm"}},
		{name: "explicit cross compiler", host: "linux", arch: "amd64", target: "darwin", cgo: "1", environment: map[string]string{"CC": "target-clang"}, want: []string{"go", "target-clang"}, absent: []string{"docker", "npm", "xcrun"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			plan := pipeline.Plan{Nodes: map[pipeline.NodeKey]pipeline.Node{"compile": {Label: "Compile", Kind: pipeline.CompileApplication, Spec: pipeline.CompileSpec{TargetOS: test.target, TargetArch: "amd64", CGOEnabled: test.cgo, Environment: test.environment, Console: test.console}}}}
			lookup := func(name string) (string, error) {
				if slices.Contains(test.installed, name) {
					return "/tools/" + name, nil
				}
				return "", errors.New("not installed")
			}
			requirements, warnings := wakeToolRequirements(plan, t.TempDir(), test.host, test.arch, lookup)
			if len(warnings) != 0 {
				t.Fatal(warnings)
			}
			names := []string{}
			for _, r := range requirements {
				names = append(names, r.Name)
			}
			for _, want := range test.want {
				if !slices.Contains(names, want) {
					t.Fatalf("missing %s: %v", want, names)
				}
			}
			for _, absent := range test.absent {
				if slices.Contains(names, absent) {
					t.Fatalf("unexpected %s: %v", absent, names)
				}
			}
		})
	}
}
func TestWakePrebuiltRequirementsSkipCompilersFrontendAndContainerProbes(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("owned shell probe uses POSIX")
	}
	root := wakeTestProject(t)
	tools := filepath.Join(root, "tools")
	if err := os.Mkdir(tools, 0700); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(root, "unexpected-tool-execution")
	for _, name := range []string{"go", "npm", "docker", "podman"} {
		script := "#!/bin/sh\nprintf invoked > '" + marker + "'\nexit 0\n"
		if err := os.WriteFile(filepath.Join(tools, name), []byte(script), 0700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", tools)
	config := wakeTestConfig()
	config.Prebuilt = true
	response := wakeRequest(t, New(), http.MethodPost, "/api/wake/setup/requirements", config)
	var result wakeRequirementsResponse
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(response.Body.String())
	}
	if len(result.Errors) != 0 {
		t.Fatal(result.Errors)
	}
	if result.Command != "wails3 build setup --skip-build" {
		t.Fatal(result.Command)
	}
	for _, r := range result.Requirements {
		if slices.Contains([]string{"go", "npm", "docker", "podman", "cc", "pkg-config"}, r.Name) {
			t.Fatalf("prebuilt requires %s", r.Name)
		}
	}
	if len(result.Requirements) != 1 || !strings.HasPrefix(result.Requirements[0].Name, "Prebuilt ") {
		t.Fatalf("missing prebuilt input: %+v", result.Requirements)
	}
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("requirement preview executed a tool")
	}
	config.Prebuilt = false
	response = wakeRequest(t, New(), http.MethodPost, "/api/wake/setup/requirements", config)
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(response.Body.String())
	}
	if len(result.Errors) != 0 {
		t.Fatal(result.Errors)
	}
	for _, r := range result.Requirements {
		if slices.Contains([]string{"npm", "docker", "podman", "cc", "pkg-config"}, r.Name) {
			t.Fatalf("headless pure Go requires %s", r.Name)
		}
	}
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("pure Go preview executed a compiler/container")
	}
	if _, err := os.Stat(filepath.Join(root, "wails.yaml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("preview wrote manifest")
	}
}
func TestWakeStandardBackendDoesNotForceCGoOnWindowsOrHeadless(t *testing.T) {
	for _, kind := range []string{"gui", "cli", "server"} {
		t.Run(kind, func(t *testing.T) {
			root := wakeTestProject(t)
			config := wakeTestConfig()
			config.ApplicationType = kind
			config.Backend = "cgo"
			config.ExperimentalPureGo = false
			config.Targets = []WakeSetupTarget{{OS: "windows", Arch: "amd64", Formats: []string{}}}
			response := wakeRequest(t, New(), http.MethodPost, "/api/wake/setup", config)
			if response.Code != 200 {
				t.Fatal(response.Body.String())
			}
			loaded, err := manifest.Load(root, "setup")
			if err != nil {
				t.Fatal(err)
			}
			plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
			if err != nil {
				t.Fatal(err)
			}
			for _, node := range plan.Nodes {
				if spec, ok := node.Spec.(pipeline.CompileSpec); ok && spec.CGOEnabled != "0" {
					t.Fatalf("standard Windows mode: %+v", spec)
				}
			}
			get := wakeRequest(t, New(), http.MethodGet, "/api/wake/setup", nil)
			var state wakeSetupResponse
			if err := json.Unmarshal(get.Body.Bytes(), &state); err != nil || state.Error != "" || state.Config.Backend != "cgo" {
				t.Fatalf("standard mode hydration: %s", get.Body.String())
			}
		})
	}
}
func TestWakeHydratesTargetBackendAndPreservesGUICompanion(t *testing.T) {
	root := wakeTestProject(t)
	if err := os.Mkdir(filepath.Join(root, "gui"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "gui/main.go"), []byte("package main\nfunc main(){}"), 0600); err != nil {
		t.Fatal(err)
	}
	original := `version: 3
project: {name: probe, product_name: Probe, identifier: org.example.probe, version: 1.0.0, kind: cli}
binaries:
  helper: {package: ./gui, kind: gui}
targets:
  linux/amd64: {environment: {CGO_ENABLED: "0"}, tags: [purego]}
profiles:
  setup:
    targets:
      linux/amd64: {formats: [deb]}
`
	if err := os.WriteFile(filepath.Join(root, "wails.yaml"), []byte(original), 0600); err != nil {
		t.Fatal(err)
	}
	get := wakeRequest(t, New(), http.MethodGet, "/api/wake/setup", nil)
	var state wakeSetupResponse
	if err := json.Unmarshal(get.Body.Bytes(), &state); err != nil {
		t.Fatal(err)
	}
	if state.Config.Backend != "purego" || !state.Config.ExperimentalPureGo {
		t.Fatalf("target backend not hydrated: %s", get.Body.String())
	}
	response := wakeRequest(t, New(), http.MethodPost, "/api/wake/setup", state.Config)
	if response.Code != 200 {
		t.Fatal(response.Body.String())
	}
	loaded, err := manifest.Load(root, "setup")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Config.Frontend.Disabled || len(loaded.Config.Binaries) != 1 || loaded.Config.Binaries[0].Kind != "gui" {
		t.Fatal("GUI companion frontend lost")
	}
	state.Config.Backend = "cgo"
	state.Config.ExperimentalPureGo = false
	response = wakeRequest(t, New(), http.MethodPost, "/api/wake/setup", state.Config)
	if response.Code != 200 {
		t.Fatal(response.Body.String())
	}
	loaded, err = manifest.Load(root, "setup")
	if err != nil {
		t.Fatal(err)
	}
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	if err != nil {
		t.Fatal(err)
	}
	foundGUI := false
	for _, node := range plan.Nodes {
		if spec, ok := node.Spec.(pipeline.CompileSpec); ok {
			if !spec.Console {
				foundGUI = true
				if spec.CGOEnabled != "1" {
					t.Fatal("standard GUI companion was silently switched to experimental pure Go")
				}
			} else if spec.CGOEnabled != "0" {
				t.Fatal("headless primary unnecessarily requires CGo")
			}
		}
	}
	if !foundGUI {
		t.Fatal("GUI companion missing from plan")
	}
}
