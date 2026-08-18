package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestSelectedProfilePlansItsDeclaredArtifacts(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(`version = 3
project {
  name = "app"
  product_name = "App"
  identifier = "com.example.app"
  version = "1.0.0"
}
frontend {
  directory = "frontend"
  install = ["npm", "install"]
  build = ["npm", "run", "build"]
  dev = ["npm", "run", "dev"]
  output = "frontend/dist"
}
profile "release" {
  target "windows/amd64" {
    formats = ["nsis"]
    sign = true
  }
  target "linux/amd64" {
    formats = ["deb"]
  }
}`), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := PlanBuild(loaded.Config, Request{Verb: "build"})
	require.NoError(t, err)
	assert.Contains(t, plan.Nodes, NodeKey("package:windows/amd64:nsis"))
	assert.Contains(t, plan.Nodes, NodeKey("package:windows/amd64:nsis:sign"))
	assert.Contains(t, plan.Nodes, NodeKey("package:linux/amd64:deb"))
	assert.NotContains(t, plan.Nodes, NodeKey("package:linux/amd64:deb:sign"))
}

func TestFrontendOutputAcceptsDirectoryRelativeAndProjectRelativePaths(t *testing.T) {
	assert.Equal(t, "frontend/dist", frontendOutputPath("frontend", "dist"))
	assert.Equal(t, "frontend/dist", frontendOutputPath("frontend", "frontend/dist"))
}

func TestHCLProfilePlansTheCompleteTargetAndFormatMatrix(t *testing.T) {
	root := t.TempDir()
	contents := `version = 3
project {
  name = "matrix"
  product_name = "Matrix"
  identifier = "com.example.matrix"
  version = "1.0.0"
}
frontend {
  directory = "frontend"
  install = ["npm", "install"]
  build = ["npm", "run", "build"]
  dev = ["npm", "run", "dev"]
  output = "frontend/dist"
}
profile "release" {
  target "windows/amd64" {
    formats = ["nsis", "msix"]
    sign = true
  }
	  target "darwin/universal" {
	    formats = ["app", "dmg"]
	    sign = true
	    notarize = true
	    destination = "dist/macos"
	  }
	  target "linux/arm64" {
	    formats = ["appimage", "deb", "rpm", "archlinux"]
	    sign = true
	  }
	  target "ios/arm64" {
	    formats = ["app"]
	    sign = true
	  }
	  target "android/amd64" {
	    formats = ["apk", "aab"]
	    sign = true
  }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(contents), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := PlanBuild(loaded.Config, Request{Verb: "build"})
	require.NoError(t, err)

	for _, artifact := range []string{
		"package:windows/amd64:nsis", "package:windows/amd64:msix",
		"package:darwin/universal:app", "package:darwin/universal:dmg",
		"package:linux/arm64:appimage", "package:linux/arm64:deb", "package:linux/arm64:rpm", "package:linux/arm64:archlinux",
		"package:ios/arm64:app", "package:android/amd64:apk", "package:android/amd64:aab",
	} {
		assert.Contains(t, plan.Nodes, NodeKey(artifact))
	}
	assert.Contains(t, plan.Nodes, NodeKey("package:windows/amd64:nsis:sign"))
	assert.Contains(t, plan.Nodes, NodeKey("package:windows/amd64:msix:sign"))
	assert.Contains(t, plan.Nodes, NodeKey("package:darwin/universal:app:sign"))
	assert.Contains(t, plan.Nodes, NodeKey("package:linux/arm64:deb:sign"))
	assert.Contains(t, plan.Nodes, NodeKey("package:ios/arm64:app:sign"))
	assert.Contains(t, plan.Nodes, NodeKey("package:android/amd64:apk:sign"))
	assert.Contains(t, plan.Nodes, NodeKey("frontend:bindings"))
	assert.Contains(t, plan.Nodes, NodeKey("frontend:build"))
	assert.Equal(t, "windows/amd64,darwin/universal,linux/arm64,ios/arm64,android/amd64", plan.Target)
	assert.Equal(t, "dist/macos", plan.Nodes[NodeKey("package:darwin/universal:app")].Spec.(PackageSpec).Variant)
}

func TestHCLTargetOverridesReachTheResolvedPipelineSpecs(t *testing.T) {
	root := t.TempDir()
	contents := `version = 3
project {
  name = "custom"
  product_name = "Custom"
  identifier = "com.example.custom"
  version = "1.0.0"
}
frontend {
  directory = "frontend"
  install = ["npm", "install"]
  build = ["npm", "run", "build"]
  dev = ["npm", "run", "dev"]
  output = "dist"
}
build {
  output = "bin"
  tags = ["release"]
  trim_path = true
  strip = true
  ldflags = ["-s"]
}
linux {
  product_name = "Custom Linux"
  identifier = "com.example.custom.linux"
  capabilities = ["network"]
}
target "linux/arm64" {
  tags = ["enterprise"]
  minimum_version = "24.04"
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(contents), 0o644))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	plan, err := PlanBuild(loaded.Config, Request{Verb: "build", TargetOS: "linux", TargetArch: "arm64"})
	require.NoError(t, err)
	compile := plan.Nodes[NodeKey("target:linux/arm64:compile")].Spec.(CompileSpec)
	assert.Equal(t, "24.04", compile.MinimumVersion)
	assert.ElementsMatch(t, []string{"release", "enterprise", "production"}, compile.Tags)
	assert.Equal(t, "bin/custom", compile.Output)
	frontend := plan.Nodes[NodeKey("frontend:build")]
	assert.Equal(t, "frontend/dist", frontend.Output)
	packaged, err := PlanBuild(loaded.Config, Request{Verb: "package", TargetOS: "linux", TargetArch: "arm64", Formats: []string{"deb"}})
	require.NoError(t, err)
	packageSpec := packaged.Nodes[NodeKey("package:linux/arm64:deb")].Spec.(PackageSpec)
	assert.Equal(t, "Custom Linux", packageSpec.Project.ProductName)
	assert.Equal(t, []string{"network"}, packageSpec.Capabilities)
}

