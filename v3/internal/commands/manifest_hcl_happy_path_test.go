//go:build !windows

package commands

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

const hclBuildFixture = `version = 3

project {
  name = "hello"
  product_name = "Hello"
  company = "Example"
  identifier = "com.example.hello"
  version = "1.2.3"
}

frontend {
  directory = "frontend"
  install = ["npm", "install"]
  build = ["npm", "run", "bundle"]
  dev = ["npm", "run", "dev"]
  output = "frontend/dist"
}

build {
  output = "bin"
  trim_path = true
  strip = true
}
`

func TestBuildPlanUsesHCLProfileAndDoesNotInvokeTaskfile(t *testing.T) {
	root := t.TempDir()
	hcl := hclBuildFixture + `
profile "release" {
  target "linux/amd64" {
    formats = ["deb"]
  }
  target "windows/amd64" {
    formats = ["nsis"]
    sign = true
  }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.yml"), []byte("this is ignored\n"), 0o644))
	t.Chdir(root)

	previous := runTaskFunc
	defer func() { runTaskFunc = previous }()
	called := false
	runTaskFunc = func(*RunTaskOptions, []string) error {
		called = true
		return nil
	}

	output := captureHCLStdout(t, func() {
		err := Build(&flags.Build{Plan: true}, []string{"release"})
		require.NoError(t, err)
	})
	assert.False(t, called, "wails.hcl must select the native route even when a Taskfile exists")
	assert.Contains(t, output, "Targets: linux/amd64 · windows/amd64")
	assert.Contains(t, output, "No files will be changed because --plan was used.")
	assert.NotContains(t, output, "Taskfile")
}

func TestBuildJSONPlanIsDeterministicAndDescribesArtifacts(t *testing.T) {
	root := t.TempDir()
	hcl := hclBuildFixture + `
profile "release" {
  target "linux/amd64" {
    formats = ["deb"]
  }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	t.Chdir(root)

	first := captureHCLStdout(t, func() {
		require.NoError(t, Build(&flags.Build{Plan: true, JSON: true}, []string{"release"}))
	})
	second := captureHCLStdout(t, func() {
		require.NoError(t, Build(&flags.Build{Plan: true, JSON: true}, []string{"release"}))
	})
	assert.Equal(t, first, second)
	var plan struct {
		SchemaVersion int `json:"schema_version"`
		Request       struct {
			Command string `json:"command"`
			Profile string `json:"profile"`
		} `json:"request"`
		Operations []struct {
			Stage    string `json:"stage"`
			Decision string `json:"decision"`
		} `json:"operations"`
		Artifacts []struct {
			Kind string `json:"kind"`
			Path string `json:"path"`
		} `json:"artifacts"`
	}
	require.NoError(t, json.Unmarshal([]byte(first), &plan))
	assert.Equal(t, 1, plan.SchemaVersion)
	assert.Equal(t, "build", plan.Request.Command)
	assert.Equal(t, "release", plan.Request.Profile)
	assert.NotEmpty(t, plan.Operations)
	assert.Contains(t, mapStrings(plan.Operations, func(v struct {
		Stage    string `json:"stage"`
		Decision string `json:"decision"`
	}) string {
		return v.Stage
	}), "package")
	for _, operation := range plan.Operations {
		assert.Equal(t, "run", operation.Decision)
	}
	assert.Contains(t, mapArtifacts(plan.Artifacts), "deb")
}

func TestBuildPlanAcceptsCommaSeparatedTargets(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hclBuildFixture), 0o644))
	t.Chdir(root)
	output := captureHCLStdout(t, func() {
		require.NoError(t, Build(&flags.Build{Plan: true, Targets: "windows/amd64,linux/arm64"}, nil))
	})
	assert.Contains(t, output, "Targets: linux/arm64 · windows/amd64")
}

func TestBuildPlanPrintsRequestedFormats(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hclBuildFixture), 0o644))
	output := captureHCLStdout(t, func() {
		require.NoError(t, Build(&flags.Build{Plan: true, Targets: "linux/amd64", Formats: "deb"}, nil))
	})
	assert.Contains(t, output, "Formats: deb")
}

func TestDevPlanUsesTheHCLBuildPipeline(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hclBuildFixture), 0o644))
	t.Chdir(root)
	output := captureHCLStdout(t, func() {
		require.NoError(t, Dev(&DevOptions{Plan: true}))
	})
	assert.Contains(t, output, "Plan:")
	assert.Contains(t, output, "Targets: "+runtime.GOOS+"/"+runtime.GOARCH)
	assert.Contains(t, output, "No files will be changed because --plan was used.")
}

func TestBuildRejectsJSONWithoutPlanBeforeLoadingProject(t *testing.T) {
	t.Chdir(t.TempDir())
	err := Build(&flags.Build{JSON: true}, nil)
	assert.ErrorContains(t, err, "--json requires --plan")
}

func TestHCLConfigCommandsUseTheResolvedManifest(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := `version = 3
project {
  name = "config"
  product_name = "Config"
  identifier = "com.example.config"
  version = "1.0.0"
}
profile "release" {
  target "linux/amd64" {}
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))

	check := captureHCLStdout(t, func() {
		require.NoError(t, ConfigCheck(&ConfigOptions{Profile: "release"}))
	})
	assert.Contains(t, check, "wails.hcl is valid (com.example.config)")

	textConfig := captureHCLStdout(t, func() {
		require.NoError(t, ConfigShow(&ConfigOptions{Profile: "release"}))
	})
	assert.Contains(t, textConfig, "name = \"config\"")
	assert.Contains(t, textConfig, "output = \"bin\"")

	jsonConfig := captureHCLStdout(t, func() {
		require.NoError(t, ConfigShow(&ConfigOptions{Profile: "release", JSON: true}))
	})
	var decoded struct {
		Profile string `json:"profile"`
		Project struct {
			Identifier string `json:"identifier"`
		} `json:"project"`
	}
	require.NoError(t, json.Unmarshal([]byte(jsonConfig), &decoded))
	assert.Equal(t, "release", decoded.Profile)
	assert.Equal(t, "com.example.config", decoded.Project.Identifier)

	require.NoError(t, Eject(&EjectOptions{}, nil))
	ejected := readTestFile(t, filepath.Join(root, manifest.EjectedFilename))
	assert.Contains(t, ejected, "Generated by Wails CLI")
}

func TestHCLSignCommandRunsTheCompleteProfilePipeline(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := `version = 3
project {
  name = "command"
  product_name = "Command"
  identifier = "com.example.command"
  version = "1.0.0"
}
frontend {
  directory = "frontend"
  install = ["npm", "install"]
  build = ["npm", "run", "bundle"]
  dev = ["npm", "run", "dev"]
  output = "frontend/dist"
}
linux {
  signing {
    identity = "release"
    certificate = "release-key"
  }
}
profile "release" {
  target "linux/amd64" {
    formats = ["deb"]
    sign = true
  }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/command\n\ngo 1.24\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "package.json"), []byte(`{"scripts":{"bundle":"bundle"}}`), 0o644))
	fakeTools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "npm"), []byte(`#!/bin/sh
if [ "$1" = "install" ]; then mkdir -p node_modules; fi
if [ "$1" = "run" ]; then mkdir -p dist; printf bundle > dist/index.html; fi
`), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "dpkg-sig"), []byte("#!/bin/sh\nexit 0\n"), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))

	require.NoError(t, SignWrapper(&flags.SignWrapper{Profile: "release"}, nil))
	assert.FileExists(t, filepath.Join(root, "bin", "linux-amd64", "command_1.0.0_amd64.deb.signed"))
	assert.FileExists(t, filepath.Join(root, "frontend", "dist", "index.html"))
}

func TestHCLPackageCommandBuildsTheRequestedTargetAndFormat(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hclBuildFixture), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/package\n\ngo 1.24\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "package.json"), []byte(`{"scripts":{"bundle":"bundle"}}`), 0o644))
	fakeTools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "npm"), []byte(`#!/bin/sh
if [ "$1" = "install" ]; then mkdir -p node_modules; fi
if [ "$1" = "run" ]; then mkdir -p dist; printf bundle > dist/index.html; fi
`), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))

	require.NoError(t, Package(&flags.Package{Targets: "linux/amd64", Formats: "deb"}, nil))
	assert.FileExists(t, filepath.Join(root, "bin", "hello_1.2.3_amd64.deb"))
}

