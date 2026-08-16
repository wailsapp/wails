package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnapshotGoAPIIgnoresMethodBodiesButTracksSignatures(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "service.go")
	write := func(source string) {
		require.NoError(t, os.WriteFile(path, []byte(source), 0o644))
	}
	store, err := OpenCache(root)
	require.NoError(t, err)
	write("package app\ntype Service struct{}\nfunc (Service) Greet(name string) string { return name }\n")
	first, err := store.SnapshotGoAPI("bindings", root, nil)
	require.NoError(t, err)
	write("package app\ntype Service struct{}\nfunc (Service) Greet(name string) string { return \"hello \" + name }\n")
	second, err := store.SnapshotGoAPI("bindings", root, nil)
	require.NoError(t, err)
	assert.Equal(t, first, second)
	write("package app\ntype Service struct{}\nfunc (Service) Greet(name string, loud bool) string { return name }\n")
	third, err := store.SnapshotGoAPI("bindings", root, nil)
	require.NoError(t, err)
	assert.NotEqual(t, second, third)
}

func TestSnapshotGoAPITracksMainRegistrationBody(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "main.go")
	store, err := OpenCache(root)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, []byte("package main\nfunc main(){ register(A{}) }\n"), 0o644))
	first, err := store.SnapshotGoAPI("bindings", root, nil)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, []byte("package main\nfunc main(){ register(B{}) }\n"), 0o644))
	second, err := store.SnapshotGoAPI("bindings", root, nil)
	require.NoError(t, err)
	assert.NotEqual(t, first, second)
}

func TestSnapshotGoAPITracksRegistrationHelperBodies(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "services.go")
	store, err := OpenCache(root)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, []byte("package main\nfunc services() any { return A{} }\n"), 0o644))
	first, err := store.SnapshotGoAPI("bindings", root, nil)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, []byte("package main\nfunc services() any { return B{} }\n"), 0o644))
	second, err := store.SnapshotGoAPI("bindings", root, nil)
	require.NoError(t, err)
	assert.NotEqual(t, first, second)
}
