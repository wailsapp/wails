package packagetemplate

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderFileDoesNotInterpretPackageModel(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "package.tmpl")
	destination := filepath.Join(root, "generated", "package.conf")
	contents := []byte("{{.Project.ProductName}}|{{.Target.OS}}/{{.Target.Arch}}|{{index .Target.Capabilities 0}}|{{.Package.Format}}|{{index .Options \"channel\"}}\n")
	require.NoError(t, os.WriteFile(source, contents, 0o644))

	err := Copy(source, destination)
	require.NoError(t, err)
	data, err := os.ReadFile(destination)
	require.NoError(t, err)
	assert.Equal(t, contents, data)
}

func TestRenderDirectoryCopiesAllOwnedFiles(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "android")
	destination := filepath.Join(root, "generated")
	require.NoError(t, os.MkdirAll(filepath.Join(source, "app"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(source, "app", "build.gradle.tmpl"), []byte("applicationId '{{.Project.Identifier}}'\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(source, "gradlew"), []byte("#!/bin/sh\n"), 0o770))
	require.NoError(t, os.Chmod(filepath.Join(source, "gradlew"), 0o770))

	err := Copy(source, destination)
	require.NoError(t, err)
	gradle, err := os.ReadFile(filepath.Join(destination, "app", "build.gradle.tmpl"))
	require.NoError(t, err)
	assert.Equal(t, "applicationId '{{.Project.Identifier}}'\n", string(gradle))
	info, err := os.Stat(filepath.Join(destination, "gradlew"))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o770), info.Mode().Perm())
	_, err = os.Stat(filepath.Join(destination, "app", "build.gradle"))
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestRenderTreatsUnknownTemplateFieldsAsOpaqueBytes(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "package.tmpl")
	destination := filepath.Join(root, "package.conf")
	contents := []byte("{{.InternalNode}}|{{.Paths.Binary}}|{{.Paths.Output}}\n")
	require.NoError(t, os.WriteFile(source, contents, 0o644))

	err := Copy(source, destination)
	require.NoError(t, err)
	data, err := os.ReadFile(destination)
	require.NoError(t, err)
	assert.Equal(t, contents, data)
}

func TestRenderRejectsDirectorySymlinks(t *testing.T) {
	root := t.TempDir()
	missing := filepath.Join(root, "missing.tmpl")
	require.NoError(t, os.WriteFile(missing, []byte("{{.InternalNode}}"), 0o644))
	source := filepath.Join(root, "tree")
	require.NoError(t, os.MkdirAll(source, 0o755))
	require.NoError(t, os.Symlink(missing, filepath.Join(source, "linked.tmpl")))
	err := Copy(source, filepath.Join(root, "generated"))
	require.ErrorContains(t, err, "unsupported symlink")
}

func TestRenderCopiesUserFileByteForByteWithoutInterpolation(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "package.tmpl")
	destination := filepath.Join(root, "generated", "package.conf")
	contents := []byte{0xff, 0xfe, '{', '{', '.', 'P', 'r', 'o', 'j', 'e', 'c', 't', '.', 'N', 'a', 'm', 'e', '}', '}', '\r', '\n', 0x00}
	require.NoError(t, os.WriteFile(source, contents, 0o751))
	require.NoError(t, os.Chmod(source, 0o751))

	require.NoError(t, Copy(source, destination))
	actual, err := os.ReadFile(destination)
	require.NoError(t, err)
	assert.Equal(t, contents, actual)
	info, err := os.Stat(destination)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o751), info.Mode().Perm())
}

func TestRenderCopiesDirectoryNamesAndContentsExactly(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "replacement")
	destination := filepath.Join(root, "generated")
	require.NoError(t, os.MkdirAll(filepath.Join(source, "app"), 0o750))
	contents := []byte("applicationId '{{.Project.Identifier}}'\r\n")
	require.NoError(t, os.WriteFile(filepath.Join(source, "app", "build.gradle.tmpl"), contents, 0o640))

	require.NoError(t, Copy(source, destination))
	actual, err := os.ReadFile(filepath.Join(destination, "app", "build.gradle.tmpl"))
	require.NoError(t, err)
	assert.Equal(t, contents, actual)
	_, err = os.Stat(filepath.Join(destination, "app", "build.gradle"))
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestCopyDirectoryFailureLeavesPreviousDestinationUntouched(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "replacement")
	destination := filepath.Join(root, "generated")
	require.NoError(t, os.MkdirAll(source, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(source, "first"), []byte("new"), 0o644))
	require.NoError(t, os.Symlink("first", filepath.Join(source, "unsupported-link")))
	require.NoError(t, os.MkdirAll(destination, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(destination, "previous"), []byte("old"), 0o644))

	err := Copy(source, destination)
	require.ErrorContains(t, err, "unsupported symlink")
	assert.Equal(t, "old", string(mustReadFile(t, filepath.Join(destination, "previous"))))
	assert.NoFileExists(t, filepath.Join(destination, "first"))
}

func mustReadFile(t testing.TB, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return data
}

func BenchmarkCopyReplacement(b *testing.B) {
	for _, size := range []int{4 << 10, 4 << 20} {
		b.Run(fmt.Sprintf("file-%d", size), func(b *testing.B) {
			root := b.TempDir()
			source := filepath.Join(root, "source")
			destination := filepath.Join(root, "destination")
			require.NoError(b, os.WriteFile(source, bytes.Repeat([]byte{0xa5}, size), 0o755))
			b.ReportAllocs()
			b.SetBytes(int64(size))
			b.ResetTimer()
			for range b.N {
				if err := Copy(source, destination); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
	b.Run("tree-500", func(b *testing.B) {
		root := b.TempDir()
		source := filepath.Join(root, "source")
		destination := filepath.Join(root, "destination")
		for index := range 500 {
			path := filepath.Join(source, fmt.Sprintf("dir-%02d", index%20), fmt.Sprintf("file-%03d", index))
			require.NoError(b, os.MkdirAll(filepath.Dir(path), 0o755))
			require.NoError(b, os.WriteFile(path, bytes.Repeat([]byte{byte(index)}, 1024), 0o644))
		}
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			if err := Copy(source, destination); err != nil {
				b.Fatal(err)
			}
		}
	})
}
