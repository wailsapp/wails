package packagetemplate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestRenderFileUsesStablePackageModel(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "package.tmpl")
	destination := filepath.Join(root, "generated", "package.conf")
	require.NoError(t, os.WriteFile(source, []byte("{{.Project.ProductName}}|{{.Target.OS}}/{{.Target.Arch}}|{{index .Target.Capabilities 0}}|{{.Package.Format}}|{{index .Options \"channel\"}}\n"), 0o644))

	err := Render(source, destination, Model{
		Version: 1,
		Project: manifest.Project{ProductName: "Example"},
		Target:  Target{OS: "linux", Arch: "arm64", Capabilities: []string{"network"}},
		Package: Package{Format: "deb"},
		Options: map[string]any{"channel": "preview"},
	})
	require.NoError(t, err)
	data, err := os.ReadFile(destination)
	require.NoError(t, err)
	assert.Equal(t, "Example|linux/arm64|network|deb|preview\n", string(data))
}

func TestRenderDirectoryRendersTemplateFilesAndCopiesOwnedFiles(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "android")
	destination := filepath.Join(root, "generated")
	require.NoError(t, os.MkdirAll(filepath.Join(source, "app"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(source, "app", "build.gradle.tmpl"), []byte("applicationId '{{.Project.Identifier}}'\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(source, "gradlew"), []byte("#!/bin/sh\n"), 0o755))

	err := Render(source, destination, Model{Version: 1, Project: manifest.Project{Identifier: "com.example.app"}})
	require.NoError(t, err)
	gradle, err := os.ReadFile(filepath.Join(destination, "app", "build.gradle"))
	require.NoError(t, err)
	assert.Equal(t, "applicationId 'com.example.app'\n", string(gradle))
	info, err := os.Stat(filepath.Join(destination, "gradlew"))
	require.NoError(t, err)
	assert.NotZero(t, info.Mode().Perm()&0o100)
	_, err = os.Stat(filepath.Join(destination, "app", "build.gradle.tmpl"))
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestRenderModelExposesResolvedInputsWithoutPipelineInternals(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "package.tmpl")
	destination := filepath.Join(root, "package.conf")
	require.NoError(t, os.WriteFile(source, []byte("{{.Paths.Binary}}|{{.Paths.Output}}|{{(index .Associations 0).Extensions}}|{{(index .Protocols 0).Scheme}}\n"), 0o644))

	err := Render(source, destination, Model{
		Paths:        Paths{Binary: "/project/bin/app", Output: "/project/bin/app.deb"},
		Associations: []manifest.Association{{Extensions: []string{"demo"}}},
		Protocols:    []manifest.Protocol{{Scheme: "example"}},
	})
	require.NoError(t, err)
	data, err := os.ReadFile(destination)
	require.NoError(t, err)
	assert.Equal(t, "/project/bin/app|/project/bin/app.deb|[demo]|example\n", string(data))
}

func TestRenderRejectsUnknownFieldsAndDirectorySymlinks(t *testing.T) {
	root := t.TempDir()
	missing := filepath.Join(root, "missing.tmpl")
	require.NoError(t, os.WriteFile(missing, []byte("{{.InternalNode}}"), 0o644))
	err := Render(missing, filepath.Join(root, "missing.out"), Model{Version: 1})
	require.ErrorContains(t, err, "can't evaluate field InternalNode")

	source := filepath.Join(root, "tree")
	require.NoError(t, os.MkdirAll(source, 0o755))
	require.NoError(t, os.Symlink(missing, filepath.Join(source, "linked.tmpl")))
	err = Render(source, filepath.Join(root, "generated"), Model{Version: 1})
	require.ErrorContains(t, err, "unsupported symlink")
}