func TestHCLLinuxPackageTemplateReachesTheOwnedPackageWorkspace(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "packaging"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "packaging", "nfpm.yaml.tmpl"), []byte("name: {{.Project.BinaryName}}\nsource: {{.Paths.Binary}}\n"), 0o644))
	hcl := `version = 3
project {
  name = "linux-template"
  product_name = "Linux Template"
  identifier = "com.example.linuxtemplate"
  version = "1.0.0"
}
package "deb" {
  template = "packaging/nfpm.yaml.tmpl"
}
profile "release" {
  target "linux/amd64" { formats = ["deb"] }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "package"})
	require.NoError(t, err)
	spec := plan.Nodes[pipeline.NodeKey("package:linux/amd64:deb")].Spec.(pipeline.PackageSpec)
	configPath, err := (&manifestHandler{root: root, config: loaded.Config}).prepareLinuxPackageConfig(spec)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, ".wails", "build", "release", "linux-amd64", "package", "deb", "nfpm.yaml"), configPath)
	data := readTestFile(t, configPath)
	assert.Contains(t, data, "name: linux-template")
	assert.Contains(t, data, filepath.Join(root, "bin", "linux-amd64", "linux-template"))
}

func TestHCLBuildExecutesFrontendAndRealGoCompile(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hclBuildFixture), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/hello\n\ngo 1.24\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.yml"), []byte("must remain untouched\n"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend"), 0o755))
	frontendPackage := filepath.Join(root, "frontend", "package.json")
	require.NoError(t, os.WriteFile(frontendPackage, []byte(`{"scripts":{"bundle":"bundle"}}`), 0o644))
	userAsset := filepath.Join(root, "frontend", "src", "user-owned.js")
	require.NoError(t, os.MkdirAll(filepath.Dir(userAsset), 0o755))
	require.NoError(t, os.WriteFile(userAsset, []byte("export const userOwned = true;\n"), 0o644))

	bin := t.TempDir()
	fakeNPM := filepath.Join(bin, "npm")
	require.NoError(t, os.WriteFile(fakeNPM, []byte("#!/bin/sh\n"+
		"printf '%s\\n' \"$*\" >> \"$NPM_LOG\"\n"+
		"if [ \"$1\" = \"install\" ]; then mkdir -p node_modules; fi\n"+
		"if [ \"$1\" = \"run\" ]; then mkdir -p dist; printf 'bundle' > dist/index.html; fi\n"), 0o755))
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	invocationLog := filepath.Join(root, ".npm-invocation")
	t.Setenv("NPM_LOG", invocationLog)
	previousTask := runTaskFunc
	taskCalled := false
	runTaskFunc = func(*RunTaskOptions, []string) error {
		taskCalled = true
		return nil
	}
	t.Cleanup(func() { runTaskFunc = previousTask })

	first, err := runManifestPipelineResult(manifestRunOptions{Verb: "build", TargetOS: runtime.GOOS, TargetArch: runtime.GOARCH})
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(root, "bin", "hello"))
	assert.FileExists(t, filepath.Join(root, "frontend", "dist", "index.html"))
	assert.Equal(t, "install\nrun bundle\n", readTestFile(t, invocationLog))
	assert.Equal(t, "must remain untouched\n", readTestFile(t, filepath.Join(root, "Taskfile.yml")))
	assert.False(t, taskCalled, "the HCL-first pipeline must not invoke Taskfile execution")
	assert.Equal(t, "export const userOwned = true;\n", readTestFile(t, userAsset))
	assert.DirExists(t, filepath.Join(root, ".wails"))

	second, err := runManifestPipelineResult(manifestRunOptions{Verb: "build", TargetOS: runtime.GOOS, TargetArch: runtime.GOARCH})
	require.NoError(t, err)
	assert.Equal(t, "install\nrun bundle\n", readTestFile(t, invocationLog), "a warm HCL build should use the cache")
	for key, result := range second.Results {
		assert.NotEqual(t, cache.LookupMiss, result.Status, key)
	}
	assert.NotEmpty(t, first.Results)

	require.NoError(t, os.WriteFile(userAsset, []byte("export const userOwned = false;\n"), 0o644))
	third, err := runManifestPipelineResult(manifestRunOptions{Verb: "build", TargetOS: runtime.GOOS, TargetArch: runtime.GOARCH})
	require.NoError(t, err)
	assert.Equal(t, "install\nrun bundle\nrun bundle\n", readTestFile(t, invocationLog), "frontend source changes must invalidate the frontend stage")
	assert.Equal(t, cache.LookupMiss, third.Results[pipeline.NodeKey("frontend:build")].Status)

	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\n\nvar changed = true\n\nfunc main() {}\n"), 0o644))
	fourth, err := runManifestPipelineResult(manifestRunOptions{Verb: "build", TargetOS: runtime.GOOS, TargetArch: runtime.GOARCH})
	require.NoError(t, err)
	assert.Equal(t, "install\nrun bundle\nrun bundle\n", readTestFile(t, invocationLog), "Go-only changes must not invalidate the frontend stage")
	assert.Equal(t, cache.LookupMiss, fourth.Results[pipeline.NodeKey("target:"+runtime.GOOS+"/"+runtime.GOARCH+":compile")].Status)

	_, err = runManifestPipelineResult(manifestRunOptions{Verb: "build", TargetOS: runtime.GOOS, TargetArch: runtime.GOARCH, Force: true})
	require.NoError(t, err)
	assert.Equal(t, "install\nrun bundle\nrun bundle\ninstall\nrun bundle\n", readTestFile(t, invocationLog), "--force must rebuild cached HCL stages")
}

func TestHCLDevBuildExecutesTheDevelopmentPipeline(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hclBuildFixture), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/dev\n\ngo 1.24\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "package.json"), []byte(`{"scripts":{"bundle":"bundle"}}`), 0o644))

	bin := t.TempDir()
	fakeNPM := filepath.Join(bin, "npm")
	require.NoError(t, os.WriteFile(fakeNPM, []byte("#!/bin/sh\nprintf '%s:%s\n' \"$*\" \"$PRODUCTION\" >> \"$NPM_LOG\"\nif [ \"$1\" = \"install\" ]; then mkdir -p node_modules; fi\nif [ \"$1\" = \"run\" ]; then mkdir -p dist; printf dev > dist/index.html; fi\n"), 0o755))
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	log := filepath.Join(root, ".npm-dev-invocation")
	t.Setenv("NPM_LOG", log)
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	_, err = runManifestDevBuild(t.Context(), &DevOptions{}, loaded, runtime.GOOS, runtime.GOARCH, "http://127.0.0.1:9999", 9999)
	require.NoError(t, err)
	assert.Equal(t, "install:\nrun bundle:false\n", readTestFile(t, log))
	assert.FileExists(t, filepath.Join(root, "bin", "hello"))
	assert.FileExists(t, filepath.Join(root, "frontend", "dist", "index.html"))
}

func TestHCLPackagePlanGeneratesOwnedPlatformAssets(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := `version = 3
project {
  name = "assets"
  product_name = "Assets"
  identifier = "com.example.assets"
  version = "1.0.0"
  description = "Generated asset test"
}
protocol "assets" {
  description = "Open an assets document"
}
file_association "assets" {
  extensions = ["asset"]
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.yml"), []byte("user-owned taskfile\n"), 0o644))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "package", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"deb", "rpm", "archlinux"}})
	require.NoError(t, err)
	node := plan.Nodes[pipeline.NodeKey("target:linux/amd64:assets")]
	require.Equal(t, pipeline.GeneratePlatformAssets, node.Kind)
	_, err = (&manifestHandler{root: root, config: loaded.Config}).Run(t.Context(), node)
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "bin"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "bin", "assets"), []byte("binary"), 0o755))
	packageNode := plan.Nodes[pipeline.NodeKey("package:linux/amd64:deb")]
	_, err = (&manifestHandler{root: root, config: loaded.Config}).Run(t.Context(), packageNode)
	require.NoError(t, err)
	for _, format := range []string{"rpm", "archlinux"} {
		packageNode = plan.Nodes[pipeline.NodeKey("package:linux/amd64:"+format)]
		_, err = (&manifestHandler{root: root, config: loaded.Config}).Run(t.Context(), packageNode)
		require.NoError(t, err, format)
	}

	assetsRoot := filepath.Join(root, ".wails", "build", "default", "linux-amd64", "assets")
	assert.FileExists(t, filepath.Join(assetsRoot, "linux", "nfpm", "nfpm.yaml"))
	assert.FileExists(t, filepath.Join(root, "bin", "assets_1.0.0_amd64.deb"))
	config := readTestFile(t, filepath.Join(assetsRoot, "config.yml"))
	assert.Contains(t, config, "assets")
	assert.Contains(t, config, "Open an assets document")
	assert.Contains(t, config, "asset")
	assert.Equal(t, "user-owned taskfile\n", readTestFile(t, filepath.Join(root, "Taskfile.yml")))
}

