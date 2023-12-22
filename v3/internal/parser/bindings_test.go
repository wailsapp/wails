package parser

import (
	"embed"
	"github.com/google/go-cmp/cmp"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

//go:embed testdata
var testdata embed.FS

func getFile(filename string) string {
	// get the file from the testdata FS
	file, err := fs.ReadFile(testdata, filename)
	if err != nil {
		panic(err)
	}
	return string(file)
}

func TestGenerateBindings(t *testing.T) {

	tests := []struct {
		dir    string
		want   map[string]map[string]string
		useIDs bool
	}{
		{
			dir: "testdata/enum",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/enum/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			dir: "testdata/enum",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/enum/frontend/bindings/main/GreetService.name.js"),
				},
			},
		},
		{
			dir: "testdata/enum_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/enum_from_imported_package/frontend/bindings/main/GreetService.name.js"),
				},
			},
		},
		{
			dir: "testdata/enum_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/enum_from_imported_package/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			dir: "testdata/function_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_single/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			dir: "testdata/function_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_single/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			dir: "testdata/function_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_imported_package/frontend/bindings/main/GreetService.name.js"),
				},
				"services": {
					"OtherService": getFile("testdata/function_from_imported_package/frontend/bindings/services/OtherService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			dir: "testdata/function_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_imported_package/frontend/bindings/main/GreetService.js"),
				},
				"services": {
					"OtherService": getFile("testdata/function_from_imported_package/frontend/bindings/services/OtherService.js"),
				},
			},
			useIDs: true,
		},
		{
			dir: "testdata/function_from_nested_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_nested_imported_package/frontend/bindings/main/GreetService.name.js"),
				},
				"services/other": {
					"OtherService": getFile("testdata/function_from_nested_imported_package/frontend/bindings/services/other/OtherService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			dir: "testdata/function_from_nested_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_nested_imported_package/frontend/bindings/main/GreetService.js"),
				},
				"services/other": {
					"OtherService": getFile("testdata/function_from_nested_imported_package/frontend/bindings/services/other/OtherService.js"),
				},
			},
			useIDs: true,
		},
		{
			dir: "testdata/struct_literal_multiple",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple/frontend/bindings/main/GreetService.js"),
					"OtherService": getFile("testdata/struct_literal_multiple/frontend/bindings/main/OtherService.js"),
				},
			},
			useIDs: true,
		},
		{
			dir: "testdata/struct_literal_multiple",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple/frontend/bindings/main/GreetService.name.js"),
					"OtherService": getFile("testdata/struct_literal_multiple/frontend/bindings/main/OtherService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			dir: "testdata/function_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_imported_package/frontend/bindings/main/GreetService.js"),
				},
				"services": {
					"OtherService": getFile("testdata/function_from_imported_package/frontend/bindings/services/OtherService.js"),
				},
			},
			useIDs: true,
		},
		{
			dir: "testdata/function_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_imported_package/frontend/bindings/main/GreetService.name.js"),
				},
				"services": {
					"OtherService": getFile("testdata/function_from_imported_package/frontend/bindings/services/OtherService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			dir: "testdata/variable_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			dir: "testdata/variable_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			dir: "testdata/variable_single_from_function",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single_from_function/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			dir: "testdata/variable_single_from_function",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single_from_function/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			dir: "testdata/variable_single_from_other_function",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single_from_other_function/frontend/bindings/main/GreetService.name.js"),
				},
				"services": {
					"OtherService": getFile("testdata/variable_single_from_other_function/frontend/bindings/services/OtherService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			dir: "testdata/variable_single_from_other_function",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single_from_other_function/frontend/bindings/main/GreetService.js"),
				},
				"services": {
					"OtherService": getFile("testdata/variable_single_from_other_function/frontend/bindings/services/OtherService.js"),
				},
			},
			useIDs: true,
		},
		{
			dir: "testdata/struct_literal_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_single/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			dir: "testdata/struct_literal_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_single/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			dir: "testdata/struct_literal_multiple_other",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple_other/frontend/bindings/main/GreetService.name.js"),
				},
				"services": {
					"OtherService": getFile("testdata/struct_literal_multiple_other/frontend/bindings/services/OtherService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			dir: "testdata/struct_literal_multiple_other",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple_other/frontend/bindings/main/GreetService.js"),
				},
				"services": {
					"OtherService": getFile("testdata/struct_literal_multiple_other/frontend/bindings/services/OtherService.js"),
				},
			},
			useIDs: true,
		},
		{
			dir: "testdata/struct_literal_non_pointer_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_non_pointer_single/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			dir: "testdata/struct_literal_non_pointer_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_non_pointer_single/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			dir: "testdata/struct_literal_multiple_files",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple_files/frontend/bindings/main/GreetService.name.js"),
					"OtherService": getFile("testdata/struct_literal_multiple_files/frontend/bindings/main/OtherService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			dir: "testdata/struct_literal_multiple_files",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple_files/frontend/bindings/main/GreetService.js"),
					"OtherService": getFile("testdata/struct_literal_multiple_files/frontend/bindings/main/OtherService.js"),
				},
			},
			useIDs: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.dir, func(t *testing.T) {
			// Run parser on directory
			absDir, err := filepath.Abs(tt.dir)
			if err != nil {
				t.Errorf("filepath.Abs() error = %v", err)
				return
			}
			project, err := ParseProject(absDir)
			if err != nil {
				t.Errorf("ParseProject() error = %v", err)
				return
			}

			project.outputDirectory = "frontend/bindings"

			// Generate Bindings
			got := project.GenerateBindings(project.BoundMethods, tt.useIDs)

			for dirName, structDetails := range got {
				// iterate the struct names in structDetails
				for name, binding := range structDetails {
					expected, ok := tt.want[dirName][name]
					if !ok {
						outFile := filepath.Join(tt.dir, project.outputDirectory, dirName, name+".got.js")
						err = os.WriteFile(outFile, []byte(binding), 0644)
						if err != nil {
							t.Errorf("os.WriteFile() error = %v", err)
							return
						}
						t.Errorf("GenerateBindings() unexpected binding = %v", name)
						return
					}
					// compare the binding

					// convert all line endings to \n
					binding = convertLineEndings(binding)
					expected = convertLineEndings(expected)

					if diff := cmp.Diff(expected, binding); diff != "" {
						outFile := filepath.Join(tt.dir, project.outputDirectory, dirName, name+".got.js")
						err = os.WriteFile(outFile, []byte(binding), 0644)
						if err != nil {
							t.Errorf("os.WriteFile() error = %v", err)
							return
						}
						t.Fatalf("GenerateBindings() mismatch (-want +got):\n%s", diff)
					}
				}
			}
		})
	}
}
