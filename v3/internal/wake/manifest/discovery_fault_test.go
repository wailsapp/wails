package manifest

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManifestDiscoveryFilesystemFaults(t *testing.T) {
	directoryInfo, err := os.Stat(t.TempDir())
	require.NoError(t, err)

	t.Run("absolute", func(t *testing.T) {
		_, _, err := discoverWithOperations("project", func(string) (string, error) { return "", fs.ErrInvalid }, os.Stat)
		require.ErrorIs(t, err, fs.ErrInvalid)
	})
	t.Run("manifest stat", func(t *testing.T) {
		root := filepath.Clean(string(filepath.Separator) + "virtual-project")
		stat := func(path string) (fs.FileInfo, error) {
			switch path {
			case root:
				return directoryInfo, nil
			case filepath.Join(root, Filename):
				return nil, fs.ErrPermission
			default:
				return nil, fs.ErrNotExist
			}
		}
		_, _, err := discoverWithOperations(root, func(value string) (string, error) { return value, nil }, stat)
		require.ErrorIs(t, err, fs.ErrPermission)
	})
	t.Run("module stat", func(t *testing.T) {
		root := filepath.Clean(string(filepath.Separator) + "virtual-project")
		stat := func(path string) (fs.FileInfo, error) {
			switch path {
			case root:
				return directoryInfo, nil
			case filepath.Join(root, "go.mod"):
				return nil, fs.ErrPermission
			default:
				return nil, fs.ErrNotExist
			}
		}
		_, _, err := discoverWithOperations(root, func(value string) (string, error) { return value, nil }, stat)
		require.ErrorIs(t, err, fs.ErrPermission)
	})
}

func TestManifestLoadFilesystemFaults(t *testing.T) {
	permission := errors.New("permission denied")
	_, err := loadWithOperations("project", "", func(string) (string, string, error) {
		return "", "", permission
	}, os.ReadFile)
	require.ErrorIs(t, err, permission)

	_, err = loadWithOperations("project", "", func(string) (string, string, error) {
		return "/project", "/project/wails.hcl", nil
	}, func(string) ([]byte, error) { return nil, permission })
	require.ErrorIs(t, err, permission)
	assert.Contains(t, err.Error(), "read "+Filename)
}

func TestManifestDocumentWritersCoverValidAndInvalidDocuments(t *testing.T) {
	root := t.TempDir()
	invalid := Document{}
	require.Error(t, WriteDocument(root, invalid))
	require.Error(t, WriteMigrationDraft(root, invalid))

	valid := NewDocument(Project{Name: "writer", ProductName: "Writer", Identifier: "com.example.writer", Version: "1.0.0"})
	require.NoError(t, WriteDocument(root, valid))
	require.NoError(t, WriteMigrationDraft(root, valid))
	assert.FileExists(t, filepath.Join(root, Filename))
	assert.FileExists(t, filepath.Join(root, MigratedFilename))
}

func TestWriteMigrationDraftAtNeverOverwritesAndPrefixesReviewComments(t *testing.T) {
	root := t.TempDir()
	doc := NewDocument(Project{Name: "writer", ProductName: "Writer", Identifier: "com.example.writer", Version: "1.0.0"})
	relative := "proposals/legacy.hcl"
	require.NoError(t, WriteMigrationDraftAt(root, relative, doc, []string{"review this", "BLOCKED: replace custom task"}))

	path := filepath.Join(root, filepath.FromSlash(relative))
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(string(data), "# review this\n# BLOCKED: replace custom task\n\n"))

	original := append([]byte(nil), data...)
	err = WriteMigrationDraftAt(root, relative, doc, nil)
	require.ErrorIs(t, err, os.ErrExist)
	after, readErr := os.ReadFile(path)
	require.NoError(t, readErr)
	assert.Equal(t, original, after)
}

func TestWriteMigrationDraftAtRejectsActiveGeneratedAndEscapingPaths(t *testing.T) {
	root := t.TempDir()
	doc := NewDocument(Project{Name: "writer", ProductName: "Writer", Identifier: "com.example.writer", Version: "1.0.0"})
	for _, output := range []string{Filename, "WAILS.HCL", EjectedFilename, "WAILS.EJECTED.HCL", ".wails/proposal.hcl", ".WAILS/proposal.hcl", "../proposal.hcl", "proposal.json"} {
		t.Run(output, func(t *testing.T) {
			err := WriteMigrationDraftAt(root, output, doc, nil)
			require.Error(t, err)
		})
	}
}

func TestDocumentConversionRejectsEmptyRequiredProjectValues(t *testing.T) {
	value := func(text string) *string { return &text }
	_, err := documentFromHCL(hclDocument{Project: &hclProject{
		Name: value(""), ProductName: value("Product"), Identifier: value("com.example.product"), Version: value("1.0.0"),
	}})
	require.Error(t, err)
}
