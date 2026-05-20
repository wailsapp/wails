package commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWindowsTaskfilePassesArchToSyso(t *testing.T) {
	data, err := buildAssets.ReadFile("build_assets/windows/Taskfile.yml")
	require.NoError(t, err)

	content := string(data)

	t.Run("build:native passes ARCH to generate:syso", func(t *testing.T) {
		inBuildNative := false
		found := false
		lines := strings.Split(content, "\n")

		for i := 0; i < len(lines); i++ {
			line := lines[i]
			if strings.Contains(line, "build:native:") {
				inBuildNative = true
				continue
			}
			if !inBuildNative {
				continue
			}
			if strings.Contains(line, "- task: generate:syso") {
				for j := i + 1; j < len(lines) && j <= i+3; j++ {
					if strings.Contains(lines[j], "ARCH:") {
						found = true
						break
					}
				}
				break
			}
		}
		assert.True(t, found, "build:native should pass ARCH to generate:syso")
	})

	t.Run("build:docker passes ARCH to generate:syso", func(t *testing.T) {
		inBuildDocker := false
		found := false
		lines := strings.Split(content, "\n")

		for i := 0; i < len(lines); i++ {
			line := lines[i]
			if strings.Contains(line, "build:docker:") {
				inBuildDocker = true
				continue
			}
			if !inBuildDocker {
				continue
			}
			if strings.Contains(line, "- task: generate:syso") {
				for j := i + 1; j < len(lines) && j <= i+3; j++ {
					if strings.Contains(lines[j], "ARCH:") {
						found = true
						break
					}
				}
				break
			}
		}
		assert.True(t, found, "build:docker should pass ARCH to generate:syso")
	})
}
