package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	taskruntime "github.com/wailsapp/task/v3"
	"github.com/wailsapp/task/v3/taskfile/ast"
	"github.com/wailsapp/wails/v3/internal/flags"
)

func TestBuildCustomizationFlags(t *testing.T) {
	fixture := newBuildCustomizationFixture(t)

	for _, platform := range []struct {
		name         string
		task         string
		platformVar  string
		productionLD string
	}{
		{"windows", "windows:build:native", "APP_TAGS_WINDOWS", "-w -s -H windowsgui"},
		{"linux", "linux:build:native", "APP_TAGS_LINUX", "-w -s"},
		{"darwin", "darwin:build:native", "APP_TAGS_DARWIN", "-w -s"},
	} {
		t.Run(platform.name, func(t *testing.T) {
			legacyDev := fixture.buildFlags(t, platform.task, map[string]any{"DEV": "true"})
			assert.Empty(t, optionValue(legacyDev, "-tags"), "empty customization must preserve the legacy dev build")
			assert.Empty(t, assignmentValue(legacyDev, "-ldflags"))

			custom := map[string]any{
				"DEV":                "true",
				"OBFUSCATED":         "true",
				"APP_TAGS":           "global-a,global-b",
				platform.platformVar: "platform",
				"EXTRA_TAGS":         "extra",
				"APP_LDFLAGS":        "-X main.Version=1.2.3",
			}
			dev := fixture.buildFlags(t, platform.task, custom)
			assert.Equal(t, "wails_obfuscated,global-a,global-b,platform,extra", optionValue(dev, "-tags"))
			assert.Equal(t, "-X main.Version=1.2.3", assignmentValue(dev, "-ldflags"))

			custom["DEV"] = "false"
			production := fixture.buildFlags(t, platform.task, custom)
			assert.Equal(t, "production,wails_obfuscated,global-a,global-b,platform,extra", optionValue(production, "-tags"))
			assert.Equal(t, platform.productionLD+" -X main.Version=1.2.3", assignmentValue(production, "-ldflags"))
		})
	}

	for _, platform := range []struct {
		name        string
		task        string
		platformVar string
		debugTags   string
		prodTags    string
	}{
		{"android", "android:build", "APP_TAGS_ANDROID", "android,debug", "production,android"},
		{"ios", "ios:build", "APP_TAGS_IOS", "ios,debug", "production,ios"},
	} {
		t.Run(platform.name, func(t *testing.T) {
			assert.Equal(t, platform.debugTags, optionValue(fixture.buildFlags(t, platform.task, nil), "-tags"))
			production := fixture.buildFlags(t, platform.task, map[string]any{
				"PRODUCTION":         "true",
				"APP_TAGS":           "global-a,global-b",
				platform.platformVar: "platform",
				"EXTRA_TAGS":         "extra",
				"APP_LDFLAGS":        "-X main.Version=1.2.3",
			})
			assert.Equal(t, platform.prodTags+",global-a,global-b,platform,extra", optionValue(production, "-tags"))
			assert.Equal(t, "-w -s -X main.Version=1.2.3", assignmentValue(production, "-ldflags"))
		})
	}

	extraOnly := fixture.buildFlags(t, "common:build:server", map[string]any{"DEV": "true", "EXTRA_TAGS": "extra"})
	assert.Equal(t, "server,extra", optionValue(extraOnly, "-tags"))
	serverProduction := fixture.buildFlags(t, "common:build:server", map[string]any{
		"DEV": "false", "OBFUSCATED": "true", "APP_TAGS": "global", "APP_TAGS_SERVER": "server-only",
		"EXTRA_TAGS": "extra", "APP_LDFLAGS": "-X main.Version=1.2.3",
	})
	assert.Equal(t, "server,production,wails_obfuscated,global,server-only,extra", optionValue(serverProduction, "-tags"))
	assert.Equal(t, "-w -s -X main.Version=1.2.3", assignmentValue(serverProduction, "-ldflags"))
}

