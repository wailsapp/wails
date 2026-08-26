package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatValidationDiagnosticsShowsSourceCaretAndHint(t *testing.T) {
	path := filepath.Join(t.TempDir(), Filename)
	require.NoError(t, os.WriteFile(path, []byte("version = 3\nfrontend { install = \"magic\" }\n"), 0o644))
	err := &ValidationError{Field: "frontend.install", Detail: `unsupported package manager "magic"`, Range: SourceRange{Filename: path, StartLine: 2, StartColumn: 12, EndLine: 2, EndColumn: 29}}

	formatted, ok := FormatValidationDiagnostics(err)
	require.True(t, ok)
	assert.Contains(t, formatted, path+":2:12")
	assert.Contains(t, formatted, `2 | frontend { install = "magic" }`)
	assert.Contains(t, formatted, "^^^^^^^^^^^^^^^^^")
	assert.Contains(t, formatted, "Hint: use npm, pnpm, yarn, or bun")
}

func TestFormatValidationDiagnosticsRendersJoinedErrorsAndRejectsOtherDomains(t *testing.T) {
	joined := errorsJoinForTest(
		&ValidationError{Field: "one", Detail: "first"},
		&ValidationError{Field: "two", Detail: "second"},
	)
	formatted, ok := FormatValidationDiagnostics(joined)
	require.True(t, ok)
	assert.Contains(t, formatted, "one: first\n\n")
	assert.Contains(t, formatted, "two: second")

	formatted, ok = FormatValidationDiagnostics(fmt.Errorf("build failed"))
	assert.False(t, ok)
	assert.Empty(t, formatted)
}

func TestFormatValidationDiagnosticsRedactsEnvironmentValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), Filename)
	line := `  environment = { "BAD-NAME" = "do-not-print-this" }`
	require.NoError(t, os.WriteFile(path, []byte("version = 3\n"+line+"\n"), 0o644))
	err := &ValidationError{
		Field:  "build.environment",
		Detail: `contains invalid variable name "BAD-NAME"`,
		Range:  SourceRange{Filename: path, StartLine: 2, StartColumn: 3, EndLine: 2, EndColumn: len(line) + 1},
	}

	formatted, ok := FormatValidationDiagnostics(err)
	require.True(t, ok)
	assert.NotContains(t, formatted, "do-not-print-this")
	assert.Contains(t, formatted, `"BAD-NAME" = "<redacted>"`)
	assert.NotContains(t, formatted, "\x1b[")
}

func TestFormatValidationDiagnosticsMatchesGoldenOutput(t *testing.T) {
	path := filepath.Join(t.TempDir(), Filename)
	require.NoError(t, os.WriteFile(path, []byte("version = 3\nprofile \"release\" { target \"linux/amd64\" { formats = [\"aab\"] } }\n"), 0o644))
	err := &ValidationError{
		Field:  `profile["release"].target["linux/amd64"].formats`,
		Detail: `format "aab" is not a production format for linux/amd64`,
		Range:  SourceRange{Filename: path, StartLine: 2, StartColumn: 47, EndLine: 2, EndColumn: 64},
	}

	formatted, ok := FormatValidationDiagnostics(err)
	require.True(t, ok)
	formatted = strings.ReplaceAll(formatted, path, "<manifest>") + "\n"
	want, readErr := os.ReadFile(filepath.Join("testdata", "semantic-diagnostic.golden"))
	require.NoError(t, readErr)
	assert.Equal(t, string(want), formatted)
}

func BenchmarkFormatValidationDiagnostics(b *testing.B) {
	err := &ValidationError{
		Field:  `profile["release"].target["linux/amd64"].formats`,
		Detail: `format "aab" is not a production format for linux/amd64`,
	}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = FormatValidationDiagnostics(err)
	}
}

func errorsJoinForTest(values ...error) error {
	return joinedValidationErrors(values)
}

type joinedValidationErrors []error

func (e joinedValidationErrors) Error() string   { return "joined" }
func (e joinedValidationErrors) Unwrap() []error { return []error(e) }
