package parser

import (
	"embed"
	"io/fs"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
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
		dir  string
		want map[string]string
	}{
		{
			"testdata/function_single",
			map[string]string{
				"main": getFile("testdata/function_single/bindings_main.js"),
			},
		},
		{
			"testdata/function_from_imported_package",
			map[string]string{
				"main":     getFile("testdata/function_from_imported_package/bindings_main.js"),
				"services": getFile("testdata/function_from_imported_package/bindings_services.js"),
			},
		},
		{
			"testdata/variable_single",
			map[string]string{
				"main": getFile("testdata/variable_single/bindings_main.js"),
			},
		},
		{
			"testdata/variable_single_from_function",
			map[string]string{
				"main": getFile("testdata/variable_single_from_function/bindings_main.js"),
			},
		},
		{
			"testdata/variable_single_from_other_function",
			map[string]string{
				"main":     getFile("testdata/variable_single_from_other_function/bindings_main.js"),
				"services": getFile("testdata/variable_single_from_other_function/bindings_services.js"),
			},
		},
		{
			"testdata/struct_literal_single",
			map[string]string{
				"main": getFile("testdata/struct_literal_single/bindings_main.js"),
			},
		},
		{
			"testdata/struct_literal_multiple",
			map[string]string{
				"main": getFile("testdata/struct_literal_multiple/bindings_main.js"),
			},
		},
		{
			"testdata/struct_literal_multiple_other",
			map[string]string{
				"main":     getFile("testdata/struct_literal_multiple_other/bindings_main.js"),
				"services": getFile("testdata/struct_literal_multiple_other/bindings_services.js"),
			},
		},
		{
			"testdata/struct_literal_multiple_files",
			map[string]string{
				"main": getFile("testdata/struct_literal_multiple_files/bindings_main.js"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.dir, func(t *testing.T) {
			// Run parser on directory
			project, err := ParseProject(tt.dir)
			if err != nil {
				t.Errorf("ParseProject() error = %v", err)
				return
			}

			// Generate Bindings
			got := GenerateBindings(project.BoundMethods)

			for name, binding := range got {
				// check if the binding is in the expected bindings
				expected, ok := tt.want[name]
				if !ok {
					err = os.WriteFile(tt.dir+"/bindings_"+name+".got.js", []byte(binding), 0644)
					if err != nil {
						t.Errorf("os.WriteFile() error = %v", err)
						return
					}
					t.Errorf("GenerateBindings() unexpected binding = %v", name)
					return
				}
				// compare the binding
				if diff := cmp.Diff(expected, binding); diff != "" {
					err = os.WriteFile(tt.dir+"/bindings_"+name+".got.js", []byte(binding), 0644)
					if err != nil {
						t.Errorf("os.WriteFile() error = %v", err)
						return
					}
					t.Fatalf("GenerateBindings() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
