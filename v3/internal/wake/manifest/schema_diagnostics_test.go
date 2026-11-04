package manifest

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const diagnosticProjectYAML = `version: 3
project:
  name: diagnostics
  product_name: Diagnostics
  identifier: com.example.diagnostics
  version: 1.0.0
`

func TestYAMLSchemaDiagnosticsCarryExactRanges(t *testing.T) {
	root := t.TempDir()
	filename := filepath.Join(root, Filename)
	tests := []struct {
		name, suffix, field, detail string
		line, column                int
	}{
		{"unknown field", "build:\n  surprise: true\n", "build.surprise", "unsupported field", 8, 3},
		{"wrong type", "build:\n  trim_path: yes\n", "build.trim_path", "a boolean is required", 8, 14},
		{"empty required string", "project:\n  name: ''\n  product_name: App\n  identifier: com.example.app\n  version: 1.0.0\n", "project.name", "must not be empty", 3, 9},
		{"format field mismatch", "packages:\n  nsis:\n    background: image.png\n", `packages["nsis"].background`, "unsupported field", 9, 5},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := diagnosticProjectYAML + test.suffix
			if test.name == "empty required string" {
				source = "version: 3\n" + test.suffix
			}
			_, err := decodeYAML(root, filename, []byte(source), "")
			require.Error(t, err)
			var validation *ValidationError
			require.ErrorAs(t, err, &validation)
			assert.Equal(t, test.field, validation.Field)
			assert.Contains(t, validation.Detail, test.detail)
			assert.Equal(t, filename, validation.Range.Filename)
			assert.Equal(t, test.line, validation.Range.StartLine)
			assert.Equal(t, test.column, validation.Range.StartColumn)
		})
	}
}

func TestYAMLSemanticDiagnosticsUseNearestExplicitField(t *testing.T) {
	root := t.TempDir()
	filename := filepath.Join(root, Filename)
	_, err := decodeYAML(root, filename, []byte(diagnosticProjectYAML+"build:\n  output: ../outside\n"), "")
	require.Error(t, err)
	var validation *ValidationError
	require.ErrorAs(t, err, &validation)
	assert.Equal(t, "build.output", validation.Field)
	assert.Equal(t, 8, validation.Range.StartLine)
	assert.Equal(t, 11, validation.Range.StartColumn)
}

func TestStrictYAMLDiagnosticsRejectAmbiguousFeatures(t *testing.T) {
	tests := []struct {
		name, suffix, detail string
	}{
		{"duplicate", "build:\n  output: one\n  output: two\n", "duplicate mapping key"},
		{"anchor", "build: &build\n  output: bin\n", "anchors are not supported"},
		{"alias", "frontend: &shared {}\nbuild: *shared\n", "aliases are not supported"},
		{"merge", "build:\n  <<: {output: bin}\n", "not supported"},
		{"custom tag", "build: !custom {output: bin}\n", "custom YAML tags are not supported"},
		{"null", "build: null\n", "null values are not supported"},
		{"multiple documents", "---\nproject: {}\n", "only one YAML document is allowed"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := decodeYAML(t.TempDir(), Filename, []byte(diagnosticProjectYAML+test.suffix), "")
			require.Error(t, err)
			assert.Contains(t, err.Error(), test.detail)
		})
	}
}

func TestYAMLValidationJoinsAllStructuralErrors(t *testing.T) {
	_, err := decodeYAML(t.TempDir(), Filename, []byte("version: two\nproject:\n  name: ''\n  surprise: true\n"), "")
	require.Error(t, err)
	var joined interface{ Unwrap() []error }
	require.True(t, errors.As(err, &joined))
	assert.GreaterOrEqual(t, len(joined.Unwrap()), 3)
}

func TestYAMLParseErrorNamesTheFileBeingDecoded(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "wails.migrated.yaml")
	_, err := decodeYAML(t.TempDir(), filename, []byte("not: [valid"), "")
	require.Error(t, err)
	var validation *ValidationError
	require.ErrorAs(t, err, &validation)
	assert.Contains(t, validation.Detail, "wails.migrated.yaml")
}
