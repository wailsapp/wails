package commands

import (
	"go/build"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinuxDevtoolsImplementationSelectedForEveryRunMode(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		for _, tags := range [][]string{nil, {"devtools"}, {"production"}, {"production", "devtools"}} {
			context := build.Default
			context.GOOS, context.GOARCH = "linux", "amd64"
			context.BuildTags = append([]string(nil), tags...)
			if legacy {
				context.BuildTags = append(context.BuildTags, "gtk3")
			}
			count := 0
			for _, file := range []string{"webview_window_linux_dev.go", "webview_window_linux_production.go"} {
				matches, err := context.MatchFile(filepath.Join("..", "..", "pkg", "application"), file)
				require.NoError(t, err)
				if matches {
					count++
				}
			}
			require.Equal(t, 1, count, "tags %v must select exactly one devtools implementation", context.BuildTags)
		}
	}
}
