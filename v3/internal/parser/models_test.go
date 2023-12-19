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
		name          string
		dir           string
		want          string
		useInterface  bool
		useTypescript bool
	}{
		{
			name:          "function single",
			dir:           "testdata/function_single",
			useTypescript: true,
		},
		{
			name:          "function from imported package",
			dir:           "testdata/function_from_imported_package",
			want:          getFile("testdata/function_from_imported_package/models.ts"),
			useTypescript: true,
		},
		{
			name: "function from imported package (Javascript)",
			dir:  "testdata/function_from_imported_package",
			want: getFile("testdata/function_from_imported_package/models.js"),
		},
		{
			name:          "variable single",
			dir:           "testdata/variable_single",
			useTypescript: true,
		},
		{
			name:          "variable single from function",
			dir:           "testdata/variable_single_from_function",
			useTypescript: true,
		},
		{
			name:          "variable single from other function",
			dir:           "testdata/variable_single_from_other_function",
			want:          getFile("testdata/variable_single_from_other_function/models.ts"),
			useTypescript: true,
		},
		{
			name:          "struct literal single",
			dir:           "testdata/struct_literal_single",
			want:          getFile("testdata/struct_literal_single/models.ts"),
			useTypescript: true,
		},
		{
			name:          "struct literal multiple",
			dir:           "testdata/struct_literal_multiple",
			useTypescript: true,
		},
		{
			name:          "struct literal multiple other",
			dir:           "testdata/struct_literal_multiple_other",
			want:          getFile("testdata/struct_literal_multiple_other/models.ts"),
			useTypescript: true,
		},
		{
			name: "struct literal multiple other (Javascript)",
			dir:  "testdata/struct_literal_multiple_other",
			want: getFile("testdata/struct_literal_multiple_other/models.js"),
		},
		{
			name:          "struct literal non pointer single (Javascript)",
			dir:           "testdata/struct_literal_non_pointer_single",
			want:          getFile("testdata/struct_literal_non_pointer_single/models.ts"),
			useTypescript: true,
		},
		{
			name: "struct literal non pointer single (Javascript)",
			dir:  "testdata/struct_literal_non_pointer_single",
			want: getFile("testdata/struct_literal_non_pointer_single/models.js"),
		},
		{
			name:          "struct literal multiple files",
			dir:           "testdata/struct_literal_multiple_files",
			useTypescript: true,
		},
		{
			name:          "enum",
			dir:           "testdata/enum",
			want:          getFile("testdata/enum/models.ts"),
			useTypescript: true,
		},
		{
			name: "enum (Javascript)",
			dir:  "testdata/enum",
			want: getFile("testdata/enum/models.js"),
		},
		{
			name:          "enum interface",
			dir:           "testdata/enum-interface",
			want:          getFile("testdata/enum-interface/models.ts"),
			useInterface:  true,
			useTypescript: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run parser on directory
			project, err := ParseProject(tt.dir)
			if err != nil {
				t.Fatalf("ParseProject() error = %v", err)
			}

			// Generate Models
			got, err := GenerateModels(project.Models, project.Types, &flags.GenerateBindingsOptions{
				UseInterfaces: tt.useInterface,
				TS:            tt.useTypescript,
			})
			if err != nil {
				t.Fatalf("GenerateModels() error = %v", err)
			}
			// convert all line endings to \n
			got = convertLineEndings(got)
			want := convertLineEndings(tt.want)
			if diff := cmp.Diff(want, got); diff != "" {
				gotFilename := "models.got.js"
				if tt.useTypescript {
					gotFilename = "models.got.ts"
				}
				err = os.WriteFile(filepath.Join(tt.dir, gotFilename), []byte(got), 0644)
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
