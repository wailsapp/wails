package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExperimentalCLICompatibility(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "wails3")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	build := exec.Command("go", "build", "-o", binary, ".")
	output, err := build.CombinedOutput()
	require.NoError(t, err, string(output))
	env := []string{}
	for _, item := range os.Environ() {
		name, _, _ := strings.Cut(item, "=")
		if name != "WAILS_EXP_USE_WAKE" && name != "WAILS_USE_WAKE" {
			env = append(env, item)
		}
	}
	run := func(dir string, value *string, args ...string) (string, error) {
		cmd := exec.Command(binary, args...)
		cmd.Dir = dir
		cmd.Env = append([]string{}, env...)
		if value != nil {
			cmd.Env = append(cmd.Env, "WAILS_EXP_USE_WAKE="+*value)
		}
		output, err := cmd.CombinedOutput()
		return string(output), err
	}
	for _, value := range []*string{nil, new(string), ptr("false"), ptr("anything")} {
		name := "unset"
		if value != nil {
			name = "set-" + *value
		}
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			enabled := value != nil
			help, err := run(dir, value, "--help")
			require.NoError(t, err, help)
			for _, command := range []string{"migrate", "eject", "clean"} {
				if enabled {
					require.Contains(t, help, command)
				} else {
					require.NotContains(t, help, command)
				}
			}
			for _, command := range []string{"build", "dev", "package", "sign"} {
				help, err := run(dir, value, command, "--help")
				require.NoError(t, err, help)
				if enabled {
					require.Contains(t, help, "profile")
				} else {
					require.NotContains(t, help, "profile")
				}
			}
			require.NoError(t, os.WriteFile(filepath.Join(dir, "Taskfile.yml"), []byte("version: '3'\ntasks:\n  build:\n    cmds:\n      - echo legacy > marker.txt\n"), 0644))
			require.NoError(t, os.WriteFile(filepath.Join(dir, "wails.hcl"), []byte("invalid manifest"), 0644))
			output, err := run(dir, value, "build")
			if enabled {
				require.Error(t, err, output)
				require.NoFileExists(t, filepath.Join(dir, "marker.txt"))
			} else {
				require.NoError(t, err, output)
				require.FileExists(t, filepath.Join(dir, "marker.txt"))
				require.NotContains(t, output, "deprecated")
				output, err = run(dir, value, "migrate")
				require.NoError(t, err, output)
				require.Contains(t, output, "Available commands:")
				require.NoFileExists(t, filepath.Join(dir, "wails.migrated.hcl"))
			}
		})
	}
}

func ptr(value string) *string { return &value }
