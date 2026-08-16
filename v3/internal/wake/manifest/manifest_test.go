package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadMinimalAndInferFrontend(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "pnpm-lock.yaml"), nil, 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "tsconfig.json"), nil, 0o644))
	require.NoError(t, WriteMinimal(root, Project{Name: "My App", ProductName: "My App", Identifier: "com.example.app", Version: "0.1.0"}))
	loaded, err := Load(root, "")
	require.NoError(t, err)
	assert.Equal(t, "my-app", loaded.Config.Project.BinaryName)
	assert.Equal(t, "pnpm", loaded.Config.Frontend.PackageManager)
	assert.True(t, loaded.Config.Frontend.Bindings.TypeScript)
	assert.True(t, loaded.Config.Build.Production)
}

func TestExplicitBindingsWinOverInference(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[frontend.bindings]
typescript=true
interfaces=true
time_type="string"
`), 0o644))
	loaded, err := Load(root, "")
	require.NoError(t, err)
	assert.True(t, loaded.Config.Frontend.Bindings.TypeScript)
	assert.True(t, loaded.Config.Frontend.Bindings.Interfaces)
	assert.Equal(t, "string", loaded.Config.Frontend.Bindings.TimeType)
}

func TestProfileIsSparseOverlay(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[build]
output_directory="out"

[profiles.debug.build]
production=false
`), 0o644))
	loaded, err := Load(root, "debug")
	require.NoError(t, err)
	assert.Equal(t, "out", loaded.Config.Build.OutputDirectory)
	assert.False(t, loaded.Config.Build.Production)
}

func TestUnknownFieldIsError(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"
surprise=true
`), 0o644))
	_, err := Load(root, "")
	require.ErrorContains(t, err, "project.surprise")
}

func TestUnknownProfileFieldIsError(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[profiles.debug.build]
surprise=true
`), 0o644))
	_, err := Load(root, "debug")
	require.ErrorContains(t, err, "profiles.debug.build.surprise")
}

func TestEjectFreezesOnceAndPreservesResolvedValues(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, WriteMinimal(root, Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}))
	require.NoError(t, Eject(root, "", "v3.0.0", false))
	loaded, err := Load(root, "")
	require.NoError(t, err)
	assert.Equal(t, "v3.0.0", loaded.Config.Wake.EjectedBy)
	assert.Equal(t, 9245, loaded.Config.Dev.Port)
	before, err := os.ReadFile(filepath.Join(root, Filename))
	require.NoError(t, err)
	err = Eject(root, "", "v3.0.1", false)
	require.ErrorIs(t, err, ErrEjectionSuggestionsUnavailable)
	after, readErr := os.ReadFile(filepath.Join(root, Filename))
	require.NoError(t, readErr)
	assert.Equal(t, before, after)
}

func TestEjectSpecificProfile(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[profiles.debug.build]
production=false
`), 0o644))
	require.NoError(t, Eject(root, "debug", "v3.0.0", false))
	data, err := os.ReadFile(filepath.Join(root, Filename))
	require.NoError(t, err)
	baseEnd := strings.Index(string(data), "[profiles.debug]")
	require.NotEqual(t, -1, baseEnd)
	assert.NotContains(t, string(data[:baseEnd]), "debounce_ms")
	loaded, err := Load(root, "debug")
	require.NoError(t, err)
	assert.False(t, loaded.Config.Build.Production)
	assert.Equal(t, "v3.0.0", loaded.Config.Wake.EjectedProfiles["debug"])
	require.ErrorIs(t, Eject(root, "debug", "v3.0.1", false), ErrEjectionSuggestionsUnavailable)
}

func TestEjectMissingProfileCopiesEffectiveBase(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, WriteMinimal(root, Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}))
	require.NoError(t, Eject(root, "release", "v3.0.0", false))
	loaded, err := Load(root, "release")
	require.NoError(t, err)
	assert.True(t, loaded.Config.Build.Production)
	assert.Equal(t, "v3.0.0", loaded.Config.Wake.EjectedProfiles["release"])
}

func TestProfileCannotOverrideTargetIdentity(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[profiles.release.targets.darwin.arm64]
build_number=42
`), 0o644))
	_, err := Load(root, "release")
	require.ErrorContains(t, err, "cannot override target identity field targets.darwin.arm64.build_number")
}

func TestEncodeDocumentPreservesFalseAgainstTrueDefault(t *testing.T) {
	doc := NewDocument(Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"})
	doc.Build.Production = false
	data, err := EncodeDocument(doc)
	require.NoError(t, err)
	assert.Contains(t, string(data), "production = false")
}

func TestEncodeDocumentPreservesZeroValuesInExtensionLists(t *testing.T) {
	doc := NewDocument(Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"})
	doc.Extensions = map[string]map[string]any{"example": {"values": []any{0, false, ""}}}
	data, err := EncodeDocument(doc)
	require.NoError(t, err)
	assert.Contains(t, string(data), `values = [0, false, ""]`)
}

func TestHookFieldsRejectWrongTypes(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[hooks.before_build]
script=42
`), 0o644))
	_, err := Load(root, "")
	require.ErrorContains(t, err, `hook field "script" must be a string`)
}

func TestPackageOptionsAreTypedUnlessOwnedByCustomTemplate(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "packaging"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "packaging", "AppxManifest.xml.tmpl"), []byte("manifest"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[package.darwin.dmg.options]
window_width=640

[package.linux.appimage.options]
categories="Development;"

[package.windows.msix]
template="packaging/AppxManifest.xml.tmpl"
[package.windows.msix.options]
channel="preview"
`), 0o644))
	loaded, err := Load(root, "")
	require.NoError(t, err)
	assert.EqualValues(t, 640, loaded.Config.Package.Darwin.DMG.Options["window_width"])
	assert.Equal(t, "preview", loaded.Config.Package.Windows.MSIX.Options["channel"])

	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"
