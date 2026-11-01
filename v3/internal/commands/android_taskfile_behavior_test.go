package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const androidPhone = "List of devices attached\nphone-123\tdevice\n"
const androidMixedDevices = androidPhone + "emulator-5556\toffline\nemulator-5554\tdevice\nemulator-5558\tdevice\n"

// Execute the actual task commands with a fake adb, so regressions in device
// selection fail without installing an SDK or touching a developer's devices.
func TestAndroidTaskfileEmulatorSelection(t *testing.T) {
	for _, source := range androidTaskfileSources(t) {
		t.Run(source.name, func(t *testing.T) {
			t.Parallel()
			for _, tt := range []struct {
				name       string
				devices    string
				later      string
				bootZeros  int
				wantSerial string
				wantError  bool
			}{
				{"phone cannot satisfy boot wait", androidPhone, "", 0, "", true},
				{"no emulator appears", "List of devices attached\n", "", 0, "", true},
				{"offline emulator times out", "List of devices attached\nemulator-5554\toffline\n", "", 200, "emulator-5554", true},
				{"offline adb errors are retried", "List of devices attached\nemulator-5554\toffline\n", "", -1, "emulator-5554", true},
				{"running emulator still needs boot completion", androidMixedDevices, "", 2, "emulator-5554", false},
				{"phone and multiple emulators", androidMixedDevices, "", 0, "emulator-5554", false},
				{"new emulator registers after phone", androidPhone, androidPhone + "emulator-5554\tdevice\n", 1, "emulator-5554", false},
			} {
				t.Run(tt.name, func(t *testing.T) {
					output, trace, err := runAndroidTaskCommands(t, source.data, "ensure-emulator", tt.devices, tt.later, tt.bootZeros)
					if tt.wantError {
						require.Error(t, err, output)
						require.Contains(t, output, "timeout")
						require.Equal(t, 120, strings.Count(trace, "sleep 1\n"), "registration and boot share one bounded polling loop")
					} else {
						require.NoError(t, err, output)
						require.Contains(t, output, "booted successfully")
					}
					require.NotContains(t, trace, "wait-for-device")
					require.NotContains(t, trace, "adb shell ")
					if tt.devices == androidPhone || tt.devices == "List of devices attached\n" {
						require.Contains(t, trace, "emulator launched", "a new emulator must outlive the task shell")
					} else {
						require.NotContains(t, trace, "emulator -list-avds", "reuse an existing emulator without requiring a local AVD")
					}
					if tt.wantSerial != "" {
						bootCall := "adb -s " + tt.wantSerial + " shell getprop sys.boot_completed"
						require.Contains(t, trace, bootCall)
						if !tt.wantError {
							require.Equal(t, tt.bootZeros+1, strings.Count(trace, bootCall))
						}
					} else {
						require.NotContains(t, trace, "shell getprop")
					}
				})
			}
		})
	}
}

func TestAndroidTaskfileScopedCommands(t *testing.T) {
	for _, source := range androidTaskfileSources(t) {
		t.Run(source.name, func(t *testing.T) {
			for _, task := range []string{"run", "deploy-emulator", "logs", "logs:all"} {
				t.Run(task, func(t *testing.T) {
					for _, tt := range []struct {
						name, devices, wantSerial string
					}{
						{"prefer online emulator", androidMixedDevices, "emulator-5554"},
						{"phone only", androidPhone, "phone-123"},
						{"multiple phones", androidPhone + "phone-456\tdevice\n", "phone-123"},
						{"offline emulator and online phone", androidPhone + "emulator-5554\toffline\n", "phone-123"},
						{"skip unavailable devices", "List of devices attached\nphone-offline\toffline\nphone-unauthorized\tunauthorized\nemulator-5554\toffline\nemulator-5556\tunauthorized\n", ""},
					} {
						t.Run(tt.name, func(t *testing.T) {
							output, trace, err := runAndroidTaskCommands(t, source.data, task, tt.devices, "", 0)
							wantSerial := tt.wantSerial
							if !strings.HasPrefix(task, "logs") && wantSerial == "phone-123" {
								wantSerial = ""
							}
							if wantSerial == "" {
								require.Error(t, err, output)
								require.Contains(t, output, "Error:")
								require.NotContains(t, trace, "uninstall")
								require.NotContains(t, trace, "logcat")
								return
							}
							require.NoError(t, err, output)
							if strings.HasPrefix(task, "logs") {
								require.Contains(t, trace, "adb -s "+wantSerial+" logcat -v time")
							} else {
								for _, call := range []string{"uninstall com.wails.test", "install bin/test.apk", "shell am start -n com.wails.test/com.wails.app.MainActivity"} {
									require.Contains(t, trace, "adb -s "+wantSerial+" "+call)
								}
							}
							for _, line := range strings.Split(trace, "\n") {
								if strings.HasPrefix(line, "adb ") && line != "adb devices" {
									require.True(t, strings.HasPrefix(line, "adb -s "+wantSerial+" "), line)
								}
							}
						})
					}
				})
			}
		})
	}
}

