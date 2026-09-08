package commands

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	// Copies of this test executable stand in for tools on Windows, where shell
	// scripts cannot be launched directly through exec.Command.
	if runtime.GOOS == "windows" && os.Getenv("WAILS_TEST_PLAN_TOOL") == "1" {
		os.Exit(0)
	}
}

func prependFakePlanTools(t *testing.T, names ...string) {
	t.Helper()
	directory := t.TempDir()
	names = append(names, "docker")
	contents := []byte("#!/bin/sh\nexit 0\n")
	suffix := ""
	if runtime.GOOS == "windows" {
		executable, err := os.Executable()
		require.NoError(t, err)
		contents, err = os.ReadFile(executable)
		require.NoError(t, err)
		suffix = ".exe"
		t.Setenv("WAILS_TEST_PLAN_TOOL", "1")
	}
	for _, name := range names {
		require.NoError(t, os.WriteFile(filepath.Join(directory, name+suffix), contents, 0o755))
	}
	t.Setenv("PATH", directory+string(os.PathListSeparator)+os.Getenv("PATH"))
}
