package commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCrossDockerfileMatchesLinuxSupportContract(t *testing.T) {
	dockerfile, err := buildAssets.ReadFile("build_assets/docker/Dockerfile.cross")
	require.NoError(t, err)

	contents := string(dockerfile)
	require.Contains(t, contents, "FROM golang:1.26-trixie")
	require.Contains(t, contents, "pkg-config --atleast-version=4.14 gtk4")

	for _, dependency := range []string{
		"libgtk-4-dev",
		"libwebkitgtk-6.0-dev",
		"libgtk-3-dev",
		"libwebkit2gtk-4.1-dev",
	} {
		require.True(t, strings.Contains(contents, dependency), "cross image must install %s", dependency)
	}
}
