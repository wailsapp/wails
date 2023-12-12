package parser

import (
	"github.com/google/go-cmp/cmp"
	"github.com/wailsapp/wails/v3/internal/flags"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateModels(t *testing.T) {

	tests := []struct {
		dir  string
		want string
	}{
		{
			dir: "testdata/function_single",
		},
		{
			dir:  "testdata/function_from_imported_package",
			want: getFile("testdata/function_from_imported_package/models.ts"),
		},
		{
			dir: "testdata/variable_single",
		},
		{
			dir: "testdata/variable_single_from_function",
		},
		{
			dir:  "testdata/variable_single_from_other_function",
			want: getFile("testdata/variable_single_from_other_function/models.ts"),
		},
		{
			dir:  "testdata/struct_literal_single",
			want: getFile("testdata/struct_literal_single/models.ts"),
		},
		{
			dir: "testdata/struct_literal_multiple",
		},
		{
			dir:  "testdata/struct_literal_multiple_other",
			want: getFile("testdata/struct_literal_multiple_other/models.ts"),
		},
		{
			dir: "testdata/struct_literal_multiple_files",
		},
		{
			dir:  "testdata/enum",
			want: getFile("testdata/enum/models.ts"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.dir, func(t *testing.T) {
			// Run parser on directory
			project, err := ParseProject(tt.dir)
			if err != nil {
				t.Fatalf("ParseProject() error = %v", err)
			}

			// Generate Models
			got, err := GenerateModels(project.Models, project.Types, &flags.GenerateBindingsOptions{})
			if err != nil {
				t.Fatalf("GenerateModels() error = %v", err)
			}
			// convert all line endings to \n
			got = convertLineEndings(got)
			want := convertLineEndings(tt.want)
			if diff := cmp.Diff(want, got); diff != "" {
				err = os.WriteFile(filepath.Join(tt.dir, "models.got.ts"), []byte(got), 0644)
				if err != nil {
					t.Errorf("os.WriteFile() error = %v", err)
					return
				}
				t.Fatalf("GenerateModels() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func convertLineEndings(str string) string {
	// replace all \r\n with \n
	return strings.ReplaceAll(str, "\r\n", "\n")
}
