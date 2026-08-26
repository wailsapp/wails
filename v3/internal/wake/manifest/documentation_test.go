package manifest

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublishedHCLGuideExamplesMatchTheManifestSchema(t *testing.T) {
	guide := repositoryDocumentationPath("experimental", "hcl-builds.mdx")
	data, err := os.ReadFile(guide)
	require.NoError(t, err)
	examples := fencedExamples(string(data), "hcl")
	require.NotEmpty(t, examples)
	for index, example := range examples {
		if runtime.GOOS == "windows" && strings.Contains(example, `hook "`) {
			// Published hook examples use Unix scripts; the Windows extension
			// contract is covered by host-specific manifest tests.
			continue
		}
		root := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend"), 0o755))
		require.NoError(t, os.MkdirAll(filepath.Join(root, "scripts"), 0o755))
		for _, script := range []string{"generate-version.sh", "check-packages.sh"} {
			require.NoError(t, os.WriteFile(filepath.Join(root, "scripts", script), []byte("#!/bin/sh\n"), 0o755))
		}
		require.NoError(t, os.WriteFile(filepath.Join(root, "version.txt"), []byte("1.0.0\n"), 0o644))
		_, err := decodeHCL(root, guide, []byte(example), "")
		require.NoErrorf(t, err, "HCL example %d", index+1)
		_, err = decodeHCL(root, guide, []byte(example), "release")
		require.NoErrorf(t, err, "HCL example %d release profile", index+1)
	}
}

func TestPublishedFieldReferenceIsCurrent(t *testing.T) {
	data, err := os.ReadFile(repositoryDocumentationPath("reference", "wails-hcl-fields.md"))
	require.NoError(t, err)
	marker := "| Field | Type | Required | Default | Example | Applies to | Description |\n"
	index := strings.Index(string(data), marker)
	require.NotEqual(t, -1, index)
	assert.Equal(t, string(SchemaReferenceMarkdown()), string(data[index:]))
}

func repositoryDocumentationPath(parts ...string) string {
	path := filepath.Join("..", "..", "..", "..", "docs", "src", "content", "docs")
	return filepath.Join(append([]string{path}, parts...)...)
}

func fencedExamples(document, language string) []string {
	marker := "```" + language + "\n"
	var result []string
	for {
		start := strings.Index(document, marker)
		if start < 0 {
			return result
		}
		document = document[start+len(marker):]
		end := strings.Index(document, "\n```")
		if end < 0 {
			return result
		}
		result = append(result, document[:end]+"\n")
		document = document[end+4:]
	}
}
