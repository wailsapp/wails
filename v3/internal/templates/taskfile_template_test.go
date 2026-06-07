package templates

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommonTaskfileUsesBinaryName(t *testing.T) {
	data, err := templates.ReadFile("_common/Taskfile.tmpl.yml")
	require.NoError(t, err, "_common/Taskfile.tmpl.yml should be present in embedded templates")
	content := string(data)
	assert.True(t, strings.Contains(content, `{{.BinaryName}}`),
		"root Taskfile template APP_NAME should use {{.BinaryName}}, not {{.ProjectName}}")
	assert.False(t, strings.Contains(content, `{{.ProjectName}}`),
		"root Taskfile template should not fall back to {{.ProjectName}} for APP_NAME")
}
