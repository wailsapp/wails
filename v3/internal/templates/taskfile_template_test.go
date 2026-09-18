package templates

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	texttemplate "text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	taskruntime "github.com/wailsapp/task/v3"
	"github.com/wailsapp/task/v3/taskfile/ast"
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

func TestCommonTaskfileAppVariablesPreserveEnvironmentAndCLIOverrides(t *testing.T) {
	appVariables := []string{
		"APP_TAGS",
		"APP_TAGS_LINUX",
		"APP_TAGS_DARWIN",
		"APP_TAGS_WINDOWS",
		"APP_TAGS_ANDROID",
		"APP_TAGS_IOS",
		"APP_TAGS_SERVER",
		"APP_LDFLAGS",
		"APP_CGO_ENABLED",
	}

	for _, variable := range appVariables {
		t.Setenv(variable, "environment-"+variable)
	}

	projectDir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(projectDir, "Taskfile.yml"),
		[]byte(renderCommonTaskfile(t)),
		0o644,
	))
	writeStubIncludedTaskfiles(t, projectDir)

	executor := taskruntime.Executor{
		Dir:                 projectDir,
		Offline:             true,
		DisableVersionCheck: true,
		Stdin:               strings.NewReader(""),
		Stdout:              io.Discard,
		Stderr:              io.Discard,
	}
	require.NoError(t, executor.Setup())

	resolved, err := executor.Compiler.GetTaskfileVariables()
	require.NoError(t, err)
	resolvedValues := resolved.ToCacheMap()
	for _, variable := range appVariables {
		assert.Equal(t, "environment-"+variable, resolvedValues[variable],
			"the generated root variable %s should self-default to its environment value", variable)
	}

	callVars := &ast.Vars{}
	callVars.Set("APP_TAGS", ast.Var{Value: "cli-tags"})
	callVars.Set("APP_CGO_ENABLED", ast.Var{Value: "0"})
	call := &ast.Call{Task: "build", Vars: callVars}
	taskDefinition, err := executor.GetTask(call)
	require.NoError(t, err)
	resolved, err = executor.Compiler.GetVariables(taskDefinition, call)
	require.NoError(t, err)
	resolvedValues = resolved.ToCacheMap()
	assert.Equal(t, "cli-tags", resolvedValues["APP_TAGS"])
	assert.Equal(t, "0", resolvedValues["APP_CGO_ENABLED"])
}

func TestCommonTaskfileAppVariablesDefaultToEmpty(t *testing.T) {
	content := renderCommonTaskfile(t)
	for _, variable := range []string{
		"APP_TAGS", "APP_TAGS_LINUX", "APP_TAGS_DARWIN", "APP_TAGS_WINDOWS",
		"APP_TAGS_ANDROID", "APP_TAGS_IOS", "APP_TAGS_SERVER", "APP_LDFLAGS", "APP_CGO_ENABLED",
	} {
		assert.Contains(t, content, variable+`: '{{.`+variable+` | default ""}}'`)
		assert.NotContains(t, content, variable+": ''")
	}
}

func renderCommonTaskfile(t *testing.T) string {
	t.Helper()
	data, err := templates.ReadFile("_common/Taskfile.tmpl.yml")
	require.NoError(t, err)

	tmpl, err := texttemplate.New("Taskfile.tmpl.yml").Parse(string(data))
	require.NoError(t, err)

	var rendered bytes.Buffer
	require.NoError(t, tmpl.Execute(&rendered, map[string]any{
		"BinaryName": "test-app",
		"Opn":        "{{",
		"Cls":        "}}",
	}))
	return rendered.String()
}

func writeStubIncludedTaskfiles(t *testing.T, projectDir string) {
	t.Helper()
	const stub = `version: '3'

tasks:
  build:
    cmds:
      - echo build
  package:
    cmds:
      - echo package
  run:
    cmds:
      - echo run
`
	for _, relativePath := range []string{
		"build/Taskfile.yml",
		"build/windows/Taskfile.yml",
		"build/darwin/Taskfile.yml",
		"build/linux/Taskfile.yml",
		"build/ios/Taskfile.yml",
		"build/android/Taskfile.yml",
	} {
		path := filepath.Join(projectDir, relativePath)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(stub), 0o644))
	}
}
