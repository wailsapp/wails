package flags

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBuildFlags(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantFlags []string
		wantErr   bool
	}{
		{
			name:      "empty string",
			input:     "",
			wantFlags: nil,
		},
		{
			name:      "single flag, multiple spaces",
			input:     "  -v ",
			wantFlags: []string{"-v"},
		},
		{
			name:      "multiple flags, complex spaces",
			input:     "  \t-v\r\n-x",
			wantFlags: []string{"-v", "-x"},
		},
		{
			name:      "empty flag (single quotes)",
			input:     `''`,
			wantFlags: []string{""},
		},
		{
			name:      "empty flag (double quotes)",
			input:     `""`,
			wantFlags: []string{""},
		},
		{
			name:      "flag with spaces (single quotes)",
			input:     `'a 	b'`,
			wantFlags: []string{"a \tb"},
		},
		{
			name:      "flag with spaces (double quotes)",
			input:     `'a 	b'`,
			wantFlags: []string{"a \tb"},
		},
		{
			name:      "mixed quoted and non-quoted flags (single quotes)",
			input:     `-v 'a b '  -x`,
			wantFlags: []string{"-v", "a b ", "-x"},
		},
		{
			name:      "mixed quoted and non-quoted flags (double quotes)",
			input:     `-v "a b "  -x`,
			wantFlags: []string{"-v", "a b ", "-x"},
		},
		{
			name:      "mixed quoted and non-quoted flags (mixed quotes)",
			input:     `-v "a b "  '-x'`,
			wantFlags: []string{"-v", "a b ", "-x"},
		},
		{
			name:      "double quote within single quotes",
			input:     `' " '`,
			wantFlags: []string{" \" "},
		},
		{
			name:      "single quote within double quotes",
			input:     `" ' "`,
			wantFlags: []string{" ' "},
		},
		{
			name:      "unmatched single quote",
			input:     `-v "a b "  '-x -y`,
			wantFlags: []string{"-v", "a b ", "-x -y"},
			wantErr:   true,
		},
		{
			name:      "unmatched double quote",
			input:     `-v "a b "  "-x -y`,
			wantFlags: []string{"-v", "a b ", "-x -y"},
			wantErr:   true,
		},
		{
			name:      "mismatched single quote",
			input:     `-v "a b "  '-x" -y`,
			wantFlags: []string{"-v", "a b ", "-x\" -y"},
			wantErr:   true,
		},
		{
			name:      "mismatched double quote",
			input:     `-v "a b "  "-x' -y`,
			wantFlags: []string{"-v", "a b ", "-x' -y"},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := GenerateBindingsOptions{
				BuildFlagsString: tt.input,
			}

			var wantErr error = nil
			if tt.wantErr {
				wantErr = ErrUnmatchedQuote
			}

			gotFlags, gotErr := options.BuildFlags()

			if diff := cmp.Diff(tt.wantFlags, gotFlags); diff != "" {
				t.Errorf("BuildFlags() unexpected result: %s\n", diff)
			}

			if diff := cmp.Diff(wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("BuildFlags() unexpected error: %s\n", diff)
			}
		})
	}
}
