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

// TestCommonTaskfileDispatchesViaGOOS guards the fix for #5615: the root
// build/package/run tasks in the template must dispatch to the platform
// Taskfile via the GOOS variable, so running them honours any customisations in
// the root Taskfile (for both native and cross builds) rather than the built-in
// {{OS}} which the CLI used to bypass with an OS-prefixed task name. The CLI
// routes `wails3 build` and `wails3 package` through this dispatch; `wails3 dev`
// is a direct command and is not affected.
func TestCommonTaskfileDispatchesViaGOOS(t *testing.T) {
	data, err := templates.ReadFile("_common/Taskfile.tmpl.yml")
	require.NoError(t, err, "_common/Taskfile.tmpl.yml should be present in embedded templates")
	// Note: the template escapes literal go-task delimiters via {{.Opn}}/{{.Cls}},
	// so the rendered output contains `{{.GOOS}}` rather than the raw text here.
	content := string(data)
	for _, verb := range []string{"build", "package", "run"} {
		assert.Contains(t, content, `{{.Opn}}.GOOS{{.Cls}}:`+verb,
			"root Taskfile %s task should dispatch via {{.GOOS}}", verb)
	}
	assert.Contains(t, content, `GOOS: '{{.Opn}}.GOOS | default OS{{.Cls}}'`,
		"root Taskfile should define a GOOS var defaulting to the host OS")
	assert.NotContains(t, content, `{{.Opn}}OS{{.Cls}}:`,
		"root Taskfile should not dispatch via the {{OS}} built-in (bypassed by the CLI)")
}