func TestAndroidTaskfileInstallFailure(t *testing.T) {
	for _, source := range androidTaskfileSources(t) {
		for _, task := range []string{"run", "deploy-emulator"} {
			t.Run(source.name+"/"+task, func(t *testing.T) {
				output, trace, err := runAndroidTaskCommands(t, source.data, task, androidMixedDevices, "", 0, "FAKE_INSTALL_FAIL=1")
				require.Error(t, err, output)
				require.Contains(t, trace, "adb -s emulator-5554 install bin/test.apk")
				require.NotContains(t, trace, "shell am start", "a failed install must not launch an old APK")
			})
		}
	}
}

type androidTaskfileSource struct {
	name string
	data []byte
}

func androidTaskfileSources(t *testing.T) []androidTaskfileSource {
	t.Helper()
	data, err := buildAssets.ReadFile("build_assets/android/Taskfile.yml")
	require.NoError(t, err)
	sources := []androidTaskfileSource{{"template", data}}
	for _, example := range []string{"android", "mobile"} {
		data, err := os.ReadFile(filepath.Join("..", "..", "examples", example, "build", "android", "Taskfile.yml"))
		require.NoError(t, err)
		sources = append(sources, androidTaskfileSource{example, data})
	}
	return sources
}

func runAndroidTaskCommands(t *testing.T, data []byte, task, devices, later string, bootZeros int, extraEnv ...string) (string, string, error) {
	t.Helper()
	shell, err := exec.LookPath("bash")
	if err != nil && runtime.GOOS == "windows" {
		// Git for Windows bundles the POSIX shell used by these task scripts.
		if git, gitErr := exec.LookPath("git"); gitErr == nil {
			shell = filepath.Join(filepath.Dir(git), "..", "bin", "bash.exe")
			_, err = os.Stat(shell)
		}
	}
	if err != nil {
		t.Skip("bash is required to execute Android task scripts")
	}
	var taskfile struct {
		Tasks map[string]struct {
			Cmds []yaml.Node `yaml:"cmds"`
		} `yaml:"tasks"`
	}
	require.NoError(t, yaml.Unmarshal(data, &taskfile))
	definition, ok := taskfile.Tasks[task]
	require.True(t, ok)
	var commands []string
	for _, cmd := range definition.Cmds {
		if cmd.Kind == yaml.ScalarNode {
			commands = append(commands, cmd.Value)
		}
	}
	require.NotEmpty(t, commands)
	tmpl, err := template.New(task).Option("missingkey=error").Parse(strings.Join(commands, "\n"))
	require.NoError(t, err, "task commands must be valid Go templates")
	var rendered bytes.Buffer
	require.NoError(t, tmpl.Execute(&rendered, map[string]string{
		"ADB": "fake_adb", "EMULATOR": "fake_emulator", "APP_ID": "com.wails.test",
		"BIN_DIR": "bin", "APP_NAME": "test", "SDK_ROOT": "/unused",
		"AVDMANAGER": "/unused", "ANDROID_ABI": "x86_64", "HOST_ARCH": "amd64",
	}))
	script := androidFakeCommands + "\n" + rendered.String()
	dir := t.TempDir()
	tracePath := filepath.Join(dir, "trace")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, shell, "-c", script)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "FAKE_DEVICES="+devices, "FAKE_LATER="+later, fmt.Sprintf("FAKE_BOOT_ZEROS=%d", bootZeros))
	cmd.Env = append(cmd.Env, extraEnv...)
	output, runErr := cmd.CombinedOutput()
	require.NoError(t, ctx.Err(), "task commands exceeded test deadline: %s", output)
	trace, err := os.ReadFile(tracePath)
	require.NoError(t, err)
	return string(output), string(trace), runErr
}

const androidFakeCommands = `set -e
fake_adb() {
    printf 'adb %s\n' "$*" >> trace
    case "$1" in
        devices)
            count=0
            if [ -f lists ]; then read -r count < lists; fi
            count=$((count + 1))
            echo "$count" > lists
            if [ -n "$FAKE_LATER" ] && [ "$count" -ge 3 ]; then
                printf '%s' "$FAKE_LATER"
            else
                printf '%s' "$FAKE_DEVICES"
            fi
            ;;
        wait-for-device) return 0 ;;
        shell) echo 1 ;;
        -s)
            shift 2
            if [ "$*" = 'shell getprop sys.boot_completed' ]; then
                if [ "$FAKE_BOOT_ZEROS" -lt 0 ]; then return 1; fi
                count=0
                if [ -f boots ]; then read -r count < boots; fi
                count=$((count + 1))
                echo "$count" > boots
                if [ "$count" -le "$FAKE_BOOT_ZEROS" ]; then echo 0; else echo 1; fi
            elif [ "$1" = 'logcat' ]; then
                echo 'Wails test output'
            elif [ "$1" = 'uninstall' ]; then
                return 1
            elif [ "$1" = 'install' ] && [ "$FAKE_INSTALL_FAIL" = '1' ]; then
                return 1
            fi
            ;;
        *) return 1 ;;
    esac
}
fake_emulator() {
    printf 'emulator %s\n' "$*" >> trace
    if [ "$1" = '-list-avds' ]; then echo test-avd; fi
}
nohup() { echo 'emulator launched' >> trace; }
sleep() { printf 'sleep %s\n' "$*" >> trace; }
`
