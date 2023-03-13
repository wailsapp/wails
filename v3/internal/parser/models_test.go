package parser

import (
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateModels(t *testing.T) {

	tests := []struct {
		dir  string
		want string
	}{
		{
			"testdata/function_single",
			getFile("testdata/function_single/models.ts"),
		},
		{
			"testdata/function_from_imported_package",
			getFile("testdata/function_from_imported_package/models.ts"),
		},
		{
			"testdata/variable_single",
			getFile("testdata/variable_single/models.ts"),
		},
		{
			"testdata/variable_single_from_function",
			getFile("testdata/variable_single_from_function/models.ts"),
		},
		{
			"testdata/variable_single_from_other_function",
			getFile("testdata/variable_single_from_other_function/models.ts"),
		},
		{
			"testdata/struct_literal_single",
			getFile("testdata/struct_literal_single/models.ts"),
		},
		{
			"testdata/struct_literal_multiple",
			getFile("testdata/struct_literal_multiple/models.ts"),
		},
		{
			"testdata/struct_literal_multiple_other",
			getFile("testdata/struct_literal_multiple_other/models.ts"),
		},
		{
			"testdata/struct_literal_multiple_files",
			getFile("testdata/struct_literal_multiple_files/models.ts"),
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
			got, err := GenerateModels(project.Models)
			if err != nil {
				t.Fatalf("GenerateModels() error = %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
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
