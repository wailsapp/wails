package manifest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHookHCLRoundTripAndValidation(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "scripts"), 0o755))
	script := testHookScriptName("prepare")
	require.NoError(t, os.WriteFile(filepath.Join(root, "scripts", script), testHookScriptContents(), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "version.txt"), []byte("1\n"), 0o644))
	source := fmt.Sprintf(`version = 3
project {
  name = "hooks"
  product_name = "Hooks"
  identifier = "com.example.hooks"
  version = "1.0.0"
}
hook "before_build" {
  script = "scripts/%s"
  directory = "scripts"
  cache = true
  inputs = ["version.txt"]
  outputs = ["generated/version.go"]
}
`, script)
	path := filepath.Join(root, Filename)
	require.NoError(t, os.WriteFile(path, []byte(source), 0o644))
	loaded, err := Load(root, "")
	require.NoError(t, err)
	assert.Equal(t, Hook{Script: "scripts/" + script, Directory: "scripts", Cache: true, Inputs: []string{"version.txt"}, Outputs: []string{"generated/version.go"}}, loaded.Config.Hooks[BeforeBuild])

	encoded, err := EncodeConfig(loaded.Config)
	require.NoError(t, err)
	assert.Contains(t, string(encoded), `hook "before_build" {`)
	second := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(second, "scripts"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(second, "scripts", script), testHookScriptContents(), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(second, "version.txt"), []byte("1\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(second, Filename), encoded, 0o644))
	roundTrip, err := Load(second, "")
	require.NoError(t, err)
	assert.Equal(t, loaded.Config.Hooks, roundTrip.Config.Hooks)
}

func TestHookValidationRejectsUnsafeAndIncompleteContracts(t *testing.T) {
	root := t.TempDir()
	script := testHookScriptName("hook")
	require.NoError(t, os.WriteFile(filepath.Join(root, script), testHookScriptContents(), 0o755))
	base := NewDocument(Project{Name: "hooks", ProductName: "Hooks", Identifier: "com.example.hooks", Version: "1.0.0"})
	tests := []struct {
		name string
		hook Hook
		want string
	}{
		{"unsupported phase", Hook{Script: script}, "unsupported hook phase"},
		{"cache needs inputs", Hook{Script: script, Cache: true, Outputs: []string{"generated/out"}}, "requires complete inputs and outputs"},
		{"declarations need cache", Hook{Script: script, Inputs: []string{"input"}}, "must be true"},
		{"multiple roots", Hook{Script: script, Cache: true, Inputs: []string{"input"}, Outputs: []string{"one/out", "two/out"}}, "non-root directory"},
		{"output contains script", Hook{Script: filepath.ToSlash(filepath.Join("generated", script)), Cache: true, Inputs: []string{"input"}, Outputs: []string{"generated/a", "generated/b"}}, "contains script or input"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			doc := base
			doc.Hooks = map[HookPhase]Hook{BeforeBuild: test.hook}
			if test.name == "unsupported phase" {
				doc.Hooks = map[HookPhase]Hook{"unknown": test.hook}
			}
			if test.name == "output contains script" {
				require.NoError(t, os.MkdirAll(filepath.Join(root, "generated"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(root, "generated", script), testHookScriptContents(), 0o755))
			}
			config := configFromDocument(root, "", doc)
			err := validateConfig(config)
			require.ErrorContains(t, err, test.want)
		})
	}
}

func TestHookOutputRoot(t *testing.T) {
	root, err := HookOutputRoot([]string{"generated/a", "generated/nested/b"})
	require.NoError(t, err)
	assert.Equal(t, "generated", root)
	root, err = HookOutputRoot([]string{"generated/result"})
	require.NoError(t, err)
	assert.Equal(t, "generated/result", root)
	root, err = HookOutputRoot([]string{"generated", "generated/nested/result"})
	require.NoError(t, err)
	assert.Equal(t, "generated", root)
	_, err = HookOutputRoot([]string{"."})
	assert.ErrorContains(t, err, "project root")
}

func TestWindowsHookScriptExtensionsAreClosed(t *testing.T) {
	for _, value := range []string{"hook.cmd", "hook.BAT", "hook.ps1"} {
		assert.True(t, validWindowsHookScript(value), value)
	}
	for _, value := range []string{"hook", "hook.sh", "hook.exe", "hook.ps1.txt"} {
		assert.False(t, validWindowsHookScript(value), value)
	}
}

func TestHookValidationCarriesTheHCLSourceRange(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows validates hook script extensions rather than Unix executable permissions")
	}
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "hook.sh"), []byte("#!/bin/sh\n"), 0o644))
	source := `version = 3
project {
  name = "hooks"
  product_name = "Hooks"
  identifier = "com.example.hooks"
  version = "1.0.0"
}
hook "before_build" {
  script = "hook.sh"
}
`
	path := filepath.Join(root, Filename)
	require.NoError(t, os.WriteFile(path, []byte(source), 0o644))
	_, err := Load(root, "")
	require.Error(t, err)
	var validation *ValidationError
	require.True(t, errors.As(err, &validation))
	assert.Equal(t, `hook["before_build"].script`, validation.Field)
	assert.Equal(t, 9, validation.Range.StartLine)
	assert.ErrorContains(t, err, "must be executable")
}

func testHookScriptName(base string) string {
	if runtime.GOOS == "windows" {
		return base + ".cmd"
	}
	return base + ".sh"
}

func testHookScriptContents() []byte {
	if runtime.GOOS == "windows" {
		return []byte("@echo off\r\n")
	}
	return []byte("#!/bin/sh\n")
}

func TestHookValidationRejectsSymlinkEscapes(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating symlinks requires elevated privileges on some Windows hosts")
	}
	root := t.TempDir()
	outside := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(outside, "hook.sh"), []byte("#!/bin/sh\n"), 0o755))
	require.NoError(t, os.Symlink(filepath.Join(outside, "hook.sh"), filepath.Join(root, "hook.sh")))
	doc := NewDocument(Project{Name: "hooks", ProductName: "Hooks", Identifier: "com.example.hooks", Version: "1.0.0"})
	doc.Hooks = map[HookPhase]Hook{BeforeBuild: {Script: "hook.sh"}}
	err := validateConfig(configFromDocument(root, "", doc))
	require.ErrorContains(t, err, "resolves outside the project")

	require.NoError(t, os.Remove(filepath.Join(root, "hook.sh")))
	require.NoError(t, os.WriteFile(filepath.Join(root, "hook.sh"), []byte("#!/bin/sh\n"), 0o755))
	require.NoError(t, os.Symlink(outside, filepath.Join(root, "generated")))
	doc.Hooks = map[HookPhase]Hook{BeforeBuild: {Script: "hook.sh", Cache: true, Inputs: []string{"input.txt"}, Outputs: []string{"generated/out.txt"}}}
	err = validateConfig(configFromDocument(root, "", doc))
	require.ErrorContains(t, err, "resolves outside the project")
}
