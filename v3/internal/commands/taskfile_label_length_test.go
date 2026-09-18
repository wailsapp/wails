package commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskfileLabelsDoNotIncludeBuildFlags(t *testing.T) {
	data, err := buildAssets.ReadFile("build_assets/Taskfile.tmpl.yml")
	require.NoError(t, err)

	content := string(data)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "label:") && strings.Contains(trimmed, "BUILD_FLAGS") {
			t.Errorf("Taskfile label should not include BUILD_FLAGS (causes >255 char checksum filenames): %s", trimmed)
		}
	}
}

func TestTaskfileLabelLength(t *testing.T) {
	data, err := buildAssets.ReadFile("build_assets/Taskfile.tmpl.yml")
	require.NoError(t, err)

	content := string(data)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "label:") {
			assert.Less(t, len(trimmed), 200,
				"Taskfile label should be short enough to avoid checksum filename length issues")
		}
	}
}
