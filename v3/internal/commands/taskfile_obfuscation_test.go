package commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildAssetsSupportObfuscation(t *testing.T) {
	tests := []struct {
		name string
		path string
		want []string
	}{
		{
			name: "common taskfile template",
			path: "build_assets/Taskfile.tmpl.yml",
			want: []string{
				"ref: .OBFUSCATED",
				" -obfuscated",
			},
		},
		{
			name: "darwin taskfile",
			path: "build_assets/darwin/Taskfile.yml",
			want: []string{
				"command -v garble",
				"garble {{.GARBLE_ARGS}} build",
				"wails_obfuscated",
				"-e OBFUSCATED=true",
				"-e GARBLE_ARGS=\"{{.GARBLE_ARGS}}\"",
			},
		},
		{
			name: "linux taskfile",
			path: "build_assets/linux/Taskfile.yml",
			want: []string{
				"command -v garble",
				"garble {{.GARBLE_ARGS}} build",
				"wails_obfuscated",
				"-e OBFUSCATED=true",
				"-e GARBLE_ARGS=\"{{.GARBLE_ARGS}}\"",
			},
		},
		{
			name: "windows taskfile",
			path: "build_assets/windows/Taskfile.yml",
			want: []string{
				"command -v garble",
				"garble {{.GARBLE_ARGS}} build",
				"wails_obfuscated",
				"-e OBFUSCATED=true",
				"-e GARBLE_ARGS=\"{{.GARBLE_ARGS}}\"",
			},
		},
		{
			name: "cross dockerfile",
			path: "build_assets/docker/Dockerfile.cross",
			want: []string{
				"go install mvdan.cc/garble",
				`if [ "$OBFUSCATED" = "true" ]; then`,
				"wails_obfuscated",
				"garble ${GARBLE_ARGS} build",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := buildAssets.ReadFile(tt.path)
			require.NoError(t, err)

			content := string(data)
			for _, want := range tt.want {
				assert.True(t, strings.Contains(content, want), "expected %q to contain %q", tt.path, want)
			}
		})
	}
}
