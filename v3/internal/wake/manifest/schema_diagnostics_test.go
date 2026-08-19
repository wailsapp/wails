package manifest

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManifestVersionMismatchReportsItsExactRange(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := []byte(`version = 4

project {
  name = "version"
  product_name = "Version"
  identifier = "com.example.version"
  version = "1.0.0"
}

`)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	_, err := LoadFile(root, path, "")
	require.Error(t, err)
	var validation *ValidationError
	require.True(t, errors.As(err, &validation))
	assert.Equal(t, "version", validation.Field)
	assert.Equal(t, SourceRange{Filename: path, StartLine: 1, StartColumn: 1, EndLine: 1, EndColumn: 12}, validation.Range)
}

func TestManifestExpressionDiagnosticUsesTheActualSourceFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "candidate.hcl")
	source := []byte(`version = 3

project {
  name = "literal"
  product_name = "Literal"
  identifier = "com.example.literal"
  version = "1.0.0"
}

build {
  environment = { RELEASE = upper("yes") }
}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	_, err := LoadFile(root, path, "")
	require.Error(t, err)
	var validation *ValidationError
	require.True(t, errors.As(err, &validation))
	assert.Equal(t, "build.environment", validation.Field)
	assert.Equal(t, SourceRange{Filename: path, StartLine: 11, StartColumn: 29, EndLine: 11, EndColumn: 41}, validation.Range)
	assert.Contains(t, validation.Detail, "only literal values are allowed")
}

func TestManifestSchemaDiagnosticsCarryExactRanges(t *testing.T) {
	tests := []struct {
		name   string
		line   string
		field  string
		start  int
		end    int
		detail string
	}{
		{name: "unknown attribute", line: `  mystery = true`, field: "build.mystery", start: 3, end: 10, detail: "Unsupported argument"},
		{name: "wrong type", line: `  trim_path = "yes"`, field: "build.trim_path", start: 15, end: 20, detail: "Unsuitable value type"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, Filename)
			source := []byte("version = 3\n\nproject {\n  name = \"diagnostic\"\n  product_name = \"Diagnostic\"\n  identifier = \"com.example.diagnostic\"\n  version = \"1.0.0\"\n}\n\nbuild {\n" + test.line + "\n}\n")
			require.NoError(t, os.WriteFile(path, source, 0o644))

			_, err := LoadFile(root, path, "")
			require.Error(t, err)
			var validation *ValidationError
			require.True(t, errors.As(err, &validation), "%T: %v", err, err)
			assert.Equal(t, test.field, validation.Field)
			assert.Equal(t, SourceRange{Filename: path, StartLine: 11, StartColumn: test.start, EndLine: 11, EndColumn: test.end}, validation.Range)
			assert.Contains(t, validation.Detail, test.detail)
		})
	}
}

func TestManifestUnknownBlockDiagnosticCarriesExactRange(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := []byte(`version = 3

project {
  name = "diagnostic"
  product_name = "Diagnostic"
  identifier = "com.example.diagnostic"
  version = "1.0.0"
}

mystery {}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	_, err := LoadFile(root, path, "")
	require.Error(t, err)
	var validation *ValidationError
	require.True(t, errors.As(err, &validation), "%T: %v", err, err)
	assert.Equal(t, "mystery", validation.Field)
	assert.Equal(t, SourceRange{Filename: path, StartLine: 10, StartColumn: 1, EndLine: 10, EndColumn: 8}, validation.Range)
	assert.Contains(t, validation.Detail, "Unsupported block type")
}