func TestHCLDarwinAppPackageUsesGeneratedPlatformAssets(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "packaging"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "packaging", "Info.plist.tmpl"), []byte("{{.Project.Identifier}}|{{.Target.OS}}|{{.Package.Format}}"), 0o644))
	_, sourceFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	iconData, err := os.ReadFile(filepath.Join(filepath.Dir(sourceFile), "build_assets", "appicon.png"))
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "build"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "build", "icon.png"), iconData, 0o644))
	hcl := `version = 3
project {
  name = "macapp"
  product_name = "Mac App"
  identifier = "com.example.macapp"
  version = "1.0.0"
  icon = "build/icon.png"
}
package "app" {
  template = "packaging/Info.plist.tmpl"
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "package", TargetOS: "darwin", TargetArch: "amd64", Formats: []string{"app"}})
	require.NoError(t, err)
	handler := &manifestHandler{root: root, config: loaded.Config}
	assetsNode := plan.Nodes[pipeline.NodeKey("target:darwin/amd64:assets")]
	_, err = handler.Run(t.Context(), assetsNode)
	require.NoError(t, err)
	packageNode := plan.Nodes[pipeline.NodeKey("package:darwin/amd64:app")]
	spec := packageNode.Spec.(pipeline.PackageSpec)
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, spec.Binary)), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, spec.Binary), []byte("binary"), 0o755))
	_, err = handler.Run(t.Context(), packageNode)
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(root, spec.Output, "Contents", "MacOS", "macapp"))
	assert.FileExists(t, filepath.Join(root, spec.Output, "Contents", "Info.plist"))
	assert.FileExists(t, filepath.Join(root, spec.Output, "Contents", "Resources", "icons.icns"))
	assert.Equal(t, "com.example.macapp|darwin|app", readTestFile(t, filepath.Join(root, spec.Output, "Contents", "Info.plist")))
}

func TestHCLDMGPackagingBuildsTheAppAndResolvesOptions(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "packaging"), 0o755))
	for _, name := range []string{"background.png", "volume.icns", "file.icns", "README.md"} {
		require.NoError(t, os.WriteFile(filepath.Join(root, "packaging", name), []byte(name), 0o644))
	}
	require.NoError(t, os.WriteFile(filepath.Join(root, "packaging", "dmg.json.tmpl"), []byte(`{"background":"packaging/background.png","window_width":700,"window_height":500}`), 0o644))
	hcl := hclBuildFixture + `
package "dmg" {
  template = "packaging/dmg.json.tmpl"
  background = "packaging/background.png"
  volume_icon = "packaging/volume.icns"
  file_icon = "packaging/file.icns"
  files = "Read Me=packaging/README.md"
  window_width = 900
  window_height = 620
}
profile "release" {
  target "darwin/arm64" { formats = ["dmg"] }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	handler := &manifestHandler{root: root, config: loaded.Config}
	assetsNode := plan.Nodes[pipeline.NodeKey("target:darwin/arm64:assets")]
	_, err = handler.Run(t.Context(), assetsNode)
	require.NoError(t, err)
	packageNode := plan.Nodes[pipeline.NodeKey("package:darwin/arm64:dmg")]
	spec := packageNode.Spec.(pipeline.PackageSpec)
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, spec.Binary)), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, spec.Binary), []byte("darwin binary"), 0o755))

	var received *flags.ToolPackage
	previousToolPackage := manifestToolPackage
	manifestToolPackage = func(options *flags.ToolPackage) error {
		received = options
		return os.WriteFile(filepath.Join(options.Out, options.ExecutableName+".dmg"), []byte("dmg"), 0o644)
	}
	t.Cleanup(func() { manifestToolPackage = previousToolPackage })
	_, err = handler.Run(t.Context(), packageNode)
	require.NoError(t, err)
	require.NotNil(t, received)
	assert.Equal(t, "dmg", received.Format)
	assert.Equal(t, filepath.Join(root, "packaging", "background.png"), received.BackgroundImage)
	assert.Equal(t, filepath.Join(root, "packaging", "volume.icns"), received.DmgVolumeIcon)
	assert.Equal(t, filepath.Join(root, "packaging", "file.icns"), received.DmgFileIcon)
	assert.Equal(t, "Read Me="+filepath.Join(root, "packaging", "README.md"), received.DmgFiles)
	assert.Equal(t, 900, received.DmgWindowWidth)
	assert.Equal(t, 620, received.DmgWindowHeight)
	assert.FileExists(t, filepath.Join(root, spec.Output))
	assert.FileExists(t, filepath.Join(root, strings.TrimSuffix(spec.Output, ".dmg")+".app", "Contents", "MacOS", "hello"))
}

func TestHCLGeneratesAssetsForEverySupportedPlatform(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := `version = 3
project {
  name = "platforms"
  product_name = "Platforms"
  identifier = "com.example.platforms"
  version = "1.0.0"
}
profile "release" {
  target "windows/amd64" { formats = ["nsis"] }
  target "darwin/amd64" { formats = ["app"] }
  target "linux/amd64" { formats = ["deb"] }
  target "ios/arm64" { formats = ["app"] }
  target "android/arm64" { formats = ["apk", "aab"] }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	handler := &manifestHandler{root: root, config: loaded.Config}
	assets := map[string]pipeline.Node{}
	for _, node := range plan.Nodes {
		if node.Kind == pipeline.GeneratePlatformAssets {
			assets[node.Spec.(pipeline.AssetsSpec).TargetOS] = node
		}
	}
	for _, platform := range []string{"windows", "darwin", "linux", "ios", "android"} {
		node, ok := assets[platform]
		require.True(t, ok, platform)
		_, err = handler.Run(t.Context(), node)
		require.NoError(t, err, platform)
	}

	assert.FileExists(t, filepath.Join(root, ".wails", "build", "release", "windows-amd64", "assets", "windows", "wails.exe.manifest"))
	assert.FileExists(t, filepath.Join(root, ".wails", "build", "release", "ios-arm64", "assets", "ios", "xcode", "overlay.json"))
	assert.FileExists(t, filepath.Join(root, ".wails", "build", "release", "android-arm64", "assets", "android", "overlay.json"))

	fakeTools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "makensis"), []byte("#!/bin/sh\nmkdir -p ../../../bin\nprintf installer > ../../../bin/platforms-amd64-installer.exe\n"), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))
	windowsPackage := plan.Nodes[pipeline.NodeKey("package:windows/amd64:nsis")]
	windowsSpec := windowsPackage.Spec.(pipeline.PackageSpec)
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, windowsSpec.Binary)), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, windowsSpec.Binary), []byte("windows binary"), 0o755))
	_, err = handler.Run(t.Context(), windowsPackage)
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(root, windowsSpec.Output))

	androidAssets := filepath.Join(root, ".wails", "build", "release", "android-arm64", "assets", "android")
	require.NoError(t, os.WriteFile(filepath.Join(androidAssets, "gradlew"), []byte("#!/bin/sh\nif [ \"$1\" = bundleRelease ]; then mkdir -p app/build/outputs/bundle/release; printf aab > app/build/outputs/bundle/release/app-release.aab; else mkdir -p app/build/outputs/apk/release; printf apk > app/build/outputs/apk/release/app-release.apk; fi\n"), 0o755))
	for _, format := range []string{"apk", "aab"} {
		androidPackage := plan.Nodes[pipeline.NodeKey("package:android/arm64:"+format)]
		androidSpec := androidPackage.Spec.(pipeline.PackageSpec)
		require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, androidSpec.Binary)), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(root, androidSpec.Binary), []byte("shared library"), 0o755))
		_, err = handler.Run(t.Context(), androidPackage)
		require.NoError(t, err, format)
		assert.FileExists(t, filepath.Join(root, androidSpec.Output))
	}
}

func TestHCLTargetBuildNumbersAndMinimumVersionsReachGeneratedAssets(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := `version = 3
project {
  name = "versions"
  product_name = "Versions"
  identifier = "com.example.versions"
  version = "1.0.0"
  build_number = 42
}
target "windows/amd64" {
  build_number = 11
}
target "darwin/amd64" {
  build_number = 13
  minimum_version = "14.0"
}
target "ios/arm64" {
  build_number = 15
  minimum_version = "17.0"
}
target "android/arm64" {
  build_number = 17
}
profile "release" {
  target "windows/amd64" {}
  target "darwin/amd64" {}
  target "ios/arm64" {}
  target "android/arm64" {}
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	handler := &manifestHandler{root: root, config: loaded.Config}
	for _, target := range []struct{ name, format string }{
		{name: "windows/amd64", format: "nsis"},
		{name: "darwin/amd64", format: "app"},
		{name: "ios/arm64", format: "app"},
		{name: "android/arm64", format: "apk"},
	} {
		parts := strings.Split(target.name, "/")
		plan, planErr := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "package", TargetOS: parts[0], TargetArch: parts[1], Formats: []string{target.format}})
		require.NoError(t, planErr, target.name)
		node := plan.Nodes[pipeline.NodeKey("target:"+target.name+":assets")]
		_, err = handler.Run(t.Context(), node)
		require.NoError(t, err, target.name)
	}
	base := filepath.Join(root, ".wails", "build", "release")
	assert.Contains(t, readTestFile(t, filepath.Join(base, "windows-amd64", "assets", "windows", "nsis", "project.nsi")), `${INFO_PRODUCTVERSION}.11`)
	darwinInfo := readTestFile(t, filepath.Join(base, "darwin-amd64", "assets", "darwin", "Info.plist"))
	assert.Contains(t, darwinInfo, "<string>13</string>")
	assert.Contains(t, darwinInfo, "<string>14.0</string>")
	iosInfo := readTestFile(t, filepath.Join(base, "ios-arm64", "assets", "ios", "xcode", "main", "Info.plist"))
	assert.Contains(t, iosInfo, "<string>15</string>")
	assert.Contains(t, iosInfo, "<string>17.0</string>")
	assert.Contains(t, readTestFile(t, filepath.Join(base, "android-arm64", "assets", "android", "app", "build.gradle")), "versionCode 17")
}

