//go:build windows

package commands

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/flags"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	if os.Getenv("HCL_TEST_SIGNER_FAILURE") == "1" {
		fmt.Fprintln(os.Stderr, "PFX rejected: test-signing-password")
		os.Exit(23)
	}
}
func TestWindowsSigningRetainsNativeFailureAndRedactsPassword(t *testing.T) {
	executable, err := os.Executable()
	require.NoError(t, err)
	data, err := os.ReadFile(executable)
	require.NoError(t, err)
	tools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tools, "signtool.exe"), data, 0755))
	t.Setenv("PATH", tools)
	t.Setenv("HCL_TEST_SIGNER_FAILURE", "1")
	err = signWindows(&flags.Sign{Input: "app.exe", Certificate: "test.pfx", Password: "test-signing-password"})
	require.ErrorContains(t, err, "PFX rejected: <redacted>")
	require.ErrorContains(t, err, "exit status 23")
	require.NotContains(t, err.Error(), "test-signing-password")
	require.NotContains(t, err.Error(), "built-in")
}