func TestManifestRejectsFieldsOutsideTheirPlatformSchema(t *testing.T) {
	tests := []struct {
		name       string
		platform   string
		body       string
		field      string
		startLine  int
		startCol   int
		endCol     int
		diagnostic string
	}{
		{name: "iOS identity on Android", platform: "android", body: `  bundle_id = "com.example.invalid"`, field: "android.bundle_id", startLine: 11, startCol: 3, endCol: 12, diagnostic: "Unsupported argument"},
		{name: "Android SDK on Windows", platform: "windows", body: `  target_sdk = 35`, field: "windows.target_sdk", startLine: 11, startCol: 3, endCol: 13, diagnostic: "Unsupported argument"},
		{name: "Apple notarization on iOS", platform: "ios", body: `  notarization {}`, field: "ios.notarization", startLine: 11, startCol: 3, endCol: 15, diagnostic: "Unsupported block type"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, Filename)
			source := []byte("version = 3\n\nproject {\n  name = \"platform-schema\"\n  product_name = \"Platform Schema\"\n  identifier = \"com.example.platform-schema\"\n  version = \"1.0.0\"\n}\n\n" + test.platform + " {\n" + test.body + "\n}\n")
			require.NoError(t, os.WriteFile(path, source, 0o644))

			_, err := LoadFile(root, path, "")
			require.Error(t, err)
			var validation *ValidationError
			require.True(t, errors.As(err, &validation), "%T: %v", err, err)
			assert.Equal(t, test.field, validation.Field)
			assert.Equal(t, SourceRange{Filename: path, StartLine: test.startLine, StartColumn: test.startCol, EndLine: test.startLine, EndColumn: test.endCol}, validation.Range)
			assert.Contains(t, validation.Detail, test.diagnostic)
		})
	}
}

func TestManifestSemanticDiagnosticsCarryTheirAttributeRange(t *testing.T) {
	tests := []struct {
		name      string
		extra     string
		field     string
		blockPath []string
		attribute string
	}{
		{
			name: "binding conflict", field: "frontend.bindings.time_type",
			extra:     "frontend {\n  bindings {\n    interfaces = true\n    time_type = \"Date\"\n  }\n}\n",
			blockPath: []string{"frontend", "bindings"}, attribute: "time_type",
		},
		{
			name: "invalid background mode", field: "ios.background_modes",
			extra:     "ios {\n  background_modes = [\"magic\"]\n}\n",
			blockPath: []string{"ios"}, attribute: "background_modes",
		},
		{
			name: "invalid environment name", field: "build.environment",
			extra:     "build {\n  environment = { \"BAD=KEY\" = \"value\" }\n}\n",
			blockPath: []string{"build"}, attribute: "environment",
		},
		{
			name: "unsafe path", field: "build.output",
			extra:     "build {\n  output = \"../release\"\n}\n",
			blockPath: []string{"build"}, attribute: "output",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, Filename)
			source := []byte("version = 3\n\nproject {\n  name = \"semantic\"\n  product_name = \"Semantic\"\n  identifier = \"com.example.semantic\"\n  version = \"1.0.0\"\n}\n\n" + test.extra)
			require.NoError(t, os.WriteFile(path, source, 0o644))

			_, err := LoadFile(root, path, "")
			require.Error(t, err)
			var validation *ValidationError
			require.True(t, errors.As(err, &validation), "%T: %v", err, err)
			assert.Equal(t, test.field, validation.Field)
			assert.Equal(t, testAttributeRange(t, source, path, test.blockPath, test.attribute), validation.Range)
		})
	}
}

func testAttributeRange(t *testing.T, source []byte, filename string, blockPath []string, attribute string) SourceRange {
	t.Helper()
	file, diagnostics := hclsyntax.ParseConfig(source, filename, hcl.InitialPos)
	require.False(t, diagnostics.HasErrors(), diagnostics.Error())
	body := file.Body.(*hclsyntax.Body)
	for _, name := range blockPath {
		var next *hclsyntax.Body
		for _, block := range body.Blocks {
			if block.Type == name {
				next = block.Body
				break
			}
		}
		require.NotNil(t, next, name)
		body = next
	}
	value := body.Attributes[attribute]
	require.NotNil(t, value, attribute)
	return sourceRange(value.Range())
}