func TestHCLLinuxSigningProducesTheSignedArtifact(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := `version = 3
project {
  name = "signed"
  product_name = "Signed"
  identifier = "com.example.signed"
  version = "1.0.0"
}
linux {
  signing {
    identity = "release"
    certificate = "release-key"
  }
}
profile "release" {
  target "linux/amd64" {
    formats = ["deb", "rpm"]
    sign = true
  }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	handler := &manifestHandler{root: root, config: loaded.Config}
	assetsNode := plan.Nodes[pipeline.NodeKey("target:linux/amd64:assets")]
	_, err = handler.Run(t.Context(), assetsNode)
	require.NoError(t, err)

	packageSpec := plan.Nodes[pipeline.NodeKey("package:linux/amd64:deb")].Spec.(pipeline.PackageSpec)
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, packageSpec.Binary)), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, packageSpec.Binary), []byte("binary"), 0o755))

	fakeTools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "dpkg-sig"), []byte("#!/bin/sh\nexit 0\n"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "rpmsign"), []byte("#!/bin/sh\nexit 0\n"), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))
	for _, format := range []string{"deb", "rpm"} {
		packageNode := plan.Nodes[pipeline.NodeKey("package:linux/amd64:"+format)]
		packageSpec = packageNode.Spec.(pipeline.PackageSpec)
		_, err = handler.Run(t.Context(), packageNode)
		require.NoError(t, err, format)
		signNode := plan.Nodes[pipeline.NodeKey("package:linux/amd64:"+format+":sign")]
		_, err = handler.Run(t.Context(), signNode)
		require.NoError(t, err, format)
		assert.FileExists(t, filepath.Join(root, packageSpec.Output+".signed"))
	}
}

func TestHCLNativeSigningSettingsReachWindowsAndDarwinArtifacts(t *testing.T) {
	tests := []struct {
		name, platform, format, hcl string
	}{
		{
			name:     "windows",
			platform: "windows",
			format:   "msix",
			hcl: `version = 3
project {
  name = "windows"
  product_name = "Windows"
  company = "Example"
  identifier = "com.example.windows"
  version = "1.0.0"
}
windows {
  signing {
    credential = "WINDOWS_PASSWORD"
    identity = "Example Publisher"
    certificate = "release.pfx"
    thumbprint = "ABC123"
    timestamp_server = "https://timestamp.example.test"
  }
}
profile "release" {
  target "windows/amd64" {
    formats = ["msix"]
    sign = true
  }
}
`,
		},
		{
			name:     "darwin",
			platform: "darwin",
			format:   "app",
			hcl: `version = 3
project {
  name = "darwin"
  product_name = "Darwin"
  identifier = "com.example.darwin"
  version = "1.0.0"
}
darwin {
  signing {
    credential = "MACOS_KEYCHAIN"
    identity = "Developer ID Application: Example"
    entitlements = "release.entitlements"
  }
  notarization {
    credential = "NOTARY_PROFILE"
  }
}
profile "release" {
  target "darwin/arm64" {
    formats = ["app"]
    sign = true
    notarize = true
  }
}
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			t.Chdir(root)
			require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(test.hcl), 0o644))
			require.NoError(t, os.WriteFile(filepath.Join(root, "release.pfx"), []byte("certificate"), 0o644))
			require.NoError(t, os.WriteFile(filepath.Join(root, "release.entitlements"), []byte("entitlements"), 0o644))
			loaded, err := manifest.Load(root, "release")
			require.NoError(t, err)
			plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
			require.NoError(t, err)
			signNode := plan.Nodes[pipeline.NodeKey("package:"+test.platform+map[string]string{"windows": "/amd64", "darwin": "/arm64"}[test.platform]+":"+test.format+":sign")]
			spec := signNode.Spec.(pipeline.SignSpec)
			input := filepath.Join(root, spec.Input)
			if test.platform == "darwin" {
				require.NoError(t, os.MkdirAll(filepath.Join(input, "Contents"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(input, "Contents", "Info.plist"), []byte("plist"), 0o644))
			} else {
				require.NoError(t, os.MkdirAll(filepath.Dir(input), 0o755))
				require.NoError(t, os.WriteFile(input, []byte("package"), 0o644))
			}

			var received flags.Sign
			previousSign := manifestSign
			manifestSign = func(options *flags.Sign) error {
				received = *options
				return nil
			}
			t.Cleanup(func() { manifestSign = previousSign })
			handler := &manifestHandler{root: root, config: loaded.Config}
			_, err = handler.Run(t.Context(), signNode)
			require.NoError(t, err)
			assert.True(t, spec.Config.Enabled)
			assert.Equal(t, input, received.Input)
			assert.Equal(t, spec.Config.Identity, received.Identity)
			assert.Equal(t, spec.Config.Certificate, received.Certificate)
			assert.Equal(t, spec.Config.Thumbprint, received.Thumbprint)
			assert.Equal(t, spec.Config.TimestampServer, received.Timestamp)
			assert.Equal(t, spec.Config.Entitlements, received.Entitlements)
			if test.platform == "darwin" {
				assert.True(t, received.Notarize)
				assert.Equal(t, spec.Config.Credential, received.KeychainProfile)
				assert.DirExists(t, input+".signed")
			} else {
				assert.FileExists(t, input+".signed")
			}
		})
	}
}

func TestHCLAppImagePackagingUsesTheConfiguredOutput(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "packaging"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "packaging", "app.desktop.tmpl"), []byte("Name={{.Project.ProductName}}\nExec={{.Project.BinaryName}}\n"), 0o644))
	hcl := `version = 3
project {
  name = "image"
  product_name = "Image"
  identifier = "com.example.image"
  version = "1.0.0"
}
package "appimage" {
  template = "packaging/app.desktop.tmpl"
}
profile "release" {
  target "linux/amd64" { formats = ["appimage"] }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	handler := &manifestHandler{root: root, config: loaded.Config}
	assetsNode := plan.Nodes[pipeline.NodeKey("target:linux/amd64:assets")]
	_, err = handler.Run(t.Context(), assetsNode)
	require.NoError(t, err)
	packageNode := plan.Nodes[pipeline.NodeKey("package:linux/amd64:appimage")]
	spec := packageNode.Spec.(pipeline.PackageSpec)
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, spec.Binary)), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, spec.Binary), []byte("binary"), 0o755))

	fakeCLI := filepath.Join(t.TempDir(), "wails-fake")
	require.NoError(t, os.WriteFile(fakeCLI, []byte("#!/bin/sh\noutput=\nwhile [ \"$#\" -gt 0 ]; do if [ \"$1\" = \"-outputdir\" ]; then output=$2; shift; fi; shift; done\nmkdir -p \"$output\"\nprintf appimage > \"$output/image-amd64.AppImage\"\n"), 0o755))
	previous := manifestExecutable
	manifestExecutable = func() (string, error) { return fakeCLI, nil }
	t.Cleanup(func() { manifestExecutable = previous })
	_, err = handler.Run(t.Context(), packageNode)
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(root, spec.Output))
	desktop := filepath.Join(root, ".wails", "build", "release", "linux-amd64", "package", "appimage", "image.desktop")
	assert.Equal(t, "Name=Image\nExec=image\n", readTestFile(t, desktop))
}

func TestHCLDefaultAppImagePackagingGeneratesDesktopMetadata(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := `version = 3
project {
  name = "default-image"
  product_name = "Default Image"
  identifier = "com.example.defaultimage"
  version = "1.0.0"
  description = "Default image package"
}
package "appimage" {
  categories = "Development;IDE;"
}
profile "release" {
  target "linux/amd64" { formats = ["appimage"] }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	handler := &manifestHandler{root: root, config: loaded.Config}
	assetsNode := plan.Nodes[pipeline.NodeKey("target:linux/amd64:assets")]
	_, err = handler.Run(t.Context(), assetsNode)
	require.NoError(t, err)
	packageNode := plan.Nodes[pipeline.NodeKey("package:linux/amd64:appimage")]
	spec := packageNode.Spec.(pipeline.PackageSpec)
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, spec.Binary)), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, spec.Binary), []byte("binary"), 0o755))
	fakeCLI := filepath.Join(t.TempDir(), "wails-fake")
	require.NoError(t, os.WriteFile(fakeCLI, []byte(`#!/bin/sh
output=
while [ "$#" -gt 0 ]; do
  if [ "$1" = "-outputdir" ]; then output=$2; shift; fi
  shift
done
mkdir -p "$output"
printf appimage > "$output/default-image-amd64.AppImage"
`), 0o755))
	previous := manifestExecutable
	manifestExecutable = func() (string, error) { return fakeCLI, nil }
	t.Cleanup(func() { manifestExecutable = previous })
	_, err = handler.Run(t.Context(), packageNode)
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(root, spec.Output))
	desktop := filepath.Join(root, ".wails", "build", "release", "linux-amd64", "package", "appimage", "default-image.desktop")
	assert.Contains(t, readTestFile(t, desktop), "Name=Default Image")
	assert.Contains(t, readTestFile(t, desktop), "Comment=Default image package")
	assert.Contains(t, readTestFile(t, desktop), "Categories=Development;IDE;")
}

