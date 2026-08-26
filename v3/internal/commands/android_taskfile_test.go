package commands

import (
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

func TestEmbeddedAndroidProjectTargetsAPI36(t *testing.T) {
	rootData, err := buildAssets.ReadFile("build_assets/android/build.gradle")
	require.NoError(t, err)
	rootGradle := string(rootData)
	assert.Contains(t, rootGradle, "version '9.0.1'")
	assert.NotContains(t, rootGradle, "version '8.7.3'")

	data, err := buildAssets.ReadFile("build_assets/android/app/build.gradle")
	require.NoError(t, err)
	gradle := string(data)
	assert.Contains(t, gradle, "compileSdk = 36")
	assert.Contains(t, gradle, "targetSdk = 36")
	assert.NotContains(t, gradle, "compileSdk 35")
	assert.NotContains(t, gradle, "targetSdk 35")
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
