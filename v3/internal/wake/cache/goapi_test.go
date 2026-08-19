package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnapshotGoAPIStableMetadataFastPathAndBuildConstraints(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "service.go")
	source := "//go:build linux\n// +build linux\n\npackage app\n\ntype Service struct{}\nfunc (Service) Greet() string { return \"hello\" }\n"
	require.NoError(t, os.WriteFile(path, []byte(source), 0o644))
	stable := time.Now().Add(-metadataFastPathWindow - time.Second)
	require.NoError(t, os.Chtimes(path, stable, stable))
	store, err := OpenCache(root)
	require.NoError(t, err)
	first, err := store.SnapshotGoAPI(SnapshotOptions{Label: "bindings", Root: "."})
	require.NoError(t, err)
	store.ResetStats()
	second, err := store.SnapshotGoAPI(SnapshotOptions{Label: "bindings", Root: "."})
	require.NoError(t, err)
	assert.Equal(t, first, second)
	assert.Equal(t, 1, store.Stats().DigestsReused)
}

func TestSnapshotGoAPIIgnoresMethodBodiesButTracksSignatures(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "service.go")
	write := func(source string) {
		require.NoError(t, os.WriteFile(path, []byte(source), 0o644))
	}
	store, err := OpenCache(root)
	require.NoError(t, err)
	write("package app\ntype Service struct{}\nfunc (Service) Greet(name string) string { return name }\n")
	first, err := store.SnapshotGoAPI(SnapshotOptions{Label: "bindings", Root: root})
	require.NoError(t, err)
	write("package app\ntype Service struct{}\nfunc (Service) Greet(name string) string { return \"hello \" + name }\n")
	second, err := store.SnapshotGoAPI(SnapshotOptions{Label: "bindings", Root: root})
	require.NoError(t, err)
	assert.Equal(t, first, second)
	write("package app\ntype Service struct{}\nfunc (Service) Greet(name string, loud bool) string { return name }\n")
	third, err := store.SnapshotGoAPI(SnapshotOptions{Label: "bindings", Root: root})
	require.NoError(t, err)
	assert.NotEqual(t, second, third)
}

func TestSnapshotGoAPITracksMainRegistrationBody(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "main.go")
	store, err := OpenCache(root)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, []byte("package main\nfunc main(){ register(A{}) }\n"), 0o644))
	first, err := store.SnapshotGoAPI(SnapshotOptions{Label: "bindings", Root: root})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, []byte("package main\nfunc main(){ register(B{}) }\n"), 0o644))
	second, err := store.SnapshotGoAPI(SnapshotOptions{Label: "bindings", Root: root})
	require.NoError(t, err)
	assert.NotEqual(t, first, second)
}

func TestSnapshotGoAPITracksRegistrationHelperBodies(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "services.go")
	store, err := OpenCache(root)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, []byte("package main\nfunc services() any { return A{} }\n"), 0o644))
	first, err := store.SnapshotGoAPI(SnapshotOptions{Label: "bindings", Root: root})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, []byte("package main\nfunc services() any { return B{} }\n"), 0o644))
	second, err := store.SnapshotGoAPI(SnapshotOptions{Label: "bindings", Root: root})
	require.NoError(t, err)
	assert.NotEqual(t, first, second)
}

func TestSnapshotGoAPIIgnoresGitIgnoredDevelopmentPackages(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "ignored"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".gitignore"), []byte("ignored/\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nfunc main() {}\n"), 0o644))
	ignoredPath := filepath.Join(root, "ignored", "service.go")
	require.NoError(t, os.WriteFile(ignoredPath, []byte("package ignored\ntype Service struct{}\n"), 0o644))
	store, err := OpenCache(root)
	require.NoError(t, err)
	options := SnapshotOptions{Label: "bindings", Root: root, UseGitIgnore: true}
	before, err := store.SnapshotGoAPI(options)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(ignoredPath, []byte("package ignored\ntype ChangedService struct{}\n"), 0o644))
	after, err := store.SnapshotGoAPI(options)
	require.NoError(t, err)
	assert.Equal(t, before, after)
}

func TestSnapshotGoAPIIgnoresRetainedNonGoAndTestInputs(t *testing.T) {
	root := t.TempDir()
	mainPath := filepath.Join(root, "main.go")
	testPath := filepath.Join(root, "main_test.go")
	readmePath := filepath.Join(root, "README.md")
	require.NoError(t, os.WriteFile(mainPath, []byte("package main\ntype Service struct{}\n"), 0o644))
	require.NoError(t, os.WriteFile(testPath, []byte("not valid Go"), 0o644))
	require.NoError(t, os.WriteFile(readmePath, []byte("first"), 0o644))

	store, err := OpenCache(root)
	require.NoError(t, err)
	options := SnapshotOptions{Label: "bindings", Root: root, IncludeAll: true}
	before, err := store.SnapshotGoAPI(options)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(testPath, []byte("still not valid Go"), 0o644))
	require.NoError(t, os.WriteFile(readmePath, []byte("second"), 0o644))
	afterIgnoredEdits, err := store.SnapshotGoAPI(options)
	require.NoError(t, err)
	assert.Equal(t, before, afterIgnoredEdits)

	require.NoError(t, os.WriteFile(mainPath, []byte("package main\ntype RenamedService struct{}\n"), 0o644))
	afterAPIEdit, err := store.SnapshotGoAPI(options)
	require.NoError(t, err)
	assert.NotEqual(t, afterIgnoredEdits, afterAPIEdit)
}