func TestHCLPlanIdentityCoversEveryPortableStage(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := `version = 3
project {
  name = "identity"
  product_name = "Identity"
  identifier = "com.example.identity"
  version = "1.0.0"
}
profile "release" {
  target "windows/amd64" { formats = ["nsis", "msix"] }
  target "darwin/amd64" { formats = ["app", "dmg"] }
  target "linux/amd64" { formats = ["appimage", "deb"] }
  target "android/arm64" { formats = ["apk", "aab"] }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	fakeTools := t.TempDir()
	for _, name := range []string{"npm", "node", "makensis", "MakeAppx.exe"} {
		require.NoError(t, os.WriteFile(filepath.Join(fakeTools, name), []byte("#!/bin/sh\nexit 0\n"), 0o755))
	}
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))
	handler := &manifestHandler{root: root, config: loaded.Config}
	for key, node := range plan.Nodes {
		identity, err := handler.Identity(t.Context(), node)
		require.NoError(t, err, key)
		assert.NotEmpty(t, identity, key)
	}
}

func TestHCLIOSSignedAppAndIPAUseDeclaredTargetSettings(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "packaging"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "packaging", "app-Info.plist.tmpl"), []byte("app {{.Project.Identifier}}"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "packaging", "ipa-Info.plist.tmpl"), []byte("ipa {{.Project.Identifier}}"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "ios.entitlements"), []byte("entitlements"), 0o644))
	hcl := `version = 3
project {
  name = "mobile"
  product_name = "Mobile"
  identifier = "com.example.mobile"
  version = "1.0.0"
}
package "app" {
  template = "packaging/app-Info.plist.tmpl"
}
package "ipa" {
  template = "packaging/ipa-Info.plist.tmpl"
}
ios {
  signing {
    identity = "Apple Distribution: Example"
    entitlements = "ios.entitlements"
  }
}
target "ios/arm64" {
  variant = "device"
  minimum_version = "17.0"
}
profile "release" {
  target "ios/arm64" {
    formats = ["app", "ipa"]
    sign = true
  }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	previousHost := manifestHostOS
	manifestHostOS = "darwin"
	t.Cleanup(func() { manifestHostOS = previousHost })
	fakeTools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "xcrun"), []byte("#!/bin/sh\nfor arg in \"$@\"; do if [ \"$arg\" = \"--show-sdk-path\" ]; then printf /fake/sdk; exit 0; fi; if [ \"$arg\" = \"--find\" ]; then printf /fake/clang; exit 0; fi; if [ \"$arg\" = \"actool\" ]; then compile=; plist=; while [ \"$#\" -gt 0 ]; do if [ \"$1\" = \"--compile\" ]; then compile=$2; shift; fi; if [ \"$1\" = \"--output-partial-info-plist\" ]; then plist=$2; shift; fi; shift; done; mkdir -p \"$compile\" \"$(dirname \"$plist\")\"; printf car > \"$compile/Assets.car\"; printf plist > \"$plist\"; exit 0; fi; done\noutput=\nwhile [ \"$#\" -gt 0 ]; do if [ \"$1\" = \"-o\" ]; then output=$2; shift; fi; shift; done\nmkdir -p \"$(dirname \"$output\")\"\nprintf binary > \"$output\"\n"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "go"), []byte("#!/bin/sh\noutput=\nwhile [ \"$#\" -gt 0 ]; do if [ \"$1\" = \"-o\" ]; then output=$2; shift; fi; shift; done\nmkdir -p \"$(dirname \"$output\")\"\nprintf archive > \"$output\"\n"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "codesign"), []byte("#!/bin/sh\nexit 0\n"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "zip"), []byte("#!/bin/sh\noutput=$2\nprintf ipa > \"$output\"\n"), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))
	handler := &manifestHandler{root: root, config: loaded.Config}
	assetsNode := plan.Nodes[pipeline.NodeKey("target:ios/arm64:assets")]
	_, err = handler.Run(t.Context(), assetsNode)
	require.NoError(t, err)
	compileNode := plan.Nodes[pipeline.NodeKey("target:ios/arm64:compile")]
	_, err = handler.Run(t.Context(), compileNode)
	require.NoError(t, err)

	for _, format := range []string{"app", "ipa"} {
		packageKey := pipeline.NodeKey("package:ios/arm64:" + format)
		packageNode := plan.Nodes[packageKey]
		spec := packageNode.Spec.(pipeline.PackageSpec)
		packageResult, packageErr := handler.Run(t.Context(), packageNode)
		err = packageErr
		if err != nil {
			t.Logf("%s package failed: %v (%s)", format, err, packageResult.Detail)
		}
		require.NoError(t, err, format)
		if format == "app" {
			assert.DirExists(t, filepath.Join(root, spec.Output))
			assert.Equal(t, "app com.example.mobile", readTestFile(t, filepath.Join(root, spec.Output, "Info.plist")))
		} else {
			assert.FileExists(t, filepath.Join(root, spec.Output))
			workspace := filepath.Join(root, ".wails", "build", "release", "ios-arm64", "package", "ipa")
			assert.Equal(t, "ipa com.example.mobile", readTestFile(t, filepath.Join(workspace, "mobile.app", "Info.plist")))
		}
		signNode := plan.Nodes[pipeline.NodeKey(string(packageKey)+":sign")]
		_, err = handler.Run(t.Context(), signNode)
		require.NoError(t, err, format)
		if format == "app" {
			assert.DirExists(t, filepath.Join(root, spec.Output+".signed"))
		} else {
			assert.FileExists(t, filepath.Join(root, spec.Output+".signed"))
		}
	}
}

