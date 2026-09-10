package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageLoadingUsesExplicitProjectDirectory(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/explicitroot\n\ngo 1.25\n"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "service.go"), []byte("package explicitroot\nfunc ExampleService() {}\n"), 0600))
	before, err := os.Getwd()
	require.NoError(t, err)
	packages, err := loadPackages(root, nil, ".")
	require.NoError(t, err)
	require.Len(t, packages, 1)
	require.Equal(t, "example.com/explicitroot", packages[0].PkgPath)
	require.Empty(t, packages[0].Errors)
	require.NotNil(t, packages[0].Types.Scope().Lookup("ExampleService"))
	patterns, err := resolvePatterns(root, nil, ".")
	require.NoError(t, err)
	require.Equal(t, []string{"example.com/explicitroot"}, patterns)
	after, err := os.Getwd()
	require.NoError(t, err)
	require.Equal(t, before, after, "loading another module must not change the process working directory")
}
