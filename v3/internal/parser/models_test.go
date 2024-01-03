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
		want          map[string]string
		useInterface  bool
		useTypescript bool
	}{
		// enum
		{
			name: "enum - Typescript",
			dir:  "testdata/enum",
			want: map[string]string{
				"main": getFile("testdata/enum/frontend/bindings/main/models.ts"),
			},
			useTypescript: true,
		},
		{
			name: "enum - Javascript",
			dir:  "testdata/enum",
			want: map[string]string{
				"main": getFile("testdata/enum/frontend/bindings/main/models.js"),
			},
			useTypescript: false,
		},
		{
			name: "enum - Typescript interfaces",
			dir:  "testdata/enum",
			want: map[string]string{
				"main": getFile("testdata/enum/frontend/bindings/main/models.interfaces.ts"),
			},
			useTypescript: true,
			useInterface:  true,
		},
		// function from imported package
		{
			name: "function from imported package - Typescript",
			dir:  "testdata/function_from_imported_package",
			want: map[string]string{
				"main":     getFile("testdata/function_from_imported_package/frontend/bindings/main/models.ts"),
				"services": getFile("testdata/function_from_imported_package/frontend/bindings/services/models.ts"),
			},
			useTypescript: true,
		},
		{
			name: "function from imported package - Typescript interfaces",
			dir:  "testdata/function_from_imported_package",
			want: map[string]string{
				"main":     getFile("testdata/function_from_imported_package/frontend/bindings/main/models.interfaces.ts"),
				"services": getFile("testdata/function_from_imported_package/frontend/bindings/services/models.interfaces.ts"),
			},
			useTypescript: true,
			useInterface:  true,
		},
		{
			name: "function from imported package - Javascript",
			dir:  "testdata/function_from_imported_package",
			want: map[string]string{
				"main":     getFile("testdata/function_from_imported_package/frontend/bindings/main/models.js"),
				"services": getFile("testdata/function_from_imported_package/frontend/bindings/services/models.js"),
			},
			useTypescript: false,
		},
		// variable single from other function
		{
			name: "variable single from other function - Typescript",
			dir:  "testdata/variable_single_from_other_function",
			want: map[string]string{
				"main":     getFile("testdata/variable_single_from_other_function/frontend/bindings/main/models.ts"),
				"services": getFile("testdata/variable_single_from_other_function/frontend/bindings/services/models.ts"),
			},
			useTypescript: true,
		},
		{
			name: "variable single from other function - Typescript interfaces",
			dir:  "testdata/variable_single_from_other_function",
			want: map[string]string{
				"main":     getFile("testdata/variable_single_from_other_function/frontend/bindings/main/models.interfaces.ts"),
				"services": getFile("testdata/variable_single_from_other_function/frontend/bindings/services/models.interfaces.ts"),
			},
			useTypescript: true,
			useInterface:  true,
		},
		{
			name: "variable single from other function - Javascript",
			dir:  "testdata/variable_single_from_other_function",
			want: map[string]string{
				"main":     getFile("testdata/variable_single_from_other_function/frontend/bindings/main/models.js"),
				"services": getFile("testdata/variable_single_from_other_function/frontend/bindings/services/models.js"),
			},
			useTypescript: false,
		},
		// struct literal single
		{
			name: "struct literal single - Typescript",
			dir:  "testdata/struct_literal_single",
			want: map[string]string{
				"main": getFile("testdata/struct_literal_single/frontend/bindings/main/models.ts"),
			},
			useTypescript: true,
		},
		{
			name: "struct literal single - Typescript interfaces",
			dir:  "testdata/struct_literal_single",
			want: map[string]string{
				"main": getFile("testdata/struct_literal_single/frontend/bindings/main/models.interfaces.ts"),
			},
			useTypescript: true,
			useInterface:  true,
		},
		{
			name: "struct literal single - Javascript",
			dir:  "testdata/struct_literal_single",
			want: map[string]string{
				"main": getFile("testdata/struct_literal_single/frontend/bindings/main/models.js"),
			},
			useTypescript: false,
		},
		// struct literal multiple other
		{
			name: "struct literal multiple other - Typescript",
			dir:  "testdata/struct_literal_multiple_other",
			want: map[string]string{
				"main":     getFile("testdata/struct_literal_multiple_other/frontend/bindings/main/models.ts"),
				"services": getFile("testdata/struct_literal_multiple_other/frontend/bindings/services/models.ts"),
			},
			useTypescript: true,
		},
		{
			name: "struct literal multiple other - Typescript interfaces",
			dir:  "testdata/struct_literal_multiple_other",
			want: map[string]string{
				"main":     getFile("testdata/struct_literal_multiple_other/frontend/bindings/main/models.interfaces.ts"),
				"services": getFile("testdata/struct_literal_multiple_other/frontend/bindings/services/models.interfaces.ts"),
			},
			useTypescript: true,
			useInterface:  true,
		},
		{
			name: "struct literal multiple other - Javascript",
			dir:  "testdata/struct_literal_multiple_other",
			want: map[string]string{
				"main":     getFile("testdata/struct_literal_multiple_other/frontend/bindings/main/models.js"),
				"services": getFile("testdata/struct_literal_multiple_other/frontend/bindings/services/models.js"),
			},
			useTypescript: false,
		},
		// struct literal non pointer single
		{
			name: "struct literal non pointer single - Typescript",
			dir:  "testdata/struct_literal_non_pointer_single",
			want: map[string]string{
				"main": getFile("testdata/struct_literal_non_pointer_single/frontend/bindings/main/models.ts"),
			},
			useTypescript: true,
		},
		{
			name: "struct literal non pointer single - Typescript interfaces",
			dir:  "testdata/struct_literal_non_pointer_single",
			want: map[string]string{
				"main": getFile("testdata/struct_literal_non_pointer_single/frontend/bindings/main/models.interfaces.ts"),
			},
			useTypescript: true,
			useInterface:  true,
		},
		{
			name: "struct literal non pointer single - Javascript",
			dir:  "testdata/struct_literal_non_pointer_single",
			want: map[string]string{
				"main": getFile("testdata/struct_literal_non_pointer_single/frontend/bindings/main/models.js"),
			},
			useTypescript: false,
		},
		// enum from imported package
		{
			name: "enum from imported package - Typescript",
			dir:  "testdata/enum_from_imported_package",
			want: map[string]string{
				"services": getFile("testdata/enum_from_imported_package/frontend/bindings/services/models.ts"),
			},
			useTypescript: true,
		},
		{
			name: "enum from imported package - Typescript interfaces",
			dir:  "testdata/enum_from_imported_package",
			want: map[string]string{
				"services": getFile("testdata/enum_from_imported_package/frontend/bindings/services/models.interfaces.ts"),
			},
			useTypescript: true,
		},
		{
			name: "enum from imported package - Javascript",
			dir:  "testdata/enum_from_imported_package",
			want: map[string]string{
				"services": getFile("testdata/enum_from_imported_package/frontend/bindings/services/models.js"),
			},
			useTypescript: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run parser on directory
			project, err := ParseProject(tt.dir)
			if err != nil {
				t.Fatalf("ParseProject() error = %v", err)
			}

			project.outputDirectory = "frontend/bindings"

			// Generate Models
			allModels, err := project.GenerateModels(project.Models, project.Types, &flags.GenerateBindingsOptions{
				UseInterfaces: tt.useInterface,
				TS:            tt.useTypescript,
			})
			if err != nil {
				t.Fatalf("GenerateModels() error = %v", err)
			}
			for pkgDir, got := range allModels {
				// convert all line endings to \n
				got = convertLineEndings(got)
				want, ok := tt.want[pkgDir]
				if !ok {
					t.Fatalf("GenerateModels() missing package: %s", pkgDir)
				}
				want = convertLineEndings(want)
				if diff := cmp.Diff(want, got); diff != "" {
					gotFilename := "models.got.js"
					if tt.useTypescript {
						gotFilename = "models.got.ts"
					}
					// Get relative package path
					//relativeBindingsDir := project.RelativeBindingsDir(project.packageCache[pkgDir])

					err = os.WriteFile(filepath.Join(tt.dir, project.outputDirectory, pkgDir, gotFilename), []byte(got), 0644)
					if err != nil {
						t.Errorf("os.WriteFile() error = %v", err)
						return
					}
					t.Fatalf("GenerateModels() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func convertLineEndings(str string) string {
	// replace all \r\n with \n
	return strings.ReplaceAll(str, "\r\n", "\n")
}