func TestManifestAccumulatesIndependentSchemaDiagnosticsDeterministically(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := []byte(`version = 3

project {
  product_name = "Diagnostics"
  identifier = "com.example.diagnostics"
  version = "1.0.0"
  mystery = true
}

build {
  trim_path = "yes"
}

unknown {}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	_, err := LoadFile(root, path, "")
	require.Error(t, err)
	var fields []string
	var ranges []SourceRange
	for _, item := range joinedErrors(err) {
		var validation *ValidationError
		if errors.As(item, &validation) {
			fields = append(fields, validation.Field)
			ranges = append(ranges, validation.Range)
		}
	}
	assert.Equal(t, []string{"project.mystery", "project.name", "build.trim_path", "unknown"}, fields)
	for _, diagnosticRange := range ranges {
		assert.Equal(t, path, diagnosticRange.Filename)
		assert.Positive(t, diagnosticRange.StartLine)
		assert.Positive(t, diagnosticRange.StartColumn)
	}
}

func joinedErrors(err error) []error {
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		var result []error
		for _, child := range joined.Unwrap() {
			result = append(result, joinedErrors(child)...)
		}
		return result
	}
	return []error{err}
}

func TestManifestRequiredStringRejectsEmptyLiteralAtItsExactRange(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := []byte(`version = 3

project {
  name = ""
  product_name = "Empty"
  identifier = "com.example.empty"
  version = "1.0.0"
}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	_, err := LoadFile(root, path, "")
	require.Error(t, err)
	var validation *ValidationError
	require.True(t, errors.As(err, &validation))
	assert.Equal(t, "project.name", validation.Field)
	assert.Equal(t, testExpressionRange(t, source, path, []string{"project"}, "name"), validation.Range)
}

func testExpressionRange(t *testing.T, source []byte, filename string, blockPath []string, attribute string) SourceRange {
	t.Helper()
	file, diagnostics := hclsyntax.ParseConfig(source, filename, hcl.InitialPos)
	require.False(t, diagnostics.HasErrors(), diagnostics.Error())
	body := file.Body.(*hclsyntax.Body)
	for _, name := range blockPath {
		var next *hclsyntax.Body
		for _, block := range body.Blocks {
			if block.Type == name {
				next = block.Body
				break
			}
		}
		require.NotNil(t, next, name)
		body = next
	}
	value := body.Attributes[attribute]
	require.NotNil(t, value, attribute)
	return sourceRange(value.Expr.Range())
}

func TestManifestRejectsFieldsOutsideTheirPackageFormatSchema(t *testing.T) {
	tests := []struct {
		name, format, attribute, field string
		endColumn                      int
	}{
		{name: "DMG visual field on NSIS", format: "nsis", attribute: `background = "background.png"`, field: `package["nsis"].background`, endColumn: 13},
		{name: "complete replacement on MSIX", format: "msix", attribute: `template = "AppxManifest.xml"`, field: `package["msix"].template`, endColumn: 11},
		{name: "AppImage metadata on DEB", format: "deb", attribute: `categories = ["Utility"]`, field: `package["deb"].categories`, endColumn: 13},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, Filename)
			source := []byte("version = 3\n\nproject {\n  name = \"package-schema\"\n  product_name = \"Package Schema\"\n  identifier = \"com.example.package-schema\"\n  version = \"1.0.0\"\n}\n\npackage \"" + test.format + "\" {\n  " + test.attribute + "\n}\n")
			require.NoError(t, os.WriteFile(path, source, 0o644))

			_, err := LoadFile(root, path, "")
			require.Error(t, err)
			var validation *ValidationError
			require.True(t, errors.As(err, &validation), "%T: %v", err, err)
			assert.Equal(t, test.field, validation.Field)
			assert.Equal(t, SourceRange{Filename: path, StartLine: 11, StartColumn: 3, EndLine: 11, EndColumn: test.endColumn}, validation.Range)
			assert.Contains(t, validation.Detail, "Unsupported argument")
		})
	}
}