func TestCGOBuildCustomization(t *testing.T) {
	fixture := newBuildCustomizationFixture(t)
	for _, target := range []struct {
		name       string
		task       string
		defaultCGO string
	}{
		{"windows", "windows:build:native", "0"},
		{"windows docker", "windows:build:docker", "0"},
		{"linux", "linux:build:native", "1"},
		{"darwin", "darwin:build:native", "1"},
		{"server native", "common:build:server", ""},
		{"server docker", "common:build:docker", "0"},
	} {
		t.Run(target.name, func(t *testing.T) {
			for _, test := range []struct {
				name string
				vars map[string]any
				want string
			}{
				{"missing", nil, target.defaultCGO},
				{"numeric zero", map[string]any{"APP_CGO_ENABLED": 0, "CGO_ENABLED": "1"}, "0"},
				{"string zero", map[string]any{"APP_CGO_ENABLED": "0", "CGO_ENABLED": "1"}, "0"},
				{"numeric one", map[string]any{"APP_CGO_ENABLED": 1, "CGO_ENABLED": "0"}, "1"},
				{"string one", map[string]any{"APP_CGO_ENABLED": "1", "CGO_ENABLED": "0"}, "1"},
			} {
				t.Run(test.name, func(t *testing.T) {
					resolved := fixture.taskVars(t, target.task, test.vars)
					assert.Equal(t, test.want, fmt.Sprint(resolved["EFFECTIVE_CGO_ENABLED"]))
				})
			}
		})
	}

	for _, target := range []struct {
		task string
		want string
	}{
		{"windows:build:native", "CGO_ENABLED=0"},
		{"linux:build:native", "CGO_ENABLED=1"},
		{"darwin:build:native", "CGO_ENABLED=1"},
	} {
		assert.Contains(t, fixture.compiledCommand(t, target.task, nil, "go build"), target.want)
	}
	assert.NotContains(t, fixture.compiledCommand(t, "common:build:server", nil, "go build"), "CGO_ENABLED=")
	assert.Contains(t, fixture.compiledCommand(t, "common:build:server", map[string]any{"APP_CGO_ENABLED": 0}, "go build"), "CGO_ENABLED=0")
	assert.Contains(t, fixture.compiledCommand(t, "ios:build", nil, "go build"), "CGO_ENABLED=1 go build")

	for _, test := range []struct {
		value any
		want  string
	}{
		{0, "windows:build:native"},
		{1, "windows:build:docker"},
	} {
		compiled := fixture.compiledTask(t, "windows:build", map[string]any{"APP_CGO_ENABLED": test.value})
		if runtime.GOOS == "windows" {
			assert.Equal(t, "windows:build:native", compiled.Cmds[0].Task)
		} else {
			assert.Equal(t, test.want, compiled.Cmds[0].Task)
		}
	}
}

func TestDockerBuildCustomizationPropagation(t *testing.T) {
	for _, platform := range []struct {
		path        string
		platformVar string
	}{
		{"build_assets/windows/Taskfile.yml", "APP_TAGS_WINDOWS"},
		{"build_assets/linux/Taskfile.yml", "APP_TAGS_LINUX"},
		{"build_assets/darwin/Taskfile.yml", "APP_TAGS_DARWIN"},
	} {
		data, err := buildAssets.ReadFile(platform.path)
		require.NoError(t, err)
		content := strings.ReplaceAll(string(data), "\r\n", "\n")
		for _, expected := range []string{
			"BUILD_FLAGS:\n            ref: .BUILD_FLAGS",
			`-e APP_TAGS={{.APP_TAGS | default "" | q}}`,
			`-e APP_PLATFORM_TAGS={{.` + platform.platformVar + ` | default "" | q}}`,
			`-e APP_LDFLAGS={{.APP_LDFLAGS | default "" | q}}`,
			`-e EXTRA_TAGS={{.EXTRA_TAGS | default "" | q}}`,
			`-e CGO_ENABLED={{.EFFECTIVE_CGO_ENABLED | toString | q}}`,
		} {
			assert.Contains(t, content, expected)
		}
	}

	cross, err := buildAssets.ReadFile("build_assets/docker/Dockerfile.cross")
	require.NoError(t, err)
	for _, expected := range []string{
		`export CGO_ENABLED="${CGO_ENABLED:-1}"`, `if [ -n "$APP_TAGS" ]`,
		`if [ -n "$APP_PLATFORM_TAGS" ]`, `if [ -n "$APP_LDFLAGS" ]`,
	} {
		assert.Contains(t, string(cross), expected)
	}

	fixture := newBuildCustomizationFixture(t)
	common, err := os.ReadFile(filepath.Join(fixture.projectDir, "build", "Taskfile.yml"))
	require.NoError(t, err)
	commonContent := strings.ReplaceAll(string(common), "\r\n", "\n")
	for _, expected := range []string{
		"BUILD_FLAGS:\n            ref: .BUILD_FLAGS",
		`--build-arg CGO_ENABLED={{.EFFECTIVE_CGO_ENABLED | toString | q}}`,
		`--build-arg APP_TAGS={{.APP_TAGS | default "" | q}}`,
		`--build-arg APP_PLATFORM_TAGS={{.APP_TAGS_SERVER | default "" | q}}`,
		`--build-arg APP_LDFLAGS={{.APP_LDFLAGS | default "" | q}}`,
		`--build-arg EXTRA_TAGS={{.EXTRA_TAGS | default "" | q}}`,
	} {
		assert.Contains(t, commonContent, expected)
	}

	server, err := buildAssets.ReadFile("build_assets/docker/Dockerfile.server")
	require.NoError(t, err)
	for _, expected := range []string{
		"ARG APP_TAGS", "ARG APP_PLATFORM_TAGS", "ARG APP_LDFLAGS", "ARG EXTRA_TAGS",
		`-tags "$TAGS"`, `-ldflags="$LDFLAGS"`,
	} {
		assert.Contains(t, string(server), expected)
	}

	bindings := fixture.rawTask(t, "common:generate:bindings", nil)
	assert.NotEmpty(t, bindings.Sources)
	assert.NotEmpty(t, bindings.Generates)
	assert.Contains(t, bindings.Cmds[0].Cmd, `-f {{.BUILD_FLAGS | default "" | q}}`)
	assert.Contains(t, fixture.compiledCommand(t, "common:generate:bindings", nil, "generate bindings"), "generate bindings -f ''")
}

