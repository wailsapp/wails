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

func TestEncodeDocumentStaysSparse(t *testing.T) {
	doc := NewDocument(Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"})
	doc.Frontend.PackageManager = "pnpm"
	doc.Wake.Migration = &Migration{CompletedBy: "v3", Complete: true, Sources: map[string]string{"Taskfile.yml": "digest"}}
	data, err := EncodeDocument(doc)
	require.NoError(t, err)
	text := string(data)
	assert.Contains(t, text, `package_manager = "pnpm"`)
	assert.NotContains(t, text, "[wake.migration]")
	assert.NotContains(t, text, "[targets]")
	assert.NotContains(t, text, "debounce_ms")
}

func TestEncodeDocumentKeepsOnlyIncompleteMigrationMarker(t *testing.T) {
	doc := NewDocument(Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"})
	doc.Wake.Migration = &Migration{CompletedBy: "v3", Complete: false, Sources: map[string]string{"Taskfile.yml": "digest"}}
	data, err := EncodeDocument(doc)
	require.NoError(t, err)
	text := string(data)
	assert.Contains(t, text, "[wake.migration]")
	assert.Contains(t, text, "complete = false")
	assert.NotContains(t, text, "completed_by")
	assert.NotContains(t, text, "[wake.migration.sources]")
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

func TestMigrationControlsActivation(t *testing.T) {
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
	active, err := Active(root)
	require.NoError(t, err)
	assert.False(t, active)
}