[package.windows.msix.options]
channel="preview"
`), 0o644))
	_, err = Load(root, "")
	require.ErrorContains(t, err, "package.windows.msix.options requires a custom template")
}

func TestPackageTemplateMustExistInsideProject(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[package.windows.msix]
template="packaging/missing.xml.tmpl"
`), 0o644))

	_, err := Load(root, "")
	require.ErrorContains(t, err, "package.windows.msix.template")
	require.ErrorContains(t, err, "does not exist")
}

func TestPackageTemplateSymlinkMustResolveInsideProject(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(outside, "template.tmpl"), []byte("outside"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "packaging"), 0o755))
	if err := os.Symlink(filepath.Join(outside, "template.tmpl"), filepath.Join(root, "packaging", "template.tmpl")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[package.windows.msix]
template="packaging/template.tmpl"
`), 0o644))

	_, err := Load(root, "")
	require.ErrorContains(t, err, "package.windows.msix.template")
	require.ErrorContains(t, err, "must resolve inside the project")
}

func TestInvalidDevWatchPatternFailsAtLoad(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[dev]
watch=["[invalid"]
`), 0o644))

	_, err := Load(root, "")
	require.ErrorContains(t, err, `dev.watch[0] "[invalid"`)
	require.ErrorContains(t, err, "invalid pattern")
}

func TestDMGFileOptionsCannotEscapeProject(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[package.darwin.dmg.options]
background="../background.png"
`), 0o644))

	_, err := Load(root, "")
	require.ErrorContains(t, err, "package.darwin.dmg.options.background must be project-relative")
}

func TestEncodeDocumentStaysSparse(t *testing.T) {
	doc := NewDocument(Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"})
	doc.Frontend.PackageManager = "pnpm"
	data, err := EncodeDocument(doc)
	require.NoError(t, err)
	text := string(data)
	assert.Contains(t, text, `package_manager = "pnpm"`)
	assert.NotContains(t, text, "[targets]")
	assert.NotContains(t, text, "debounce_ms")
}

func TestCachedHookLongFormLoadsInputsAndOutputs(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[hooks.before_build]
script="scripts/generate.sh"
cache=true
inputs=["version.txt"]
outputs=["generated/version.go"]
`), 0o644))
	loaded, err := Load(root, "")
	require.NoError(t, err)
	assert.Equal(t, []string{"version.txt"}, loaded.Config.Hooks.BeforeBuild.Inputs)
	assert.Equal(t, []string{"generated/version.go"}, loaded.Config.Hooks.BeforeBuild.Outputs)
}

func TestCachedHookOutputsMustHaveOneSafeArtifactRoot(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[hooks.before_build]
script="generate"
cache=true
inputs=["version.txt"]
outputs=["generated/version.go", "dist/version.txt"]
`), 0o644))
	_, err := Load(root, "")
	require.ErrorContains(t, err, "multiple outputs must share a top-level directory")

	output, err := HookOutputRoot([]string{"generated/version.go", "generated/version.txt"})
	require.NoError(t, err)
	assert.Equal(t, "generated", output)
}

func TestHookCacheDeclarationsRequireOptInAndCannotCaptureInputs(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[hooks.before_build]
script="scripts/generate.sh"
inputs=["version.txt"]
outputs=["generated/version.go"]
`), 0o644))
	_, err := Load(root, "")
	require.ErrorContains(t, err, "inputs and outputs require cache = true")

	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[hooks.before_build]
script="generated/generate.sh"
cache=true
inputs=["version.txt"]
outputs=["generated/version.go", "generated/metadata.json"]
`), 0o644))
	_, err = Load(root, "")
	require.ErrorContains(t, err, `output root "generated" contains input "generated/generate.sh"`)
}

func TestProjectPathsCannotEscape(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[build]
output_directory="..\\outside"
`), 0o644))
	_, err := Load(root, "")
	require.ErrorContains(t, err, "build.output_directory must be a project-relative path")
}

func TestSigningAndRegistrationPathsCannotEscape(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[signing.darwin]
entitlements="../outside.plist"
`), 0o644))
	_, err := Load(root, "")
	require.ErrorContains(t, err, "signing.darwin.entitlements must be project-relative")

	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[[associations]]
extensions=["demo"]
icon="../outside.ico"
`), 0o644))
	_, err = Load(root, "")
	require.ErrorContains(t, err, "associations[0].icon must be project-relative")
}

func TestRegistrationPlatformsAreValidated(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"

[[protocols]]
scheme="example"
platforms=["plan9"]
`), 0o644))
	_, err := Load(root, "")
	require.ErrorContains(t, err, `protocols[0].platforms contains unsupported platform "plan9"`)
}

func TestLegacyMigrationFieldsAreIgnoredAndNotReencoded(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(`[project]
name="app"
product_name="App"
identifier="com.example.app"
version="1.0.0"
[wake.migration]
completed_by="v3"
complete=false
[wake.migration.sources]
"Taskfile.yml"="digest"
`), 0o644))
	loaded, err := Load(root, "")
	require.NoError(t, err)
	data, err := EncodeDocument(loaded.Document)
	require.NoError(t, err)
	assert.NotContains(t, string(data), "migration")
}