type buildCustomizationFixture struct {
	executor   *taskruntime.Executor
	projectDir string
}

func newBuildCustomizationFixture(t *testing.T) *buildCustomizationFixture {
	t.Helper()
	for _, variable := range []string{
		"APP_TAGS", "APP_TAGS_LINUX", "APP_TAGS_DARWIN", "APP_TAGS_WINDOWS", "APP_TAGS_ANDROID",
		"APP_TAGS_IOS", "APP_TAGS_SERVER", "APP_LDFLAGS", "APP_CGO_ENABLED", "CGO_ENABLED",
	} {
		t.Setenv(variable, "")
	}

	projectDir := t.TempDir()
	require.NoError(t, GenerateBuildAssets(&BuildAssetsOptions{
		Dir: filepath.Join(projectDir, "build"), Name: "Test App", BinaryName: "testapp", Silent: true,
	}))
	const rootTaskfile = `version: '3'
vars:
  APP_NAME: testapp
  BIN_DIR: bin
  APP_TAGS: '{{.APP_TAGS | default ""}}'
  APP_TAGS_LINUX: '{{.APP_TAGS_LINUX | default ""}}'
  APP_TAGS_DARWIN: '{{.APP_TAGS_DARWIN | default ""}}'
  APP_TAGS_WINDOWS: '{{.APP_TAGS_WINDOWS | default ""}}'
  APP_TAGS_ANDROID: '{{.APP_TAGS_ANDROID | default ""}}'
  APP_TAGS_IOS: '{{.APP_TAGS_IOS | default ""}}'
  APP_TAGS_SERVER: '{{.APP_TAGS_SERVER | default ""}}'
  APP_LDFLAGS: '{{.APP_LDFLAGS | default ""}}'
  APP_CGO_ENABLED: '{{.APP_CGO_ENABLED | default ""}}'
includes:
  common: ./build/Taskfile.yml
  windows: ./build/windows/Taskfile.yml
  linux: ./build/linux/Taskfile.yml
  darwin: ./build/darwin/Taskfile.yml
  android: ./build/android/Taskfile.yml
  ios: ./build/ios/Taskfile.yml
`
	require.NoError(t, os.WriteFile(filepath.Join(projectDir, "Taskfile.yml"), []byte(rootTaskfile), 0o644))
	executor := &taskruntime.Executor{
		Dir: projectDir, Offline: true, DisableVersionCheck: true,
		Stdin: strings.NewReader(""), Stdout: io.Discard, Stderr: io.Discard,
	}
	require.NoError(t, executor.Setup())
	return &buildCustomizationFixture{executor: executor, projectDir: projectDir}
}

func (fixture *buildCustomizationFixture) call(task string, values map[string]any) *ast.Call {
	vars := &ast.Vars{}
	for key, value := range values {
		vars.Set(key, ast.Var{Value: value})
	}
	return &ast.Call{Task: task, Vars: vars}
}

func (fixture *buildCustomizationFixture) rawTask(t *testing.T, task string, values map[string]any) *ast.Task {
	t.Helper()
	result, err := fixture.executor.GetTask(fixture.call(task, values))
	require.NoError(t, err)
	return result
}

func (fixture *buildCustomizationFixture) taskVars(t *testing.T, task string, values map[string]any) map[string]any {
	t.Helper()
	call := fixture.call(task, values)
	resolved, err := fixture.executor.Compiler.FastGetVariables(fixture.rawTask(t, task, values), call)
	require.NoError(t, err)
	return resolved.ToCacheMap()
}

func (fixture *buildCustomizationFixture) compiledTask(t *testing.T, task string, values map[string]any) *ast.Task {
	t.Helper()
	compiled, err := fixture.executor.FastCompiledTask(fixture.call(task, values))
	require.NoError(t, err)
	return compiled
}

func (fixture *buildCustomizationFixture) compiledCommand(t *testing.T, task string, values map[string]any, fragment string) string {
	t.Helper()
	for _, command := range fixture.compiledTask(t, task, values).Cmds {
		if strings.Contains(command.Cmd, fragment) {
			return command.Cmd
		}
	}
	require.FailNow(t, "task command not found", "no command in %s contained %q", task, fragment)
	return ""
}

func (fixture *buildCustomizationFixture) buildFlags(t *testing.T, task string, values map[string]any) []string {
	t.Helper()
	buildFlags := fmt.Sprint(fixture.taskVars(t, task, values)["BUILD_FLAGS"])
	options := flags.GenerateBindingsOptions{BuildFlagsString: buildFlags}
	fields, err := options.BuildFlags()
	require.NoError(t, err)
	return fields
}

func optionValue(fields []string, option string) string {
	for index, field := range fields {
		if field == option && index+1 < len(fields) {
			return fields[index+1]
		}
	}
	return ""
}

func assignmentValue(fields []string, name string) string {
	prefix := name + "="
	for _, field := range fields {
		if strings.HasPrefix(field, prefix) {
			return strings.TrimPrefix(field, prefix)
		}
	}
	return ""
}