func TestHCLIOSSignedDefaultSimulatorAppUsesGeneratedInfoPlist(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := `version = 3
project {
  name = "simulator-app"
  product_name = "Simulator App"
  identifier = "com.example.simulatorapp"
  version = "1.0.0"
}
profile "release" {
  target "ios/arm64" { formats = ["app"] }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	previousHost := manifestHostOS
	manifestHostOS = "darwin"
	t.Cleanup(func() { manifestHostOS = previousHost })
	fakeTools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "xcrun"), []byte(`#!/bin/sh
for arg in "$@"; do
  if [ "$arg" = "--show-sdk-path" ]; then printf /fake/iphonesimulator.sdk; exit 0; fi
  if [ "$arg" = "actool" ]; then
    compile=
    plist=
    while [ "$#" -gt 0 ]; do
      if [ "$1" = "--compile" ]; then compile=$2; shift; fi
      if [ "$1" = "--output-partial-info-plist" ]; then plist=$2; shift; fi
      shift
    done
    mkdir -p "$compile" "$(dirname "$plist")"
    printf car > "$compile/Assets.car"
    printf plist > "$plist"
    exit 0
  fi
done
output=
while [ "$#" -gt 0 ]; do
  if [ "$1" = "-o" ]; then output=$2; shift; fi
  shift
done
mkdir -p "$(dirname "$output")"
printf executable > "$output"
`), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "codesign"), []byte("#!/bin/sh\nexit 0\n"), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))
	handler := &manifestHandler{root: root, config: loaded.Config}
	assetsNode := plan.Nodes[pipeline.NodeKey("target:ios/arm64:assets")]
	_, err = handler.Run(t.Context(), assetsNode)
	require.NoError(t, err)
	packageNode := plan.Nodes[pipeline.NodeKey("package:ios/arm64:app")]
	spec := packageNode.Spec.(pipeline.PackageSpec)
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, spec.Binary)), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, spec.Binary), []byte("archive"), 0o755))
	_, err = handler.Run(t.Context(), packageNode)
	require.NoError(t, err)
	assert.DirExists(t, filepath.Join(root, spec.Output))
	assert.Contains(t, readTestFile(t, filepath.Join(root, spec.Output, "Info.plist")), "com.example.simulatorapp")
}

func TestHCLDarwinUniversalCompileUsesBothArchitectures(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := hclBuildFixture + `
profile "release" {
  target "darwin/universal" {}
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)

	fakeTools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "go"), []byte(`#!/bin/sh
output=
while [ "$#" -gt 0 ]; do
  if [ "$1" = "-o" ]; then output=$2; shift; fi
  shift
done
mkdir -p "$(dirname "$output")"
printf darwin > "$output"
`), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))

	var lipoInputs []string
	previousLipo := manifestLipo
	manifestLipo = func(options *flags.Lipo) error {
		lipoInputs = append([]string(nil), options.Inputs...)
		return os.WriteFile(options.Output, []byte("universal"), 0o644)
	}
	t.Cleanup(func() { manifestLipo = previousLipo })

	handler := &manifestHandler{root: root, config: loaded.Config}
	compileNode := plan.Nodes[pipeline.NodeKey("target:darwin/universal:compile")]
	_, err = handler.Run(t.Context(), compileNode)
	require.NoError(t, err)
	assert.Len(t, lipoInputs, 2)
	assert.Contains(t, lipoInputs[0], "darwin-universal/binaries/hello-")
	assert.Contains(t, lipoInputs[1], "darwin-universal/binaries/hello-")
	assert.FileExists(t, filepath.Join(root, compileNode.Output))
}

func TestHCLCompileUsesConfiguredGoFlagsAndTargetMinimum(t *testing.T) {
	tests := []struct {
		name, target, hcl, wantEnv string
	}{
		{
			name:   "linux flags",
			target: "linux/amd64",
			hcl: `version = 3
project {
  name = "flags"
  product_name = "Flags"
  identifier = "com.example.flags"
  version = "1.0.0"
}
build {
  tags = ["enterprise"]
  compiler_flags = ["all=-l"]
  ldflags = ["-X example/build.version=2.0.0"]
}
profile "release" {
  target "linux/amd64" {}
}
`,
			wantEnv: "linux|amd64|||",
		},
		{
			name:   "darwin minimum",
			target: "darwin/amd64",
			hcl: `version = 3
project {
  name = "minimum"
  product_name = "Minimum"
  identifier = "com.example.minimum"
  version = "1.0.0"
}
target "darwin/amd64" {
  minimum_version = "13.0"
}
profile "release" {
  target "darwin/amd64" {}
}
`,
			wantEnv: "darwin|amd64|13.0|-mmacosx-version-min=13.0|-mmacosx-version-min=13.0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			t.Chdir(root)
			require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(test.hcl), 0o644))
			loaded, err := manifest.Load(root, "release")
			require.NoError(t, err)
			plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
			require.NoError(t, err)
			node := plan.Nodes[pipeline.NodeKey("target:"+test.target+":compile")]

			fakeTools := t.TempDir()
			log := filepath.Join(root, "go.log")
			require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "go"), []byte(`#!/bin/sh
printf '%s|%s|%s|%s|%s\n' "$GOOS" "$GOARCH" "$MACOSX_DEPLOYMENT_TARGET" "$CGO_CFLAGS" "$CGO_LDFLAGS" >> "$GO_LOG"
printf '%s\n' "$*" >> "$GO_LOG"
output=
while [ "$#" -gt 0 ]; do
  if [ "$1" = "-o" ]; then output=$2; shift; fi
  shift
done
mkdir -p "$(dirname "$output")"
printf binary > "$output"
`), 0o755))
			t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))
			t.Setenv("GO_LOG", log)

			_, err = (&manifestHandler{root: root, config: loaded.Config}).Run(t.Context(), node)
			require.NoError(t, err)
			data := readTestFile(t, log)
			assert.Contains(t, data, test.wantEnv)
			if test.name == "linux flags" {
				assert.Contains(t, data, "-tags enterprise,production")
				assert.Contains(t, data, "-gcflags all=-l")
				assert.Contains(t, data, "-ldflags -X example/build.version=2.0.0 -w -s")
			}
		})
	}
}

func TestHCLAndroidCompileDiscoversNDKForEachHostToolchain(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := `version = 3
project {
  name = "android-host"
  product_name = "Android Host"
  identifier = "com.example.androidhost"
  version = "1.0.0"
}
profile "release" {
  target "android/arm64" {}
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)

	sdk := t.TempDir()
	ndk := filepath.Join(sdk, "ndk", "26.3.11579264")
	fakeTools := t.TempDir()
	compileLog := filepath.Join(root, "android-host.log")
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "go"), []byte(`#!/bin/sh
printf '%s|%s|%s\n' "$GOARCH" "$CC" "$WAILS_ANDROID_JNI" >> "$ANDROID_HOST_LOG"
output=
while [ "$#" -gt 0 ]; do
  if [ "$1" = "-o" ]; then output=$2; shift; fi
  shift
done
mkdir -p "$(dirname "$output")"
printf shared-library > "$output"
`), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("ANDROID_HOME", sdk)
	t.Setenv("ANDROID_NDK_HOME", "")
	t.Setenv("ANDROID_HOST_LOG", compileLog)

	previousHost := manifestHostOS
	t.Cleanup(func() { manifestHostOS = previousHost })
	handler := &manifestHandler{root: root, config: loaded.Config}
	for _, test := range []struct {
		host, tag, suffix string
	}{
		{host: "linux", tag: "linux-x86_64"},
		{host: "darwin", tag: "darwin-x86_64"},
		{host: "windows", tag: "windows-x86_64", suffix: ".cmd"},
	} {
		toolchain := filepath.Join(ndk, "toolchains", "llvm", "prebuilt", test.tag, "bin")
		require.NoError(t, os.MkdirAll(toolchain, 0o755))
		compiler := filepath.Join(toolchain, "aarch64-linux-android21-clang"+test.suffix)
		require.NoError(t, os.WriteFile(compiler, []byte("#!/bin/sh\nexit 0\n"), 0o755))
		manifestHostOS = test.host
		node := plan.Nodes[pipeline.NodeKey("target:android/arm64:compile")]
		_, err = handler.Run(t.Context(), node)
		require.NoError(t, err, test.host)
	}
	log := readTestFile(t, compileLog)
	assert.Contains(t, log, "arm64|"+filepath.Join(ndk, "toolchains", "llvm", "prebuilt", "linux-x86_64", "bin", "aarch64-linux-android21-clang")+"|arm64-v8a")
	assert.Contains(t, log, "arm64|"+filepath.Join(ndk, "toolchains", "llvm", "prebuilt", "darwin-x86_64", "bin", "aarch64-linux-android21-clang")+"|arm64-v8a")
	assert.Contains(t, log, "arm64|"+filepath.Join(ndk, "toolchains", "llvm", "prebuilt", "windows-x86_64", "bin", "aarch64-linux-android21-clang.cmd")+"|arm64-v8a")
}

func TestHCLIOSSImulatorCompileUsesTheDefaultVariantAndMinimum(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := `version = 3
project {
  name = "simulator"
  product_name = "Simulator"
  identifier = "com.example.simulator"
  version = "1.0.0"
}
profile "release" {
  target "ios/arm64" {}
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)

	fakeTools := t.TempDir()
	compileLog := filepath.Join(root, "ios-simulator.log")
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "xcrun"), []byte(`#!/bin/sh
for arg in "$@"; do
  if [ "$arg" = "--show-sdk-path" ]; then printf /fake/iphonesimulator.sdk; exit 0; fi
  if [ "$arg" = "--find" ]; then printf /fake/clang; exit 0; fi
done
output=
while [ "$#" -gt 0 ]; do
  if [ "$1" = "-o" ]; then output=$2; shift; fi
  shift
done
mkdir -p "$(dirname "$output")"
printf executable > "$output"
`), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "go"), []byte(`#!/bin/sh
printf '%s|%s|%s|%s\n' "$GOOS" "$GOARCH" "$CGO_CFLAGS" "$CGO_LDFLAGS" >> "$IOS_SIMULATOR_LOG"
output=
while [ "$#" -gt 0 ]; do
  if [ "$1" = "-o" ]; then output=$2; shift; fi
  shift
done
mkdir -p "$(dirname "$output")"
printf archive > "$output"
`), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("IOS_SIMULATOR_LOG", compileLog)
	previousHost := manifestHostOS
	manifestHostOS = "darwin"
	t.Cleanup(func() { manifestHostOS = previousHost })

	node := plan.Nodes[pipeline.NodeKey("target:ios/arm64:compile")]
	_, err = (&manifestHandler{root: root, config: loaded.Config}).Run(t.Context(), node)
	require.NoError(t, err)
	assert.Contains(t, readTestFile(t, compileLog), "ios|arm64|-isysroot /fake/iphonesimulator.sdk -target arm64-apple-ios15.0-simulator -mios-simulator-version-min=15.0|-isysroot /fake/iphonesimulator.sdk -target arm64-apple-ios15.0-simulator")
}

