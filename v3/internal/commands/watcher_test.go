package commands

import (
	"testing"

	"github.com/atterpac/refresh/engine"
	"github.com/atterpac/refresh/process"
	"github.com/stretchr/testify/require"
)

func TestEnsureIgnored(t *testing.T) {
	t.Run("adds pattern when not present", func(t *testing.T) {
		list := []string{".gitignore", ".DS_Store"}
		ensureIgnored(&list, "*_test.go")
		require.Contains(t, list, "*_test.go")
		require.Len(t, list, 3)
	})

	t.Run("does not duplicate pattern when already present", func(t *testing.T) {
		list := []string{".gitignore", "*_test.go"}
		ensureIgnored(&list, "*_test.go")
		require.Contains(t, list, "*_test.go")
		require.Len(t, list, 2)
	})

	t.Run("adds to empty list", func(t *testing.T) {
		var list []string
		ensureIgnored(&list, "*_test.go")
		require.Contains(t, list, "*_test.go")
		require.Len(t, list, 1)
	})
}

func TestEnsureFrontendDevServerReadyTask(t *testing.T) {
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://localhost:9245")
	config := engine.Config{ExecStruct: []process.Execute{
		{Cmd: "frontend", Type: process.Background},
		{Cmd: "application", Type: process.Primary},
	}}

	ensureFrontendDevServerReadyTask(&config)

	require.Equal(t, []process.Execute{
		{Cmd: "frontend", Type: process.Background},
		{Cmd: frontendDevServerReadyCommand, Type: process.Once},
		{Cmd: "application", Type: process.Primary},
	}, config.ExecStruct)
}

func TestEnsureFrontendDevServerReadyTaskIsIdempotent(t *testing.T) {
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://localhost:9245")
	config := engine.Config{ExecStruct: []process.Execute{
		{Cmd: frontendDevServerReadyCommand, Type: process.Once},
		{Cmd: "application", Type: process.Primary},
	}}

	ensureFrontendDevServerReadyTask(&config)

	require.Len(t, config.ExecStruct, 2)
}

func TestEnsureFrontendDevServerReadyTaskRequiresDevServer(t *testing.T) {
	t.Setenv("FRONTEND_DEVSERVER_URL", "")
	config := engine.Config{ExecStruct: []process.Execute{
		{Cmd: "application", Type: process.Primary},
	}}

	ensureFrontendDevServerReadyTask(&config)

	require.Len(t, config.ExecStruct, 1)
}
