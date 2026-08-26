package manifest

import (
	"fmt"
	"os"
	"path/filepath"
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

func errorsJoinForTest(values ...error) error {
	return joinedValidationErrors(values)
}

type joinedValidationErrors []error

func (e joinedValidationErrors) Error() string   { return "joined" }
func (e joinedValidationErrors) Unwrap() []error { return []error(e) }