func TestHCLProjectIconIsCopiedIntoGeneratedPlatformAssets(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	_, sourceFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	source := filepath.Join(filepath.Dir(sourceFile), "build_assets", "appicon.png")
	require.NoError(t, os.MkdirAll(filepath.Join(root, "build"), 0o755))
	iconData, err := os.ReadFile(source)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(root, "build", "icon.png"), iconData, 0o644))
	hcl := `version = 3
project {
  name = "icon"
  product_name = "Icon"
  identifier = "com.example.icon"
  version = "1.0.0"
  icon = "build/icon.png"
}
profile "release" {
  target "windows/amd64" {}
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	node := plan.Nodes[pipeline.NodeKey("target:windows/amd64:assets")]
	_, err = (&manifestHandler{root: root, config: loaded.Config}).Run(t.Context(), node)
	require.NoError(t, err)
	assets := filepath.Join(root, ".wails", "build", "release", "windows-amd64", "assets")
	assert.FileExists(t, filepath.Join(assets, "appicon.png"))
	assert.FileExists(t, filepath.Join(assets, "windows", "icon.ico"))
}

func TestHCLAndroidCompileUsesConfiguredNDKForBothArchitectures(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := hclBuildFixture + `
profile "release" {
  target "android/arm64" {}
  target "android/amd64" {}
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)

	ndk := t.TempDir()
	toolchain := filepath.Join(ndk, "toolchains", "llvm", "prebuilt", "linux-x86_64", "bin")
	require.NoError(t, os.MkdirAll(toolchain, 0o755))
	for _, compiler := range []string{"aarch64-linux-android21-clang", "x86_64-linux-android21-clang"} {
		require.NoError(t, os.WriteFile(filepath.Join(toolchain, compiler), []byte("#!/bin/sh\nexit 0\n"), 0o755))
	}
	t.Setenv("ANDROID_NDK_HOME", ndk)

	fakeTools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "go"), []byte(`#!/bin/sh
printf '%s:%s:%s:%s\n' "$GOARCH" "$CGO_ENABLED" "$CC" "$WAILS_ANDROID_JNI" >> "$ANDROID_COMPILE_LOG"
output=
while [ "$#" -gt 0 ]; do
  if [ "$1" = "-o" ]; then output=$2; shift; fi
  shift
done
mkdir -p "$(dirname "$output")"
printf shared-library > "$output"
`), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))
	compileLog := filepath.Join(root, "android-compile.log")
	t.Setenv("ANDROID_COMPILE_LOG", compileLog)

	handler := &manifestHandler{root: root, config: loaded.Config}
	for _, target := range []string{"android/arm64", "android/amd64"} {
		node := plan.Nodes[pipeline.NodeKey("target:"+target+":compile")]
		_, err = handler.Run(t.Context(), node)
		require.NoError(t, err, target)
		assert.FileExists(t, filepath.Join(root, node.Output))
	}
	log := readTestFile(t, compileLog)
	assert.Contains(t, log, "arm64:1:")
	assert.Contains(t, log, "amd64:1:")
	assert.Contains(t, log, ":arm64-v8a")
	assert.Contains(t, log, ":x86_64")
}

func TestHCLWindowsCompileGeneratesTheResourceOverlay(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := hclBuildFixture + `
profile "release" {
  target "windows/amd64" {}
  target "windows/arm64" {}
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)

	fakeTools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "go"), []byte(`#!/bin/sh
output=
while [ "$#" -gt 0 ]; do
  if [ "$1" = "-o" ]; then output=$2; shift; fi
  shift
done
mkdir -p "$(dirname "$output")"
printf windows > "$output"
	`), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))
	handler := &manifestHandler{root: root, config: loaded.Config}
	for _, arch := range []string{"amd64", "arm64"} {
		assetsNode := plan.Nodes[pipeline.NodeKey("target:windows/"+arch+":assets")]
		_, err = handler.Run(t.Context(), assetsNode)
		require.NoError(t, err, arch)
		compileNode := plan.Nodes[pipeline.NodeKey("target:windows/"+arch+":compile")]
		_, err = handler.Run(t.Context(), compileNode)
		require.NoError(t, err, arch)
		assert.FileExists(t, filepath.Join(root, compileNode.Output))
		assert.FileExists(t, filepath.Join(root, ".wails", "generated", "windows", arch, "wails_windows_"+arch+".syso"))
		assert.FileExists(t, filepath.Join(root, ".wails", "generated", "windows", arch, "overlay.json"))
	}
}

func TestHCLObfuscatedBuildUsesGarbleAndConfiguredArguments(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := `version = 3
project {
  name = "obfuscated"
  product_name = "Obfuscated"
  identifier = "com.example.obfuscated"
  version = "1.0.0"
}
build {
  obfuscated = true
  garble_args = ["-tiny", "-literals"]
  tags = ["enterprise"]
}
profile "release" {
  target "linux/amd64" {}
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	compileNode := plan.Nodes[pipeline.NodeKey("target:linux/amd64:compile")]
	compile := compileNode.Spec.(pipeline.CompileSpec)
	assert.True(t, compile.Obfuscated)
	assert.Contains(t, compile.Tags, "wails_obfuscated")
	assert.Equal(t, []string{"-tiny", "-literals"}, compile.GarbleArgs)

	fakeTools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "garble"), []byte(`#!/bin/sh
printf '%s\n' "$*" >> "$GARBLE_LOG"
output=
while [ "$#" -gt 0 ]; do
  if [ "$1" = "-o" ]; then output=$2; shift; fi
  shift
done
mkdir -p "$(dirname "$output")"
printf obfuscated > "$output"
`), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))
	garbleLog := filepath.Join(root, "garble.log")
	t.Setenv("GARBLE_LOG", garbleLog)
	handler := &manifestHandler{root: root, config: loaded.Config}
	_, err = handler.Run(t.Context(), compileNode)
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(root, compileNode.Output))
	assert.Contains(t, readTestFile(t, garbleLog), "-tiny -literals build -tags enterprise,production,wails_obfuscated")
}

func TestHCLLinuxCompileSupportsArmAnd386Targets(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := hclBuildFixture + `
profile "release" {
  target "linux/arm" {}
  target "linux/386" {}
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	fakeTools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "go"), []byte(`#!/bin/sh
output=
while [ "$#" -gt 0 ]; do
  if [ "$1" = "-o" ]; then output=$2; shift; fi
  shift
done
mkdir -p "$(dirname "$output")"
printf linux > "$output"
`), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))
	handler := &manifestHandler{root: root, config: loaded.Config}
	for _, arch := range []string{"arm", "386"} {
		node := plan.Nodes[pipeline.NodeKey("target:linux/"+arch+":compile")]
		_, err = handler.Run(t.Context(), node)
		require.NoError(t, err, arch)
		assert.FileExists(t, filepath.Join(root, node.Output))
	}
}

func TestHCLMSIXPackagingUsesMakeAppxOnTheWindowsTarget(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "packaging"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "packaging", "AppxManifest.xml.tmpl"), []byte("<Package Identity=\"{{.Project.Identifier}}\" />"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "release.pfx"), []byte("certificate"), 0o644))
	hcl := hclBuildFixture + `
windows {
  signing {
    identity = "CN=Example Publisher"
    certificate = "release.pfx"
  }
}
package "msix" {
  template = "packaging/AppxManifest.xml.tmpl"
}
file_association "document" {
  extensions = ["docx"]
  description = "Example document"
  platforms = ["windows"]
}
protocol "example" {
  description = "Example protocol"
  platforms = ["windows"]
}
profile "release" {
  target "windows/amd64" { formats = ["msix"] }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)

	fakeTools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "where"), []byte("#!/bin/sh\nexit 0\n"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "signtool.exe"), []byte("#!/bin/sh\nexit 0\n"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "MakeAppx.exe"), []byte(`#!/bin/sh
output=
while [ "$#" -gt 0 ]; do
  if [ "$1" = "/p" ]; then output=$2; shift; fi
  shift
done
mkdir -p "$(dirname "$output")"
printf msix > "$output"
`), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))
	previousHost := manifestHostOS
	manifestHostOS = "windows"
	t.Cleanup(func() { manifestHostOS = previousHost })

	handler := &manifestHandler{root: root, config: loaded.Config}
	assetsNode := plan.Nodes[pipeline.NodeKey("target:windows/amd64:assets")]
	_, err = handler.Run(t.Context(), assetsNode)
	require.NoError(t, err)
	packageNode := plan.Nodes[pipeline.NodeKey("package:windows/amd64:msix")]
	spec := packageNode.Spec.(pipeline.PackageSpec)
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, spec.Binary)), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, spec.Binary), []byte("windows binary"), 0o755))
	_, err = handler.Run(t.Context(), packageNode)
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(root, spec.Output))
	workspace := filepath.Join(root, ".wails", "build", "release", "windows-amd64", "package", "msix")
	assert.Contains(t, readTestFile(t, filepath.Join(workspace, "config.json")), "docx")
	assert.Contains(t, readTestFile(t, filepath.Join(workspace, "config.json")), "example")
	assert.Equal(t, "<Package Identity=\"com.example.hello\" />", readTestFile(t, filepath.Join(workspace, "AppxManifest.xml")))
}