func TestHCLProfileDestinationsReachEverySupportedTarget(t *testing.T) {
	root := t.TempDir()
	contents := `version = 3
project {
  name = "destinations"
  product_name = "Destinations"
  identifier = "com.example.destinations"
  version = "1.0.0"
}
profile "release" {
  target "windows/amd64" { destination = "windows-output" }
  target "darwin/amd64" { destination = "darwin-output" }
  target "linux/amd64" { destination = "linux-output" }
  target "ios/arm64" { destination = "ios-device" }
  target "android/arm64" { destination = "android-output" }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(contents), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := PlanBuild(loaded.Config, Request{Verb: "build"})
	require.NoError(t, err)
	for _, target := range []struct {
		name, variant string
	}{
		{"windows/amd64", "windows-output"},
		{"darwin/amd64", "darwin-output"},
		{"linux/amd64", "linux-output"},
		{"ios/arm64", "ios-device"},
		{"android/arm64", "android-output"},
	} {
		node, ok := plan.Nodes[NodeKey("target:"+target.name+":compile")]
		require.True(t, ok, target.name)
		spec, ok := node.Spec.(CompileSpec)
		require.True(t, ok, target.name)
		assert.Equal(t, target.variant, spec.Variant, target.name)
	}
}

func TestHCLProfileSigningWithoutFormatsUsesPlatformDefaults(t *testing.T) {
	root := t.TempDir()
	contents := `version = 3
project {
  name = "signed-default"
  product_name = "Signed Default"
  identifier = "com.example.signed-default"
  version = "1.0.0"
}
profile "release" {
  target "linux/amd64" { sign = true }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(contents), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := PlanBuild(loaded.Config, Request{Verb: "build"})
	require.NoError(t, err)
	assert.Contains(t, plan.Nodes, NodeKey("package:linux/amd64:appimage"))
	assert.Contains(t, plan.Nodes, NodeKey("package:linux/amd64:appimage:sign"))
}

func TestHCLProfileDestinationsReachArmAnd386Targets(t *testing.T) {
	root := t.TempDir()
	contents := `version = 3
project {
  name = "destinations"
  product_name = "Destinations"
  identifier = "com.example.destinations"
  version = "1.0.0"
}
profile "release" {
  target "linux/arm" { destination = "arm-output" }
  target "windows/386" { destination = "x86-output" }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(contents), 0o644))
	loaded, err := manifest.Load(root, "release")
	require.NoError(t, err)
	plan, err := PlanBuild(loaded.Config, Request{Verb: "build"})
	require.NoError(t, err)
	assert.Equal(t, "arm-output", plan.Nodes[NodeKey("target:linux/arm:compile")].Spec.(CompileSpec).Variant)
	assert.Equal(t, "x86-output", plan.Nodes[NodeKey("target:windows/386:compile")].Spec.(CompileSpec).Variant)
}

func TestHCLPlanTracksLocalGoModuleReplacements(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "app")
	dependency := filepath.Join(parent, "dependency")
	require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend"), 0o755))
	require.NoError(t, os.MkdirAll(dependency, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(`version = 3
project {
  name = "app"
  product_name = "App"
  identifier = "com.example.app"
  version = "1.0.0"
}
`), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example/app\n\ngo 1.24\n\nreplace example/dependency => ../dependency\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dependency, "go.mod"), []byte("module example/dependency\n"), 0o644))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	plan, err := PlanBuild(loaded.Config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	compile := plan.Nodes[NodeKey("target:linux/amd64:compile")]
	found := false
	for _, input := range compile.Inputs {
		if input.Root == dependency {
			found = true
			break
		}
	}
	assert.True(t, found, "HCL-first plans must snapshot local replacement modules")
}

func TestHCLPlanTracksGoWorkspaceModules(t *testing.T) {
	workspace := t.TempDir()
	root := filepath.Join(workspace, "app")
	dependency := filepath.Join(workspace, "dependency")
	require.NoError(t, os.MkdirAll(root, 0o755))
	require.NoError(t, os.MkdirAll(dependency, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(workspace, "go.work"), []byte("go 1.24\n\nuse (\n\t./app\n\t./dependency\n)\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(`version = 3
project {
  name = "app"
  product_name = "App"
  identifier = "com.example.app"
  version = "1.0.0"
}
`), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dependency, "go.mod"), []byte("module example/dependency\n"), 0o644))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	plan, err := PlanBuild(loaded.Config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	compile := plan.Nodes[NodeKey("target:linux/amd64:compile")]
	found := false
	for _, input := range compile.Inputs {
		if input.Root == dependency {
			found = true
			break
		}
	}
	assert.True(t, found, "HCL-first plans must snapshot modules listed by a parent go.work")
}

func TestHCLMultiTargetExecutorSharesProjectStages(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend", "node_modules"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "package.json"), []byte("{}"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "package-lock.json"), []byte("{}"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/multi\n\ngo 1.24\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nfunc main() {}\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(`version = 3
project {
  name = "multi"
  product_name = "Multi"
  identifier = "com.example.multi"
  version = "1.0.0"
}
`), 0o644))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	plan, err := PlanBuild(loaded.Config, Request{Verb: "build", Targets: []Target{{OS: "linux", Arch: "amd64"}, {OS: "linux", Arch: "arm64"}}})
	require.NoError(t, err)
	handler := &fakeHandler{root: root}
	results, err := (Executor{Handler: handler}).Execute(context.Background(), plan, ExecuteOptions{Root: root, Reporter: report.Nop{}})
	require.NoError(t, err)
	assert.Len(t, results, len(plan.Nodes))
	assert.Equal(t, 1, countNodeRuns(handler.runs, NodeKey("frontend:install")))
	assert.Equal(t, 1, countNodeRuns(handler.runs, NodeKey("frontend:bindings")))
	assert.Equal(t, 1, countNodeRuns(handler.runs, NodeKey("frontend:build")))
	assert.Equal(t, 1, countNodeRuns(handler.runs, NodeKey("target:linux/amd64:compile")))
	assert.Equal(t, 1, countNodeRuns(handler.runs, NodeKey("target:linux/arm64:compile")))
}

func countNodeRuns(runs []NodeKey, want NodeKey) int {
	count := 0
	for _, run := range runs {
		if run == want {
			count++
		}
	}
	return count
}
