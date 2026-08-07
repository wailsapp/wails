package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveMCPRootDefaultsToWorkingDirectory(t *testing.T) {
	root, err := resolveMCPRoot("")
	require.NoError(t, err)
	want, err := os.Getwd()
	require.NoError(t, err)
	want, err = filepath.EvalSymlinks(want)
	require.NoError(t, err)
	require.Equal(t, want, root)
}

func TestMCPProjectPathRejectsEscapes(t *testing.T) {
	root := t.TempDir()
	server := &mcpServer{root: root, token: "test"}

	_, err := server.projectPath("../outside")
	require.Error(t, err)

	link := filepath.Join(root, "link")
	require.NoError(t, os.Symlink(t.TempDir(), link))
	_, err = server.projectPath("link")
	require.Error(t, err)
}

func TestMCPAuthorizeUsesExactConstantTimeToken(t *testing.T) {
	server := &mcpServer{token: "correct-token"}
	require.NoError(t, server.authorize("correct-token"))
	require.Error(t, server.authorize("correct-token-extra"))
	require.Error(t, server.authorize(""))
}

func TestNewMCPTokenHasSufficientEntropyAndIsURLSafe(t *testing.T) {
	one, err := newMCPToken()
	require.NoError(t, err)
	two, err := newMCPToken()
	require.NoError(t, err)
	require.Len(t, one, 43)
	require.NotEqual(t, one, two)
	require.NotContains(t, one, "+")
	require.NotContains(t, one, "/")
	require.NotContains(t, one, "=")
}

func TestValidateMCPArgsBoundsInput(t *testing.T) {
	require.NoError(t, validateMCPArgs([]string{"--tags", "mcp"}))
	require.Error(t, validateMCPArgs(make([]string, maxMCPArgs+1)))
	require.Error(t, validateMCPArgs([]string{"bad\x00arg"}))
	require.Error(t, validateMCPArgs([]string{string(make([]byte, 4097))}))
}

func TestMCPInitRequiresExplicitExternalApproval(t *testing.T) {
	server := &mcpServer{root: t.TempDir(), token: "secret"}
	_, _, err := server.init(t.Context(), nil, mcpInitInput{
		mcpAuthInput: mcpAuthInput{Token: "secret"},
		Name:         "demo",
		Template:     "https://example.invalid/template",
	})
	require.ErrorContains(t, err, "allowExternal")
}