func TestHCLNSISPackagingUsesTheARM64InstallerDefinition(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "packaging"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "packaging", "project.nsi.tmpl"), []byte("# custom {{.Target.Arch}}"), 0o644))
	hcl := hclBuildFixture + `
package "nsis" {
  template = "packaging/project.nsi.tmpl"
}
profile "release" {
  target "windows/arm64" { formats = ["nsis"] }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	fakeTools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "makensis"), []byte(`#!/bin/sh
mkdir -p ../../../bin
printf installer > ../../../bin/hello-arm64-installer.exe
`), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))
	handler := &manifestHandler{root: root, config: loaded.Config}
	assetsNode := plan.Nodes[pipeline.NodeKey("target:windows/arm64:assets")]
	_, err = handler.Run(t.Context(), assetsNode)
	require.NoError(t, err)
	packageNode := plan.Nodes[pipeline.NodeKey("package:windows/arm64:nsis")]
	spec := packageNode.Spec.(pipeline.PackageSpec)
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, spec.Binary)), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, spec.Binary), []byte("windows binary"), 0o755))
	_, err = handler.Run(t.Context(), packageNode)
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(root, spec.Output))
	assert.Equal(t, "# custom arm64", readTestFile(t, filepath.Join(root, ".wails", "build", "release", "windows-arm64", "package", "nsis", "assets", "windows", "nsis", "project.nsi")))
}

func TestHCLAndroidSigningProducesAPKAndAABArtifacts(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := hclBuildFixture + `
android {
  signing {
    certificate = "debug.keystore"
    identity = "release"
    credential = "ANDROID_PASSWORD"
  }
}
profile "release" {
  target "android/arm64" {
    formats = ["apk", "aab"]
    sign = true
  }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "debug.keystore"), []byte("keystore"), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	t.Setenv("ANDROID_PASSWORD", "secret")

	fakeTools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "apksigner"), []byte(`#!/bin/sh
output=
while [ "$#" -gt 0 ]; do
  if [ "$1" = "--out" ]; then output=$2; shift; fi
  shift
done
mkdir -p "$(dirname "$output")"
printf signed-apk > "$output"
`), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(fakeTools, "jarsigner"), []byte(`#!/bin/sh
output=
while [ "$#" -gt 0 ]; do
  if [ "$1" = "-signedjar" ]; then output=$2; shift; fi
  shift
done
mkdir -p "$(dirname "$output")"
printf signed-aab > "$output"
`), 0o755))
	t.Setenv("PATH", fakeTools+string(os.PathListSeparator)+os.Getenv("PATH"))

	handler := &manifestHandler{root: root, config: loaded.Config}
	assetsNode := plan.Nodes[pipeline.NodeKey("target:android/arm64:assets")]
	_, err = handler.Run(t.Context(), assetsNode)
	require.NoError(t, err)
	assetsRoot := filepath.Join(root, ".wails", "build", "release", "android-arm64", "assets", "android")
	require.NoError(t, os.WriteFile(filepath.Join(assetsRoot, "gradlew"), []byte(`#!/bin/sh
if [ "$1" = "bundleRelease" ]; then
  mkdir -p app/build/outputs/bundle/release
  printf aab > app/build/outputs/bundle/release/app-release.aab
else
  mkdir -p app/build/outputs/apk/release
  printf apk > app/build/outputs/apk/release/app-release.apk
fi
`), 0o755))

	for _, format := range []string{"apk", "aab"} {
		packageNode := plan.Nodes[pipeline.NodeKey("package:android/arm64:"+format)]
		spec := packageNode.Spec.(pipeline.PackageSpec)
		require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, spec.Binary)), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(root, spec.Binary), []byte("shared library"), 0o755))
		_, err = handler.Run(t.Context(), packageNode)
		require.NoError(t, err, format)
		signNode := plan.Nodes[pipeline.NodeKey("package:android/arm64:"+format+":sign")]
		_, err = handler.Run(t.Context(), signNode)
		require.NoError(t, err, format)
		assert.FileExists(t, filepath.Join(root, spec.Output+".signed"))
	}
}

func TestHCLAndroidAMD64PackagingUsesTheX86ABI(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	hcl := hclBuildFixture + `
profile "release" {
  target "android/amd64" { formats = ["apk"] }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	handler := &manifestHandler{root: root, config: loaded.Config}
	assetsNode := plan.Nodes[pipeline.NodeKey("target:android/amd64:assets")]
	_, err = handler.Run(t.Context(), assetsNode)
	require.NoError(t, err)
	assetsRoot := filepath.Join(root, ".wails", "build", "release", "android-amd64", "assets", "android")
	require.NoError(t, os.WriteFile(filepath.Join(assetsRoot, "gradlew"), []byte(`#!/bin/sh
mkdir -p app/build/outputs/apk/release
printf apk > app/build/outputs/apk/release/app-release.apk
`), 0o755))
	packageNode := plan.Nodes[pipeline.NodeKey("package:android/amd64:apk")]
	spec := packageNode.Spec.(pipeline.PackageSpec)
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, spec.Binary)), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, spec.Binary), []byte("shared library"), 0o755))
	_, err = handler.Run(t.Context(), packageNode)
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(root, spec.Output))
	workspace := filepath.Join(root, ".wails", "build", "release", "android-amd64", "package", "apk")
	assert.FileExists(t, filepath.Join(workspace, "app", "src", "main", "jniLibs", "x86_64", "libwails.so"))
}

func TestHCLAndroidPackagingRendersAUserTemplateDirectory(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	templateRoot := filepath.Join(root, "packaging", "android")
	require.NoError(t, os.MkdirAll(filepath.Join(templateRoot, "app", "src", "main"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(templateRoot, "gradlew"), []byte(`#!/bin/sh
mkdir -p app/build/outputs/apk/release
printf template-apk > app/build/outputs/apk/release/app-release.apk
`), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(templateRoot, "app", "src", "main", "AndroidManifest.xml.tmpl"), []byte("<manifest package=\"{{.Project.Identifier}}\" />"), 0o644))
	hcl := hclBuildFixture + `
package "apk" {
  template = "packaging/android"
}
profile "release" {
  target "android/arm64" { formats = ["apk"] }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)
	handler := &manifestHandler{root: root, config: loaded.Config}
	packageNode := plan.Nodes[pipeline.NodeKey("package:android/arm64:apk")]
	spec := packageNode.Spec.(pipeline.PackageSpec)
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, spec.Binary)), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, spec.Binary), []byte("shared library"), 0o755))
	_, err = handler.Run(t.Context(), packageNode)
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(root, spec.Output))
	workspace := filepath.Join(root, ".wails", "build", "release", "android-arm64", "package", "apk")
	assert.Equal(t, "<manifest package=\"com.example.hello\" />", readTestFile(t, filepath.Join(workspace, "app", "src", "main", "AndroidManifest.xml")))
	assert.FileExists(t, filepath.Join(workspace, "app", "src", "main", "jniLibs", "arm64-v8a", "libwails.so"))
}

func captureHCLStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = writer
	done := make(chan string, 1)
	go func() {
		var buffer bytes.Buffer
		_, _ = io.Copy(&buffer, reader)
		done <- buffer.String()
	}()
	defer func() {
		os.Stdout = old
		_ = reader.Close()
	}()
	fn()
	require.NoError(t, writer.Close())
	return <-done
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(data)
}

func mapStrings[T any](items []T, selectValue func(T) string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, selectValue(item))
	}
	return result
}

func mapArtifacts(items []struct {
	Kind string `json:"kind"`
	Path string `json:"path"`
}) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, item.Kind)
	}
	return result
}
