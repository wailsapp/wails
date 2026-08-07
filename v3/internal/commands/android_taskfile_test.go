package commands

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestAndroidTaskfileDeployDevice(t *testing.T) {
	// Given
	deployTask := androidTaskYAML(t, "deploy-device")

	// Then
	wantSnippets := []string{
		"DEVICE_ID",
		`DEVICE="${DEVICE_ID:-}"`,
		`"{{.ADB}}" devices`,
		`!~ /^emulator-/`,
		`- task: package`,
		`ARCH: arm64`,
		`"{{.ADB}}" -s "$DEVICE" uninstall {{.APP_ID}}`,
		`"{{.ADB}}" -s "$DEVICE" install "{{.BIN_DIR}}/{{.APP_NAME}}.apk"`,
		`"{{.ADB}}" -s "$DEVICE" shell am start -n {{.APP_ID}}/com.wails.app.MainActivity`,
	}
	for _, want := range wantSnippets {
		assert.Contains(t, deployTask, want)
	}
	assert.NotContains(t, deployTask, "ensure-emulator")
}

func TestAndroidTaskfileRunDevice(t *testing.T) {
	// Given
	runDeviceTask := androidTaskYAML(t, "run:device")

	// Then
	wantSnippets := []string{
		"DEVICE_ID",
		`DEVICE="${DEVICE_ID:-}"`,
		`"{{.ADB}}" devices`,
		`!~ /^emulator-/`,
		`- task: build`,
		`ARCH: arm64`,
		`"{{.ADB}}" -s "$DEVICE" uninstall {{.APP_ID}}`,
		`"{{.ADB}}" -s "$DEVICE" install "{{.BIN_DIR}}/{{.APP_NAME}}.apk"`,
		`"{{.ADB}}" -s "$DEVICE" shell am start -n {{.APP_ID}}/com.wails.app.MainActivity`,
	}
	for _, want := range wantSnippets {
		assert.Contains(t, runDeviceTask, want)
	}
	assert.NotContains(t, runDeviceTask, "ensure-emulator")
}

func TestAndroidTaskfileBuildUsesAndroidBindingContext(t *testing.T) {
	// Given
	buildTask := androidTaskYAML(t, "build")
	var task struct {
		Deps []struct {
			Task string         `yaml:"task"`
			Vars map[string]any `yaml:"vars"`
		} `yaml:"deps"`
	}
	require.NoError(t, yaml.Unmarshal([]byte(buildTask), &task))

	// Then
	assert.Contains(t, buildTask, "task: common:build:frontend")
	assert.Contains(t, buildTask, "GOOS: android")
	assert.NotContains(t, buildTask, "generate:android:bindings")

	for _, dependency := range task.Deps {
		if dependency.Task == "common:build:frontend" {
			require.Equal(t, "android", dependency.Vars["GOOS"])
			require.Equal(t, "0", dependency.Vars["CGO_ENABLED"])
			return
		}
	}
	t.Fatal("Android build should depend on common:build:frontend")
}

func TestCommonTaskfileForwardsBindingContext(t *testing.T) {
	// Given
	data, err := buildAssets.ReadFile("build_assets/Taskfile.tmpl.yml")
	require.NoError(t, err)
	taskfile := string(data)

	// Then
	assert.Contains(t, taskfile, "ref: .GOOS")
	assert.Contains(t, taskfile, "ref: .CGO_ENABLED")
	assert.Contains(t, taskfile, `GOOS: '{{.Opn}}.GOOS | default OS{{.Cls}}'`)
	assert.Contains(t, taskfile, `CGO_ENABLED: '{{.Opn}}.CGO_ENABLED | default ""{{.Cls}}'`)
	assert.Contains(t, taskfile, `label: build:frontend (DEV={{.Opn}}.DEV{{.Cls}} RUNNER={{.Opn}}.PACKAGE_MANAGER{{.Cls}} GOOS={{.Opn}}.GOOS | default OS{{.Cls}} CGO_ENABLED={{.Opn}}.CGO_ENABLED | default ""{{.Cls}})`)
	assert.Contains(t, taskfile, `label: generate:bindings (GOOS={{.Opn}}.GOOS | default OS{{.Cls}} CGO_ENABLED={{.Opn}}.CGO_ENABLED | default ""{{.Cls}})`)
}

func TestAndroidExamplesUseCommonBindingContext(t *testing.T) {
	for _, path := range []string{
		"../../examples/android/build/android/Taskfile.yml",
		"../../examples/mobile/build/android/Taskfile.yml",
	} {
		data, err := os.ReadFile(path)
		require.NoError(t, err)

		var taskfile struct {
			Tasks map[string]yaml.Node `yaml:"tasks"`
		}
		require.NoError(t, yaml.Unmarshal(data, &taskfile), path)
		buildNode, ok := taskfile.Tasks["build"]
		require.True(t, ok, "%s should include a build task", path)
		var buildTask struct {
			Deps []struct {
				Task string         `yaml:"task"`
				Vars map[string]any `yaml:"vars"`
			} `yaml:"deps"`
		}
		require.NoError(t, buildNode.Decode(&buildTask), path)
		assert.NotContains(t, string(data), "generate:android:bindings", path)

		foundFrontendBuild := false
		for _, dependency := range buildTask.Deps {
			if dependency.Task != "common:build:frontend" {
				continue
			}
			foundFrontendBuild = true
			assert.Equal(t, "android", dependency.Vars["GOOS"], path)
			assert.Equal(t, "0", dependency.Vars["CGO_ENABLED"], path)
		}
		assert.True(t, foundFrontendBuild, "%s build task should depend on common:build:frontend", path)
	}

	for _, path := range []string{
		"../../examples/android/build/Taskfile.yml",
		"../../examples/mobile/build/Taskfile.yml",
	} {
		data, err := os.ReadFile(path)
		require.NoError(t, err)
		taskfile := string(data)

		assert.Contains(t, taskfile, `GOOS={{.GOOS | default OS}}`, path)
		assert.Contains(t, taskfile, `CGO_ENABLED={{.CGO_ENABLED | default ""}}`, path)
	}
}

func androidTaskYAML(t *testing.T, name string) string {
	t.Helper()

	data, err := buildAssets.ReadFile("build_assets/android/Taskfile.yml")
	require.NoError(t, err)

	// When
	var taskfile struct {
		Tasks map[string]yaml.Node `yaml:"tasks"`
	}
	require.NoError(t, yaml.Unmarshal(data, &taskfile))
	taskNode, ok := taskfile.Tasks[name]
	require.True(t, ok, "android Taskfile should include %s", name)
	taskData, err := yaml.Marshal(&taskNode)
	require.NoError(t, err)
	return string(taskData)
}
