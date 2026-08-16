package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestLegacyDevRejectsManifestOnlyFlags(t *testing.T) {
	t.Chdir(t.TempDir())
	err := Dev(&DevOptions{Profile: "release"})
	assert.ErrorContains(t, err, "require an active wails.toml")
}

func TestDevPathExcludedAlwaysSkipsGeneratedAndDependencyTrees(t *testing.T) {
	config := manifest.Config{Build: manifest.Build{OutputDirectory: "out"}, Frontend: manifest.Frontend{Directory: "frontend"}}
	for _, path := range []string{".wails/cache", ".git/objects", "out/app", "frontend/src/main.ts", "tools/node_modules/pkg"} {
		assert.True(t, devPathExcluded(config, path), path)
	}
	assert.False(t, devPathExcluded(config, "internal/service.go"))
}

func TestMatchesDevWatchDoubleStar(t *testing.T) {
	patterns := []string{"**/*.go", "wails.toml"}
	assert.True(t, matchesDevWatch(patterns, "internal/service.go"))
	assert.True(t, matchesDevWatch(patterns, "main.go"))
	assert.True(t, matchesDevWatch(patterns, "wails.toml"))
	assert.False(t, matchesDevWatch(patterns, "README.md"))
}

func TestFrontendSessionChanged(t *testing.T) {
	base := manifest.Config{Frontend: manifest.Frontend{Directory: "frontend", PackageManager: "npm", DevCommand: "dev"}}
	assert.False(t, frontendSessionChanged(base, 9245, base, 9245))
	next := base
	next.Frontend.DevCommand = "serve"
	assert.True(t, frontendSessionChanged(base, 9245, next, 9245))
	assert.True(t, frontendSessionChanged(base, 9245, base, 5173))
}
