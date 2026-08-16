package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

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
