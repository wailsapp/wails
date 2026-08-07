package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadPackagesAndroidWithoutCGO(t *testing.T) {
	t.Setenv("GOOS", "android")
	t.Setenv("CGO_ENABLED", "0")

	packages, err := LoadPackages(
		[]string{"-tags", "android,debug"},
		"github.com/wailsapp/wails/v3/pkg/application",
	)
	require.NoError(t, err)
	require.NotEmpty(t, packages, "LoadPackages returned no packages")

	for _, pkg := range packages {
		require.Empty(t, pkg.Errors, "package %s has loading errors", pkg.PkgPath)
	}
}
