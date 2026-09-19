package commands

import (
	"testing"

	"github.com/atterpac/refresh/engine"
	"github.com/atterpac/refresh/process"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureIgnored(t *testing.T) {
	t.Run("adds pattern when not present", func(t *testing.T) {
		list := []string{".gitignore", ".DS_Store"}
		ensureIgnored(&list, "*_test.go")
		assert.Contains(t, list, "*_test.go")
		assert.Len(t, list, 3)
	})

	t.Run("does not duplicate pattern when already present", func(t *testing.T) {
		list := []string{".gitignore", "*_test.go"}
		ensureIgnored(&list, "*_test.go")
		assert.Contains(t, list, "*_test.go")
		assert.Len(t, list, 2)
	})

	t.Run("adds to empty list", func(t *testing.T) {
		var list []string
		ensureIgnored(&list, "*_test.go")
		assert.Contains(t, list, "*_test.go")
		assert.Len(t, list, 1)
	})
}

func TestFrontendReadiness(t *testing.T) {
	for _, tc := range []struct{ name, frontendURL, address string }{
		{"default", "http://localhost:9245", "localhost:9245"},
		{"custom port", "http://localhost:17321", "localhost:17321"},
		{"secure", "https://localhost:17322", "localhost:17322"},
		{"IPv6", "http://[::1]:17323", "[::1]:17323"},
		{"HTTP default port", "http://localhost", "localhost:80"},
		{"HTTPS default port", "https://localhost", "localhost:443"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := engine.Config{ExecStruct: []process.Execute{{Cmd: "wails3 task common:dev:frontend", Type: process.Background}}}
			require.NoError(t, applyFrontendReadiness(&config, tc.frontendURL))
			require.NotNil(t, config.ExecStruct[0].Readiness)
			assert.Equal(t, tc.address, config.ExecStruct[0].Readiness.TCP)
			assert.Equal(t, "60s", config.ExecStruct[0].Readiness.Timeout)
		})
	}
	t.Run("preserve custom configuration", func(t *testing.T) {
		explicit := &process.Readiness{TCP: "custom:8000", Timeout: "2m"}
		config := engine.Config{ExecStruct: []process.Execute{
			{Cmd: "wails3 task common:dev:frontend", Type: process.Background, Readiness: explicit},
			{Cmd: "npm run dev", Type: process.Background},
			{Cmd: "wails3 task common:dev:frontend && echo ready", Type: process.Background},
			{Cmd: "wails3 task common:dev:frontend", Type: process.Primary},
		}}
		require.NoError(t, applyFrontendReadiness(&config, "http://localhost:9245"))
		assert.Same(t, explicit, config.ExecStruct[0].Readiness)
		for _, step := range config.ExecStruct[1:] {
			assert.Nil(t, step.Readiness)
		}
	})
	t.Run("no frontend", func(t *testing.T) {
		config := engine.Config{ExecStruct: []process.Execute{{Cmd: "wails3 task common:dev:frontend", Type: process.Background}}}
		require.NoError(t, applyFrontendReadiness(&config, ""))
		assert.Nil(t, config.ExecStruct[0].Readiness)
	})
	t.Run("invalid URL", func(t *testing.T) {
		for _, target := range []string{"not-a-url", "file:///tmp/frontend", "http://localhost:bad"} {
			config := engine.Config{ExecStruct: []process.Execute{{Cmd: "wails3 task common:dev:frontend", Type: process.Background}}}
			assert.Error(t, applyFrontendReadiness(&config, target))
		}
	})
}
