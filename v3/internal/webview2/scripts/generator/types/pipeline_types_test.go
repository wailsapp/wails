package types_test

// These tests drive the full generator pipeline (Process → Generate →
// GenerateTests) through the real participle parser over every bundled IDL.
// They live in the external test package so they can import updater/generator
// (which imports types) without an import cycle, giving the types package its
// own end-to-end coverage rather than relying on the generator package's tests.

import (
	"os"
	"path/filepath"
	"testing"

	"updater/generator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func bundledIDLs(t *testing.T) []string {
	t.Helper()
	matches, err := filepath.Glob("../../WebView2.*.idl")
	require.NoError(t, err)
	require.NotEmpty(t, matches)
	return matches
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return data
}

func TestFullPipelineOverBundledIDLs(t *testing.T) {
	for _, path := range bundledIDLs(t) {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			data := readFile(t, path)

			files, err := generator.ParseIDL(data)
			require.NoError(t, err)
			require.NotEmpty(t, files)
			for _, f := range files {
				assert.Positive(t, f.Content.Len(), "%s empty", f.FileName)
			}

			withTests, err := generator.ParseIDLWithTests(data)
			require.NoError(t, err)
			assert.Greater(t, len(withTests), len(files))
		})
	}
}
