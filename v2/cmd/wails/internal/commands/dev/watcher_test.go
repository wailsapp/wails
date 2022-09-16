package dev

import (
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wailsapp/wails/v2/internal/fs"
)

func Test_processDirectories(t *testing.T) {
	tests := []struct {
		name       string
		dirs       []string
		ignoreDirs []string
		want       []string
	}{
		{
			name:       "Should ignore .git",
			ignoreDirs: []string{".git"},
			dirs:       []string{".git", "some/path/to/nested/.git", "some/path/to/nested/.git/CHANGELOG"},
			want:       []string{},
		},
		{
			name:       "Should ignore node_modules",
			ignoreDirs: []string{"node_modules"},
			dirs:       []string{"node_modules", "path/to/node_modules", "path/to/node_modules/some/other/path"},
			want:       []string{},
		},
		{
			name:       "Should ignore dirs starting with .",
			ignoreDirs: []string{".*"},
			dirs:       []string{".test", ".gitignore", ".someother", "valid"},
			want:       []string{"valid"},
		},
		{
			name:       "Should ignore dirs in ignoreDirs",
			dirs:       []string{"build", "build/my.exe", "build/my.app"},
			ignoreDirs: []string{"build"},
			want:       []string{},
		},
		{
			name:       "Should ignore subdirectories",
			dirs:       []string{"build", "some/path/to/build", "some/path/to/CHANGELOG", "some/other/path"},
			ignoreDirs: []string{"some/**/*"},
			want:       []string{"build"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processDirectories(tt.dirs, tt.ignoreDirs)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processDirectories() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetIgnoreDirs(t *testing.T) {

	// Remove testdir if it exists
	_ = os.RemoveAll("testdir")

	tests := []struct {
		name      string
		files     []string
		want      []string
		shouldErr bool
	}{
		{
			name:  "Should have defaults",
			files: []string{},
			want:  []string{"testdir/build/*", ".*", "node_modules"},
		},
		{
			name:  "Should ignore dotFiles",
			files: []string{".test1", ".wailsignore"},
			want:  []string{"testdir/build/*", ".*", "node_modules"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			err := fs.Mkdir("testdir")
			require.NoError(t, err)
			defer func() {
				err := os.RemoveAll("testdir")
				require.NoError(t, err)
			}()
			for _, file := range tt.files {
				fs.MustWriteString(filepath.Join("testdir", file), "")
			}

			got := getIgnoreDirs("testdir")

			got = lo.Map(got, func(s string, _ int) string {
				return filepath.ToSlash(s)
			})

			if (err != nil) != tt.shouldErr {
				t.Errorf("initialiseWatcher() error = %v, shouldErr %v", err, tt.shouldErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initialiseWatcher() got = %v, want %v", got, tt.want)
			}
		})
	}
}
