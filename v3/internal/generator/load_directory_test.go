package generator

import (
	"go/types"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageLoadingUsesIsolatedCGoEnvironment(t *testing.T) {
	root := t.TempDir()
	for name, content := range map[string]string{
		"go.mod":  "module example.com/backend\n\ngo 1.26.0\n",
		"pure.go": "//go:build !cgo\n\npackage backend\nconst Backend = \"purego\"\n",
		"cgo.go":  "//go:build cgo\n\npackage backend\nconst Backend = \"cgo\"\n",
	} {
		require.NoError(t, os.WriteFile(filepath.Join(root, name), []byte(content), 0o644))
	}
	t.Setenv("CGO_ENABLED", "1")
	t.Setenv("CC", filepath.Join(root, "no-c-compiler"))
	environment := append(os.Environ(), "CGO_ENABLED=0", "GOWORK=off")
	loader := packageLoader{environment: environment}
	loaded, err := loader.loadPackages(root, nil, ".")
	require.NoError(t, err)
	require.Len(t, loaded, 1)
	require.Empty(t, loaded[0].Errors)
	backend, ok := loaded[0].Types.Scope().Lookup("Backend").(*types.Const)
	require.True(t, ok)
	require.Equal(t, `"purego"`, backend.Val().ExactString())
	require.Equal(t, "1", os.Getenv("CGO_ENABLED"))
	paths, err := loader.resolvePatterns(root, nil, ".")
	require.NoError(t, err)
	require.Equal(t, []string{"example.com/backend"}, paths)
}

func TestPackageLoadingUsesExplicitProjectDirectory(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/explicitroot\n\ngo 1.25\n"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "service.go"), []byte("package explicitroot\nfunc ExampleService() {}\n"), 0600))
	before, err := os.Getwd()
	require.NoError(t, err)
	packages, err := (packageLoader{}).loadPackages(root, nil, ".")
	require.NoError(t, err)
	require.Len(t, packages, 1)
	require.Equal(t, "example.com/explicitroot", packages[0].PkgPath)
	require.Empty(t, packages[0].Errors)
	require.NotNil(t, packages[0].Types.Scope().Lookup("ExampleService"))
	patterns, err := (packageLoader{}).resolvePatterns(root, nil, ".")
	require.NoError(t, err)
	require.Equal(t, []string{"example.com/explicitroot"}, patterns)
	after, err := os.Getwd()
	require.NoError(t, err)
	require.Equal(t, before, after, "loading another module must not change the process working directory")
}
