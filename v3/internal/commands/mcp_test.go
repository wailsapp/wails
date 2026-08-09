package commands

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
	root, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	server := &mcpServer{root: root, token: "test"}

	inside := filepath.Join(root, "inside")
	require.NoError(t, os.Mkdir(inside, 0o755))
	got, err := server.projectPath("inside")
	require.NoError(t, err)
	require.Equal(t, inside, got)

	_, err = server.projectPath("../outside")
	require.Error(t, err)

	link := filepath.Join(root, "link")
	require.NoError(t, os.Symlink(t.TempDir(), link))
	_, err = server.projectPath("link")
	require.Error(t, err)
}

func TestMCPProjectPathRejectsSymlinkedParentOfMissingChild(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	server := &mcpServer{root: root, token: "test"}
	require.NoError(t, os.Symlink(t.TempDir(), filepath.Join(root, "link")))

	_, err = server.projectPath("link/new-project")
	require.Error(t, err)
}

func TestMCPAuthorizeRequiresExactToken(t *testing.T) {
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
	root, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	server := &mcpServer{root: root, token: "test"}
	require.NoError(t, server.validateMCPArgs(root, []string{"--tags", "mcp"}))
	require.NoError(t, server.validateMCPArgs(root, []string{"--output", "group"}))
	require.NoError(t, server.validateMCPArgs(root, []string{"--config", "build/config.yml"}))
	require.Error(t, server.validateMCPArgs(root, make([]string, maxMCPArgs+1)))
	require.Error(t, server.validateMCPArgs(root, []string{"bad\x00arg"}))
	require.Error(t, server.validateMCPArgs(root, []string{strings.Repeat("a", 4097)}))
	require.Error(t, server.validateMCPArgs(root, []string{"--taskfile", "../outside/Taskfile.yml"}))
	require.Error(t, server.validateMCPArgs(root, []string{"--obfuscated-output=../outside"}))
	require.Error(t, server.validateMCPArgs(root, []string{"--config="}))
	require.Error(t, server.validateMCPArgs(root, []string{"--config", ""}))
}

func TestMCPJobStopReportsStopping(t *testing.T) {
	job, err := newMCPJob(t.TempDir(), func() {})
	require.NoError(t, err)

	job.requestStop()
	require.Equal(t, "stopping", job.snapshot().State)
}

func TestDecodeMCPJSON(t *testing.T) {
	require.Equal(t, map[string]any{"ok": true}, decodeMCPJSON(`{"ok":true}`))
	require.Equal(t, "not json", decodeMCPJSON("not json"))
}

func TestMCPRequestIsLoopback(t *testing.T) {
	tests := []struct {
		name   string
		host   string
		origin string
		want   bool
	}{
		{name: "ipv4", host: "127.0.0.1:1234", origin: "http://127.0.0.1:1234", want: true},
		{name: "ipv6", host: "[::1]:1234", origin: "http://localhost:1234", want: true},
		{name: "missing origin", host: "localhost:1234", want: true},
		{name: "non-loopback host", host: "192.0.2.1:1234", origin: "http://192.0.2.1:1234", want: false},
		{name: "non-loopback origin", host: "127.0.0.1:1234", origin: "http://example.com:1234", want: false},
		{name: "https origin", host: "127.0.0.1:1234", origin: "https://127.0.0.1:1234", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "http://127.0.0.1:1234/mcp", nil)
			r.Host = tt.host
			if tt.origin != "" {
				r.Header.Set("Origin", tt.origin)
			}
			require.Equal(t, tt.want, mcpRequestIsLoopback(r))
		})
	}
}

func TestMCPRejectsConflictingTransportsBeforeRootResolution(t *testing.T) {
	err := MCP(&MCPOptions{HTTP: true, Stdio: true, Root: filepath.Join(t.TempDir(), "missing")})
	require.EqualError(t, err, "choose at most one of --http and --stdio")
}

func TestMCPJobsReapFinishedJobs(t *testing.T) {
	jobs := newMCPJobs()
	all := make([]*mcpJob, 0, maxMCPJobs)
	for range maxMCPJobs {
		job, err := newMCPJob(t.TempDir(), func() {})
		require.NoError(t, err)
		require.True(t, jobs.add(job))
		all = append(all, job)
	}

	blocked, err := newMCPJob(t.TempDir(), func() {})
	require.NoError(t, err)
	require.False(t, jobs.add(blocked))

	all[0].finish(nil)
	require.True(t, jobs.add(blocked))
	_, ok := jobs.get(all[0].id)
	require.False(t, ok)
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
