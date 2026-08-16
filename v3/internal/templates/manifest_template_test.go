package templates

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommonTemplateUsesMinimalManifest(t *testing.T) {
	data, err := templates.ReadFile("_common/wails.tmpl.toml")
	require.NoError(t, err)
	content := string(data)
	assert.Contains(t, content, `[project]`)
	assert.Contains(t, content, `name = "{{.ProjectName}}"`)
	assert.Contains(t, content, `# binary_name = "{{.BinaryName}}"`)
	assert.False(t, strings.Contains(content, "[build]"))
	_, err = templates.ReadFile("_common/Taskfile.tmpl.yml")
	assert.Error(t, err)
}
